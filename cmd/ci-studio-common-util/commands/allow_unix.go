// +build !windows

package commands

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/chef/ci-studio-common/internal/pkg/system"
)

func grantSuperUserPermission(cmd *cobra.Command, args []string) error {
	err := system.AddSudoPermission("/bin/hab", args[0])
	if err != nil {
		return errors.Wrapf(err, "failed to grant sudo permissions to /bin/hab for %s", args[0])
	}

	err = system.AddSudoPermission("/usr/bin/ci-studio-common-util", args[0])
	if err != nil {
		return errors.Wrapf(err, "failed to grant sudo permissions to /bin/bin/ci-studio-common-uril for %s", args[0])
	}

	err = system.AddSudoPermission("/usr/bin/install-habitat", args[0])
	if err != nil {
		return errors.Wrapf(err, "failed to grant sudo permissions to /usr/bin/install-habitat for %s", args[0])
	}

	return nil
}
