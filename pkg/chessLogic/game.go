package chesslogic

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type GameStatus int

const (
	Waiting GameStatus = iota
	Active
	Finished
)

type Player struct {
	ID       int
	Username string
	Color    Color
	TimeLeft time.Duration 
}

type Game struct {
	ID string 
	Players [2]*Player
	Board *Board
	Moves []Move
	Status GameStatus
	TurnColor Color
	Winner Color
	Mut sync.Mutex 
}

func NewGame(playerWhite, playerBlack *Player) *Game {
	return &Game {
		ID: generateGameId(),
		Players: [2]*Player{playerWhite,playerBlack},
		Board: NewBoard(),
		Moves: []Move{},
		Status: Active,
		TurnColor: White,
	}
} 

func (g *Game) MakeMove(move Move) error {
	g.Mut.Lock()
	defer g.Mut.Unlock()

	if g.Status != Active {
		return fmt.Errorf("game is over")
	}

	figFile, figRank, err := squareToCoords(move.From)
	if err != nil {
		return fmt.Errorf("invalid move")
	}
	capFile, capRank, err := squareToCoords(move.To)
	if err != nil {
		return fmt.Errorf("invalid move")
	}
	fig := g.Board.Squares[figRank][figFile]
	if fig == nil {
		return fmt.Errorf("no piece on origin square")
	}

	if g.TurnColor != fig.Color {
		return fmt.Errorf("not your turn")
	}

	if !g.Board.IsLegalMove(*fig,move) {
		return fmt.Errorf("illegal move")
	}

	captured := g.Board.Squares[capRank][capFile]
	origType := fig.Type

	if err := g.Board.ApplyMove(move); err != nil {
		return err
	}

	if g.Board.IsCheck(fig.Color) {
		g.Board.Squares[figRank][figFile] = fig
		g.Board.Squares[capRank][capFile] = captured
		fig.Type = origType
		return fmt.Errorf("move leaves king in check")
	}

	g.Moves = append(g.Moves, move)

	opp := OppositeColor(g.TurnColor)
	if g.Board.IsCheckMate(opp) {
		g.Status = Finished
		g.Winner = g.TurnColor
		return nil
	}
	
	if g.Board.IsStaleMate(opp) {
		g.Status = Finished
		return nil
	}

	g.TurnColor = opp
	return nil
}

func (g *Game) IsGameOver() bool {
	return g.Status == Finished
}

func generateGameId() string {
	id := uuid.New()
	return id.String()
}

