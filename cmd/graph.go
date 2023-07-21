package main

import (
	"math"
	"sort"
)

type graph struct {
	demandMap    map[int]float64
	adjacencyMap map[int]map[int]float64
	pheromoneMap map[int]map[int]float64
	cfg          config
}

func newGraph(graphData map[int]point, demand map[int]float64, cfg config) graph {
	return graph{
		cfg:          cfg,
		adjacencyMap: createAdjacencyMap(graphData),
		pheromoneMap: createPheromoneMap(keys(graphData), cfg),
		demandMap:    demand,
	}
}

func keys(m map[int]point) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func createAdjacencyMap(graphData map[int]point) map[int]map[int]float64 {
	adjacencyMap := make(map[int]map[int]float64)
	for cityID, city := range graphData {
		if adjacencyMap[cityID] == nil {
			adjacencyMap[cityID] = make(map[int]float64)
		}
		for otherCityID, otherCity := range graphData {
			if cityID == otherCityID {
				continue
			}
			dist := getDistance(city, otherCity)
			if adjacencyMap[otherCityID] == nil {
				adjacencyMap[otherCityID] = make(map[int]float64)
			}
			adjacencyMap[cityID][otherCityID] = dist
			adjacencyMap[otherCityID][cityID] = dist
		}
	}
	return adjacencyMap
}

func createPheromoneMap(nodes []int, cfg config) map[int]map[int]float64 {
	pheromoneMap := make(map[int]map[int]float64)
	for _, cityID := range nodes {
		if pheromoneMap[cityID] == nil {
			pheromoneMap[cityID] = make(map[int]float64)
		}

		for _, otherCityID := range nodes {
			if cityID == otherCityID {
				continue
			}
			if pheromoneMap[otherCityID] == nil {
				pheromoneMap[otherCityID] = make(map[int]float64)
			}

			pheromoneMap[cityID][otherCityID] = cfg.initialPheromone
			pheromoneMap[otherCityID][cityID] = cfg.initialPheromone
		}
	}
	return pheromoneMap
}

func (g *graph) updatePheromoneMap(solutions []solution, bestSolution solution) {
	nodes := sorted(keys2(g.pheromoneMap))
	for i, n := range nodes {
		for _, n2 := range nodes[i+1:] {
			newVal := math.Max((1-g.cfg.rho)*g.pheromoneMap[n][n2], 1e-10)
			g.pheromoneMap[n][n2] = newVal
			g.pheromoneMap[n2][n] = newVal
		}
	}
	type edge struct {
		from, to int
	}

	if !g.cfg.rank {
		for _, solution := range solutions {
			pi := g.cfg.initialPheromone / solution.cost
			for _, route := range solution.routes {
				var edges []edge
				for i := 0; i < len(route)-1; i++ {
					edges = append(edges, edge{route[i], route[i+1]})
				}
				for _, edge := range edges {
					g.pheromoneMap[edge.from][edge.to] += pi
					g.pheromoneMap[edge.to][edge.from] += pi
				}

			}
		}
	} else {
		sort.Slice(solutions, func(i, j int) bool {
			return solutions[i].cost < solutions[j].cost
		})
		for mi, solution := range solutions[:g.cfg.sigma-1] {
			pi := float64(g.cfg.sigma-mi-1) * g.cfg.initialPheromone / solution.cost
			for _, route := range solution.routes {
				var edges []edge
				for i := 0; i < len(route)-1; i++ {
					edges = append(edges, edge{route[i], route[i+1]})
				}
				for _, edge := range edges {
					g.pheromoneMap[edge.from][edge.to] += pi
					g.pheromoneMap[edge.to][edge.from] += pi
				}
			}

		}
	}
	if g.cfg.elite || g.cfg.rank {
		for _, route := range bestSolution.routes {
			var edges []edge
			for i := 0; i < len(route)-1; i++ {
				edges = append(edges, edge{route[i], route[i+1]})
			}
			for _, edge := range edges {
				g.pheromoneMap[edge.from][edge.to] += float64(g.cfg.sigma) * g.cfg.initialPheromone / bestSolution.cost
				g.pheromoneMap[edge.to][edge.from] += float64(g.cfg.sigma) * g.cfg.initialPheromone / bestSolution.cost
			}
		}
	}
}

func getDistance(p1, p2 point) float64 {
	return math.Sqrt(math.Pow(p1.x-p2.x, 2) + math.Pow(p1.y-p2.y, 2))
}

func sorted(nodes []int) []int {
	sort.Ints(nodes)
	return nodes
}
