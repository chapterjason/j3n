/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	if workflow != "multi_branch" {
		return errors.New("only multi_branch workflow is supported")
	}

	if directory == "" {
		wd, err := os.Getwd()

		if err != nil {
			return err
		}

		directory = wd
	}

	_, err := os.Stat(directory)

	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(directory, os.ModePerm)

			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		f, err := os.ReadDir(directory)

		if err != nil {
			return err
		}

		if len(f) > 0 {
			return ErrDirectoryNotEmpty
		}
	}

	r, err := git.PlainInit(directory, false)

	if err != nil {
		return err
	}

	cfg, err := r.Config()

	if err != nil {
		return err
	}

	cfg.Init.DefaultBranch = "release/0.1"

	err = r.SetConfig(cfg)

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	viper.Set("$schema", "https://raw.githubusercontent.com/chapterjason/j3n/release/0.1/resources/schema/all.json")
	viper.Set("version", "0.1.0-DEV")
	viper.Set("release.workflow", workflow)

	err = viper.SafeWriteConfigAs(path.Join(directory, "j3n.json"))

	if err != nil {
		return err
	}

	w, err := r.Worktree()

	if err != nil {
		return err
	}

	_, err = w.Commit(
		"feat: Add initial set of files",
		&git.CommitOptions{
			All: true,
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringP("directory", "d", "", "Directory to initialize in")
	initCmd.Flags().String("workflow", "multi_branch", "Release workflow to use")
}
