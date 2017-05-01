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
	// edges map[int32]*edge
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
  fmt.Println("Adding edge", u, v)
	glock.Lock()
	defer glock.Unlock()

	if v1, ok := g.nodes[u]; ok {
		if v2, ok2 := g.nodes[v]; ok2 {
			// e := edge{v1, v2, trans}
			// g.edges[trans] = &e
      g.nodes[u].neighbors[v] = v2
			// g.nodes[v].neighbors[u] = v1
			g.inUse[trans] = append(g.inUse[trans], &edge{v1, v2, trans})
			fmt.Printf("%+v\n", g)
			fmt.Printf("%+v, %+v\n", g.nodes[u], g.nodes[v])
			return nil
		} else {
			fmt.Println("YOU FUCKED UP BOY", v)
			return fmt.Errorf("Vertex %v not in graph", v)
		}
	}
	fmt.Println("YOU FUCKED UP BOY", u)
	return fmt.Errorf("Vertex %v not in graph\n", u)
}

func (g *Graph) RemoveEdge(u, v string) error {
	fmt.Println("Removing edge")
	glock.Lock()
	defer glock.Unlock()

	// if e, ok := g.edges[trans]; ok {
		// delete(g.edges, trans)
		delete(g.nodes[u].neighbors, v)
		// delete(g.nodes[v].neighbors, u)
	// } else {
		// return fmt.Errorf("Edge %v not in graph", trans)
	// }

	fmt.Println("Edge removed")
	return nil
}

func (g *Graph) RemoveTransaction(trans int32) {
	// glock.Lock()
	// defer glock.Unlock()
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
  fmt.Println("Detecting cycle")
	// glock.RLock()
	// defer glock.RUnlock()

	other := g.CopyGraph()
	// fmt.Println("Old news baby", g.edges[trans], other.edges[trans])
	// edge := other.edges[trans]
	other.RemoveEdge(u, v)
	fmt.Println(other)
	fmt.Printf("%+v, %+v in dc\n", other.nodes[u], other.nodes[v])
	// fmt.Println(other.edges)
	// fmt.Println(edge.start, edge.end)
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
		fmt.Println(stack)
		fmt.Println(seen)
		u := stack[len(stack)-1]
		fmt.Println("Looking at", u.id)
		fmt.Printf("%+v\n", u.neighbors)
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
				fmt.Println("HEY NOW", stack)
			}
		}
	}

	fmt.Println("I AINT SEEN SHIT BOI")
	return false
}
