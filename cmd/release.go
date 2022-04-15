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
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"

	"github.com/chapterjason/j3n/mod/release"
	"github.com/chapterjason/j3n/mod/version"
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release [version]",
	Short: "Create a new release of a project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		v, err := version.Parse(args[0])

		if err != nil {
			return err
		}

		wd, err := os.Getwd()

		if err != nil {
			return err
		}

		repo, err := git.PlainOpen(wd)

		if err != nil {
			return err
		}

		has, err := release.HasTag(repo, v)

		if err != nil {
			return err
		}

		if has {
			return fmt.Errorf("tag %s already exists", v)
		}

		if v.Patch != 0 {
			// is a patch release
			// ensure tag does not already exist
			// ensure release branch exist
			// checkout release branch
			// pre-release hook
			// create tag
			// post-release hook
		} else if v.Minor != 0 {
			// is a minor release
		} else {
			// is a major release
		}

		// rbn := release.DefaultGitReleaseBranchFormatter(v)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(releaseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// releaseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// releaseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
