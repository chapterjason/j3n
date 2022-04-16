/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/chapterjason/j3n/mod/release"
	"github.com/chapterjason/j3n/mod/version"
)

var (
	ErrDirectoryNotEmpty = errors.New("directory not empty")
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new project",
	RunE: func(cmd *cobra.Command, args []string) error {
		workflow, err := cmd.Flags().GetString("workflow")

		if err != nil {
			return err
		}

		directory, err := cmd.Flags().GetString("directory")

		if err != nil {
			return err
		}

		return initProject(directory, workflow)
	},
}

// initProject initializes a new project
// - ensure the directory exists and is empty
// - initialize the git repository
// - set the default branch to release/0.1
// - set default config values
// - save config
// - commit config
func initProject(directory string, workflow string) error {
	v := version.MustParse("0.1.0-DEV")

	if workflow != "multi_branch" {
		return errors.New("only multi_branch workflow is supported")
	}

	if directory == "" {
		wd, err := os.Getwd()

		if err != nil {
			return errors.Wrap(err, "failed to get working directory")
		}

		directory = wd
	}

	_, err := os.Stat(directory)

	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(directory, os.ModePerm)

			if err != nil {
				return errors.Wrap(err, "failed to create directory")
			}
		} else {
			return errors.Wrap(err, "failed to stat directory")
		}
	} else {
		f, err := os.ReadDir(directory)

		if err != nil {
			return errors.Wrap(err, "failed to read directory")
		}

		if len(f) > 0 {
			return ErrDirectoryNotEmpty
		}
	}

	r, err := git.PlainInit(directory, false)

	if err != nil {
		return errors.Wrap(err, "failed to initialize git repository")
	}

	w, err := r.Worktree()

	if err != nil {
		return errors.Wrap(err, "failed to get worktree")
	}

	viper.Set("$schema", "https://raw.githubusercontent.com/chapterjason/j3n/release/0.1/resources/schema/all.json")
	viper.Set("version.current", v.String())
	viper.Set("release.workflow.type", workflow)

	err = viper.SafeWriteConfigAs(path.Join(directory, "j3n.json"))

	if err != nil {
		return errors.Wrap(err, "failed to write config")
	}

	err = w.AddWithOptions(
		&git.AddOptions{
			All: true,
		},
	)

	if err != nil {
		return errors.Wrap(err, "failed to add files")
	}

	_, err = w.Commit(
		"feat: Add initial set of files",
		&git.CommitOptions{
			All: true,
		},
	)

	if err != nil {
		return errors.Wrap(err, "failed to commit files")
	}

	h, err := r.Head()

	if err != nil {
		return errors.Wrap(err, "failed to get head")
	}

	rbs := release.GitReleaseBranchFormatter(v)
	ref := plumbing.NewHashReference(plumbing.NewBranchReferenceName(rbs), h.Hash())

	err = r.Storer.SetReference(ref)

	if err != nil {
		return errors.Wrapf(err, "failed to create %s branch", rbs)
	}

	err = r.CreateBranch(
		&config.Branch{
			Name:  rbs,
			Merge: ref.Name(),
		},
	)

	if err != nil {
		return err
	}

	err = w.Checkout(
		&git.CheckoutOptions{
			Branch: ref.Name(),
		},
	)

	if err != nil {
		return errors.Wrap(err, "failed to checkout branch")
	}

	err = r.Storer.RemoveReference(plumbing.Master)

	if err != nil {
		return errors.Wrap(err, "failed to delete branch")
	}

	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringP("directory", "d", "", "Directory to initialize in")
	initCmd.Flags().String("workflow", "multi_branch", "Release workflow to use")
}
