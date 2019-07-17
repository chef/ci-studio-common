// +build !windows

package install

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnixInstallBinPath(t *testing.T) {
	install := DefaultInstall()
	actual := install.BinPath("foo")
	assert.Equal(t, "/opt/ci-utils/bin/foo", actual)
}
