package chesslogic

import (
	"fmt"
	"strings"
)

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
    //capPiece := b.Squares[toRank][toFile] <- this is for castle but never used
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
    if piece.Type == King && abs(toFile-fromFile) == 2  { //TODO: in the future castle can be executed if player clicked on rook 
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
	//but i'm not pretty sure it will be useful
	//i could make it though but if only i found someone who could make frontend..
	//surely this is will be good project but that's okay. I tried so hard every day to keep myself focused
	// in these cold days when everything falls apart. When pleasure and happiness are so far away from me
	// when i feel so alone
	//Why am i keep writing it in the code? i think it's on mind every time i think about something important
	// like my goals or my life. Or even this project that nobody will ever see except me
	// and yet. I'll keep trying. and keep feel sadness
	// if anyone will ever read this - i wish you (no matter who you are or where you're from) lots of love and joy in your life
	// and i will go play chess - because it's the only thing that brings me joy. Nothing else.

    // it's horrible that i made such a note here. But whatever. 
    //These days brings me nothing. I keep trying to study to work...Every day i search for permanence - a place that i can go
    // to make me whole
    //And yet there's nothing but a silence. I think i sink in silence. It's 28 August. After a few days my study will start.
    //If only it brought me some joy.
    //28 August - the most lonely day of my life. The most lonely week. Month.Year. My life is like the coldest place in the world
    // Where no man lives.
    //It's strange that i keep writing these things in the code. but i hope if someone would ever read this..May you be blessed
    // It's August 28. Tommorow will be August 29. I'm so tired of climbing up the walls that i build for myself.
    // But whatever. I'll keep trying. 
    //It can't be real actually. I'm not a man who can sit there and write code. I must've lost myself. I'm not a backend developer
    // Not a developer at all. I think i'm not a human but a shadow of myself. I was..joyful? I guess yes. But now i don't feel pleasure
    // but pain in all parts of the body. Only pain occupies my life.
    // Don't think i'll kill myself. It's disgusting. And i have something that brings me hope. Note from my groupmate (i guess she doesn't even remember me)
    // "I wish you lots of love and joy in the upcoming year! May all your dreams come true". And a heart emoticon below the note.
    // None of my dreams have come true. But still i'm trying.
}

func parsePromotion(s string) (PieceType, error) {
    switch strings.ToLower(s) {
    case "q", "queen":  return Queen, nil
    case "r", "rook":   return Rook, nil
    case "b", "bishop": return Bishop, nil
    case "n", "knight": return Knight, nil
    case "":            return Queen, nil 
    default:
        return Queen, fmt.Errorf("unknown promotion: %q", s)
    }
}