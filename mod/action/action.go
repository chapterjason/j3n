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
	"github.com/pkg/errors"
)

var (
	ErrStepNotFound = errors.New("step not found")
)

type Action struct {
	Dependencies []string         `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
	Steps        map[string]*Step `json:"steps" yaml:"steps"`
}

func (a *Action) HasStep(step string) bool {
	_, ok := a.Steps[step]
	return ok
}

func (a *Action) GetStep(step string) (*Step, error) {
	s, ok := a.Steps[step]

	if !ok {
		return s, ErrStepNotFound
	}

	return s, nil
}

func (a *Action) MustGetStep(step string) *Step {
	s, err := a.GetStep(step)

	if err != nil {
		panic(err)
	}

	return s
}
