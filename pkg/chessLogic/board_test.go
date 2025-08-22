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