package lib

import(
	"fmt"
	"os"
)

// Check will throw a panic if the err is not nil
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// AddSudoPermission will ensure that the user can run the given command with sudo without a password
func AddSudoPermission(command string, user string) {
	sudoersFilePath := fmt.Sprintf("/etc/sudoers.d/%s", user)
	newLine := fmt.Sprintf("%s ALL=NOPASSWD:SETENV: %s\n", user, command)

	err := AppendIfMissing(sudoersFilePath, newLine)
	Check(err)

	os.Chmod(sudoersFilePath, 0440)
}