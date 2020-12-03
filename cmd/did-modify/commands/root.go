package commands

import (
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/ryanuber/go-glob"
	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type rootCmdOptions struct {
	gitref string
	globs  []string
}

var (
	rootCmd = &cobra.Command{
		Use:   "did-modify",
		Short: `Prints "true" to STDOUT if any files matching GLOBS were modified between HEAD and GITREF. Otherwise, prints "false".`,
		RunE:  didModifyE,
	}

	rootCmdOpts = &rootCmdOptions{}
)

func init() {
	rootCmd.Flags().StringVar(&rootCmdOpts.gitref, "git-ref", "HEAD~1", "A valid Git reference (e.g. HEAD, master, origin/master, etc).")
	rootCmd.Flags().StringSliceVar(&rootCmdOpts.globs, "globs", []string{"*"}, "Comma-separated list of glob patterns to inspect to determine if there are changes.")
}

// Execute handles the execution of child commands and flags.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func didModifyE(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "failed to get current working directory")
	}

	repo, err := git.PlainOpen(cwd)
	if err != nil {
		return errors.Wrapf(err, "failed to open git repository at %s", cwd)
	}

	head, err := repo.Head()
	if err != nil {
		return errors.Wrap(err, "failed to fetch HEAD of repository")
	}

	headCommit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return errors.Wrap(err, "failed to fetch commit sha for HEAD of repository")
	}

	gitRefRev, err := repo.ResolveRevision(plumbing.Revision(rootCmdOpts.gitref))
	if err != nil {
		return errors.Wrapf(err, "failed to fetch the fully qualified revision for %s", rootCmdOpts.gitref)
	}

	gitRefCommit, err := repo.CommitObject(*gitRefRev)
	if err != nil {
		return errors.Wrap(err, "failed to fetch git commit for revision")
	}

	patch, err := headCommit.Patch(gitRefCommit)
	if err != nil {
		return errors.Wrap(err, "failed to get git patch for git commit")
	}

	for _, fileStat := range patch.Stats() {
		for _, globPattern := range rootCmdOpts.globs {
			if glob.Glob(globPattern, fileStat.Name) {
				cmd.Print("true")

				return nil
			}
		}
	}

	cmd.Print("false")

	return nil
}
