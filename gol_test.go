package main

import (
	"testing"
	"io/ioutil"
	"strings"
	"errors"
	"fmt"
	"reflect"
)

func TestAllBoards(t *testing.T) {
	files, err := ioutil.ReadDir("test_boards")
	if err != nil {
		t.Error("Error loading test boards")
	}
	for _, file := range files {
		println(fmt.Sprintf("Testing %s", file.Name()))
		testTransition(file.Name(), t)
	}
}

func testTransition(boardName string, t *testing.T) {
	before, expectedAfter, err := readBoards(boardName)
	if err != nil {
		t.Error(err)
	}

	actualAfter := Step(before)

	if !reflect.DeepEqual(actualAfter, expectedAfter) {
		t.Error(fmt.Sprintf("Stepping the board did not produce the expected outcome.\nExpected:\n%s\nActual:\n%s", BoardToString(expectedAfter), BoardToString(actualAfter)))
	}
}

func readBoards(boardName string) (*BoardState, *BoardState, error) {
	body, err := ioutil.ReadFile("test_boards/" + boardName)
	if (err != nil) {
		return nil, nil, errors.New(fmt.Sprintf("Error loading board: %s", boardName))
	}

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

