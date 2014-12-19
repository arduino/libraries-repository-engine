package config

// Configuration part
// TODO: Move into his own file
var gh_auth_token = "e282d0ab13c38a4303e65620aeab13c2beba3385"
var gh_user = "arlibs"
var gh_callbackUrl = "http://kungfu.bug.st:8088/github/event/"
var localGitFolder = "git/"
var librariesIndexFile = "db.json"

func GithubAuthToken() string {
	return gh_auth_token
}

func GithubUser() string {
	return gh_user
}

func GithubCallbackURL() string {
	return gh_callbackUrl
}

func LocalGitFolder() string {
	return localGitFolder
}

func LibraryDBFile() string {
	return librariesIndexFile
}
