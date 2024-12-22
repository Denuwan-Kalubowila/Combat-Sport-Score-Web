package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Player struct {
	Color    string
	Score    int
	IsWinner bool
}

type Match struct {
	RedPlayed  Player
	BluePlayer Player
}

type Score struct {
	PlayerColor string `json:"playerColor"`
	Points      int    `json:"points"`
}

type Winner struct {
	PlayerColor string `json:"playerColor"`
	IsWinner    bool   `json:"isWinner"`
}

var (
	redPlayer  = Player{Color: "red", Score: 0, IsWinner: false}
	bluePlayer = Player{Color: "blue", Score: 0, IsWinner: false}
)

var matches = make(map[string]*Match)

func main() {

	// Update score by sending a POST request to /updateScore
	http.HandleFunc("/updateScore", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		var scoreUpdate Score
		err := json.NewDecoder(r.Body).Decode(&scoreUpdate)
		if err != nil {
			http.Error(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest)
			return
		}

		switch scoreUpdate.PlayerColor {
		case "red":
			redPlayer.Score += scoreUpdate.Points
		case "blue":
			bluePlayer.Score += scoreUpdate.Points
		default:
			http.Error(w, "Invalid player color", http.StatusBadRequest)
			return
		}
		fmt.Println(redPlayer.Score)
		fmt.Println(bluePlayer.Score)

		// Check if there is a winner
		checkWinner()

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Score updated"))
	})

	// Update winner by sending a POST request to /updateWinner
	http.HandleFunc("/updateWinner", func(w http.ResponseWriter, r *http.Request) {
		var Winner Winner
		err := json.NewDecoder(r.Body).Decode(&Winner)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if Winner.PlayerColor == "red" {
			redPlayer.IsWinner = Winner.IsWinner
		} else if Winner.PlayerColor == "blue" {
			bluePlayer.IsWinner = Winner.IsWinner
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Winner updated"))
	})

	// Stream score by sending a GET request to /streamScore
	http.HandleFunc("/streamScore", streamScoreHandler)

	http.ListenAndServe(":8080", nil)

}

// Stream score
func streamScoreHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		if redPlayer.IsWinner || bluePlayer.IsWinner {
			winner := fmt.Sprintf("data: {\"red\": %d, \"blue\": %d, \"winner\": \"%s\"}\n\n", redPlayer.Score, bluePlayer.Score, func() string {
				if redPlayer.IsWinner {
					return "red"
				} else {
					return "blue"
				}
			}())
			w.Write([]byte(winner))
			flusher.Flush()
			return
		}

		content := fmt.Sprintf("data: {\"red\": %d, \"blue\": %d}\n\n", redPlayer.Score, bluePlayer.Score)
		w.Write([]byte(content))
		flusher.Flush()
		time.Sleep(2 * time.Second)
	}
}

// Check if there is a winner
func checkWinner() {
	if redPlayer.Score >= 10 {
		redPlayer.IsWinner = true
	} else if bluePlayer.Score >= 10 {
		bluePlayer.IsWinner = true
	}
}
