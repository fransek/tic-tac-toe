package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type TileState string

const (
	Empty  TileState = ""
	Player TileState = "X"
	AI     TileState = "O"
)

type BoardState map[int]TileState

func main() {
	fmt.Println("http://localhost:8080/")
	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/move", move)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func move(c *gin.Context) {
	var body BoardState

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tileIndex := calculateBestMove(&body)

	c.JSON(http.StatusOK, gin.H{
		"tileIndex": tileIndex,
	})
}

func checkForThreeInARow(
	boardState *BoardState,
	startingTileIndex int,
	increment int,
	player TileState,
) bool {
	tile1 := (*boardState)[startingTileIndex]
	tile2 := (*boardState)[startingTileIndex+increment]
	tile3 := (*boardState)[startingTileIndex+2*increment]

	return tile1 == player && tile2 == player && tile3 == player
}

func checkForWinner(boardState *BoardState, player TileState) bool {
	for i := 0; i < 3; i++ {
		if checkForThreeInARow(boardState, i*3, 1, player) {
			return true
		}
		if checkForThreeInARow(boardState, i, 3, player) {
			return true
		}
	}
	if checkForThreeInARow(boardState, 0, 4, player) {
		return true
	}
	if checkForThreeInARow(boardState, 2, 2, player) {
		return true
	}
	return false
}

func aiWins(boardState *BoardState) bool {
	return checkForWinner(boardState, AI)
}

func playerWins(boardState *BoardState) bool {
	return checkForWinner(boardState, Player)
}

func calculateBestMove(boardState *BoardState) int {
	return -1
}

func getEmptyTiles(boardState *BoardState) []int {
	emptyTiles := []int{}
	for i := 0; i < 9; i++ {
		if (*boardState)[i] == Empty {
			emptyTiles = append(emptyTiles, i)
		}
	}
	return emptyTiles
}

func getWinningOutcomes(boardState *BoardState) []int {
	winningOutcomes := []int{}
	for _, index := range getEmptyTiles(boardState) {
		if checkForWinner(boardState, Player) {
			winningOutcomes = append(winningOutcomes, index)
		}
	}
	return winningOutcomes
}
