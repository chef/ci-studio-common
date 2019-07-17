// +build windows

package install

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWinInstallBinPath(t *testing.T) {
	install := DefaultInstall()
	actual := install.BinPath("foo")
	assert.Equal(t, "C:\\ci-utils\\bin\\foo.exe", actual)
}
