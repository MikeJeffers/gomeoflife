package main

import (
	"errors"
	"math/rand"
)

/*
Defines a cell in a cellular automata system
For now this will implement a vanilla conways game of life
0 = DEAD
1 = ALIVE
TBD other states and rules w/ some elaborate config driven system
*/
type Cell struct {
	x, y, state int
}

func (c *Cell) drawState() (int, int, int) {
	return c.x, c.y, c.state
}

func (c *Cell) nextState(neighbors []Cell) {
	liveNeighbors := 0
	for _, neighbor := range neighbors {
		liveNeighbors += neighbor.state
	}
	// Typical Conway ruleset for moores neighborhood
	if c.state == 1 {
		if liveNeighbors < 2 || liveNeighbors > 3 {
			c.state = 0
		}
	} else {
		if liveNeighbors == 3 {
			c.state = 1
		}
	}
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
			if rand.Float32() > 0.7 {
				initialState = 1
			}
			g.cells = append(g.cells, Cell{x, y, initialState})
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
	cellsCopy := make([]Cell, len(g.cells))
	copy(cellsCopy, g.cells)
	// This is the embarrassingly parallelizable part, lets implement serial first
	for i, cell := range g.cells {
		neighbors := g.getNeighbors(cell.x, cell.y)
		cellsCopy[i].nextState(neighbors)
	}
	g.cells = cellsCopy
}
