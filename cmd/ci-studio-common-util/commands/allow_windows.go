// +build windows

package commands

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func grantSuperUserPermission(cmd *cobra.Command, args []string) error {
	return errors.New("this command is not supported on Windows")
}
