package commands

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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

const (
	remoteEtag     = "abcdef0123456789"
	archiveContent = "I am really a tarball"
)

var (
	// Buffer that captures output of execCommand
	commandBuf bytes.Buffer

	// Buffer that captures content of archiver
	archiverBuf bytes.Buffer

	// The command that will be run next
	expectedNextCmd string
)

func TestUpdate(t *testing.T) {
	var err error

	defaultUpdateCmdOpts := &updateCmdOptions{
		channel: "stable",
		force:   false,
		phase:   1,
	}

	ciutils = install.DefaultInstall()

	cmd := &cobra.Command{}
	args := []string{}
	output := new(bytes.Buffer)
	cmd.SetOutput(output)

	t.Run("failed to get lockfile", func(t *testing.T) {
		fs = filesystem.NewMemFs()
		updateCmdOpts = defaultUpdateCmdOpts
		fslock = &filesystem.MemMapLock{
			Locks: make(map[string]*filesystem.MemLock),
			Err:   errors.New("some error"),
		}

		// Assertions!
		err = updateE(cmd, args)
		assert.NotNil(t, err)
		assert.Empty(t, output.String())
		assert.Regexp(t, "could not create lockfile", err.Error())
		assert.Empty(t, commandBuf.Bytes())

		// Cleanup!
		commandBuf.Reset()
		output.Reset()
	})

	t.Run("error occurred trying to acquire lock", func(t *testing.T) {
		fs = filesystem.NewMemFs()
		updateCmdOpts = defaultUpdateCmdOpts
		fslock = &filesystem.MemMapLock{
			Locks:   make(map[string]*filesystem.MemLock),
			LockErr: errors.New("some error"),
		}

		// Assertions!
		err = updateE(cmd, args)
		assert.NotNil(t, err)
		assert.Empty(t, output.String())
		assert.Regexp(t, "could not acquire file lock", err.Error())
		assert.Empty(t, commandBuf.Bytes())

		// Cleanup!
		commandBuf.Reset()
		output.Reset()
	})

	t.Run("could not acquire lock (install already in progress)", func(t *testing.T) {
		fs = filesystem.NewMemFs()
		updateCmdOpts = defaultUpdateCmdOpts
		fslock = &filesystem.MemMapLock{
			Locks: make(map[string]*filesystem.MemLock),
		}

		lock, _ := fslock.GetLock(ciutils.LockPath("upgrade-ci-utils"))
		_ = lock.Lock()

		// Assertions!
		err = updateE(cmd, args)
		assert.Nil(t, err)
		assert.Equal(t, "ci-utils upgrade already in progress -- skipping\n", output.String())
		assert.Empty(t, commandBuf.Bytes())

		// Cleanup!
		commandBuf.Reset()
		output.Reset()
	})

	t.Run("an unsupported phase is given", func(t *testing.T) {
		fs = filesystem.NewMemFs()
		updateCmdOpts = &updateCmdOptions{
			phase: 4,
		}

		// Assertions!
		err = updateE(cmd, args)
		assert.NotNil(t, err)
		assert.Regexp(t, "4 is an unsupported phase", err.Error())
		assert.Empty(t, output.String())
		assert.Empty(t, commandBuf.Bytes())

		// Cleanup!
		commandBuf.Reset()
		output.Reset()
	})
}

