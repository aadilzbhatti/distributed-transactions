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
	inUse map[int32][]*edge
}

var glock = &sync.RWMutex{}

func NewGraph() *Graph {
	return &Graph{make(map[string]*vertex, 0), make(map[int32]([]*edge), 0)}
}

func (g *Graph) CopyGraph() *Graph {
	res := NewGraph()
	for k, _ := range g.nodes {
		res.AddVertex(k)
		for n, v := range g.nodes[k].neighbors {
			res.nodes[k].neighbors[n] = v
		}
	}
	return res
}

func (g *Graph) AddVertex(id string) {
	if !g.IsVertexInGraph(id) {
		v := vertex{id, make(map[string]*vertex, 0)}
		glock.Lock()
		g.nodes[id] = &v
		glock.Unlock()
		fmt.Println(g)
	}
}

func (g *Graph) AddEdge(u string, v string, trans int32) error {
	glock.Lock()
	defer glock.Unlock()

	if v1, ok := g.nodes[u]; ok {
		if v2, ok2 := g.nodes[v]; ok2 {
      g.nodes[u].neighbors[v] = v2
			g.inUse[trans] = append(g.inUse[trans], &edge{v1, v2, trans})
			return nil
		} else {
			return fmt.Errorf("Vertex %v not in graph", v)
		}
	}
	return fmt.Errorf("Vertex %v not in graph\n", u)
}

func (g *Graph) RemoveEdge(u, v string) error {
	fmt.Println("Removing edge")
	glock.Lock()
	defer glock.Unlock()

	delete(g.nodes[u].neighbors, v)
	return nil
}

func (g *Graph) RemoveTransaction(trans int32) {
	for _, k := range g.inUse[trans] {
		g.RemoveEdge(k.start.id, k.end.id)
	}
	delete(g.inUse, trans)
}

func (g *Graph) IsVertexInGraph(name string) bool {
	glock.RLock()
	defer glock.RUnlock()

	_, ok := g.nodes[name]; return ok
}

func (g *Graph) DetectCycle(u, v string) bool {
	other := g.CopyGraph()
	other.RemoveEdge(u, v)
	fmt.Println(other)
	fmt.Printf("%+v, %+v in dc\n", other.nodes[u], other.nodes[v])
	return other.cycleHelper(other.nodes[u], other.nodes[v]) || other.cycleHelper(other.nodes[v], other.nodes[u])
}

func (g *Graph) cycleHelper(start *vertex, end *vertex) bool {
	if start.id == end.id {
		return false
	}
	seen := make(map[string]bool)
	stack := make([]vertex, 0)

	stack = append(stack, *start)
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:(len(stack) - 1)]
		if u.id == end.id {
			return true
		}
		if _, ok := seen[u.id]; !ok {
			seen[u.id] = true
			for _, v := range u.neighbors {
				stack = append(stack, *v)
			}
		}
	}

	return false
}
