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

package action

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/pkg/errors"

	"github.com/chapterjason/j3n/modx/slicex"
)

func init() {
	MustRegister(
		"exec", func(input any, params map[string]any) (any, error) {
			command := params["command"].(string)

			continueOnError := false
			ignoreExitCodes := []float64{}
			printStdout := false
			printStderr := false

			var args []string

			var stdout, stderr bytes.Buffer

			if params["args"] != nil {
				args = slicex.ToString(params["args"])
			}

			cmd := exec.Command(command, args...)

			if params["directory"] != nil {
				dir := params["directory"].(string)

				if !path.IsAbs(dir) {
					wd, err := os.Getwd()

					if err != nil {
						return nil, errors.Wrap(err, "failed to get working directory")
					}

					dir = path.Join(wd, dir)
				}

				cmd.Dir = dir
			}

			if params["env"] != nil {
				cmd.Env = slicex.ToString(params["env"])
			}

			if params["continue_on_error"] != nil {
				continueOnError = params["continueOnError"].(bool)
			}

			if params["ignore_exit_codes"] != nil {
				ignoreExitCodes = slicex.ToFloat(params["ignore_exit_codes"])
			}

			if params["print_stdout"] != nil {
				printStdout = params["print_stdout"].(bool)
			}

			if params["print_stderr"] != nil {
				printStderr = params["print_stderr"].(bool)
			}

			if input != nil {
				is := input.(string)
				cmd.Stdin = strings.NewReader(is)
			}

			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Start()

			if err != nil {
				return nil, errors.Wrap(err, "failed to start command")
			}

			err = cmd.Wait()

			if err != nil {
				exitCode := cmd.ProcessState.ExitCode()

				if !slicex.Contains(ignoreExitCodes, float64(exitCode)) && !continueOnError {
					message := fmt.Sprintf("command \"%s\" failed\n", cmd.String())
					message += fmt.Sprintf("    exit code: %d\n", exitCode)
					message += fmt.Sprintf("    stderr: %s\n", strings.TrimSpace(stderr.String()))
					message += fmt.Sprintf("    stdout: %s\n", strings.TrimSpace(stdout.String()))

					return nil, errors.New(message)
				}
			}

			if printStdout {
				fmt.Print(stdout.String())
			}

			if printStderr {
				fmt.Print(stderr.String())
			}

			return stdout.String() + stderr.String(), nil
		},
	)
}
