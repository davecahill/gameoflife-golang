package main

import (
	"net/http"
	"math/rand"
	"io/ioutil"
    "encoding/json"
	"log"
	"strings"
	"strconv"
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

func CreateEmptySquareBoard(boardSize int) *BoardState {
	deadCellChooser := CellChooser(func() cellState {
		return Dead
	})
	return createSquareBoard(boardSize, deadCellChooser)
}

func CreateRandomSquareBoard(boardSize int) *BoardState {
	randomCellChooser := CellChooser(func() cellState {
		return rand.Int() % 2 == 0
	})
	return createSquareBoard(boardSize, randomCellChooser)
}

func createSquareBoard(boardSize int, cellChooser CellChooser) *BoardState {
	boardState := &BoardState{}

	for i := 0; i < boardSize; i++ {
		rowState := []cellState{}
		for j := 0; j < boardSize; j++ {
			rowState = append(rowState, cellChooser())
		}
		boardState.States = append(boardState.States, rowState)
	}
	return boardState
}

// Board stepping

func step(oldBoard *BoardState) error {
	for i := 0; i < len(oldBoard.States); i++ {
		for j := 0; j < len(oldBoard.States[i]); j++ {
			oldBoard.States[i][j] = !oldBoard.States[i][j];
		}
	}
	return nil
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
		step(boardState)
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
