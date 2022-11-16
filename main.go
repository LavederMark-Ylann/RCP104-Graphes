package main

import (
	"fmt"
	"log"
	rand "math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	charts "github.com/go-echarts/go-echarts/v2/charts"
	components "github.com/go-echarts/go-echarts/v2/components"
	opts "github.com/go-echarts/go-echarts/v2/opts"
	types "github.com/go-echarts/go-echarts/v2/types"
)

type Graph struct {
	Nodes []Node
	Edges []Edge
}

type Node struct {
	Name string
}

type Edge struct {
	Source      Node
	Destination Node
	Weight      int
}

func main() {
	randomGraph := generateRandomGraph()
	printGraph(randomGraph)
}

func generateRandomGraph() (graph Graph) {
	// Seed pour avoir un semi-random
	rand.Seed(time.Now().UnixNano())
	minNodes, maxNodes := 4, 7
	nNodes := rand.Intn(maxNodes-minNodes+1) + minNodes
	// nombre max d'aretes = (nNodes * (nNodes - 1)) /
	minEdges, maxEdges := nNodes, (nNodes*(nNodes-1))/2
	nEdges := rand.Intn(maxEdges-minEdges+1) + minEdges
	// Création de nNodes sommets
	for i := 0; i < nNodes; i++ {
		graph.Nodes = append(graph.Nodes, Node{Name: fmt.Sprintf("%d", i+1)})
	}
	// Création de nEdges arêtes valides
	for i := 0; i < nEdges; i++ {
		var source, destination Node
		for {
			source = graph.Nodes[rand.Intn(len(graph.Nodes))]
			destination = graph.Nodes[rand.Intn(len(graph.Nodes))]
			if !relationExists(graph, source, destination) {
				break
			}
		}
		i, err := strconv.Atoi(source.Name)
		if err != nil {
			panic(err)
		}
		j, err := strconv.Atoi(destination.Name)
		if err != nil {
			panic(err)
		}
		// inversion des deux sommets si i > j
		if i > j {
			temp := source
			source = destination
			destination = temp
		}
		graph.Edges = append(graph.Edges, Edge{
			Source:      source,
			Destination: destination,
			Weight:      rand.Intn(6) + 1,
		})
	}
	return graph
}

func printGraph(g Graph) {
	graph := charts.NewGraph()
	graph.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Graph rendered"}),
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
	)
	convertedNodes := make([]opts.GraphNode, 0)
	for _, node := range g.Nodes {
		convertedNodes = append(convertedNodes, opts.GraphNode{Name: node.Name})
	}
	convertedEdges := make([]opts.GraphLink, 0)
	for _, edge := range g.Edges {
		convertedEdges = append(convertedEdges, opts.GraphLink{
			Source: edge.Source.Name,
			Target: edge.Destination.Name,
			Value:  float32(edge.Weight),
		})
		fmt.Println("Source: " + edge.Source.Name + " Target: " + edge.Destination.Name + " Value: " + strconv.Itoa(edge.Weight))
	}

	graph.AddSeries("graph", convertedNodes, convertedEdges,
		charts.WithGraphChartOpts(
			opts.GraphChart{
				Force:  &opts.GraphForce{Repulsion: 5000},
				Layout: "circular",
				Roam:   true,
			}),
		charts.WithLabelOpts(opts.Label{Show: true, Position: "right"}),
		charts.WithMarkPointStyleOpts(opts.MarkPointStyle{Label: &opts.Label{Show: true}}),
	)

	page := components.NewPage()
	page.AddCharts(graph)
	f, err := os.Create("graph.html")

	if err != nil {
		panic(err)
	}
	page.Render(f)
	openBrowser("graph.html")
}

func relationExists(graph Graph, node1 Node, node2 Node) (value bool) {
	for _, edge := range graph.Edges {
		if (edge.Source == node1 && edge.Destination == node2) || (edge.Source == node2 && edge.Destination == node1) {
			return true
		}
	}
	return false
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
