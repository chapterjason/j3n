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
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	"github.com/chapterjason/j3n/mod/action"
)

// actionCmd represents the action command
var actionCmd = &cobra.Command{
	Use:   "action [name]",
	Short: "Run an action",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		configFilePath, err := cmd.Flags().GetString("config")

		if err != nil {
			return err
		}

		file, err := os.OpenFile(configFilePath, os.O_RDONLY, 0600)

		if err != nil {
			return err
		}

		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				panic(err)
			}
		}(file)

		b, err := ioutil.ReadAll(file)

		if err != nil {
			return err
		}

		var m action.Map

		err = json.Unmarshal(b, &m)

		if err != nil {
			return err
		}

		ep := action.NewExecuter(&m)

		err = ep.Execute(args[0])

		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(actionCmd)

	actionCmd.Flags().StringP("config", "c", "actions.json", "Path to config file")
	actionCmd.MarkFlagFilename("config", "json")
}
