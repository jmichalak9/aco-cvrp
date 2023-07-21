package main

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

type point struct {
	x, y float64
}

var (
	capRG = regexp.MustCompile(`CAPACITY : (\d+)`)
	dimRG = regexp.MustCompile(`DIMENSION : (\d+)`)
	// TODO: capture float
	optRG    = regexp.MustCompile(`COMMENT : (\d+)`)
	graphRG  = regexp.MustCompile(`(\d+) (\d+) (\d+)`)
	demandRG = regexp.MustCompile(`(\d+) (\d+)`)
)

func getData(filepath string) (cap, opt, cities int, demand map[int]float64, graph map[int]point) {
	graph = make(map[int]point)
	demand = make(map[int]float64)
	dat, err := os.ReadFile(filepath)
	check(err)
	content := string(dat)

	x, err := strconv.Atoi(capRG.FindStringSubmatch(content)[1])
	check(err)
	cap = x

	x, err = strconv.Atoi(dimRG.FindStringSubmatch(content)[1])
	check(err)
	cities = x

	x, err = strconv.Atoi(optRG.FindStringSubmatch(content)[1])
	check(err)
	opt = x

	for _, match := range graphRG.FindAllStringSubmatch(content, -1) {
		numm := strings.Split(match[0], " ")
		var cityID int
		var point point
		cityID, _ = strconv.Atoi(numm[0])
		point.x, _ = strconv.ParseFloat(numm[1], 64)
		point.y, _ = strconv.ParseFloat(numm[2], 64)
		graph[cityID] = point
	}
	for _, match := range demandRG.FindAllStringSubmatch(content, -1) {
		numm := strings.Split(match[0], " ")
		var cityID int
		var demandd float64
		cityID, _ = strconv.Atoi(numm[0])
		demandd, _ = strconv.ParseFloat(numm[1], 64)
		demand[cityID] = demandd
	}
	return
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