func TestUpdate_phaseOne(t *testing.T) {
	var err error

	defaultUpdateCmdOpts := &updateCmdOptions{
		channel: "stable",
		force:   false,
		phase:   1,
	}

	ciutils = install.DefaultInstall()

	cmd := &cobra.Command{}
	args := []string{}
	output := new(bytes.Buffer)
	cmd.SetOutput(output)

	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()

	t.Run("failed to get HEAD for remote asset", func(t *testing.T) {
		fs = filesystem.NewMemFs()
		updateCmdOpts = defaultUpdateCmdOpts
		fslock = &filesystem.MemMapLock{
			Locks: make(map[string]*filesystem.MemLock),
		}

		// Mock out the HTTP response
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("HEAD", `=~^https://packages.chef.io/files/stable/ci-utils/latest/.+\z`,
			func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("some error")
			})

		// Assertions!
		err = updateE(cmd, args)
		assert.NotNil(t, err)
		assert.Regexp(t, "unable to download URL", err.Error())
		assert.Empty(t, output.Bytes())
		assert.Empty(t, commandBuf.Bytes())
		assertUpgradeStateReleased(t, ciutils, fs, fslock)

		// Cleanup!
		commandBuf.Reset()
		output.Reset()
	})

	t.Run("non-200 remote asset response", func(t *testing.T) {
		fs = filesystem.NewMemFs()
		updateCmdOpts = defaultUpdateCmdOpts
		fslock = &filesystem.MemMapLock{
			Locks: make(map[string]*filesystem.MemLock),
		}

		// Mock out the HTTP response
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("HEAD", `=~^https://packages.chef.io/files/stable/ci-utils/latest/.+\z`,
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(404, ""), nil
			})

		// Assertions!
		err := updateE(cmd, args)
		assert.NotNil(t, err)
		assert.Regexp(t, "unable to download URL", err.Error())
		assert.Empty(t, output.Bytes())
		assert.Empty(t, commandBuf.Bytes())
		assertUpgradeStateReleased(t, ciutils, fs, fslock)

		// Cleanup!
		commandBuf.Reset()
		output.Reset()
	})

	t.Run("failed to download the remote asset", func(t *testing.T) {
		fs = filesystem.NewMemFs()
		updateCmdOpts = defaultUpdateCmdOpts
		fslock = &filesystem.MemMapLock{
			Locks: make(map[string]*filesystem.MemLock),
		}

		// Mock out the HTTP response
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("HEAD", `=~^https://packages.chef.io/files/stable/ci-utils/latest/.+\z`,
			func(req *http.Request) (*http.Response, error) {
				response := httpmock.NewStringResponse(200, "")
				response.Header.Set("etag", remoteEtag)

				return response, nil
			})

		// Force the error by returning an error on the GET
		httpmock.RegisterResponder("GET", `=~^https://packages.chef.io/files/stable/ci-utils/latest/.+\z`,
			func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("some error")
			})

		// Assertions!
		err := updateE(cmd, args)
		assert.NotNil(t, err)
		assert.Regexp(t, "failed to download asset tarball", err.Error())
		expectedOutput := heredoc.Doc(`
			--> {1} Upgrade available
			--> {1} Downloading latest release
		`)
		assert.Equal(t, expectedOutput, output.String())
		assert.Empty(t, commandBuf.Bytes())
		assertUpgradeStateReleased(t, ciutils, fs, fslock)

		// Cleanup!
		commandBuf.Reset()
		output.Reset()
	})

	t.Run("failed to backup the current installation", func(t *testing.T) {
		fs = filesystem.NewMemFs()
		updateCmdOpts = defaultUpdateCmdOpts
		fslock = &filesystem.MemMapLock{
			Locks: make(map[string]*filesystem.MemLock),
		}

		// Mock out the HTTP response
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("HEAD", `=~^https://packages.chef.io/files/stable/ci-utils/latest/.+\z`,
			func(req *http.Request) (*http.Response, error) {
				response := httpmock.NewStringResponse(200, "")
				response.Header.Set("etag", remoteEtag)

				return response, nil
			})

		httpmock.RegisterResponder("GET", `=~^https://packages.chef.io/files/stable/ci-utils/latest/.+\z`,
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(200, archiveContent), nil
			})

		// Assertions!
		err := updateE(cmd, args)
		assert.NotNil(t, err)
		assert.Regexp(t, "failed to backup existing installation", err.Error())
		expectedOutput := heredoc.Doc(`
			--> {1} Upgrade available
			--> {1} Downloading latest release
			--> {1} Release downloaded -- backing up current installation
		`)
		assert.Equal(t, expectedOutput, output.String())
		assert.Empty(t, commandBuf.Bytes())
		assertUpgradeStateReleased(t, ciutils, fs, fslock)

		// Cleanup!
		commandBuf.Reset()
		output.Reset()
	})

	t.Run("remote etag == local etag (install up-to-date)", func(t *testing.T) {
		fs = filesystem.NewMemFs()
		updateCmdOpts = defaultUpdateCmdOpts
		fslock = &filesystem.MemMapLock{
			Locks: make(map[string]*filesystem.MemLock),
		}

		// Mock out the on-disk etag
		err = fs.WriteFile(ciutils.SettingsPath("etag"), []byte(remoteEtag), 0644)
		assert.Nil(t, err)

		// Mock out the HTTP response
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("HEAD", `=~^https://packages.chef.io/files/stable/ci-utils/latest/.+\z`,
			func(req *http.Request) (*http.Response, error) {
				response := httpmock.NewStringResponse(200, "")
				response.Header.Set("etag", remoteEtag)

				return response, nil
			})

		// Assertions!
		err = updateE(cmd, args)
		assert.Nil(t, err)
		assert.Equal(t, "ci-utils is up-to-date\n", output.String())
		assert.Empty(t, commandBuf.Bytes())

		// Cleanup!
		commandBuf.Reset()
		output.Reset()
	})

	t.Run("install up-to-date but force is specified", func(t *testing.T) {
		fs = filesystem.NewMemFs()
		updateCmdOpts = &updateCmdOptions{
			channel: "stable",
			force:   true,
			phase:   1,
		}
		fslock = &filesystem.MemMapLock{
			Locks: make(map[string]*filesystem.MemLock),
		}

		// Mock out the on-disk etag
		err = fs.WriteFile(ciutils.SettingsPath("etag"), []byte(remoteEtag), 0644)
		assert.Nil(t, err)

		// Create the existing installation
		err = fs.MkdirAll(ciutils.RootDir(), 0755)
		assert.Nil(t, err)

		// Mock out the HTTP response
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("HEAD", `=~^https://packages.chef.io/files/stable/ci-utils/latest/.+\z`,
			func(req *http.Request) (*http.Response, error) {
				response := httpmock.NewStringResponse(200, "")
				response.Header.Set("etag", remoteEtag)

				return response, nil
			})

		httpmock.RegisterResponder("GET", `=~^https://packages.chef.io/files/stable/ci-utils/latest/.+\z`,
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(200, archiveContent), nil
			})

		// Assertions!
		err = updateE(cmd, args)
		assert.Nil(t, err)
		expectedOutput := heredoc.Doc(`
			--> {1} Upgrade available
			--> {1} Downloading latest release
			--> {1} Release downloaded -- backing up current installation
		`)
		assert.Equal(t, expectedOutput, output.String())

		expectedNextCmd = fmt.Sprintf("%s update --phase 2", ciutils.BackupBinPath("ci-utils"))
		assert.Equal(t, expectedNextCmd, commandBuf.String())

		inProgressFileExits, err := fs.Exists(ciutils.SettingsPath("upgrade-in-progress"))
		assert.Nil(t, err)
		assert.Equal(t, true, inProgressFileExits)

		backupDirExists, err := fs.Exists(ciutils.BackupDir())
		assert.Nil(t, err)
		assert.Equal(t, true, backupDirExists)

		// Cleanup!
		commandBuf.Reset()
		output.Reset()
	})

	t.Run("remote etag != local etag (new version)", func(t *testing.T) {
		fs = filesystem.NewMemFs()
		updateCmdOpts = defaultUpdateCmdOpts
		fslock = &filesystem.MemMapLock{
			Locks: make(map[string]*filesystem.MemLock),
		}

		// Mock out the on-disk etag
		err = fs.WriteFile(ciutils.SettingsPath("etag"), []byte("old-etag-value"), 0644)
		assert.Nil(t, err)

		// Create the existing installation
		err = fs.MkdirAll(ciutils.RootDir(), 0755)
		assert.Nil(t, err)

		// Mock out the HTTP response
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("HEAD", `=~^https://packages.chef.io/files/stable/ci-utils/latest/.+\z`,
			func(req *http.Request) (*http.Response, error) {
				response := httpmock.NewStringResponse(200, "")
				response.Header.Set("etag", remoteEtag)

				return response, nil
			})

		httpmock.RegisterResponder("GET", `=~^https://packages.chef.io/files/stable/ci-utils/latest/.+\z`,
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(200, archiveContent), nil
			})

		// Assertions!
		err = updateE(cmd, args)
		assert.Nil(t, err)
		expectedOutput := heredoc.Doc(`
			--> {1} Upgrade available
			--> {1} Downloading latest release
			--> {1} Release downloaded -- backing up current installation
		`)
		assert.Equal(t, expectedOutput, output.String())

		expectedNextCmd = fmt.Sprintf("%s update --phase 2", ciutils.BackupBinPath("ci-utils"))
		assert.Equal(t, expectedNextCmd, commandBuf.String())

		inProgressFileExits, err := fs.Exists(ciutils.SettingsPath("upgrade-in-progress"))
		assert.Nil(t, err)
		assert.Equal(t, true, inProgressFileExits)

		backupDirExists, err := fs.Exists(ciutils.BackupDir())
		assert.Nil(t, err)
		assert.Equal(t, true, backupDirExists)

		// Cleanup!
		commandBuf.Reset()
		output.Reset()
	})

	t.Run("new version with specified channel", func(t *testing.T) {
		fs = filesystem.NewMemFs()
		updateCmdOpts = &updateCmdOptions{
			channel: "current",
			force:   false,
			phase:   1,
		}
		fslock = &filesystem.MemMapLock{
			Locks: make(map[string]*filesystem.MemLock),
		}

		// Mock out the on-disk etag
		err = fs.WriteFile(ciutils.SettingsPath("etag"), []byte("old-etag-value"), 0644)
		assert.Nil(t, err)

		// Create the existing installation
		err = fs.MkdirAll(ciutils.RootDir(), 0755)
		assert.Nil(t, err)

		// Mock out the HTTP response
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("HEAD", `=~^https://packages.chef.io/files/current/ci-utils/latest/.+\z`,
			func(req *http.Request) (*http.Response, error) {
				response := httpmock.NewStringResponse(200, "")
				response.Header.Set("etag", remoteEtag)

				return response, nil
			})

		httpmock.RegisterResponder("GET", `=~^https://packages.chef.io/files/current/ci-utils/latest/.+\z`,
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(200, archiveContent), nil
			})

		// Assertions!
		err = updateE(cmd, args)
		assert.Nil(t, err)
		expectedOutput := heredoc.Doc(`
			--> {1} Upgrade available
			--> {1} Downloading latest release
			--> {1} Release downloaded -- backing up current installation
		`)
		assert.Equal(t, expectedOutput, output.String())

		expectedNextCmd = fmt.Sprintf("%s update --phase 2", ciutils.BackupBinPath("ci-utils"))
		assert.Equal(t, expectedNextCmd, commandBuf.String())

		inProgressFileExits, err := fs.Exists(ciutils.SettingsPath("upgrade-in-progress"))
		assert.Nil(t, err)
		assert.Equal(t, true, inProgressFileExits)

		backupDirExists, err := fs.Exists(ciutils.BackupDir())
		assert.Nil(t, err)
		assert.Equal(t, true, backupDirExists)

		// Cleanup!
		commandBuf.Reset()
		output.Reset()
	})
}

