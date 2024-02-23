package main

// Welcome to
// __________         __    __  .__                               __
// \______   \_____ _/  |__/  |_|  |   ____   ______ ____ _____  |  | __ ____
//  |    |  _/\__  \\   __\   __\  | _/ __ \ /  ___//    \\__  \ |  |/ // __ \
//  |    |   \ / __ \|  |  |  | |  |_\  ___/ \___ \|   |  \/ __ \|    <\  ___/
//  |________/(______/__|  |__| |____/\_____>______>___|__(______/__|__\\_____>
//
// This file can be a nice home for your Battlesnake logic and helper functions.
//
// To get you started we've included code to prevent your Battlesnake from moving backwards.
// For more info see docs.battlesnake.com

import (
	"log"
	"math"
	"math/rand"
	"slices"
)

// info is called when you create your Battlesnake on play.battlesnake.com
// and controls your Battlesnake's appearance
// TIP: If you open your Battlesnake URL in a browser you should see this data
func info() BattlesnakeInfoResponse {
	log.Println("INFO")

	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "BumbleBee", // TODO: Your Battlesnake username
		Color:      "#9c27b0",   // TODO: Choose color
		Head:       "bee",       // TODO: Choose head
		Tail:       "ladybug",   // TODO: Choose tail
	}
}

// start is called when your Battlesnake begins a game
func start(state GameState) {
	log.Println("GAME START")
}

// end is called when your Battlesnake finishes a game
func end(state GameState) {
	log.Printf("GAME OVER\n\n")
}

// move is called on every turn and returns your next move
// Valid moves are "up", "down", "left", or "right"
// See https://docs.battlesnake.com/api/example-move for available data
func move(state GameState) BattlesnakeMoveResponse {

	isMoveSafe := map[string]bool{
		"up":    true,
		"down":  true,
		"left":  true,
		"right": true,
	}

	// We've included code to prevent your Battlesnake from moving backwards
	myHead := state.You.Body[0] // Coordinates of your head
	myNeck := state.You.Body[1] // Coordinates of your "neck"

	// TODO: Step 1 - Prevent your Battlesnake from moving out of bounds
	boardWidth := state.Board.Width
	boardHeight := state.Board.Height

	if myNeck.X < myHead.X { // Neck is left of head, don't move left
		isMoveSafe["left"] = false

	} else if myNeck.X > myHead.X { // Neck is right of head, don't move right
		isMoveSafe["right"] = false

	} else if myNeck.Y < myHead.Y { // Neck is below head, don't move down
		isMoveSafe["down"] = false

	} else if myNeck.Y > myHead.Y { // Neck is above head, don't move up
		isMoveSafe["up"] = false
	}

	if state.You.Head.X-1 == -1 {
		isMoveSafe["left"] = false
	} else if state.You.Head.X+1 == boardWidth {
		isMoveSafe["right"] = false
	} else if state.You.Head.Y-1 == -1 { // Neck is below head, don't move down
		isMoveSafe["down"] = false
	} else if state.You.Head.Y+1 == boardHeight { // Neck is above head, don't move up
		isMoveSafe["up"] = false
	}

	isMoveSafe = checkBody(state, state.You.Body, isMoveSafe)

	// Are there any safe moves left?
	safeMoves := []string{}
	for move, isSafe := range isMoveSafe {
		if isSafe {
			safeMoves = append(safeMoves, move)
		}
	}

	// TODO: Step 3 - Prevent your Battlesnake from colliding with other Battlesnakes
	opponents := state.Board.Snakes

	for _, opponent := range opponents {
		bo := []Coord{opponent.Head}
		bo = append(bo, opponent.Body...)
		isMoveSafe = checkBody(state, bo, isMoveSafe)
	}

	if len(safeMoves) == 0 {
		log.Printf("MOVE %d: No safe moves detected! Moving down\n", state.Turn)
		return BattlesnakeMoveResponse{Move: "down"}
	}

	// TODO: Step 4 - Move towards food instead of random, to regain health and survive longer
	food := state.Board.Food

	nextMove := "down"
	if len(food) == 0 {
		nextMove = safeMoves[rand.Intn(len(safeMoves))]
		log.Printf("len(food) == 0 nextMove %s\n", nextMove)
	} else {
		nextMove = getDirection(findPath(state.You.Head, food))
		log.Printf("getDirection nextMove %s\n", nextMove)
	}

	log.Printf("MOVE %d: %s\n", state.Turn, nextMove)
	return BattlesnakeMoveResponse{Move: nextMove}
}

func getDirection(d Coord) string {
	log.Printf("preGetDirection%v\n", d)

	switch {
	case d.X == 0 && d.Y == 1:
		return "up"
	case d.X == 0 && d.Y == -1:
		return "down"
	case d.X == -1 && d.Y == 0:
		return "left"
	case d.X == 1 && d.Y == 0:
		return "right"
	default:
		return "down" // Вернуть пустое направление в случае неверных координат
	}
}

func findPath(current Coord, targets []Coord) Coord {
	minDist := math.Inf(1)
	var closest Coord

	for _, target := range targets {
		dist := math.Sqrt(math.Pow(float64(target.X-current.X), 2) + math.Pow(float64(target.Y-current.Y), 2))
		if dist < minDist {
			minDist = dist
			closest = target
		}
	}

	log.Printf("findPath %v\n", closest)

	return closest
}

func checkBody(state GameState, body []Coord, isMoveSafe map[string]bool) map[string]bool {
	// TODO: Step 2 - Prevent your Battlesnake from colliding with itself
	mybodyX := []int{}
	mybodyY := []int{}

	for _, bodyPart := range body {
		mybodyX = append(mybodyX, bodyPart.X)
		mybodyY = append(mybodyY, bodyPart.Y)
	}

	for move := range isMoveSafe {
		if move == "left" && slices.Contains(mybodyX, state.You.Head.X-1) {
			isMoveSafe[move] = false
		}
		if move == "right" && slices.Contains(mybodyX, state.You.Head.X+1) {
			isMoveSafe[move] = false
		}
		if move == "up" && slices.Contains(mybodyY, state.You.Head.Y+1) {
			isMoveSafe[move] = false
		}
		if move == "down" && slices.Contains(mybodyY, state.You.Head.Y-1) {
			isMoveSafe[move] = false
		}
	}

	return isMoveSafe
}

func main() {
	RunServer()
}
