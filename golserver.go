package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
	"strings"
	"strconv"
)

func writeBoard(boardState *BoardState, w http.ResponseWriter) {
	data, err := json.Marshal(boardState)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(data)
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
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

func main() {
	http.HandleFunc("/new/", newHandler)
	http.ListenAndServe(":8080", nil)
}