func TestUpdate_phaseTwo(t *testing.T) {
	var err error

	defaultUpdateCmdOpts := &updateCmdOptions{
		channel: "stable",
		force:   false,
		phase:   2,
	}

	ciutils = install.DefaultInstall()

	cmd := &cobra.Command{}
	args := []string{}
	output := new(bytes.Buffer)
	cmd.SetOutput(output)

	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()

	t.Run("failed to unarchive asset", func(t *testing.T) {
		fs = filesystem.NewMemFs()
		updateCmdOpts = defaultUpdateCmdOpts
		fslock = &filesystem.MemMapLock{
			Locks: make(map[string]*filesystem.MemLock),
		}

		// Seed phase one state
		err := fs.WriteFile(assetOnDisk(), []byte(archiveContent), 0644)
		assert.Nil(t, err)

		err = fs.WriteFile(ciutils.SettingsPath("upgrade-in-progress"), []byte(""), 0644)
		assert.Nil(t, err)

		// Unsuccessful untar
		archiver = &fakeUnarchiver{
			Logger:   &archiverBuf,
			Response: errors.New("woops"),
		}

		// Assertions!
		err = updateE(cmd, args)
		assert.NotNil(t, err)
		assert.Regexp(t, "failed to untar new installation into place", err.Error())
		expectedOutput := heredoc.Doc(`
			--> {2} Removing current installation
			--> {2} Moving new installation into place
		`)
		assert.Equal(t, expectedOutput, output.String())
		assert.Empty(t, commandBuf.Bytes())
		assertUpgradeStateReleased(t, ciutils, fs, fslock)

		// Cleanup!
		archiverBuf.Reset()
		commandBuf.Reset()
		output.Reset()
	})

	t.Run("happy path", func(t *testing.T) {
		fs = filesystem.NewMemFs()
		updateCmdOpts = defaultUpdateCmdOpts
		fslock = &filesystem.MemMapLock{
			Locks: make(map[string]*filesystem.MemLock),
		}

		assetOnDisk := filepath.Join(ciutils.RootDirParent(), "ci-utils.tar.gz")

		// Seed phase one state
		err = fs.WriteFile(assetOnDisk, []byte(archiveContent), 0644)
		assert.Nil(t, err)

		err = fs.WriteFile(ciutils.SettingsPath("upgrade-in-progress"), []byte("foo"), 0644)
		assert.Nil(t, err)

		// Successful un-tar
		archiver = &fakeUnarchiver{
			Logger:   &archiverBuf,
			Response: nil,
		}

		// Assertions!
		err = updateE(cmd, args)
		assert.Nil(t, err)
		expectedOutput := heredoc.Doc(`
			--> {2} Removing current installation
			--> {2} Moving new installation into place
		`)
		assert.Equal(t, expectedOutput, output.String())

		expectedNextCmd = fmt.Sprintf("%s update --phase 3", ciutils.BinPath("ci-utils"))
		assert.Equal(t, expectedNextCmd, commandBuf.String())

		expectedArchiveLog := fmt.Sprintf("unarchive %s into %s", assetOnDisk, ciutils.RootDirParent())
		assert.Equal(t, expectedArchiveLog, archiverBuf.String())

		inProgressFileExits, err := fs.Exists(ciutils.SettingsPath("upgrade-in-progress"))
		assert.Nil(t, err)
		assert.Equal(t, true, inProgressFileExits)

		// Cleanup!
		archiverBuf.Reset()
		commandBuf.Reset()
		output.Reset()
	})
}

