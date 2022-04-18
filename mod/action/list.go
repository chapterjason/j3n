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
	"errors"

	"github.com/chapterjason/j3n/mod/topology"
)

var ErrActionNotFound = errors.New("action not found")

type List struct {
	Actions map[string]*Action `json:"actions" yaml:"actions"`
}

func (l *List) HasAction(actionName string) bool {
	_, ok := l.Actions[actionName]

	return ok
}

func (l *List) GetAction(actionName string) (*Action, error) {
	if !l.HasAction(actionName) {
		return nil, ErrActionNotFound
	}

	return l.Actions[actionName], nil
}

func (l *List) GetGraph() *topology.DependencyGraph {
	graph := topology.NewDependencyGraph()

	for actionName, action := range l.Actions {
		graph.AddNode(actionName)

		for _, dep := range action.Dependencies {
			graph.AddEdge(actionName, dep)
		}
	}

	return graph
}
