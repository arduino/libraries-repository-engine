package file

// SCCSFiles is a map of folder names used internally by source code control systems.
var SCCSFiles = map[string]bool{
	"CVS":  true,
	"RCS":  true,
	".git": true,
	".svn": true,
	".hg":  true,
	".bzr": true}

// IsSCCS returns whether the given string is a folder name used internally by source code control systems.
func IsSCCS(name string) bool {
	return SCCSFiles[name]
}
