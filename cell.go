package main

import (
	"errors"
	"math/rand"
	"sync"
)

type State interface {
	nextState([]Cell) State
	getValue() int
}
type AliveState struct {
	value int
}
type DeadState struct {
	value int
}
type AnotherState struct {
	value int
}

var DEAD = &DeadState{0}
var ALIVE = &AliveState{1}
var ANOTHER = &AnotherState{2}

func getStateByValue(value int) State {
	switch value {
	case 0:
		return DEAD
	case 1:
		return ALIVE
	case 2:
		return ANOTHER
	default:
		return DEAD
	}
}
func countMatchingState(neighbors []Cell, state State) int {
	liveNeighbors := 0
	for _, neighbor := range neighbors {
		if neighbor.state.getValue() == state.getValue() {
			liveNeighbors += 1
		}
	}
	return liveNeighbors
}
func mapNeighborStateCounts(neighbors []Cell) map[State]int {
	m := make(map[State]int)
	m[ALIVE] = 0
	m[DEAD] = 0
	m[ANOTHER] = 0
	for _, neighbor := range neighbors {
		m[neighbor.state] += 1
	}
	return m
}

func (s *AliveState) getValue() int {
	return s.value
}
func (s *DeadState) getValue() int {
	return s.value
}
func (s *AnotherState) getValue() int {
	return s.value
}

func (s *AliveState) nextState(neighbors []Cell) State {
	counts := mapNeighborStateCounts(neighbors)
	if counts[ALIVE] < 2 || counts[ALIVE] > 3 {
		return DEAD
	}
	return ALIVE
}

func (s *DeadState) nextState(neighbors []Cell) State {
	counts := mapNeighborStateCounts(neighbors)
	if counts[ALIVE] == 3 {
		return ALIVE
	}
	return DEAD
}

func (s *AnotherState) nextState(neighbors []Cell) State {
	return ANOTHER
}

/*
Defines a cell in a cellular automata system
*/
type Cell struct {
	x, y  int
	state State
}

func (c *Cell) drawState() (int, int, int) {
	return c.x, c.y, c.state.getValue()
}

func (c *Cell) nextState(neighbors []Cell) {
	c.state = c.state.nextState(neighbors)
}

/*
Represents a game state internally as a 1D array with length=cols*rows
*/
type Grid struct {
	rows, cols int
	cells      []Cell
}

func (g *Grid) init(cols, rows int) {
	//TODO bounds check input
	g.rows = rows
	g.cols = cols
	for y := 0; y < g.rows; y++ {
		for x := 0; x < g.cols; x++ {
			initialState := 0
			// Random seed
			r := rand.Float32()
			if r > 0.7 {
				initialState = 1
			} else if r < 0.05 {
				initialState = 2
			}
			g.cells = append(g.cells, Cell{x, y, getStateByValue(initialState)})
		}
	}
}

func (g *Grid) getCell(x, y int) (Cell, error) {
	if x >= g.cols || y >= g.rows || x < 0 || y < 0 {
		return Cell{}, errors.New("out of bounds")
	}
	return g.cells[(y*g.cols)+x], nil
}

func (g *Grid) getNeighbors(x, y int) []Cell {
	var neighbors []Cell
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			cell, err := g.getCell(x+i, y+j)
			if err != nil {
				continue
			}
			neighbors = append(neighbors, cell)
		}
	}
	return neighbors
}

func (g *Grid) nextState() {
	// Store a copy that will be the next state
	// Leaving original array unmutated
	cellsCopy := make([]Cell, len(g.cells))
	copy(cellsCopy, g.cells)
	var wg sync.WaitGroup
	wg.Add(len(g.cells))
	for i, cell := range g.cells {
		go func(i int, cell Cell) {
			defer wg.Done()
			neighbors := g.getNeighbors(cell.x, cell.y)
			cellsCopy[i].nextState(neighbors)
		}(i, cell)
	}
	wg.Wait()
	g.cells = cellsCopy
}
