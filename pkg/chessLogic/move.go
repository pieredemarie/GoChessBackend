package chesslogic

import (
	"fmt"
	"strings"
)

//TODO
//

type Move struct {
	From      string
	To        string
	Promotion string
}

func (b *Board) ApplyMove(m Move) error {
    fromFile, fromRank, _ := squareToCoords(m.From)
    toFile, toRank, _ := squareToCoords(m.To)

    piece := b.Squares[fromRank][fromFile]
    b.Squares[fromRank][fromFile] = nil

    // en passant
    if piece.Type == Pawn && fromFile != toFile && b.Squares[toRank][toFile] == nil {
        captureRank := fromRank
        if piece.Color == White {
            captureRank = toRank - 1
        } else {
            captureRank = toRank + 1
        }
        b.Squares[captureRank][toFile] = nil
    }

    // castle
    if piece.Type == King && abs(toFile-fromFile) == 2 {
        if toFile > fromFile { // short castling
            rook := b.Squares[fromRank][7]
            b.Squares[fromRank][7] = nil
            b.Squares[fromRank][fromFile+1] = rook
        } else { // long castling
            rook := b.Squares[fromRank][0]
            b.Squares[fromRank][0] = nil
            b.Squares[fromRank][fromFile-1] = rook
        }
    }

    // pawn promotion
    if piece.Type == Pawn && (toRank == 7 || toRank == 0) {
        pt, err := parsePromotion(m.Promotion)
    	if err != nil {
        	pt = Queen
   	 	}
    	piece.Type = pt
    }

    b.Squares[toRank][toFile] = piece

    if piece.Type == Pawn && abs(toRank-fromRank) == 2 {
        midRank := (toRank + fromRank) / 2
        b.EnPassantSquare = string(rune('a'+fromFile)) + string(rune('1'+midRank))
    } else {
        b.EnPassantSquare = ""
    }
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

func parsePromotion(s string) (PieceType, error) {
    switch strings.ToLower(s) {
    case "q", "queen":  return Queen, nil
    case "r", "rook":   return Rook, nil
    case "b", "bishop": return Bishop, nil
    case "n", "knight": return Knight, nil
    case "":            return Queen, nil // по умолчанию ферзь
    default:
        return Queen, fmt.Errorf("unknown promotion: %q", s)
    }
}