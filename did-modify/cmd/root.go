package cmd

import (
	"fmt"
	"os"

	"github.com/chef/ci-studio-common/lib"
	"github.com/ryanuber/go-glob"
	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

var (
	rootCmd = &cobra.Command{
		Use:   "did-modify",
		Short: `Prints "true" to STDOUT if any files matching GLOBS were modified between HEAD and GITREF. Otherwise, prints "false".`,
		Run:   detectModifiedFiles,
	}

	gitref string
	globs  []string
)

func init() {
	rootCmd.Flags().StringVar(&gitref, "git-ref", "HEAD~1", "A valid Git reference (e.g. HEAD, master, origin/master, etc).")
	rootCmd.Flags().StringSliceVar(&globs, "globs", []string{"*"}, "Comma-separated list of glob patterns to inspect to determine if there are changes.")
}

// Execute runs the command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func detectModifiedFiles(cmd *cobra.Command, args []string) {
	cwd, _ := os.Getwd()

	repo, err := git.PlainOpen(cwd)
	lib.Check(err)

	head, err := repo.Head()
	lib.Check(err)
	headCommit, err := repo.CommitObject(head.Hash())

	gitRefRev, err := repo.ResolveRevision(plumbing.Revision(gitref))
	lib.Check(err)

	gitRefCommit, err := repo.CommitObject(*gitRefRev)
	lib.Check(err)

	patch, err := headCommit.Patch(gitRefCommit)
	lib.Check(err)

	for _, fileStat := range patch.Stats() {
		for _, globPattern := range globs {
			if glob.Glob(globPattern, fileStat.Name) {
				fmt.Print("true")
                                return
			}
		}
	}

	fmt.Print("false")
}
