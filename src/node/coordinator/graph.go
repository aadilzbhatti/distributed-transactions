package coordinator

import (
	"fmt"
	"sync"
)

type vertex struct {
  id string
	neighbors map[string]*vertex
}

type edge struct {
	start *vertex
	end   *vertex
	trans int32
}

type Graph struct {
	nodes map[string]*vertex
	edges map[int32]*edge
}

var glock = &sync.RWMutex{}

func NewGraph() *Graph {
	return &Graph{make(map[string]*vertex, 0), make(map[int32]*edge, 0)}
}

func (g *Graph) CopyGraph() *Graph {
	res := NewGraph()
	for k, _ := range g.nodes {
		res.AddVertex(k)
	}
	for _, e := range g.edges {
		res.AddEdge(e.start.id, e.end.id, e.trans)
	}
	return res
}

func (g *Graph) AddVertex(id string) {
	if !g.IsVertexInGraph(id) {
	  fmt.Println("Adding vertex", id)
		v := vertex{id, make(map[string]*vertex, 0)}
		glock.Lock()
		g.nodes[id] = &v
		glock.Unlock()
		fmt.Println(g)
	}
}

func (g *Graph) AddEdge(u string, v string, trans int32) error {
  fmt.Println("Adding edge")
	glock.Lock()
	defer glock.Unlock()

	if v1, ok := g.nodes[u]; ok {
		fmt.Println("OK")
		if v2, ok2 := g.nodes[v]; ok2 {
			fmt.Println("OK2")
			e := edge{v1, v2, trans}
			fmt.Println("MADE EDGE", e)
			g.edges[trans] = &e
			fmt.Println("ADDED EDGE TO GRAPH")
      g.nodes[u].neighbors[v] = v2
			fmt.Println("SET NEIGHBOR of u to v")
			g.nodes[v].neighbors[u] = v1
			fmt.Println("SET NEIGHBOR OF v to u")
			return nil
		} else {
			fmt.Println("YOU FUCKED UP BOY", v)
			return fmt.Errorf("Vertex %v not in graph\n", v)
		}
	}
	fmt.Println("YOU FUCKED UP BOY", u)
	return fmt.Errorf("Vertex %v not in graph\n", u)
}

func (g *Graph) RemoveEdge(trans int32) error {
	glock.Lock()
	defer glock.Unlock()

	if e, ok := g.edges[trans]; ok {
		delete(g.edges, trans)
		delete(g.nodes[e.start.id].neighbors, e.end.id)
		delete(g.nodes[e.end.id].neighbors, e.start.id)
	} else {
		return fmt.Errorf("Edge %v not in graph\n", trans)
	}

	return nil
}

func (g *Graph) IsVertexInGraph(name string) bool {
	glock.RLock()
	defer glock.RUnlock()

	_, ok := g.nodes[name]; return ok
}

func (g *Graph) DetectCycle(trans int32) bool {
  fmt.Println("Detecting cycle")
	glock.RLock()
	defer glock.RUnlock()

	other := g.CopyGraph()
	other.RemoveEdge(trans)
	fmt.Println(other)
	edge := g.edges[trans]
	return other.cycleHelper(edge.start, edge.end)
}

func (g *Graph) cycleHelper(start *vertex, end *vertex) bool {
	fmt.Println(start, end)
	if end.id == start.id {
		fmt.Println("FOUND SOMETHING 1")
		return false
	}
	fmt.Println(start.neighbors, end.neighbors)
	for _, v := range start.neighbors {
    fmt.Println(v)
		if end.id == v.id {
			fmt.Println("FOUND SOMETHING 2")
			return true
		}
		if g.cycleHelper(v, end) {
			fmt.Println("FOUND SOMETHING 3")
			return true
		}
	}
	fmt.Println("FOUND SOMETHING 4")
	return false
}
