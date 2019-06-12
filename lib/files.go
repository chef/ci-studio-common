package lib

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// AppendIfMissing appends the content string to the given file
func AppendIfMissing(filePath string, content string) error {
	fileToModify, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	Check(err)

	fileContents, err := ioutil.ReadFile(fileToModify.Name())
	Check(err)

	if !strings.Contains(string(fileContents), content) {
		_, err = fileToModify.Write([]byte(content))
	}

	fileToModify.Close()

	return err
}

// RenameDir renames a directory while handling cross-link issues
// credit https://github.com/rawlingsj/jx/blob/master/pkg/util/files.go
func RenameDir(src string, dst string, force bool) (err error) {
	err = CopyDir(src, dst, force)
	if err != nil {
		return fmt.Errorf("failed to copy source dir %s to %s: %s", src, dst, err)
	}
	err = os.RemoveAll(src)
	if err != nil {
		return fmt.Errorf("failed to cleanup source dir %s: %s", src, err)
	}
	return nil
}

// RenameFile renames a file while handling cross-link issues
// credit https://github.com/rawlingsj/jx/blob/master/pkg/util/files.go
func RenameFile(src string, dst string) (err error) {
	if src == dst {
		return nil
	}
	err = CopyFile(src, dst)
	if err != nil {
		return fmt.Errorf("failed to copy source file %s to %s: %s", src, dst, err)
	}
	err = os.RemoveAll(src)
	if err != nil {
		return fmt.Errorf("failed to cleanup source file %s: %s", src, err)
	}
	return nil
}

// CopyDir copies the entire contents of a directory
// credit https://gist.github.com/r0l1/92462b38df26839a3ca324697c8cba04
func CopyDir(src string, dst string, force bool) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)
	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}
	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		if force {
			os.RemoveAll(dst)
		} else {
			return fmt.Errorf("destination already exists")
		}
	}
	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath, force)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}
			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}
	return
}

// CopyFile copies file while preseving permissions
// credit https://gist.github.com/r0l1/92462b38df26839a3ca324697c8cba04
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()
	_, err = io.Copy(out, in)
	if err != nil {
		return
	}
	err = out.Sync()
	if err != nil {
		return
	}
	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}
	return
}

// FileExists returns whether or not the given file exists
func FileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	resp.Body.Close()
	out.Close()

	return err
}

// FindAndReplace finds the given regular expression in the given file and replaces it with the new string
func FindAndReplace(filePath string, regexStr string, newString string) error {
	r := regexp.MustCompile(regexStr)

	fileContents, err := ioutil.ReadFile(filePath)
	Check(err)

	newContents := r.ReplaceAll(fileContents, []byte(newString))

	err = ioutil.WriteFile(filePath, newContents, 0644)

	return err
}
