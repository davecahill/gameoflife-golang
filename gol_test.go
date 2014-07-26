package main

import (
	"testing"
	"io/ioutil"
	"strings"
	"errors"
	"fmt"
)

func TestStep(t *testing.T) {
	testTransition("middleDot", t)
}

func testTransition(boardName string, t *testing.T) {
	before, expectedAfter, err := readBoards(boardName)
	if err != nil {
		t.Error(err)
	}

	Step(before)

	if before != expectedAfter {
		t.Error(fmt.Sprintf("Stepping the board did not produce the expected outcome.\nActual:\n%s\nExpected:\n%s", BoardToString(before), BoardToString(expectedAfter)))
	}
}

func readBoards(boardName string) (*BoardState, *BoardState, error) {
	body, _ := ioutil.ReadFile("test_boards/" + boardName +  ".board")
	lines := strings.Split(string(body), "\n")

	// Number of lines should be uneven - x rows for board one, one empty middle row,
	if len(lines) % 2 == 0 {
		return nil, nil, errors.New("Badly formatted file; should have uneven number of lines")
	}

	middleLineIndex := len(lines) / 2

	if len(lines[middleLineIndex]) > 0 {
		return nil, nil, errors.New("Badly formatted file; middle line should be empty")
	}

	beforeBoard, err := readBoard(lines[0:middleLineIndex])
	if err != nil {
		return nil, nil, errors.New("First board is invalid: " + err.Error())
	}

	afterBoard, err := readBoard(lines[middleLineIndex + 1:])
	if err != nil {
		return nil, nil, errors.New("Second board is invalid: " + err.Error())
	}

	if !DimensionsEqual(beforeBoard, afterBoard) {
		return nil, nil, errors.New("Boards have different dimensions")
	}

	return beforeBoard, afterBoard, nil
}

func readBoard(boardLines []string) (*BoardState, error) {
	board, err := TextToBoard(boardLines)
	if err != nil {
		return nil, errors.New("Error reading board: " + err.Error())
	}

	return board, nil
}

