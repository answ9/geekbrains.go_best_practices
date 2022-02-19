// Package files is represented by two structs: File and UniqueFiles.
package files

// Struct File contains Path and Name of a file.
// Path can be used for fast access of a definite file
// Name can be used mostly for printing it to user
type File struct {
	Path, Name string
}

// Use method NewFile(path, name string) to create a new object
func NewFile(path, name string) File {
	return File{path, name}
}
