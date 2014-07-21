package main

func CreateEmptySquareBoard(boardSize int32) BoardState {
	boardState := BoardState{}
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

func main() {
	println("hi")
	boardState := CreateEmptySquareBoard(4)
	println(len(boardState.RowStates))
}
