// +build !windows

package commands

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/jarcoal/httpmock"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/chef/ci-studio-common/internal/pkg/filesystem"
	"github.com/chef/ci-studio-common/internal/pkg/install"
)

func TestRootE(t *testing.T) {
	var err error
	var commandBuf bytes.Buffer

	defaultRootCmdOpts := &rootCmdOptions{
		channel: "stable",
		version: globalHabValue,
	}

	fs = filesystem.NewMemFs()
	ciutils = install.DefaultInstall()

	cmd := &cobra.Command{}
	args := []string{}
	output := new(bytes.Buffer)
	cmd.SetOutput(output)

	t.Run("failed to create lockfile", func(t *testing.T) {
		fs = filesystem.NewMemFs()
		rootCmdOpts = defaultRootCmdOpts
		fslock = &filesystem.MemMapLock{
			Locks: make(map[string]*filesystem.MemLock),
			Err:   errors.New("some error"),
		}

		// Assertions!
		err = rootE(cmd, args)
		assert.NotNil(t, err)
		assert.Empty(t, output.String())
		assert.Regexp(t, "could not create .+/install-habitat.lock", err.Error())
		assert.Empty(t, commandBuf.Bytes())

		// Cleanup!
		commandBuf.Reset()
		output.Reset()
	})

	t.Run("failed to acquire lock (other error)", func(t *testing.T) {
		fs = filesystem.NewMemFs()
		rootCmdOpts = defaultRootCmdOpts
		fslock = &filesystem.MemMapLock{
			Locks:   make(map[string]*filesystem.MemLock),
			LockErr: errors.New("some error"),
		}

		// Assertions!
		err = rootE(cmd, args)
		assert.NotNil(t, err)
		assert.Empty(t, output.String())
		assert.Regexp(t, "could not acquire file lock", err.Error())
		assert.Empty(t, commandBuf.Bytes())

		// Cleanup!
		commandBuf.Reset()
		output.Reset()
	})

	t.Run("failed to acquire lock (install in progress)", func(t *testing.T) {
		fs = filesystem.NewMemFs()
		rootCmdOpts = defaultRootCmdOpts
		fslock = &filesystem.MemMapLock{
			Locks: make(map[string]*filesystem.MemLock),
		}

		lock, _ := fslock.GetLock(ciutils.LockPath("install-habitat"))
		_ = lock.Lock()

		// Assertions!
		err = rootE(cmd, args)
		assert.Nil(t, err)
		assert.Equal(t, "Chef Habitat install already in progress -- skipping\n", output.String())
		assert.Empty(t, commandBuf.Bytes())

		// Cleanup!
		commandBuf.Reset()
		output.Reset()
	})

	t.Run("failed to determine current version", func(t *testing.T) {
		fs = filesystem.NewMemFs()
		rootCmdOpts = defaultRootCmdOpts
		fslock = &filesystem.MemMapLock{
			Locks: make(map[string]*filesystem.MemLock),
		}

		execCommand = func(command string, args ...string) *exec.Cmd {
			cs := []string{"-test.run=TestHabitatVersionNotFound", "--", command}
			cs = append(cs, args...)
			cmd := exec.Command(os.Args[0], cs...)
			cmd.Env = []string{
				"GO_WANT_HELPER_PROCESS=1",
				fmt.Sprintf("COMMAND=%s %s", command, strings.Join(args, " ")),
			}

			if fmt.Sprintf("%s %s", command, strings.Join(args, " ")) != "hab --version" {
				if _, err := fmt.Fprintf(&commandBuf, fmt.Sprintf("%s %s", command, strings.Join(args, " "))); err != nil {
					panic(err)
				}
			}

			return cmd
		}
		defer func() { execCommand = exec.Command }()

		// Mock out the HTTP response
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", HabitatInstallScriptURL,
			func(req *http.Request) (*http.Response, error) {
				response := httpmock.NewStringResponse(200, "install-content")

				return response, nil
			})

		// Assertions!
		err = rootE(cmd, args)
		assert.Nil(t, err)

		expectedOutput := heredoc.Docf(`
			Going to install the x86_64-%s build of Chef Habitat %s
		`, runtime.GOOS, globalHabValue)
		assert.Equal(t, expectedOutput, output.String())

		expectedCommand := fmt.Sprintf("%s -v %s -c stable -t x86_64-%s", ciutils.SettingsPath("install-habitat.sh"), globalHabValue, runtime.GOOS)
		assert.Equal(t, expectedCommand, commandBuf.String())

		lock, _ := fslock.GetLock(ciutils.LockPath("install-habitat"))
		_, err = lock.TryLock()
		assert.Nil(t, err)

		// Cleanup!
		commandBuf.Reset()
		output.Reset()
	})

	t.Run("failed to download install script", func(t *testing.T) {
		fs = filesystem.NewMemFs()
		rootCmdOpts = defaultRootCmdOpts
		fslock = &filesystem.MemMapLock{
			Locks: make(map[string]*filesystem.MemLock),
		}

		execCommand = func(command string, args ...string) *exec.Cmd {
			cs := []string{"-test.run=TestHabitatVersionNotFound", "--", command}
			cs = append(cs, args...)
			cmd := exec.Command(os.Args[0], cs...)
			cmd.Env = []string{
				"GO_WANT_HELPER_PROCESS=1",
				fmt.Sprintf("COMMAND=%s %s", command, strings.Join(args, " ")),
			}

			if fmt.Sprintf("%s %s", command, strings.Join(args, " ")) != "hab --version" {
				if _, err := fmt.Fprintf(&commandBuf, fmt.Sprintf("%s %s", command, strings.Join(args, " "))); err != nil {
					panic(err)
				}
			}

			return cmd
		}
		defer func() { execCommand = exec.Command }()

		// Mock out the HTTP response
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", HabitatInstallScriptURL,
			func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("some error")
			})

		// Assertions!
		err = rootE(cmd, args)
		assert.NotNil(t, err)
		assert.Regexp(t, "failed to download Chef Habitat install file", err.Error())

		expectedOutput := heredoc.Docf(`
			Going to install the x86_64-%s build of Chef Habitat %s
		`, runtime.GOOS, globalHabValue)
		assert.Equal(t, expectedOutput, output.String())
		assert.Empty(t, commandBuf.String())

		lock, _ := fslock.GetLock(ciutils.LockPath("install-habitat"))
		_, err = lock.TryLock()
		assert.Nil(t, err)

		// Cleanup!
		commandBuf.Reset()
		output.Reset()
	})

	t.Run("current version == specified version (install up-to-date)", func(t *testing.T) {
		fs = filesystem.NewMemFs()
		rootCmdOpts = defaultRootCmdOpts
		fslock = &filesystem.MemMapLock{
			Locks: make(map[string]*filesystem.MemLock),
		}

		execCommand = func(command string, args ...string) *exec.Cmd {
			cs := []string{"-test.run=TestCurrentHabitatVersionInstalled", "--", command}
			cs = append(cs, args...)
			cmd := exec.Command(os.Args[0], cs...)
			cmd.Env = []string{
				"GO_WANT_HELPER_PROCESS=1",
				fmt.Sprintf("COMMAND=%s %s", command, strings.Join(args, " ")),
			}

			return cmd
		}
		defer func() { execCommand = exec.Command }()

		// Assertions!
		err = rootE(cmd, args)
		assert.Nil(t, err)

		expectedOutput := heredoc.Docf(`
			Chef Habitat is already up-to-date (%s)
		`, globalHabValue)
		assert.Equal(t, expectedOutput, output.String())

		assert.Empty(t, commandBuf.Bytes())

		// Cleanup!
		commandBuf.Reset()
		output.Reset()
	})

	t.Run("current version != specified version", func(t *testing.T) {
		fs = filesystem.NewMemFs()
		rootCmdOpts = defaultRootCmdOpts
		fslock = &filesystem.MemMapLock{
			Locks: make(map[string]*filesystem.MemLock),
		}

		execCommand = func(command string, args ...string) *exec.Cmd {
			cs := []string{"-test.run=TestOldHabitatVersionInstalled", "--", command}
			cs = append(cs, args...)
			cmd := exec.Command(os.Args[0], cs...)
			cmd.Env = []string{
				"GO_WANT_HELPER_PROCESS=1",
				fmt.Sprintf("COMMAND=%s %s", command, strings.Join(args, " ")),
			}

			if fmt.Sprintf("%s %s", command, strings.Join(args, " ")) != "hab --version" {
				if _, err := fmt.Fprintf(&commandBuf, fmt.Sprintf("%s %s", command, strings.Join(args, " "))); err != nil {
					panic(err)
				}
			}

			return cmd
		}
		defer func() { execCommand = exec.Command }()

		// Mock out the HTTP response
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", HabitatInstallScriptURL,
			func(req *http.Request) (*http.Response, error) {
				response := httpmock.NewStringResponse(200, "install-content")

				return response, nil
			})

		// Assertions!
		err = rootE(cmd, args)
		assert.Nil(t, err)

		expectedOutput := heredoc.Docf(`
			Going to install the x86_64-%s build of Chef Habitat %s
		`, runtime.GOOS, globalHabValue)
		assert.Equal(t, expectedOutput, output.String())

		expectedCommand := fmt.Sprintf("%s -v %s -c stable -t x86_64-%s", ciutils.SettingsPath("install-habitat.sh"), globalHabValue, runtime.GOOS)
		assert.Equal(t, expectedCommand, commandBuf.String())

		lock, _ := fslock.GetLock(ciutils.LockPath("install-habitat"))
		_, err = lock.TryLock()
		assert.Nil(t, err)

		// Cleanup!
		commandBuf.Reset()
		output.Reset()
	})

	t.Run("install script already exists", func(t *testing.T) {
		fs = filesystem.NewMemFs()
		rootCmdOpts = defaultRootCmdOpts
		fslock = &filesystem.MemMapLock{
			Locks: make(map[string]*filesystem.MemLock),
		}

		execCommand = func(command string, args ...string) *exec.Cmd {
			cs := []string{"-test.run=TestOldHabitatVersionInstalled", "--", command}
			cs = append(cs, args...)
			cmd := exec.Command(os.Args[0], cs...)
			cmd.Env = []string{
				"GO_WANT_HELPER_PROCESS=1",
				fmt.Sprintf("COMMAND=%s %s", command, strings.Join(args, " ")),
			}

			if fmt.Sprintf("%s %s", command, strings.Join(args, " ")) != "hab --version" {
				if _, err := fmt.Fprintf(&commandBuf, fmt.Sprintf("%s %s", command, strings.Join(args, " "))); err != nil {
					panic(err)
				}
			}

			return cmd
		}
		defer func() { execCommand = exec.Command }()

		err = fs.WriteFile(filepath.Join(ciutils.SettingsPath("install-habitat.sh")), []byte("install script content"), 0777)
		assert.Nil(t, err)

		// Assertions!
		err = rootE(cmd, args)
		expectedOutput := heredoc.Docf(`
			Going to install the x86_64-%s build of Chef Habitat %s
		`, runtime.GOOS, globalHabValue)
		assert.Nil(t, err)
		assert.Equal(t, expectedOutput, output.String())

		expectedCommand := fmt.Sprintf("%s -v %s -c stable -t x86_64-%s", ciutils.SettingsPath("install-habitat.sh"), globalHabValue, runtime.GOOS)
		assert.Equal(t, expectedCommand, commandBuf.String())

		lock, _ := fslock.GetLock(ciutils.LockPath("install-habitat"))
		_, err = lock.TryLock()
		assert.Nil(t, err)

		// Cleanup!
		commandBuf.Reset()
		output.Reset()
	})
}

func TestHabitatVersionNotFound(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	if os.Getenv("COMMAND") == "hab --version" {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func TestCurrentHabitatVersionInstalled(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	if os.Getenv("COMMAND") == "hab --version" {
		fmt.Printf("hab %s/20190101125959\n", globalHabValue)
		os.Exit(0)
	} else {
		os.Exit(0)
	}
}

func TestOldHabitatVersionInstalled(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	if os.Getenv("COMMAND") == "hab --version" {
		fmt.Print("hab 0.82.0/20190101125959\n")
		os.Exit(0)
	} else {
		os.Exit(0)
	}
}
