package main

import (
	"net/http"
	"math/rand"
	"io/ioutil"
    "encoding/json"
	"log"
)

// Define board and cells
type cellState bool

const Alive cellState = true
const Dead cellState = true

type BoardState struct {
	States [][]cellState
}

// Board creation

type CellChooser func() cellState

func CreateEmptySquareBoard(boardSize int32) *BoardState {
	deadCellChooser := CellChooser(func() cellState {
		return Dead
	})
	return createSquareBoard(boardSize, deadCellChooser)
}

func CreateRandomSquareBoard(boardSize int32) *BoardState {
	randomCellChooser := CellChooser(func() cellState {
		return rand.Int() % 2 == 0
	})
	return createSquareBoard(boardSize, randomCellChooser)
}

func createSquareBoard(boardSize int32, cellChooser CellChooser) *BoardState {
	boardState := &BoardState{}

	for i := int32(0); i < boardSize; i++ {
		rowState := []cellState{}
		for j := int32(0); j < boardSize; j++ {
			rowState = append(rowState, cellChooser())
		}
		boardState.States = append(boardState.States, rowState)
	}
	return boardState
}

// Web server

func newHandler(w http.ResponseWriter, r *http.Request) {
	boardState := CreateRandomSquareBoard(4)
	data, err := json.Marshal(boardState)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(data)
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
