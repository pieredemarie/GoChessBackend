package chesslogic

import "fmt"

//TODO
// 

type Move struct {
	From      string
	To        string
	Promotion *PieceType
}

func (b *Board) ApplyMove(m Move) error {
	fromFile, fromRank, err := squareToCoords(m.From)
	if err != nil {
		return fmt.Errorf("invalid from square")
	}
	toFile, toRank, err := squareToCoords(m.To)
	if err != nil {
		return fmt.Errorf("invalid to square")
	}

	piece := b.Squares[fromRank][fromFile]
	if piece == nil {
		return fmt.Errorf("no piece at the selected square")
	}

	if !b.IsLegalMove(*piece,m) {
		return fmt.Errorf("illegal move")
	}
	b.Squares[toRank][toFile] = nil
	
	return nil
}

func (b *Board) UndoMove(m Move) {
	//TODO: restore previos move function
}

func (b *Board) GenerateAllPossibleMoves(piece Piece) {
	//Here will be function that generate all possible moves for figure
	//but i'm not pretty sure that it will be useful
	//i could make it though but if only i found someone who could make frontend..
	//surely this is will be good project but that's okay. I tried so hard every day to keep myself focused
	// in this cold days when everything falls apart. When pleasure and happines is so far away from me
	// when i feel so alone
	//Why am i keep writing it in the code? i think it's on mind every time i think about something important
	// like my goals or my life. Or even this project that nobody will ever see except me
	// and yet. I'll keep trying. and keep feel sadness
	// if anyone will ever read this - i wish you (no matter who you are or where you're from) lots of love and joy in your life
	// and i will go play chess - because it's the only thing that brings me joy. Nothing else.
}