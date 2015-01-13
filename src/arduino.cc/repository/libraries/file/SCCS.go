package file

var SCCSFiles map[string]bool = map[string]bool{"CVS":true, "RCS":true, ".git":true, ".svn":true, ".hg":true, ".bzr":true}

func IsSCCS(name string) bool {
	return SCCSFiles[name]
}
