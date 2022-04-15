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
	"strings"

	"github.com/pkg/errors"
)

var (
	ReferenceSeparator  = "."
	ErrInvalidReference = errors.New("invalid reference")
)

type reference struct {
	action string
	step   string
}

func newReference(action string, step string) reference {
	return reference{
		action: action,
		step:   step,
	}
}

func parseReference(ref string) (reference, error) {
	parts := strings.Split(ref, ReferenceSeparator)

	if len(parts) != 2 {
		return reference{}, ErrInvalidReference
	}

	return newReference(parts[0], parts[1]), nil
}

func resolveReference(s string, action string) (reference, error) {
	if strings.Contains(s, ReferenceSeparator) {
		ref, err := parseReference(s)

		if err != nil {
			return reference{}, err
		}

		return ref, nil
	} else {
		return newReference(action, s), nil
	}
}

func (r *reference) String() string {
	return r.action + ReferenceSeparator + r.step
}

func (r *reference) compare(other reference) bool {
	return r.action == other.action && r.step == other.step
}