func TestUpdate_phaseThree(t *testing.T) {
	var err error

	defaultUpdateCmdOpts := &updateCmdOptions{
		channel: "stable",
		force:   false,
		phase:   3,
	}

	ciutils = install.DefaultInstall()

	cmd := &cobra.Command{}
	args := []string{}
	output := new(bytes.Buffer)
	cmd.SetOutput(output)

	t.Run("failed to rename etag file", func(t *testing.T) {
		fs = filesystem.NewMemFs()
		updateCmdOpts = defaultUpdateCmdOpts
		fslock = &filesystem.MemMapLock{
			Locks: make(map[string]*filesystem.MemLock),
		}

		etagPath := ciutils.SettingsPath("etag")

		// Mock out the on-disk etags
		err := fs.WriteFile(etagPath, []byte("old-etag-value"), 0644)
		assert.Nil(t, err)

		// Missing new etag-new -- this will cause error

		err = fs.WriteFile(ciutils.SettingsPath("upgrade-in-progress"), []byte(""), 0644)
		assert.Nil(t, err)

		// Assertions!
		err = updateE(cmd, args)
		assert.NotNil(t, err)
		assert.Regexp(t, "failed to move new etag into place on disk", err.Error())

		expectedOutput := heredoc.Doc(`
			--> {3} Upgrade complete -- cleaning up
		`)
		assert.Equal(t, expectedOutput, output.String())
		assert.Empty(t, commandBuf.String())
		assertUpgradeStateReleased(t, ciutils, fs, fslock)

		// Cleanup!
		output.Reset()
	})

	t.Run("happy path", func(t *testing.T) {
		fs = filesystem.NewMemFs()
		updateCmdOpts = defaultUpdateCmdOpts
		fslock = &filesystem.MemMapLock{
			Locks: make(map[string]*filesystem.MemLock),
		}

		// Mock out install state
		err = fs.WriteFile(ciutils.SettingsPath("etag"), []byte("old-etag-value"), 0644)
		assert.Nil(t, err)

		err = fs.WriteFile(ciutils.SettingsPath("etag-new"), []byte(remoteEtag), 0644)
		assert.Nil(t, err)

		err = fs.WriteFile(ciutils.SettingsPath("upgrade-in-progress"), []byte(""), 0644)
		assert.Nil(t, err)

		// Assertions!
		err = updateE(cmd, args)
		assert.Nil(t, err)
		expectedOutput := heredoc.Doc(`
			--> {3} Upgrade complete -- cleaning up
		`)
		assert.Equal(t, expectedOutput, output.String())
		assert.Empty(t, commandBuf.String())
		assertUpgradeStateReleased(t, ciutils, fs, fslock)

		// Cleanup!
		output.Reset()
	})
}

