package filesystem

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// FileSystem is an interface that wraps a subset of afero calls but also adds a few of our own
type FileSystem interface {
	AppendIfMissing(name string, content []byte, mode os.FileMode) error
	Chmod(name string, mode os.FileMode) error
	CopyDir(src string, dst string) error
	CopyFile(src string, dst string) error
	Chtimes(name string, atime time.Time, mtime time.Time) error
	Create(name string) (afero.File, error)
	DownloadRemoteFile(url string, name string) error
	Exists(path string) (bool, error)
	FileContainsBytes(filename string, subslice []byte) (bool, error)
	MkdirAll(path string, perm os.FileMode) error
	ReadFile(filename string) ([]byte, error)
	Rename(oldname string, newname string) error
	Remove(name string) error
	RemoveAll(path string) error
	Stat(path string) (os.FileInfo, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
}

// A wrapper around the real afero.Fs we want
type Fs struct {
	AferoFs afero.Fs
}

func NewOsFs() *Fs {
	return &Fs{
		AferoFs: afero.NewOsFs(),
	}
}

func NewMemFs() *Fs {
	return &Fs{
		AferoFs: afero.NewMemMapFs(),
	}
}

// AppendIfMissing appends the content string to the given file
func (f *Fs) AppendIfMissing(filePath string, content []byte, mode os.FileMode) error {
	exists, err := f.Exists(filePath)
	if err != nil {
		return errors.Wrapf(err, "coult not confirm that file %s exists", filePath)
	}

	contentToWrite := append(content, []byte("\n")...)

	if !exists {
		err := f.WriteFile(filePath, contentToWrite, mode)
		if err != nil {
			return errors.Wrapf(err, "could not create file %s", filePath)
		}

		return nil
	}

	contains, err := f.FileContainsBytes(filePath, content)
	if err != nil {
		return errors.Wrapf(err, "could not check file %s for content", filePath)
	}

	if !contains {
		fileContents, err := f.ReadFile(filePath)
		if err != nil {
			return errors.Wrapf(err, "failed to read file %s", filePath)
		}

		newContents := append(fileContents, contentToWrite...)
		err = f.WriteFile(filePath, newContents, mode)
		if err != nil {
			return errors.Wrapf(err, "failed to append contents to file %s", filePath)
		}
	}

	return nil
}

// Chmod modified the file permissions
func (f *Fs) Chmod(name string, mode os.FileMode) error {
	return f.AferoFs.Chmod(name, mode)
}

// Chtimes changes the access and modification times of the named file
func (f *Fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return f.AferoFs.Chtimes(name, atime, mtime)
}

// CopyDir copies the contents of the given directory to another directory on the same filesystem
func (f *Fs) CopyDir(src string, dst string) error {
	srcInfo, err := f.AferoFs.Stat(src)
	if err != nil {
		return err
	}

	dir, err := f.AferoFs.Open(src)
	if err != nil {
		return err
	}
	defer dir.Close()

	entries, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	dstExists, err := f.Exists(dst)
	if err != nil {
		return err
	}

	if !dstExists {
		if err = f.MkdirAll(dst, srcInfo.Mode()); err != nil {
			return err
		}
	}

	for _, e := range entries {
		srcFullPath := filepath.Join(src, e.Name())
		dstFullPath := filepath.Join(dst, e.Name())

		if e.IsDir() {
			if err = f.CopyDir(srcFullPath, dstFullPath); err != nil {
				return err
			}
		} else {
			if err = f.CopyFile(srcFullPath, dstFullPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// CopyFile copies the specified file to the the given destination
func (f *Fs) CopyFile(src string, dst string) error {
	srcFile, err := f.AferoFs.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	dstFile, err := f.AferoFs.OpenFile(dst, os.O_RDWR|os.O_CREATE, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}

// Create creates a file
func (f *Fs) Create(name string) (afero.File, error) {
	return f.AferoFs.Create(name)
}

// DownloadRemoteFile downloads file from the internet onto disk
func (f *Fs) DownloadRemoteFile(url string, name string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := f.AferoFs.Create(name)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	return err
}

// Exists returns whether or not the file exists
func (f *Fs) Exists(path string) (bool, error) {
	return afero.Exists(f.AferoFs, path)
}

// FileContainsBytes returns whether or not the given file contains the subslice, otherwise an error
func (f *Fs) FileContainsBytes(filename string, subslice []byte) (bool, error) {
	return afero.FileContainsBytes(f.AferoFs, filename, subslice)
}

// MkdirAll creates a directory path and all parents that does not exist yet.
func (f *Fs) MkdirAll(path string, perm os.FileMode) error {
	return f.AferoFs.MkdirAll(path, perm)
}

// ReadFile returns the contents of the file as a slice, otherwise error
func (f *Fs) ReadFile(filename string) ([]byte, error) {
	return afero.ReadFile(f.AferoFs, filename)
}

// Rename returns an error if there was an issue renaming the given path
func (f *Fs) Rename(oldname string, newname string) error {
	return f.AferoFs.Rename(oldname, newname)
}

// Remove removes a file identified by name, returning an error, if any happens.
func (f *Fs) Remove(name string) error {
	return f.AferoFs.Remove(name)
}

// RemoveAll removes a directory path and any children it contains. It
// does not fail if the path does not exist (return nil).
func (f *Fs) RemoveAll(path string) error {
	return f.AferoFs.RemoveAll(path)
}

// Stat returns a FileInfo describing the named file, or an error, if any
// happens.
func (f *Fs) Stat(path string) (os.FileInfo, error) {
	return f.AferoFs.Stat(path)
}

// WriteFile writes the byte slice to the givn file
func (f *Fs) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return afero.WriteFile(f.AferoFs, filename, data, perm)
}
