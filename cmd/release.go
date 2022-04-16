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

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/chapterjason/j3n/mod/release"
	"github.com/chapterjason/j3n/mod/version"
	"github.com/chapterjason/j3n/modx/viperx"
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release [version]",
	Short: "Create a new release of a project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rs := viper.Get("release")

		var rc release.Config

		if err := viperx.Transcode(rs, &rc); err != nil {
			return err
		}

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

		return release.Release(repo, v, rc.Workflow)
	},
}

func init() {
	rootCmd.AddCommand(releaseCmd)
}
