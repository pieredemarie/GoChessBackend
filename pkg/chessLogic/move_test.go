package chesslogic

import "testing"

func TestKnightMoves(t *testing.T) {
    b := NewBoard() 

    knight := Piece{Type: Knight, Color: White}
    b.Squares[0][1] = &knight // putting knight on b1

    moves := []struct {
        from, to string
        want     bool
    }{
        {"b1", "c3", true},  // correct
        {"b1", "a3", true},  // correct
        {"b1", "b3", false}, // invalid
    }

    for _, m := range moves {
        if got := b.isLegalKnightMove(knight, Move{From: m.from, To: m.to}); got != m.want {
            t.Errorf("knight %s->%s = %v, want %v", m.from, m.to, got, m.want)
        }
    }
}
func TestBishopMove(t *testing.T) {
	b := NewBoard()

	// pawn from e2 removed so bishop from f1 can move
	b.Squares[1][4] = nil 

	bishop := *b.Squares[0][5] // bishop on f1
	move := Move{From: "f1", To: "c4"} // diagonal move from f1-c4

	if !b.IsLegalMove(bishop, move) {
		t.Errorf("Expected bishop move f1->c4 to be legal")
	}

	// Invalid move test from f1 to f3 bishop can't move
	move2 := Move{From: "f1", To: "f3"}
	if b.IsLegalMove(bishop, move2) {
		t.Errorf("Expected bishop move f1->f3 to be illegal")
	}
}

func TestQueenMove(t *testing.T) {
	b := NewBoard()
	
	//also removing pawn from e2 so queen can actually move
	b.Squares[1][4] = nil 
	queen := *b.Squares[0][3]
	moves := []struct {
        from, to string
        want     bool
    }{
        {"d1", "f3", true},  // correct
        {"d1", "h5", true},  // correct
        {"d1", "b3", false}, // invalid
    }
	for _, m := range moves {
        if got := b.isLegalQueenMove(queen, Move{From: m.from, To: m.to}); got != m.want {
            t.Errorf("queen %s->%s = %v, want %v", m.from, m.to, got, m.want)
        }
    }
}

func TestRookMove(t *testing.T) {
	b := NewBoard()

	//removing pawn from a2 so the file can be opened for rook
	b.Squares[1][0] = nil 
	rook := *b.Squares[0][0]
	moves := []struct {
        from, to string
        want     bool
    }{
        {"a1", "a5", true},  // correct
        {"a1", "a7", true},  // correct
        {"a1", "a8", false}, // invalid
		{"a1", "b1", false},
    }
	for _, m := range moves {
        if got := b.isLegalRookMove(rook, Move{From: m.from, To: m.to}); got != m.want {
            t.Errorf("queen %s->%s = %v, want %v", m.from, m.to, got, m.want)
        }
    }
}

func TestPawnMove(t *testing.T) {
	b := NewBoard()

	pawn := *b.Squares[1][0] //taking as an example pawn from a2
	b.Squares[2][1] = NewPiece(Pawn,Black) // let's put enemy pawn here on b3
	moves := []struct { // as a test we'll see can pawn moves 2 squares away from start position
        from, to string
        want     bool
    }{
        {"a2", "a3", true},  // correct
        {"a2", "a4", true},  // correct
        {"a2", "b3", true}, // invalid
		{"a2", "b4", false}, // invalid
    }
	for _, m := range moves {
        if got := b.isLegalPawnMove(pawn, Move{From: m.from, To: m.to}); got != m.want {
            t.Errorf("queen %s->%s = %v, want %v", m.from, m.to, got, m.want)
        }
    }
}