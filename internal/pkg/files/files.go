package files

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// AppendIfMissing appends the content string to the given file
func AppendIfMissing(filePath string, content string) error {
	fileToModify, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to open file %s", filePath)
	}
	defer fileToModify.Close()

	fileContents, err := ioutil.ReadFile(fileToModify.Name())
	if err != nil {
		return errors.Wrapf(err, "failed to read file %s", fileToModify.Name())
	}

	if !strings.Contains(string(fileContents), content) {
		_, err = fileToModify.WriteString(fmt.Sprintf("%s\n", content))

		if err != nil {
			return errors.Wrapf(err, "failed to write contents to file %s", fileToModify.Name())
		}
	}

	return nil
}

// RenameFile renames a file while handling cross-link issues
func RenameFile(src string, dst string) error {
	if src == dst {
		return nil
	}

	dstDir := filepath.Dir(dst)
	dstFile := filepath.Base(dst)
	dstTmpFile := filepath.Join(dstDir, fmt.Sprintf("%s.tmp", dstFile))

	err := CopyFile(src, dstTmpFile)
	if err != nil {
		return errors.Wrapf(err, "failed to copy %s to %s as part of file rename", src, dstTmpFile)
	}

	err = os.Rename(dstTmpFile, dst)
	if err != nil {
		return errors.Wrapf(err, "failed to rename %s to %s as part of file rename", dstTmpFile, dst)
	}

	err = os.RemoveAll(src)
	if err != nil {
		return errors.Wrapf(err, "failed to remove %s as part of file rename", src)
	}

	return nil
}

// CopyDir copies the entire contents of a directory
// credit https://gist.github.com/r0l1/92462b38df26839a3ca324697c8cba04
func CopyDir(src string, dst string, force bool) error {
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
		return err
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
		return err
	}
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath, force)
			if err != nil {
				return err
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}
			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CopyFile copies file while preseving permissions
// credit https://gist.github.com/r0l1/92462b38df26839a3ca324697c8cba04
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	err = out.Sync()
	if err != nil {
		return err
	}
	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return err
	}
	return nil
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
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	return err
}

// FindAndReplace finds the given regular expression in the given file and replaces it with the new string
func FindAndReplace(filePath string, regexStr string, newString string) error {
	r := regexp.MustCompile(regexStr)

	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return errors.Wrapf(err, "failed to read contents of file %s", filePath)
	}

	newContents := r.ReplaceAll(fileContents, []byte(newString))

	err = ioutil.WriteFile(filePath, newContents, 0644)

	return err
}
