package main

import (
	"fmt"
)

type vertex struct {
	tid       int32
  id string
	neighbors map[string]vertex
}

type edge struct {
	start vertex
	end   vertex
}

type Graph struct {
	nodes map[string]vertex
	edges []edge
}

func New() Graph {
	return Graph{make(map[string]vertex, 0), make([]edge, 0)}
}

func (g Graph) AddVertex(tid int32, id string) {
	v := vertex{tid, id, make(map[string]vertex, 0)}
	g.nodes[id] = v
}

func (g Graph) AddEdge(u string, v string) error {
	if v1, ok := g.nodes[u]; ok {
		if v2, ok2 := g.nodes[v]; ok2 {
			e := edge{v1, v2}
			g.edges = append(g.edges, e)
      g.nodes[u].neighbors[v] = v2
			return nil
		} else {
			return fmt.Errorf("Vertex %v not in graph\n", v)
		}
	}
	return fmt.Errorf("Vertex %v not in graph\n", u)
}

func (g Graph) DetectCycle() {
	// TODO cycle detection
}
