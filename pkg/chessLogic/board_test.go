package chesslogic

import "testing"

func TestIsAttackingCell(t *testing.T) {
	b := NewBoard() 
	
	//test 1: knight should attack c3
	if !b.IsCellAttacked("c3",Black) {
		t.Errorf("Expected c3 is attacked by knight on b1")
	}

	//test 2: nothing attacks e5
	if b.IsCellAttacked("e5", Black) {
		t.Errorf("Expected nothing attacks e5")
	}
}

func TestChecks(t *testing.T) {
	b := NewBoard()

	//test 1: queen on h4 should put king in check (with f2 pawn removed from board) 
	b.Squares[3][7] = NewPiece(Queen,Black)
	b.Squares[1][5] = nil
	if !b.IsCheck(White) {
		t.Errorf("Expected check by Queen on h4")
	}
}

func TestCheckMate(t *testing.T) {
	b := NewBoard()

	//clearing the board but it's much better to crete NewEmptyBoard without arranging figures on it
	for rank := 0;rank < 8;rank++ {
		for file := 0;file < 8;file++ {
			b.Squares[rank][file] = nil
		}
	}

	// test 1: King on h8 Queen on g7 with the King on f6
	b.Squares[7][7] = NewPiece(King, Black) // h8
    b.Squares[6][6] = NewPiece(Queen, White) // g7
    b.Squares[5][5] = NewPiece(King, White) // f6

    // checking that black king is checkmated
    if !b.IsCheckMate(Black) {
        t.Errorf("Expected checkmate for Black")
    }
}

func TestStaleMate(t *testing.T) {
	b := NewBoard()

	//clearing the board but it's much better to crete NewEmptyBoard without arranging figures on it
	for rank := 0;rank < 8;rank++ {
		for file := 0;file < 8;file++ {
			b.Squares[rank][file] = nil
		}
	}
	// black king has no legal moves so this is a stalemate for sure
    b.Squares[7][7] = NewPiece(King, Black) // h8
    b.Squares[5][5] = NewPiece(King, White) // f6
    b.Squares[5][6] = NewPiece(Queen, White) // g6

    // checking if there's a stalemate
    if !b.IsStaleMate(Black) {
        t.Errorf("Expected stalemate for Black")
    }
}

func TestCastlingShort(t *testing.T) {
    // creating new Board
    b := NewBoard()
    // putting white king and rook 
    b.Squares[0][4] = &Piece{Type: King, Color: White}
    b.Squares[0][7] = &Piece{Type: Rook, Color: White}
    b.WhiteKingMoved = false
    b.WhiteRookHMoved = false

    // clearing the path
    b.Squares[0][5] = nil
    b.Squares[0][6] = nil

    // castling kingside
    move := Move{From: "e1", To: "g1"}
    err := b.ApplyMove(move)
    if err != nil {
        t.Fatalf("castling failed: %v", err)
    }

    // making sure that castle is done
    if b.Squares[0][6] == nil || b.Squares[0][6].Type != King {
        t.Errorf("expected king on g1")
    }
    if b.Squares[0][5] == nil || b.Squares[0][5].Type != Rook {
        t.Errorf("expected rook on f1")
    }
    if b.Squares[0][4] != nil || b.Squares[0][7] != nil {
        t.Errorf("old squares should be empty")
    }
}

func TestCastlingLo(t *testing.T) {
b := NewBoard()
    // putting white king and rook 
    b.Squares[0][4] = &Piece{Type: King, Color: White}
    b.Squares[0][0] = &Piece{Type: Rook, Color: White}
    b.WhiteKingMoved = false
    b.WhiteRookHMoved = false

    // clearing the path
    b.Squares[0][3] = nil
    b.Squares[0][2] = nil
	b.Squares[0][1] = nil

    // castling queenside
    move := Move{From: "e1", To: "c1"}
    err := b.ApplyMove(move)
    if err != nil {
        t.Fatalf("castling failed: %v", err)
    }

    // making sure that castle is done
    if b.Squares[0][2] == nil || b.Squares[0][2].Type != King {
        t.Errorf("expected king on c1")
    }
    if b.Squares[0][3] == nil || b.Squares[0][3].Type != Rook {
        t.Errorf("expected rook on d1")
    }
    if b.Squares[0][4] != nil || b.Squares[0][0] != nil {
        t.Errorf("old squares should be empty")
    }
}