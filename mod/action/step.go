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

type Step struct {
	Type            string         `json:"type" yaml:"type"`
	Dependencies    []string       `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
	Input           string         `json:"input,omitempty" yaml:"input,omitempty"`
	Output          bool           `json:"output,omitempty" yaml:"output,omitempty"`
	ContinueOnError bool           `json:"continue_on_error,omitempty" yaml:"continue_on_error,omitempty"`
	IgnoreExitCodes []int          `json:"ignore_exit_codes,omitempty" yaml:"ignore_exit_codes,omitempty"`
	Params          map[string]any `json:"params,omitempty" yaml:"params,omitempty"`
}
