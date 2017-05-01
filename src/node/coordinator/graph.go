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
	fmt.Println("Called! with", id)
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
		if v2, ok2 := g.nodes[v]; ok2 {
			e := edge{v1, v2, trans}
			g.edges[trans] = &e
      g.nodes[u].neighbors[v] = v2
			g.nodes[v].neighbors[u] = v1
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
	fmt.Println("Removing edge")
	glock.Lock()
	defer glock.Unlock()

	if e, ok := g.edges[trans]; ok {
		delete(g.edges, trans)
		delete(g.nodes[e.start.id].neighbors, e.end.id)
		delete(g.nodes[e.end.id].neighbors, e.start.id)
	} else {
		return fmt.Errorf("Edge %v not in graph\n", trans)
	}

	fmt.Println("Edge removed")
	return nil
}

func (g *Graph) IsVertexInGraph(name string) bool {
	glock.RLock()
	defer glock.RUnlock()

	_, ok := g.nodes[name]; return ok
}

func (g *Graph) DetectCycle(trans int32) bool {
  fmt.Println("Detecting cycle")
	// glock.RLock()
	// defer glock.RUnlock()

	other := g.CopyGraph()
	fmt.Println("Old news baby", g.edges[trans], other.edges[trans])
	other.RemoveEdge(trans)
	fmt.Println(other)
	edge := g.edges[trans]
	fmt.Println(other.edges)
	return other.cycleHelper(edge.start, edge.end)
}

func (g *Graph) cycleHelper(start *vertex, end *vertex) bool {
	seen := make(map[string]bool)
	stack := make([]vertex, 0)

	stack = append(stack, *start)
	for len(stack) > 0 {
		fmt.Println(stack)
		fmt.Println(seen)
		u := stack[len(stack)-1]
		fmt.Println("Looking at", u.id)
		stack = stack[:(len(stack) - 1)]
		if u.id == end.id {
			fmt.Println("WE FOUND ONE GUYS!")
			return true
		}
		if _, ok := seen[u.id]; !ok {
			fmt.Println("WE ARE IN THE MONEY")
			seen[u.id] = true
			for _, v := range u.neighbors {
				stack = append(stack, *v)
			}
		}
	}

	fmt.Println("I AINT SEEN SHIT BOI")
	return false
}
