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
	"fmt"
	"testing"
	"time"

	"github.com/arduino/go-paths-helper"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveTag(t *testing.T) {
	// Prepare the file system for the test repository.
	repositoryPath, err := paths.TempDir().MkTempDir("gitutils-TestResolveTag-repo")
	require.Nil(t, err)

	// Create test repository.
	repository, err := git.PlainInit(repositoryPath.String(), false)
	require.Nil(t, err)

	testTables := []struct {
		objectTypeName string
		objectHash     plumbing.Hash
		annotated      bool
		errorAssertion assert.ErrorAssertionFunc
	}{
		{
			objectTypeName: "Commit",
			objectHash:     makeCommit(t, repository, repositoryPath),
			errorAssertion: assert.NoError,
		},
		{
			objectTypeName: "Tree",
			objectHash:     getTreeHash(t, repository),
			errorAssertion: assert.Error,
		},
		{
			objectTypeName: "Blob",
			objectHash:     getBlobHash(t, repository),
			errorAssertion: assert.Error,
		},
	}

	for _, testTable := range testTables {
		for _, annotationConfig := range []struct {
			annotated  bool
			descriptor string
		}{
			{
				annotated:  true,
				descriptor: "Annotated",
			},
			{
				annotated:  false,
				descriptor: "Lightweight",
			},
		} {
			testName := fmt.Sprintf("%s, %s", testTable.objectTypeName, annotationConfig.descriptor)
			tag := makeTag(t, repository, testName, testTable.objectHash, annotationConfig.annotated)
			resolvedTag, err := ResolveTag(tag, repository)
			testTable.errorAssertion(t, err, fmt.Sprintf("%s tag resolution error", testName))
			if err == nil {
				assert.Equal(t, testTable.objectHash, *resolvedTag, fmt.Sprintf("%s tag resolution", testName))
			}
		}
	}
}

// makeCommit creates a test commit in the given repository and returns its plumbing.Hash object.
func makeCommit(t *testing.T, repository *git.Repository, repositoryPath *paths.Path) plumbing.Hash {
	_, err := paths.WriteToTempFile([]byte{}, repositoryPath, "gitutils-makeCommit-tempfile")
	require.Nil(t, err)

	worktree, err := repository.Worktree()
	require.Nil(t, err)
	_, err = worktree.Add(".")
	require.Nil(t, err)

	signature := &object.Signature{
		Name:  "Jane Developer",
		Email: "janedeveloper@example.com",
		When:  time.Now(),
	}

	commit, err := worktree.Commit(
		"Test commit message",
		&git.CommitOptions{
			Author: signature,
		},
	)
	require.Nil(t, err)

	return commit
}

// getTreeHash returns the plumbing.Hash object for an arbitrary Git tree object.
func getTreeHash(t *testing.T, repository *git.Repository) plumbing.Hash {
	trees, err := repository.TreeObjects()
	require.Nil(t, err)
	tree, err := trees.Next()
	require.Nil(t, err)
	return tree.ID()
}

// getTreeHash returns the plumbing.Hash object for an arbitrary Git blob object.
func getBlobHash(t *testing.T, repository *git.Repository) plumbing.Hash {
	blobs, err := repository.BlobObjects()
	require.Nil(t, err)
	blob, err := blobs.Next()
	require.Nil(t, err)
	return blob.ID()
}

// makeTag creates a Git tag in the given repository and returns its *plumbing.Reference object.
func makeTag(t *testing.T, repository *git.Repository, name string, hash plumbing.Hash, annotated bool) *plumbing.Reference {
	var tag *plumbing.Reference
	var err error
	if annotated {
		signature := &object.Signature{
			Name:  "Jane Developer",
			Email: "janedeveloper@example.com",
			When:  time.Now(),
		}

		tag, err = repository.CreateTag(
			name,
			hash,
			&git.CreateTagOptions{
				Tagger:  signature,
				Message: name,
			},
		)
	} else {
		tag, err = repository.CreateTag(name, hash, nil)
	}
	require.Nil(t, err)

	return tag
}