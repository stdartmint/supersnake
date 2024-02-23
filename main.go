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
		Color:      "#ff5722",   // TODO: Choose color
		Head:       "pig",       // TODO: Choose head
		Tail:       "pig",       // TODO: Choose tail
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
	log.Printf("MOVE %d: No safe moves detected! Moving down\n", state.Turn)

	if myNeck.X < myHead.X { // Neck is left of head, don't move left
		isMoveSafe["left"] = false

	} else if myNeck.X > myHead.X { // Neck is right of head, don't move right
		isMoveSafe["right"] = false

	} else if myNeck.Y < myHead.Y { // Neck is below head, don't move down
		isMoveSafe["down"] = false

	} else if myNeck.Y > myHead.Y { // Neck is above head, don't move up
		isMoveSafe["up"] = false
	}

	if state.You.Head.X-1 < 0 {
		isMoveSafe["left"] = false
	}
	if state.You.Head.X+1 >= boardWidth {
		isMoveSafe["right"] = false
	}
	if state.You.Head.Y-1 < 0 {
		isMoveSafe["down"] = false
	}
	if state.You.Head.Y+1 >= boardHeight {
		isMoveSafe["up"] = false
	}

	isMoveSafe = checkBody(state, state.You.Body, isMoveSafe)
	if !isInHazard(state.You.Head, state.Board.Hazards) {
		isMoveSafe = checkBody(state, state.Board.Hazards, isMoveSafe)
	}

	// Are there any safe moves left?
	safeMoves := []string{}

	opponents := state.Board.Snakes
	for _, opponent := range opponents {
		bo := []Coord{opponent.Head}
		bo = append(bo, opponent.Body...)
		isMoveSafe = checkBody(state, bo, isMoveSafe)
	}

	food := state.Board.Food
	nextMove := getDirection(state.You.Head, findPath(state.You.Head, food))

	for move, isSafe := range isMoveSafe {
		if isSafe {
			safeMoves = append(safeMoves, move)
		}
	}

	if len(safeMoves) == 0 {
		log.Printf("MOVE %d: No safe moves detected! Moving down\n", state.Turn)
		return BattlesnakeMoveResponse{Move: "down"}
	}

	access := false
	for !access {
		if !slices.Contains(safeMoves, nextMove) {
			nextMove = safeMoves[rand.Intn(len(safeMoves))]
		}
		access = true
	}

	log.Printf("MOVE %d: %s\n", state.Turn, nextMove)
	return BattlesnakeMoveResponse{Move: nextMove}
}

func getDirection(me, target Coord) string {
	log.Printf("preGetDirection%v\n", target)

	dx := target.X - me.X
	dy := target.Y - me.Y

	switch {
	case dx == 0 && dy < 0:
		return "up"
	case dx == 0 && dy > 0:
		return "down"
	case dx < 0 && dy == 0:
		return "left"
	case dx > 0 && dy == 0:
		return "right"
	default:
		return ""
	}
}

func isInHazard(head Coord, hazards []Coord) bool {
	for _, hazard := range hazards {
		if head.X == hazard.X && head.Y == hazard.Y {
			return true
		}
	}
	return false
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

//func checkBody(state GameState, body []Coord, isMoveSafe map[string]bool) map[string]bool {
//	headX := state.You.Head.X
//	headY := state.You.Head.Y
//	myLength := len(state.You.Body)
//
//	for _, part := range body {
//		if part.X == headX-1 && part.Y == headY {
//			isMoveSafe["left"] = false
//		}
//		if part.X == headX+1 && part.Y == headY {
//			isMoveSafe["right"] = false
//		}
//		if part.Y == headY-1 && part.X == headX {
//			isMoveSafe["down"] = false
//		}
//		if part.Y == headY+1 && part.X == headX {
//			isMoveSafe["up"] = false
//		}
//	}
//
//	for _, snake := range state.Board.Snakes {
//		if snake.ID == state.You.ID {
//			continue
//		}
//		opponentLength := len(snake.Body)
//
//		if myLength <= opponentLength {
//			possibleMoves := []Coord{
//				{X: snake.Head.X + 1, Y: snake.Head.Y}, // right
//				{X: snake.Head.X - 1, Y: snake.Head.Y}, // left
//				{X: snake.Head.X, Y: snake.Head.Y + 1}, // up
//				{X: snake.Head.X, Y: snake.Head.Y - 1}, // down
//			}
//
//			for _, move := range possibleMoves {
//				if move.X == headX && move.Y == headY {
//					if snake.Head.X < headX {
//						isMoveSafe["right"] = false
//					} else if snake.Head.X > headX {
//						isMoveSafe["left"] = false
//					} else if snake.Head.Y < headY {
//						isMoveSafe["up"] = false
//					} else if snake.Head.Y > headY {
//						isMoveSafe["down"] = false
//					}
//				}
//			}
//		}
//	}
//
//	return isMoveSafe
//}

func checkBody(state GameState, body []Coord, isMoveSafe map[string]bool) map[string]bool {
	headX := state.You.Head.X
	headY := state.You.Head.Y

	for _, part := range body {
		if part.X == headX-1 && part.Y == headY {
			isMoveSafe["left"] = false
		}
		if part.X == headX+1 && part.Y == headY {
			isMoveSafe["right"] = false
		}
		if part.Y == headY-1 && part.X == headX {
			isMoveSafe["down"] = false
		}
		if part.Y == headY+1 && part.X == headX {
			isMoveSafe["up"] = false
		}
	}

	return isMoveSafe
}

func main() {
	RunServer()
}
