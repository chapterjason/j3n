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
	"fmt"

	"github.com/pkg/errors"
)

type Executer struct {
	m       *Map
	outputs map[string]any
}

func NewExecuter(m *Map) *Executer {
	return &Executer{
		m:       m,
		outputs: make(map[string]any),
	}
}

func (r *Executer) hasOutput(ref reference) bool {
	_, ok := r.outputs[ref.String()]

	return ok
}

func (r *Executer) Execute(s string) error {
	ep := newExecutionPlan(r.m)
	err := ep.plan(s)

	if err != nil {
		return errors.Wrapf(err, "failed to plan execution of %s", s)
	}

	for _, ref := range ep.refs {
		action := r.m.MustGetAction(ref.action)
		step := action.MustGetStep(ref.step)

		err := r.executeStep(ref, step)

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Executer) executeStep(ref reference, s *Step) error {
	var input any

	if s.Input != "" {
		inputRef, err := resolveReference(s.Input, ref.action)

		if err != nil {
			return errors.Wrapf(err, "failed to resolve input reference %s", s.Input)
		}

		if !r.hasOutput(inputRef) {
			return fmt.Errorf("input %s not found", s.Input)
		}

		input = r.mustGetOutput(inputRef)
	}

	sr, ok := Steps[s.Type]

	if !ok {
		return fmt.Errorf("no runner for step type %s", s.Type)
	}

	out, err := sr(input, s.Params)

	if err != nil {
		return err
	}

	if s.Output {
		if out == nil {
			return fmt.Errorf("output of step %s is nil", ref.String())
		}

		r.outputs[ref.String()] = out
	}

	return nil
}

func (r *Executer) getOutput(ref reference) (any, error) {
	rs := ref.String()

	if !r.hasOutput(ref) {
		return nil, fmt.Errorf("output %s not found", rs)
	}

	return r.outputs[rs], nil
}

func (r *Executer) mustGetOutput(ref reference) any {
	out, err := r.getOutput(ref)

	if err != nil {
		panic(err)
	}

	return out
}
