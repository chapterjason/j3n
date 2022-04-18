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

package topology

import (
	"github.com/chapterjason/j3n/modx/mapx"
	"github.com/chapterjason/j3n/modx/slicex"
)

type DependencyGraph struct {
	nodes map[string][]string
}

func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		nodes: make(map[string][]string),
	}
}

func (dg *DependencyGraph) AddNode(name string) {
	if _, ok := dg.nodes[name]; !ok {
		dg.nodes[name] = []string{}
	}
}

func (dg *DependencyGraph) AddEdge(from, to string) {
	dg.nodes[from] = append(dg.nodes[from], to)
}

func (dg *DependencyGraph) IsCyclic() bool {
	visited := make(map[string]bool)
	recursionStack := make(map[string]bool)

	for key := range dg.nodes {
		visited[key] = false
		recursionStack[key] = false
	}

	for key := range dg.nodes {
		if !mapx.HasKey(visited, key) && dg.isCyclic(key, &visited, &recursionStack) {
			return true
		}
	}

	return false
}

func (dg *DependencyGraph) isCyclic(key string, visited *map[string]bool, recursionStack *map[string]bool) bool {
	if !mapx.HasKey(*visited, key) {
		(*visited)[key] = true
		(*recursionStack)[key] = true

		for _, dep := range dg.nodes[key] {
			if !mapx.HasKey(*visited, dep) {
				if dg.isCyclic(dep, visited, recursionStack) {
					return true
				}
			} else if mapx.HasKey(*recursionStack, dep) {
				return true
			}
		}
	}

	(*recursionStack)[key] = false

	return false
}

func (dg *DependencyGraph) Iterate() <-chan []string {
	ch := make(chan []string)

	go func() {
		nodes := dg.nodes

		for {
			var keys []string

			for key, deps := range nodes {
				if len(deps) == 0 {
					keys = append(keys, key)
				}
			}

			for _, key := range keys {
				delete(nodes, key)

				for i, m := range nodes {
					nodes[i] = slicex.RemoveString(m, key)
				}
			}

			if len(keys) > 0 {
				ch <- keys
			} else {
				close(ch)

				break
			}
		}
	}()

	return ch
}

func (dg *DependencyGraph) GetKeys() []string {
	keys := []string{}

	for key := range dg.nodes {
		keys = append(keys, key)
	}

	return keys
}

func (dg *DependencyGraph) GetDependencies(key string) []string {
	return dg.nodes[key]
}

func (dg *DependencyGraph) Add(d *DependencyGraph, key string) {
	deps := d.GetDependencies(key)

	dg.AddNode(key)

	for _, dep := range deps {
		dg.AddNode(dep)
		dg.AddEdge(key, dep)
	}

	for _, dep := range deps {
		dg.Add(d, dep)
	}
}
