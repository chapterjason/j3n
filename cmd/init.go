/*
 * Copyright © 2022 Jason Schilling
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */
package cmd

import (
	"os"
	"path"

	"github.com/gogs/git-module"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/chapterjason/j3n/mod/release"
	"github.com/chapterjason/j3n/mod/version"
	"github.com/chapterjason/j3n/modx/gitx"
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

	err = git.Init(
		directory, git.InitOptions{
			Bare: false,
		},
	)

	if err != nil {
		return errors.Wrap(err, "failed to initialize git repository")
	}

	r, err := git.Open(directory)

	if err != nil {
		return errors.Wrap(err, "failed to open git repository")
	}

	viper.Set("$schema", "https://raw.githubusercontent.com/chapterjason/j3n/release/0.1/resources/schema/all.json")
	viper.Set("version.current", v.String())
	viper.Set("release.workflow.type", workflow)

	err = viper.SafeWriteConfigAs(path.Join(directory, "j3n.json"))

	if err != nil {
		return errors.Wrap(err, "failed to write config")
	}

	err = r.Add(
		git.AddOptions{
			All: true,
		},
	)

	if err != nil {
		return errors.Wrap(err, "failed to add files to git repository")
	}

	sig, err := gitx.GetSignature(r)

	if err != nil {
		return errors.Wrap(err, "failed to get signature")
	}

	err = r.Commit(sig, "feat: Add initial set of files")

	if err != nil {
		return errors.Wrap(err, "failed to commit")
	}

	headRevision, err := r.RevParse("HEAD")

	if err != nil {
		return errors.Wrap(err, "failed to get head")
	}

	branches, err := r.Branches()

	if err != nil {
		return errors.Wrap(err, "failed to get branches")
	}

	rbs := release.BranchFormatter(v)

	err = r.Checkout(rbs, git.CheckoutOptions{BaseBranch: headRevision})

	if err != nil {
		return errors.Wrap(err, "failed to create first release branch")
	}

	for _, branch := range branches {
		err := r.DeleteBranch(
			branch, git.DeleteBranchOptions{
				Force: true,
			},
		)

		if err != nil {
			return errors.Wrap(err, "failed to delete branch")
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringP("directory", "d", "", "Directory to initialize in")
	initCmd.Flags().String("workflow", "multi_branch", "Release workflow to use")
}
