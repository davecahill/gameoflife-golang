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

// Board stepping

func Step(oldBoard *BoardState) error {
	for i := 0; i < len(oldBoard.States); i++ {
		for j := 0; j < len(oldBoard.States[i]); j++ {
			oldBoard.States[i][j] = !oldBoard.States[i][j];
		}
	}
	return nil
}

func DimensionsEqual(a *BoardState, b *BoardState) bool {
	return len(a.States) == len(b.States) && (len(a.States) == 0 || (len(a.States[0]) == len(b.States[0])))
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
		Step(boardState)
		writeBoard(boardState, w)
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
