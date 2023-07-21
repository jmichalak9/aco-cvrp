package main

import (
	"fmt"
	wr "github.com/mroth/weightedrand"
	"math"
	"math/rand"
)

type ant struct {
	startingCap float64
	graph       graph
	cfg         config

	citiesLeft    []int
	cap           float64
	totalPathCost float64
	routes        []route
}

type route []int

func newAnt(graph graph, cap float64, cfg config) ant {
	a := ant{
		graph:       graph,
		startingCap: cap,
		cfg:         cfg,
	}
	a.resetState()
	return a
}

func (a *ant) getAvailableCities() []int {
	var cities []int
	for _, city := range a.citiesLeft {
		if a.cap >= a.graph.demandMap[city] {
			cities = append(cities, city)
		}
	}

	return cities
}

func (a *ant) selectFirstCity() int {
	cities := a.getAvailableCities()

	return cities[rand.Intn(len(cities))]
}

func (a *ant) resetState() {
	a.cap = a.startingCap
	a.citiesLeft = keys2(a.graph.adjacencyMap)
	a.removeFromCitiesLeft(1)
	a.routes = []route{}
	a.totalPathCost = 0
}

func keys2(m map[int]map[int]float64) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (a *ant) findSolution() solution {
	a.startNewRoute()

	for len(a.citiesLeft) != 0 {
		l := len(a.routes) - 1
		ll := len(a.routes[l]) - 1
		current := a.routes[l][ll]
		next := a.selectNextCity(current)
		if next == -1 {
			a.moveToCity(current, 1)
			a.startNewRoute()
		} else {
			a.moveToCity(current, next)
		}
	}
	l := len(a.routes) - 1
	ll := len(a.routes[l]) - 1
	a.moveToCity(a.routes[l][ll], 1)
	return solution{
		routes: a.routes,
		cost:   a.totalPathCost,
	}
}

func (a *ant) moveToCity(current, next int) {
	l := len(a.routes) - 1
	a.routes[l] = append(a.routes[l], next)
	if next != 1 {
		a.removeFromCitiesLeft(next)
	}
	a.cap -= a.graph.demandMap[next]
	a.totalPathCost += a.graph.adjacencyMap[current][next]
}

func (a *ant) removeFromCitiesLeft(next int) {
	for i, v := range a.citiesLeft {
		if v == next {
			a.citiesLeft = append(a.citiesLeft[:i], a.citiesLeft[i+1:]...)
			return
		}
	}
}

func (a *ant) startNewRoute() {
	a.cap = a.startingCap
	a.routes = append(a.routes, route{1})
	first := a.selectFirstCity()

	a.moveToCity(1, first)
}

func (a *ant) selectNextCity(current int) int {
	available := a.getAvailableCities()

	if len(available) == 0 {
		return -1
	}
	var scores []float64

	for _, city := range available {

		if _, ok := a.graph.adjacencyMap[current][city]; city == current && ok {
			return city
		}

		score := math.Pow(a.graph.pheromoneMap[current][city], a.cfg.alpha) *
			math.Pow(1.0/a.graph.adjacencyMap[current][city], a.cfg.beta)
		if math.IsInf(score, 1) {
			return city
		}
		scores = append(scores, score)
	}

	denom := sum(scores)
	var probs []float64
	for _, v := range scores {
		probs = append(probs, v/denom)
	}
	var choices []wr.Choice
	for i := range probs {
		choices = append(choices, wr.Choice{
			Item:   available[i],
			Weight: uint(probs[i] * 1_000_000_000),
		})
	}
	chooser, err := wr.NewChooser(choices...)
	if err != nil {
		fmt.Println(scores)
	}
	check(err)
	next := chooser.Pick().(int)
	return next
}

func sum(slice []float64) float64 {
	var result float64
	for _, numb := range slice {
		result += numb
	}
	return result
}
