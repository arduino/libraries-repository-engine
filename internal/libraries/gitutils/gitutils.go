// This file is part of libraries-repository-engine.
//
// Copyright 2021 ARDUINO SA (http://www.arduino.cc/)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

package gitutils

import (
	"sort"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// ResolveTag returns the commit hash associated with a tag.
func ResolveTag(tag *plumbing.Reference, repository *git.Repository) (*plumbing.Hash, error) {
	// Annotated tags have their own hash, different from the commit hash, so the tag must be resolved to get the has for
	// the associated commit.
	// Tags may point to any Git object. Although not common, this can include tree and blob objects in addition to commits.
	// Resolving non-commit objects results in an error.
	return repository.ResolveRevision(plumbing.Revision(tag.Hash().String()))
}

// SortedCommitTags returns the repository's commit object tags sorted by their chronological order in the current branch's history.
// Tags for commits not in the branch's history are returned in lexicographical order relative to their adjacent tags.
func SortedCommitTags(repository *git.Repository) ([]*plumbing.Reference, error) {
	/*
		Given a repository tag structure like so (I've omitted 1.0.3-1.0.9 as irrelevant):

			* HEAD -> main, tag: 1.0.11
			* tag: 1.0.10
			* tag: 1.0.2
			| * tag: 1.0.2-rc2, development-branch
			| * tag: 1.0.2-rc1
			|/
			* tag: 1.0.1
			* tag: 1.0.0

			The raw tags order is lexicographical:
			1.0.0
			1.0.1
			1.0.10
			1.0.2
			1.0.2-rc1
			1.0.2-rc2

			This order is not meaningful. More meaningful would be to order the tags according to the chronology of the
			branch's commit history:
			1.0.0
			1.0.1
			1.0.2
			1.0.10

			This leaves the question of how to handle tags from other branches, which is likely why a sensible sorting
			capability was not provided. However, even if the sorting of those tags is not optimal, a meaningful sort of the
			current branch's tags will be a significant improvement over the default behavior.
	*/

	headRef, err := repository.Head()
	if err != nil {
		return nil, err
	}

	headCommit, err := repository.CommitObject(headRef.Hash())
	if err != nil {
		return nil, err
	}

	commits := object.NewCommitIterCTime(headCommit, nil, nil) // Iterator for the head commit and parents in reverse chronological commit time order.
	commitMap := make(map[plumbing.Hash]int)                   // commitMap associates each hash with its chronological position in the branch history.
	var commitIndex int
	for { // Iterate over all commits.
		commit, err := commits.Next()
		if err != nil {
			// Reached end of commits
			break
		}
		commitMap[commit.Hash] = commitIndex

		commitIndex-- // Decrement to reflect reverse chronological order.
	}

	tags, err := repository.Tags() // Get an iterator of the refs of the repository's tags. These are returned in a useless lexicographical order (e.g, 1.0.10 < 1.0.2), so it's necessary to cross-reference them against the commits, which are in a meaningful order.

	type tagDataType struct {
		tag      *plumbing.Reference
		position int
	}
	var tagData []tagDataType
	associatedCommitIndex := commitIndex // Initialize to index of oldest commit in case the first tags aren't in the branch.
	var tagIndex int
	for { // Iterate over all tag refs.
		tag, err := tags.Next()
		if err != nil {
			// Reached end of tags
			break
		}

		// Annotated tags have their own hash, different from the commit hash, so tags must be resolved before
		// cross-referencing against the commit hashes.
		resolvedTag, err := ResolveTag(tag, repository)
		if err != nil {
			// Non-commit object tags are not included in the sorted list.
			continue
		}

		commitIndex, ok := commitMap[*resolvedTag]
		if ok {
			// There is a commit in the branch associated with the tag.
			associatedCommitIndex = commitIndex
		}

		tagData = append(
			tagData,
			tagDataType{
				tag:      tag,
				position: associatedCommitIndex*10000 + tagIndex, // Leave intervals between positions to allow the insertion of unassociated tags in the existing lexicographical order relative to the last associated tag.
			},
		)

		tagIndex++
	}

	// Sort the tags according to the branch's history where possible.
	sort.SliceStable(
		tagData,
		// "less" function
		func(thisIndex, otherIndex int) bool {
			return tagData[thisIndex].position < tagData[otherIndex].position
		},
	)

	var sortedTags []*plumbing.Reference
	for _, tagDatum := range tagData {
		sortedTags = append(sortedTags, tagDatum.tag)
	}

	return sortedTags, nil
}
