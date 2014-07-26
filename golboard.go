package main

import (
	"math/rand"
	"strings"
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

type cellChooser func() CellState

func CreateEmptySquareBoard(boardSize int) *BoardState {
	deadCellChooser := cellChooser(func() CellState {
		return Dead
	})
	return createSquareBoard(boardSize, deadCellChooser)
}

func CreateRandomSquareBoard(boardSize int) *BoardState {
	randomCellChooser := cellChooser(func() CellState {
		return rand.Int() % 2 == 0
	})
	return createSquareBoard(boardSize, randomCellChooser)
}

func createSquareBoard(boardSize int, cc cellChooser) *BoardState {
	boardState := &BoardState{}

	for i := 0; i < boardSize; i++ {
		rowState := []CellState{}
		for j := 0; j < boardSize; j++ {
			rowState = append(rowState, cc())
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
	heightA, widthA := GetDimensions(a)
	heightB, widthB := GetDimensions(b)
	return heightA == heightB && widthA == widthB
}

func GetDimensions(board *BoardState) (int, int) {
	height := len(board.States)
	width := 0
	if height > 0 {
		width = len(board.States[0])
	}
	return height, width
}
