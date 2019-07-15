// +build !windows

package system

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/chef/ci-studio-common/internal/pkg/files"
)

// SuperUser returns the name of the superuser
const SuperUser string = "root"

// AddSudoPermission will ensure that the user can run the given command with sudo without a password
func AddSudoPermission(command string, user string) error {
	sudoersFilePath := fmt.Sprintf("/etc/sudoers.d/%s", user)
	newLine := fmt.Sprintf("%s ALL=NOPASSWD:SETENV: %s\n", user, command)

	err := files.AppendIfMissing(sudoersFilePath, newLine)
	if err != nil {
		return errors.Wrapf(err, "failed append sudo permission for %s to sudoers file", user)
	}

	err = os.Chmod(sudoersFilePath, 0440)
	if err != nil {
		return errors.Wrap(err, "failed to change sudoers file permissions to 0440")
	}

	return nil
}
