package filesystem

import (
	"github.com/mholt/archiver"
)

// Archiver defines what this package has available.
type Archiver interface {
	Unarchive(source string, destination string) error
}

// Unarchiver has access to the methods of this package.
type Unarchiver struct{}

// Unarchive - unarchive the source.
func (a *Unarchiver) Unarchive(source string, destination string) error {
	return archiver.Unarchive(source, destination)
}
