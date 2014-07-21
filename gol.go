package main

import (
	"net/http"
	"code.google.com/p/goprotobuf/proto"
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

func newHandler(w http.ResponseWriter, r *http.Request) {
	boardState := CreateEmptySquareBoard(4)
	data, _ := proto.Marshal(boardState)
	w.Write(data)
}

func main() {
	http.HandleFunc("/new/", newHandler)
	http.ListenAndServe(":8080", nil)
}
