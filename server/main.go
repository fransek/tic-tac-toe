package main

import (
	"fmt"
	"math"
	"net/http"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type TileState string

const (
	Empty    TileState = ""
	Opponent TileState = "X"
	AI       TileState = "O"
)

type BoardState map[int]TileState

func main() {
	fmt.Println("http://localhost:8080/")
	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/move", move)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	// test()
}

func test() {
	exampleBoardState := &BoardState{
		0: AI,
		1: AI,
		2: Empty,
		3: Opponent,
		4: Empty,
		5: Opponent,
		6: Empty,
		7: Empty,
		8: Empty,
	}

	fmt.Println(getScores(exampleBoardState))
}

func move(c *gin.Context) {
	var body BoardState

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tileIndex": calculateBestMove(&body),
		"scores":    getScores(&body),
		"board":     body,
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

func opponentWins(boardState *BoardState) bool {
	return checkForWinner(boardState, Opponent)
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

func copyBoardState(boardState *BoardState) BoardState {
	newBoardState := BoardState{}
	for i := 0; i < 9; i++ {
		newBoardState[i] = (*boardState)[i]
	}
	return newBoardState
}

func simulateMove(boardState *BoardState, player TileState, index int) BoardState {
	newBoardState := copyBoardState(boardState)
	newBoardState[index] = player
	return newBoardState
}

func permutations(nums []int) [][]int {
	var result [][]int
	permute(nums, 0, len(nums)-1, &result)
	return result
}

func permute(nums []int, l, r int, result *[][]int) {
	if l == r {
		// Create a copy of the current permutation and append it to the result
		permutation := make([]int, len(nums))
		copy(permutation, nums)
		*result = append(*result, permutation)
	} else {
		for i := l; i <= r; i++ {
			// Swap the current element with the first element
			// Then recursively permute the remaining elements
			nums[l], nums[i] = nums[i], nums[l]
			permute(nums, l+1, r, result)
			// Undo the swap to restore the original order
			nums[l], nums[i] = nums[i], nums[l]
		}
	}
}

func createScoreMap(arr []int) map[int]int {
	result := make(map[int]int)

	for _, num := range arr {
		result[num] = 0
	}

	return result
}

func getScores(boardState *BoardState) map[int]int {
	emptyTiles := getEmptyTiles(boardState)
	sequences := permutations(emptyTiles)
	scores := createScoreMap(emptyTiles)

	// Create a channel to receive scores from goroutines
	scoreChan := make(chan map[int]int)

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Launch a goroutine for each sequence
	for _, sequence := range sequences {
		wg.Add(1)
		go func(sequence []int) {
			defer wg.Done()
			maxScore := len(sequence)

			turns, result, firstMove := getOutcome(boardState, sequence)
			scoreMap := createScoreMap(emptyTiles)

			if result == Win {
				scoreMap[firstMove] += maxScore - turns
			} else if result == Loss {
				scoreMap[firstMove] -= maxScore - turns
			}

			// Send the score map to the channel
			scoreChan <- scoreMap
		}(sequence)
	}

	// Start a goroutine to collect scores from the channel and merge them
	go func() {
		for scoreMap := range scoreChan {
			for key, value := range scoreMap {
				scores[key] += value
			}
		}
	}()

	// Wait for all goroutines to finish
	wg.Wait()

	// Close the score channel
	close(scoreChan)

	return scores
}

type Result string

const (
	Win  Result = "Win"
	Loss Result = "Loss"
	Draw Result = "Draw"
)

func getOutcome(boardState *BoardState, sequence []int) (turns int, result Result, firstMove int) {
	simulatedBoardState := copyBoardState(boardState)

	for index, tile := range sequence {

		if index%2 == 0 {
			simulatedBoardState = simulateMove(&simulatedBoardState, AI, tile)

			if aiWins(&simulatedBoardState) {
				return index + 1, Win, sequence[0]
			}

		} else {
			simulatedBoardState = simulateMove(&simulatedBoardState, Opponent, tile)

			if opponentWins(&simulatedBoardState) {
				return index + 1, Loss, sequence[0]
			}
		}
	}

	return len(sequence), Draw, sequence[0]
}

func findKeyWithGreatestValue(m map[int]int) int {
	maxValue := math.MinInt32
	maxKey := 0

	for key, value := range m {
		if value > maxValue {
			maxValue = value
			maxKey = key
		}
	}

	return maxKey
}

func calculateBestMove(boardState *BoardState) int {
	scores := getScores(boardState)
	return findKeyWithGreatestValue(scores)
}
