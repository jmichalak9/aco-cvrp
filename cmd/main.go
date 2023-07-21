package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	rand.Seed(420)
	cfg := config{
		iterations:       3000,
		elite:            false,
		rank:             false,
		initialPheromone: 80,
		alpha:            2,
		beta:             5,
		rho:              0.2,
		sigma:            6,
	}
	runACO(cfg)
}

type solution struct {
	cost   float64
	routes []route
}

func runACO(cfg config) (solution, []solution) {
	cap, opt, cities, demand, graphData := getData("data/C11.txt")
	graph := newGraph(graphData, demand, cfg)
	var ants []ant
	for i := 0; i < cities; i++ {
		ants = append(ants, newAnt(graph, float64(cap), cfg))
	}
	cfg.sigma = cities
	if cfg.rank {
		cfg.sigma = 6
	}
	bestSolution := solution{cost: -1}
	var candidates []solution
	for i := 1; i < cfg.iterations+1; i++ {
		var solutions []solution
		ch := make(chan solution, 1000)
		var wg sync.WaitGroup
		for _, _ant := range ants {
			wg.Add(1)
			go func(ant ant) {
				defer wg.Done()
				ant.resetState()
				solution := ant.findSolution()
				//solutions = append(solutions, solution)
				ch <- solution
			}(_ant)
		}
		wg.Wait()
		close(ch)
		for s := range ch {
			solutions = append(solutions, s)
		}
		candidateBestSolution := min(solutions)
		candidates = append(candidates, candidateBestSolution)
		if bestSolution.cost == -1 || candidateBestSolution.cost < bestSolution.cost {
			bestSolution = candidateBestSolution
		}
		if i == 1 || i%10 == 0 {
			fmt.Printf("%v\n", bestSolution.cost)
		}
		graph.updatePheromoneMap(solutions, bestSolution)
	}
	fmt.Printf("opt: %d", opt)
	return bestSolution, candidates
}

func worker(solutions []solution, ch <-chan solution) {
	for s := range ch {
		solutions = append(solutions, s)
	}
}

func min(solutions []solution) solution {
	bestSolution := solution{cost: -1}
	for _, solution := range solutions {
		if bestSolution.cost == -1 || solution.cost < bestSolution.cost {
			bestSolution = solution
		}
	}
	return bestSolution
}

type config struct {
	iterations       int
	elite            bool
	rank             bool
	initialPheromone float64

	alpha float64
	beta  float64
	rho   float64
	sigma int
}
