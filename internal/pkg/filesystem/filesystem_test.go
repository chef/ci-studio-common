package filesystem

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestAppendIfMissing(t *testing.T) {
	aferoFs := afero.NewMemMapFs()
	fs := &Fs{AferoFs: aferoFs}
	newContent := []byte("new content!")

	tempFile, err := afero.TempFile(aferoFs, "", "append-test")
	assert.Nil(t, err)

	t.Run("new file", func(t *testing.T) {
		err := fs.AppendIfMissing("/tmp/newfile", newContent, 0644)
		assert.Nil(t, err)

		contents, err := afero.ReadFile(aferoFs, "/tmp/newfile")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, fmt.Sprintf("%s\n", newContent), string(contents))
	})

	t.Run("existing file, no match", func(t *testing.T) {
		err := afero.WriteFile(aferoFs, tempFile.Name(), []byte("some existing content\n"), 0644)
		if err != nil {
			t.Fatal(err)
		}

		err = fs.AppendIfMissing(tempFile.Name(), newContent, 0644)
		assert.Nil(t, err)

		expected := []byte("some existing content\nnew content!\n")

		contents, err := afero.ReadFile(aferoFs, tempFile.Name())
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, expected, contents)
	})

	t.Run("existing file, match", func(t *testing.T) {
		err := afero.WriteFile(aferoFs, tempFile.Name(), []byte("some existing new content!\n"), 0644)
		if err != nil {
			t.Fatal(err)
		}

		err = fs.AppendIfMissing(tempFile.Name(), newContent, 0644)
		assert.Nil(t, err)

		expected := []byte("some existing new content!\n")

		contents, err := afero.ReadFile(aferoFs, tempFile.Name())
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, expected, contents)
	})
}

func TestCopyFile(t *testing.T) {
	aferoFs := afero.NewMemMapFs()
	fs := &Fs{AferoFs: aferoFs}
	fileContent := []byte("I am some content in a file")

	srcFile, err := afero.TempFile(aferoFs, "", "copy-test-src")
	assert.Nil(t, err)

	dstFile, err := afero.TempFile(aferoFs, "", "copy-test-dst")
	assert.Nil(t, err)

	err = afero.WriteFile(aferoFs, srcFile.Name(), fileContent, 0644)
	assert.Nil(t, err)

	err = fs.CopyFile(srcFile.Name(), dstFile.Name())
	assert.Nil(t, err)

	contents, err := afero.ReadFile(aferoFs, dstFile.Name())
	assert.Nil(t, err)
	assert.Equal(t, fileContent, contents)
}

func TestCopyDir(t *testing.T) {
	aferoFs := afero.NewMemMapFs()
	fs := &Fs{AferoFs: aferoFs}
	fileContent := []byte("I am some content in a file")

	var err error
	type pathSpec struct {
		path string
		dir  bool
	}

	paths := []*pathSpec{
		{"/file1", false},
		{"/dir1", true},
		{"/dir1/file1", false},
		{"/dir1/dir2", true},
		{"/dir1/dir2/file2", false},
	}

	for _, p := range paths {
		if p.dir {
			if err = aferoFs.MkdirAll(p.path, 0755); err != nil {
				t.Error(err.Error())
			}
		} else {
			if err = afero.WriteFile(aferoFs, p.path, fileContent, 0644); err != nil {
				t.Error(err.Error())
			}
		}
	}

	destDirPrefix := "/subdir1"
	for _, p := range paths {
		if destExists, err := afero.Exists(aferoFs, filepath.Join(destDirPrefix, p.path)); err != nil {
			t.Error(err.Error())
		} else if destExists {
			t.Errorf("Dest path '%s' should not exist", p.path)
		}
	}

	if err = fs.CopyDir("/", destDirPrefix); err != nil {
		t.Error(err.Error())
	}

	for _, p := range paths {
		if destExists, err := afero.Exists(aferoFs, filepath.Join(destDirPrefix, p.path)); err != nil {
			t.Error(err.Error())
		} else if !destExists {
			t.Errorf("Dest path '%s' should exist", p.path)
		}
	}
}

func TestDownloadRemoteFile(t *testing.T) {
	aferoFs := afero.NewMemMapFs()
	fs := &Fs{AferoFs: aferoFs}
	downloadContent := "I am some content you download from the internet"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, downloadContent)
	}))
	defer ts.Close()

	tempFile, err := afero.TempFile(aferoFs, "", "http-download-test")
	assert.Nil(t, err)

	err = fs.DownloadRemoteFile(ts.URL, tempFile.Name())
	assert.Nil(t, err)

	expected := fmt.Sprintf("%s\n", downloadContent)
	fileContent, err := afero.ReadFile(aferoFs, tempFile.Name())
	assert.Nil(t, err)
	assert.Equal(t, expected, string(fileContent))
}
