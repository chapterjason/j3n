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
)

type executionPlan struct {
	m *Map

	refs []reference
}

func newExecutionPlan(m *Map) *executionPlan {
	return &executionPlan{
		m: m,
	}
}

func (p *executionPlan) plan(name string) error {
	if !p.m.HasAction(name) {
		return fmt.Errorf("action %s not found", name)
	}

	p.refs = []reference{}

	err := p.planAction(name, p.m.MustGetAction(name))

	if err != nil {
		return err
	}

	return nil
}

func (p *executionPlan) planAction(name string, a *Action) error {
	for _, dependency := range a.Dependencies {
		if !p.m.HasAction(dependency) {
			return fmt.Errorf("action %s depends on %s, but %s is not defined (actions can only dependend on other actions)", name, dependency, dependency)
		}

		err := p.planAction(dependency, p.m.MustGetAction(dependency))

		if err != nil {
			return err
		}
	}

	for stepName, step := range a.Steps {
		err := p.planStep(newReference(name, stepName), step)

		if err != nil {
			return err
		}
	}

	return nil
}

func (p *executionPlan) planStep(ref reference, s *Step) error {
	deps := s.Dependencies
	deps = append(deps, s.Input)

	for _, dependency := range s.Dependencies {
		dependencyRef, err := resolveReference(ref.action, dependency)

		if err != nil {
			return err
		}

		if !p.hasReference(dependencyRef) {
			action, err := p.m.GetAction(dependencyRef.action)

			if err != nil {
				return err
			}

			if !action.HasStep(dependencyRef.step) {
				return fmt.Errorf("step %s depends on %s, but %s is not defined", ref.step, dependency, dependency)
			}

			dependencyStep, err := action.GetStep(dependencyRef.step)

			if err != nil {
				return err
			}

			err = p.planStep(dependencyRef, dependencyStep)

			if err != nil {
				return err
			}
		} else {
			// ignore dependency if it's already planned
		}
	}

	p.refs = append(p.refs, ref)

	return nil
}

func (p *executionPlan) hasReference(ref reference) bool {
	for _, r := range p.refs {
		if r.compare(ref) {
			return true
		}
	}

	return false
}
