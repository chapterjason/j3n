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
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func init() {
	MustRegister(
		"exec", func(input any, params map[string]any) (any, error) {
			command := params["command"].(string)
			var args []string

			if params["args"] != nil {
				args = params["args"].([]string)
			}

			cmd := exec.Command(command, args...)

			if params["directory"] != nil {
				cmd.Dir = params["directory"].(string)
			}

			if params["env"] != nil {
				cmd.Env = params["env"].([]string)
			}

			if input != nil {
				is := input.(string)
				cmd.Stdin = strings.NewReader(is)
			}

			b, err := cmd.CombinedOutput()

			if err != nil {
				return nil, errors.Wrap(err, "failed to run command")
			}

			return string(b), nil
		},
	)
}