// Fake archiver
type fakeUnarchiver struct {
	Logger   *bytes.Buffer
	Response error
}

func (u *fakeUnarchiver) Unarchive(source string, destination string) error {
	_, err := fmt.Fprintf(u.Logger, "unarchive %s into %s", source, destination)
	if err != nil {
		return err
	}

	return u.Response
}

// Fake Shell Commands
func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}

	if _, err := fmt.Fprintf(&commandBuf, fmt.Sprintf("%s %s", command, strings.Join(args, " "))); err != nil {
		panic(err)
	}

	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	os.Exit(0)
}

func assertUpgradeStateReleased(t *testing.T, install install.Install, fs filesystem.FileSystem, fslock filesystem.Locker) {
	inProgressFileExits, err := fs.Exists(ciutils.SettingsPath("upgrade-in-progress"))
	assert.Nil(t, err)
	assert.False(t, inProgressFileExits)

	backupDirExists, err := fs.Exists(ciutils.BackupDir())
	assert.Nil(t, err)
	assert.False(t, backupDirExists)

	assetExists, err := fs.Exists(assetOnDisk())
	assert.Nil(t, err)
	assert.False(t, assetExists)

	lock, _ := fslock.GetLock(ciutils.LockPath("upgrade-ci-utils"))
	_, err = lock.TryLock()
	assert.Nil(t, err)
}
