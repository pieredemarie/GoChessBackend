package chesslogic

import (
	"fmt"
)

//TODO
//in game.go Check and hasAttackedCell functions!!

func (b *Board) IsLegalMove(piece Piece,m Move) bool {
	switch piece.Type {
	case Pawn:
		return b.isLegalPawnMove(piece,m) 
	case Knight:
		return b.isLegalKnightMove(piece,m)
	case Bishop: 
		return b.isLegalBishopMove(piece,m)
	case Rook:
		return b.isLegalRookMove(piece,m)
	case Queen:
		return b.isLegalQueenMove(piece,m) 
	case King:
		return b.isLegalKingMove(piece,m)
	default:
		return false
	}
}

func (b *Board) isLegalKnightMove(piece Piece,m Move) bool {
	fromFile, fromRank, err := squareToCoords(m.From)
    if err != nil {
        return false
    }
    toFile, toRank, err := squareToCoords(m.To)
    if err != nil {
        return false
    }

	pieceToEnd := b.GetPiece(m.To)
	if pieceToEnd != nil && pieceToEnd.Color == piece.Color { // Horses can't jump on their friends, otherwise you have a weird horse
		return false
	}

	rowsDiff := abs(fromRank - toRank)
	colsDiff := abs(fromFile - toFile)

	return (rowsDiff == 2 && colsDiff == 1) || (rowsDiff == 1 && colsDiff == 2)
}

func (b *Board) isLegalKingMove(piece Piece,m Move) bool {
	fromFile, fromRank, err := squareToCoords(m.From)
    if err != nil {
        return false
    }
    toFile, toRank, err := squareToCoords(m.To)
    if err != nil {
        return false
    }

	king := b.Squares[fromRank][fromFile]

	dRank := abs(toRank - fromRank)
	dFile := abs(toFile - fromFile)

	if dFile <= 1 && dRank <= 1 { 
		if b.Squares[toRank][toFile] == nil || b.Squares[toRank][toFile].Color != piece.Color {
			target := b.Squares[toRank][toFile]

			//temporarily moving the figure on the square
			b.Squares[fromRank][fromFile] = nil 
			b.Squares[toRank][toFile] = king

			//checking if king is under attack
			kingPos := string(rune('a'+toFile)) + string(rune('1'+toRank))
			isKingAttacked := b.IsCellAttacked(kingPos,king.Color)

			//undo the move
			b.Squares[fromRank][fromFile] = king
			b.Squares[toRank][toFile] = target

			return !isKingAttacked
		}	
	}
	return false
}

func (b *Board) isLegalBishopMove(piece Piece,m Move) bool {
	fromFile, fromRank, err := squareToCoords(m.From)
	if err != nil {
		return false
	}

	toFile, toRank, err := squareToCoords(m.To)
	if err != nil {
		return false
	}

	if abs(toFile-fromFile) != abs(toRank-fromRank) {
        return false
    }

	stepR := sign(toRank - fromRank)
	stepL := sign(toFile - fromFile)
	
	r := fromRank + stepR
	l := fromFile + stepL
	for r != toRank || l != toFile {
		if !isValidCoordinate(l,r) {
			return false
		}

		if b.Squares[r][l] != nil {
			if b.Squares[r][l].Color == piece.Color {
				return false
			}

			return r == toRank && l == toFile
		}
		r += stepR
		l += stepL
	}

	if b.Squares[toRank][toFile] == nil || b.Squares[toRank][toFile].Color != piece.Color {
		return true
	}

	return false
}

func (b *Board) isLegalRookMove(piece Piece,m Move) bool {
	fromFile, fromRank, err := squareToCoords(m.From)
	if err != nil {
		return false
	}

	toFile, toRank, err := squareToCoords(m.To)
	if err != nil {
		return false
	}
	if (fromRank != toRank && fromFile != toFile) {
		return false
	}

	stepR := sign(toRank - fromRank)
	stepL := sign(toFile - fromFile)
	
	currRank := fromRank + stepR
	currFile := fromFile + stepL

	for currRank != toRank || currFile != toFile {
		if !isValidCoordinate(currFile,currRank) {
			return false
		}

		if b.Squares[currRank][currFile] != nil {
			if b.Squares[currRank][currFile].Color != piece.Color {
				return currRank == toRank && currFile == toFile
			}
			return false
		}

		currRank += stepR
		currFile += stepL
	}
	
	if b.Squares[toRank][toFile] == nil || b.Squares[currRank][currFile].Color != piece.Color {
		return true
	}
	return false
}

