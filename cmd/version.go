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

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/chapterjason/j3n/mod/version"
	"github.com/chapterjason/j3n/modx/viperx"
)

var (
	ErrUnknownStrategy = errors.New("unknown versioning strategy")
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Manage the version of a project",
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func initVersionConfig() {
	wd, err := os.Getwd()
	cobra.CheckErr(err)

	pkgPath := path.Join(wd, "package.json")

	if _, err := os.Stat(pkgPath); err == nil {
		strategy := version.NewNpmStrategy(wd)

		version.Setters = append(version.Setters, strategy)
		version.Getters = append(version.Getters, strategy)
	}

	j3nPath := path.Join(wd, "j3n.json")

	if _, err := os.Stat(j3nPath); err == nil {
		strategy := version.NewVersionStrategy(wd)

		version.Setters = append(version.Setters, strategy)
		version.Getters = append(version.Getters, strategy)
	}

	sv := viper.Get("version.strategies")

	if sv != nil {
		strategies := sv.([]interface{})

		for _, strategy := range strategies {
			type typ struct {
				Type string `json:"type"`
			}

			t := typ{}

			err := viperx.Transcode(strategy, &t)
			cobra.CheckErr(err)

			switch t.Type {
			case "npm":
				ns := version.NpmStrategy{}
				err := viperx.Transcode(strategy, &ns)
				cobra.CheckErr(err)

				version.Setters = append(version.Setters, &ns)
				version.Getters = append(version.Getters, &ns)
			case "expression":
				es := version.ExpressionStrategy{}
				err := viperx.Transcode(strategy, &es)
				cobra.CheckErr(err)

				version.Setters = append(version.Setters, &es)
			default:
				cobra.CheckErr(ErrUnknownStrategy)
			}
		}
	}
}
