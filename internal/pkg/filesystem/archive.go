package filesystem

import (
	"github.com/mholt/archiver"
)

type Archiver interface {
	Unarchive(source string, destination string) error
}

type Unarchiver struct{}

func (a *Unarchiver) Unarchive(source string, destination string) error {
	return archiver.Unarchive(source, destination)
}
