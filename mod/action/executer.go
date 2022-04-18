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
	"sync"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/chapterjason/j3n/mod/topology"
)

var ErrOutputNotFound = errors.New("output not found")

type Executer struct {
	list    *List
	storage map[string]any
}

func NewExecuter(list *List) *Executer {
	return &Executer{
		list:    list,
		storage: make(map[string]any),
	}
}

func (e *Executer) Execute(actionName string) (map[string]map[string]error, error) {
	if _, ok := e.list.Actions[actionName]; !ok {
		return nil, fmt.Errorf("action %s not found", actionName)
	}

	ldg := e.list.GetGraph()

	if ldg.IsCyclic() {
		return nil, errors.New("cyclic dependency detected")
	}

	adg := topology.NewDependencyGraph()
	adg.Add(ldg, actionName)

	results := map[string]map[string]error{}

	for items := range adg.Iterate() {
		wg := sync.WaitGroup{}

		actionNames := []string{}

		for _, item := range items {
			wg.Add(1)

			actionNames = append(actionNames, item)

			go func(actionName string) {
				if ers := e.ExecuteAction(actionName); ers != nil {
					results[actionName] = ers
				}

				wg.Done()
			}(item)
		}

		wg.Wait()

		ers := map[string]map[string]error{}

		for _, actionName := range actionNames {
			for outputName, err := range results[actionName] {
				if ers[actionName] == nil {
					ers[actionName] = map[string]error{}
				}

				ers[actionName][outputName] = err
			}
		}

		if len(ers) > 0 {
			return ers, nil
		}
	}

	return nil, nil
}

func (e *Executer) ExecuteStep(action *Action, stepName string) error {
	step, err := action.GetStep(stepName)

	if err != nil {
		return err
	}

	log.Infof("executing step %s", stepName)

	var input any

	if step.Input != "" {
		var err error

		input, err = e.GetOutput(step.Input)

		if err != nil {
			return errors.Wrapf(err, "failed to get input %s", step.Input)
		}
	}

	stepRunner, ok := Steps[step.Type]

	if !ok {
		return fmt.Errorf("no runner for step %s and type %s", stepName, step.Type)
	}

	out, err := stepRunner(input, step.Params)

	if err != nil {
		return err
	}

	if step.Output != "" {
		if out == nil {
			return fmt.Errorf("output of step %s is nil", stepName)
		}

		e.storage[step.Output] = out
	}

	log.Debugf("step %s executed", stepName)

	return nil
}

func (e *Executer) GetOutput(key string) (any, error) {
	v, ok := e.storage[key]

	if !ok {
		return nil, ErrOutputNotFound
	}

	return v, nil
}

func (e *Executer) ExecuteAction(actionName string) map[string]error {
	log.Infof("executing action %s", actionName)

	action, err := e.list.GetAction(actionName)

	if err != nil {
		return map[string]error{actionName: err}
	}

	sdg := action.GetGraph()

	if sdg.IsCyclic() {
		return map[string]error{actionName: errors.New("cyclic dependency")}
	}

	results := map[string]error{}

	for items := range sdg.Iterate() {
		wg := sync.WaitGroup{}

		stepNames := []string{}

		for _, item := range items {
			wg.Add(1)

			stepNames = append(stepNames, item)

			go func(action *Action, stepName string) {
				if err := e.ExecuteStep(action, stepName); err != nil {
					results[stepName] = err
				}

				wg.Done()
			}(action, item)
		}

		wg.Wait()

		ers := map[string]error{}

		for _, stepName := range stepNames {
			if err, ok := results[stepName]; ok {
				ers[stepName] = err
			}
		}

		if len(ers) > 0 {
			return ers
		}
	}

	log.Debugf("action %s executed", actionName)

	return nil
}
