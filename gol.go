package main

import (
	"net/http"
	"math/rand"
	"io/ioutil"
    "encoding/json"
	"log"
	"strings"
	"strconv"
	"errors"
	"fmt"
	"bytes"
)

// Define board and cells
type CellState bool

const Alive CellState = true
const Dead CellState = false

type BoardState struct {
	States [][]CellState
}

// Board creation

type CellChooser func() CellState

func CreateEmptySquareBoard(boardSize int) *BoardState {
	deadCellChooser := CellChooser(func() CellState {
		return Dead
	})
	return createSquareBoard(boardSize, deadCellChooser)
}

func CreateRandomSquareBoard(boardSize int) *BoardState {
	randomCellChooser := CellChooser(func() CellState {
		return rand.Int() % 2 == 0
	})
	return createSquareBoard(boardSize, randomCellChooser)
}

func createSquareBoard(boardSize int, cellChooser CellChooser) *BoardState {
	boardState := &BoardState{}

	for i := 0; i < boardSize; i++ {
		rowState := []CellState{}
		for j := 0; j < boardSize; j++ {
			rowState = append(rowState, cellChooser())
		}
		boardState.States = append(boardState.States, rowState)
	}
	return boardState
}

func TextToBoard(lines []string) (*BoardState, error) {
	boardState := &BoardState{}

	for i := 0; i < len(lines); i++ {
		rowState := []CellState{}
		for j := 0; j < len(lines[i]); j++ {
			cellState, err := charToCellState(lines[i][j])
			if err != nil {
				return nil, errors.New("Bad character encountered in input: " + err.Error())
			}
			rowState = append(rowState, cellState)
		}
		boardState.States = append(boardState.States, rowState)
	}
	return boardState, nil
}

func BoardToString(board *BoardState) string {
	return strings.Join(BoardToText(board), "\n")
}

func BoardToText(board *BoardState) []string {
	text := []string{}

	for i := 0; i < len(board.States); i++ {
		var buffer bytes.Buffer
		for j := 0; j < len(board.States[i]); j++ {
			c := cellStateToChar(board.States[i][j])
			buffer.WriteByte(c)
		}
		text = append(text, buffer.String())
	}
	return text
}

func cellStateToChar(b CellState) byte {
	if (b == Alive) {
		return 'x'
	} else {
		return '-'
	}
}

func charToCellState(c byte) (CellState, error) {
	switch {
	case c == '-':
		return Dead, nil
	case c == 'x':
		return Alive, nil
	}
	return Alive, errors.New(fmt.Sprintf("Invalid character %s", c))
}

func DimensionsEqual(a *BoardState, b *BoardState) bool {
	heightA, widthA := getDimensions(a)
	heightB, widthB := getDimensions(b)
	return heightA == heightB && widthA == widthB
}

// Returns height, width
func getDimensions(board *BoardState) (int, int) {
	height := len(board.States)
	width := 0
	if height > 0 {
		width = len(board.States[0])
	}
	return height, width
}

// Board stepping

func Step(oldBoard *BoardState) *BoardState {
	newBoard := deepCopy(oldBoard)

	for i := 0; i < len(oldBoard.States); i++ {
		for j := 0; j < len(oldBoard.States[i]); j++ {
			// Calculate number of neighbors
			oldCellState := oldBoard.States[i][j]
			oldNumLiveNeighbors := getNumLiveNeighbors(oldBoard, i, j)
			newCellState := getNextState(oldCellState, oldNumLiveNeighbors)
			newBoard.States[i][j] = newCellState
		}
	}
	return newBoard
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

type addr struct {
	row int
	col int
}

/*
The universe of the Game of Life is an infinite two-dimensional orthogonal grid
of square cells, each of which is in one of two possible states, alive or dead.

Every cell interacts with its eight neighbours, which are the cells that are
horizontally, vertically, or diagonally adjacent.
*/
func getNumLiveNeighbors(board *BoardState, i int, j int) int {
	height, width := getDimensions(board)

	below := (i + 1) % height
	above := (i - 1 + height) % height
	right := (j + 1) % width
	left := (j - 1 + width) % width

	neighborAddresses := []addr{{below, left}, {below, j}, {below, right}, {i, left}, {i, right}, {above, left}, {above, j}, {above, right}}

	numLiveNeighbors := 0
	for idx := 0; idx < len(neighborAddresses); idx++ {
		numLiveNeighbors += getNumLiveCells(board, neighborAddresses[idx])
	}
	return numLiveNeighbors;
}

func getNumLiveCells(board *BoardState, a addr) int {
	if board.States[a.row][a.col] == Alive {
		return 1
	} else {
		return 0
	}
}

// Web server

func writeBoard(boardState *BoardState, w http.ResponseWriter) {
	data, err := json.Marshal(boardState)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(data)
}

func newHandler(w http.ResponseWriter, r *http.Request) {

	switch {
	case r.Method == "POST":
		body, _ := ioutil.ReadAll(r.Body)
		boardState := &BoardState{}
		json.Unmarshal(body, boardState)
		newBoardState := Step(boardState)
		writeBoard(newBoardState, w)
	case r.Method == "GET":
		end := strings.TrimPrefix(r.URL.Path, "/new/")
		num, _ := strconv.Atoi(end)
		boardState := CreateRandomSquareBoard(num)
		writeBoard(boardState, w)
	}
}

func baseHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadFile("index.html")
	w.Write(body)
}

func main() {
	http.HandleFunc("/new/", newHandler)
	http.HandleFunc("/", baseHandler)
	http.ListenAndServe(":8080", nil)
}
