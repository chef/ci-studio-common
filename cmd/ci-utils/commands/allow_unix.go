// +build !windows

package commands

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func allowE(cmd *cobra.Command, args []string) error {
	err := addSudoPermission("/bin/hab", args[0])
	if err != nil {
		return errors.Wrapf(err, "failed to grant sudo permissions to /bin/hab for %s", args[0])
	}

	err = addSudoPermission("/usr/bin/ci-utils", args[0])
	if err != nil {
		return errors.Wrapf(err, "failed to grant sudo permissions to /usr/bin/ci-utils for %s", args[0])
	}

	err = addSudoPermission("/usr/bin/install-habitat", args[0])
	if err != nil {
		return errors.Wrapf(err, "failed to grant sudo permissions to /usr/bin/install-habitat for %s", args[0])
	}

	return nil
}

// addSudoPermission will ensure that the user can run the given command with sudo without a password
func addSudoPermission(command string, user string) error {
	sudoersFilePath := fmt.Sprintf("/etc/sudoers.d/%s", user)
	newLine := fmt.Sprintf("%s ALL=NOPASSWD:SETENV: %s", user, command)

	err := fs.AppendIfMissing(sudoersFilePath, []byte(newLine), 0440)
	if err != nil {
		return errors.Wrapf(err, "failed to append sudo permission for %s to %s", user, sudoersFilePath)
	}

	return nil
}
