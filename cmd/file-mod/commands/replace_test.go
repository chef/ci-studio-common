package commands

import (
	"bytes"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/chef/ci-studio-common/internal/pkg/filesystem"
)

func TestFindAndReplaceE(t *testing.T) {
	fs = filesystem.NewMemFs()

	tempFileName := "/tmp/test-find-and-replace"

	existingContent := heredoc.Doc(`
		the quick brown fox
		jumps over
		the lazy dog
	`)

	cmd := &cobra.Command{}
	output := new(bytes.Buffer)
	cmd.SetOutput(output)

	t.Run("no match", func(t *testing.T) {
		err := fs.WriteFile(tempFileName, []byte(existingContent), 0644)
		assert.Nil(t, err)

		err = findAndReplaceE(cmd, []string{"/foo.+/", "bar", tempFileName})
		assert.Nil(t, err)

		actual, err := fs.ReadFile(tempFileName)
		assert.Nil(t, err)
		assert.Equal(t, existingContent, string(actual))
	})

	t.Run("single match", func(t *testing.T) {
		err := fs.WriteFile(tempFileName, []byte(existingContent), 0644)
		assert.Nil(t, err)

		err = findAndReplaceE(cmd, []string{"br.+n", "purple", tempFileName})
		assert.Nil(t, err)

		expected := heredoc.Doc(`
			the quick purple fox
			jumps over
			the lazy dog
		`)

		actual, err := fs.ReadFile(tempFileName)
		assert.Nil(t, err)
		assert.Equal(t, expected, string(actual))
	})

	t.Run("single match w/ replace", func(t *testing.T) {
		err := fs.WriteFile(tempFileName, []byte(existingContent), 0644)
		assert.Nil(t, err)

		err = findAndReplaceE(cmd, []string{`(?m)brown\s+(.+)$`, "pink $1", tempFileName})
		assert.Nil(t, err)

		expected := heredoc.Doc(`
			the quick pink fox
			jumps over
			the lazy dog
		`)

		actual, err := fs.ReadFile(tempFileName)
		assert.Nil(t, err)
		assert.Equal(t, expected, string(actual))
	})

	t.Run("multiple matches", func(t *testing.T) {
		err := fs.WriteFile(tempFileName, []byte(existingContent), 0644)
		assert.Nil(t, err)

		err = findAndReplaceE(cmd, []string{"the", "da", tempFileName})
		assert.Nil(t, err)

		expected := heredoc.Doc(`
			da quick brown fox
			jumps over
			da lazy dog
		`)

		actual, err := fs.ReadFile(tempFileName)
		assert.Nil(t, err)
		assert.Equal(t, expected, string(actual))
	})

	t.Run("single multi-line match", func(t *testing.T) {
		err := fs.WriteFile(tempFileName, []byte(existingContent), 0644)
		assert.Nil(t, err)

		err = findAndReplaceE(cmd, []string{"(?ms)fox(.+)over", "fox\nhops over", tempFileName})
		assert.Nil(t, err)

		expected := heredoc.Doc(`
			the quick brown fox
			hops over
			the lazy dog
		`)

		actual, err := fs.ReadFile(tempFileName)
		assert.Nil(t, err)
		assert.Equal(t, expected, string(actual))
	})

	t.Run("multiple multi-line match", func(t *testing.T) {
		err := fs.WriteFile(tempFileName, []byte(existingContent), 0644)
		assert.Nil(t, err)

		err = findAndReplaceE(cmd, []string{"(?ms).{1}$.{2}", "1\n2", tempFileName})
		assert.Nil(t, err)

		expected := heredoc.Doc(`
			the quick brown fo1
			2umps ove1
			2he lazy dog
		`)

		actual, err := fs.ReadFile(tempFileName)
		assert.Nil(t, err)
		assert.Equal(t, expected, string(actual))
	})
}
