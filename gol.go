package main

import (
	"net/http"
	"code.google.com/p/goprotobuf/proto"
	"math/rand"
	"io/ioutil"
)

func CreateEmptySquareBoard(boardSize int32) *BoardState {
	boardState := &BoardState{}
	for i := int32(0); i < boardSize; i++ {
		rowState := &BoardState_RowState{};
		for j := int32(0); j < boardSize; j++ {
			cellState := CellState_DEAD
			rowState.CellStates = append(rowState.CellStates, cellState)
		}
		boardState.RowStates = append(boardState.RowStates, rowState)
	}
	return boardState
}

func CreateRandomSquareBoard(boardSize int32) *BoardState {
	boardState := &BoardState{}
	for i := int32(0); i < boardSize; i++ {
		rowState := &BoardState_RowState{};
		for j := int32(0); j < boardSize; j++ {
			if rand.Int() % 2 == 0 {
				rowState.CellStates = append(rowState.CellStates, CellState_DEAD)
			} else {
				rowState.CellStates = append(rowState.CellStates, CellState_ALIVE)
			}
		}
		boardState.RowStates = append(boardState.RowStates, rowState)
	}
	return boardState
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	boardState := CreateRandomSquareBoard(4)
	data, _ := proto.Marshal(boardState)
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
