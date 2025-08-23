package chesslogic

import "testing"

func TestMakeMoveAndApplyMove(t *testing.T) {
    // 1. moving knight
    g := NewGame(&Player{Username: "White", Color: White}, &Player{Username: "Black", Color: Black})
    
    move1 := Move{From: "b1", To: "c3"}
    if err := g.MakeMove(move1); err != nil {
        t.Errorf("expected knight move to be legal, got error: %v", err)
    }

    // 2. moving black pieces  while it's not their turn
    move2 := Move{From: "c8", To: "g4"}
    if err := g.MakeMove(move2); err == nil {
        t.Errorf("expected error for moving wrong color, got nil")
    }

    // 3. checks,mate and stalemate tests
    g = &Game{
        Board:     NewBoard(),
        TurnColor: White,
        Status:    Active,
    }
    
    // clearing position
    for i := 0; i < 8; i++ {
        for j := 0; j < 8; j++ {
            g.Board.Squares[i][j] = nil
        }
    }
    g.Board.Squares[7][7] = NewPiece(King, Black)    // h8
    g.Board.Squares[6][5] = NewPiece(Queen, White)   // f7
    g.Board.Squares[5][5] = NewPiece(King, White)    // f6

    move3 := Move{From: "f7", To: "g7"}
    if err := g.MakeMove(move3); err != nil {
        t.Errorf("expected checkmate move to be accepted, got error: %v", err)
    }

    if g.Status != Finished {
        t.Errorf("expected game to be finished after checkmate, got %v", g.Status)
    }
    if g.Winner != White {
        t.Errorf("expected winner white, got %v", g.Winner)
    }
}