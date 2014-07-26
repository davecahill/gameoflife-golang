package main

func Step(oldBoard *BoardState) *BoardState {
	newBoard := deepCopy(oldBoard)

	for i := 0; i < len(oldBoard.States); i++ {
		for j := 0; j < len(oldBoard.States[i]); j++ {
			oldCellState := oldBoard.States[i][j]
			oldNumLiveNeighbors := getNumLiveNeighbors(oldBoard, i, j)
			newCellState := getNextState(oldCellState, oldNumLiveNeighbors)
			newBoard.States[i][j] = newCellState
		}
	}
	return newBoard
}

func getNextState(oldCellState CellState, oldNumLiveNeighbors int) CellState {

	if oldCellState == Alive {
		if oldNumLiveNeighbors < 2 {
			// Any live cell with fewer than two live neighbours dies, as if caused by under-population.
			return Dead
		} else if oldNumLiveNeighbors == 2 || oldNumLiveNeighbors == 3 {
			// Any live cell with two or three live neighbours lives on to the next generation.
			return Alive
		} else {
			// Any live cell with more than three live neighbours dies, as if by overcrowding.
			return Dead
		}
	} else {
		if oldNumLiveNeighbors == 3 {
			// Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.
			return Alive
		} else {
			return Dead
		}
	}
}

func deepCopy(oldBoard *BoardState) *BoardState {
	newBoard := &BoardState{}

	for i := 0; i < len(oldBoard.States); i++ {
		rowState := []CellState{}
		for j := 0; j < len(oldBoard.States[i]); j++ {
			rowState = append(rowState, oldBoard.States[i][j])
		}
		newBoard.States = append(newBoard.States, rowState)
	}
	return newBoard
}

/*
The universe of the Game of Life is an infinite two-dimensional orthogonal grid
of square cells, each of which is in one of two possible states, alive or dead.

Every cell interacts with its eight neighbours, which are the cells that are
horizontally, vertically, or diagonally adjacent.
*/
func getNumLiveNeighbors(board *BoardState, i int, j int) int {
	height, width := GetDimensions(board)

	below := (i + 1) % height
	above := (i - 1 + height) % height
	right := (j + 1) % width
	left := (j - 1 + width) % width

	neighborAddresses := []cellLocation{{below, left}, {below, j}, {below, right}, {i, left}, {i, right}, {above, left}, {above, j}, {above, right}}

	numLiveNeighbors := 0
	for idx := 0; idx < len(neighborAddresses); idx++ {
		numLiveNeighbors += getNumLiveCells(board, neighborAddresses[idx])
	}
	return numLiveNeighbors;
}

func getNumLiveCells(board *BoardState, a cellLocation) int {
	if board.States[a.row][a.col] == Alive {
		return 1
	} else {
		return 0
	}
}

type cellLocation struct {
	row int
	col int
}
