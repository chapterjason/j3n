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

	"github.com/erikgeiser/promptkit/selection"
	"github.com/gogs/git-module"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/chapterjason/j3n/mod/release"
	"github.com/chapterjason/j3n/mod/version"
	"github.com/chapterjason/j3n/modx/slicex"
	"github.com/chapterjason/j3n/modx/viperx"
)

// releaseNextCmd represents the releaseNext command
var releaseNextCmd = &cobra.Command{
	Use:   "next [major|minor]",
	Short: "Prepare the next release branch",
	Args:  cobra.MaximumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			if !slicex.Contains([]string{"major", "minor"}, args[0]) {
				return errors.Errorf("invalid release type: %s", args[0])
			}
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		rs := viper.Get("release")

		var rc release.Config

		if err := viperx.Transcode(rs, &rc); err != nil {
			return err
		}

		_, isMulti := rc.Workflow.(*release.MultiBranchWorkflow)

		if !isMulti {
			return errors.New(cmd.Use + " requires a multi-branch workflow")
		}

		var v version.Version
		var err error
		vf := cmd.Flag("version")

		if vf != nil && vf.Value.String() != "" {
			v, err = version.Parse(vf.Value.String())

			if err != nil {
				return errors.Wrap(err, "failed to parse version")
			}
		} else {
			v, err = version.Get()

			if err != nil {
				return errors.Wrap(err, "failed to get current version")
			}
		}

		if len(args) == 0 {
			t, err := askForNextReleaseType()

			if err != nil {
				return errors.Wrap(err, "failed to ask for next release type")
			}

			args = []string{t}
		}

		next := args[0]

		nv := v

		if next == "major" {
			nv.Patch = 0
			nv.Minor = 0
			nv.Major++
		} else if next == "minor" {
			nv.Patch = 0
			nv.Minor++
		} else {
			return errors.Errorf("invalid release type: %s", next)
		}

		nv.Prerelease = []string{"DEV"}

		wd, err := os.Getwd()

		if err != nil {
			return err
		}

		r, err := git.Open(wd)

		if err != nil {
			return err
		}

		err = release.Next(r, v, nv, rc.Workflow)

		if err != nil {
			return errors.Wrap(err, "failed to create next release branch")
		}

		return nil
	},
}

func askForNextReleaseType() (string, error) {
	items := selection.Choices(
		[]string{
			"minor",
			"major",
		},
	)

	sel := selection.New("What type of release will be the next one?", items)
	item, err := sel.RunPrompt()

	if err != nil {
		return "", err
	}

	return item.String, nil
}

func init() {
	releaseCmd.AddCommand(releaseNextCmd)

	releaseNextCmd.Flags().StringP("version", "v", "", "version")
}