func (b *Board) isLegalQueenMove(piece Piece,m Move) bool {
	fromFile, fromRank, err := squareToCoords(m.From)
	if err != nil {
		return false
	}

	toFile, toRank, err := squareToCoords(m.To)
	if err != nil {
		return false
	}

	if fromRank != toRank && fromFile != toFile && abs(fromRank-toRank) != abs(fromFile-toFile) {
    	return false // non horizontal vertical or diagonal move
	} 

	stepR := sign(toRank - fromRank)
	stepL := sign(toFile - fromFile)
	
	currRank := fromRank + stepR
	currFile := fromFile + stepL

	for currRank != toRank || currFile != toFile {
		if !isValidCoordinate(currFile, currRank) {
			return false
		}
		
		if b.Squares[currRank][currFile] != nil {
			if b.Squares[currRank][currFile].Color == piece.Color {
				return false
			}
			return currRank == toRank && currFile == toFile 
		}

		currRank += stepR
		currFile += stepL
	}

	if b.Squares[toRank][toFile] == nil || b.Squares[currRank][currFile].Color != piece.Color {
		return true
	}
	return false
}

func (b *Board) isLegalPawnMove(piece Piece,m Move) bool {
	fromFile, fromRank, err := squareToCoords(m.From)
	if err != nil {
		return false
	}

	toFile, toRank, err := squareToCoords(m.To)
	if err != nil {
		return false
	}

	dRank := toRank - fromRank
	dFile := toFile - fromFile

	if dRank == 0 {
		return false
	}
	if dFile == 0 {
		if piece.Color == White {
			if dRank == 1 && b.Squares[toRank][toFile] == nil {
				return true
			}

			if fromRank == 1 && dRank == 2 && b.Squares[toRank-1][toFile] == nil && b.Squares[toRank][toFile] == nil { // move from only start position (2 sqaures away)
				return true
				//if this happened a pawn can be captured en passat
			}
		} else if piece.Color == Black {
			if dRank == -1 && b.Squares[toRank][toFile] == nil {
				return true
			}
			if fromRank == 6 && dRank == -2 && b.Squares[toRank+1][toFile] == nil && b.Squares[toRank][toFile] == nil { // move from only start position (2 sqaures away)
				return true
				//if this happened a pawn can be captured en passat
			}
		}
	} else if abs(dFile) == 1 { // capturing
		if piece.Color == White && dRank == 1 { 
			//capture check
			if b.Squares[toRank][toFile] != nil && b.Squares[toRank][toFile].Color != piece.Color {
				return true
			}
			//en passat check
			if enemy := b.Squares[fromRank][toFile];enemy != nil && enemy.Type == Pawn && enemy.Color != piece.Color {
				return true
			}  
		} else if piece.Color == Black && dRank == -1 {
			//capture check
			if b.Squares[toRank][toFile] != nil && b.Squares[toRank][toFile].Color != piece.Color {
				return true
			}
			//en passat check
			if enemy := b.Squares[fromRank][toFile]; enemy.Type == Pawn && enemy.Color != piece.Color {
				return true
			}  
		}
	}
	return false
}

func abs(x int) int {
	if x < 0 {return -x}
	return x
}

func sign(x int) int {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	}
	return 0
}

func isValidCoordinate(file, rank int) bool {
    return file >= 0 && file < 8 && rank >= 0 && rank < 8
}

func OppositeColor(c Color) Color {
	if c == White { return Black }
	return White
}

func squareToCoords(sq string) (int, int, error) {
	if len(sq) != 2 {
		return 0, 0, fmt.Errorf("invalid square: %s", sq)
	}
	file := int(sq[0] - 'a') 
    rank := int(sq[1] - '1') 
    if file < 0 || file > 7 || rank < 0 || rank > 7 {
        return 0, 0, fmt.Errorf("invalid square: %s", sq)
    }
    return file, rank, nil // rank, file, nil
}