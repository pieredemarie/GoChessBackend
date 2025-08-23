package chesslogic

type Board struct {
	Squares [8][8]*Piece
	Turn    Color
	EnPassantSquare string
}

type CapablancaBoard struct { //In the future will be done this type of chess
	Squares [10][8]*Piece
	Turn    Color
	EnPassantSquare string
}

func NewBoard() *Board {
	b := &Board{}
	b.ArrangeFigures()
	return b
}

func (b *Board) ArrangeFigures() {
	//pawns
	for cols := 0; cols < 8; cols++ {
		b.Squares[1][cols] = NewPiece(Pawn, White)
		b.Squares[6][cols] = NewPiece(Pawn, Black)
	}

	//Rooks
	b.Squares[0][0] = NewPiece(Rook, White)
	b.Squares[0][7] = NewPiece(Rook, White)
	b.Squares[7][0] = NewPiece(Rook, Black)
	b.Squares[7][7] = NewPiece(Rook, Black)

	//Knights
	b.Squares[0][1] = NewPiece(Knight, White)
	b.Squares[0][6] = NewPiece(Knight, White)
	b.Squares[7][1] = NewPiece(Knight, Black)
	b.Squares[7][6] = NewPiece(Knight, Black)

	// Bishops
	b.Squares[0][2] = NewPiece(Bishop, White)
	b.Squares[0][5] = NewPiece(Bishop, White)
	b.Squares[7][2] = NewPiece(Bishop, Black)
	b.Squares[7][5] = NewPiece(Bishop, Black)

	// Queens
	b.Squares[0][3] = NewPiece(Queen, White)
	b.Squares[7][3] = NewPiece(Queen, Black)

	// Kings
	b.Squares[0][4] = NewPiece(King, White)
	b.Squares[7][4] = NewPiece(King, Black)
}

func (b *Board) GetPiece(pos string) *Piece {
	row, col, ok := squareToIndex(pos)
	if !ok {
		return nil
	}

	return b.Squares[row][col]
}

// for no reason i thought that color means the attacker's color but it's wrong. I hate myself.
// Here color means the defender's color (the side that owns the square).
func (b *Board) IsCellAttacked(sussyPos string, color Color) bool { //sussy yeah
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			if b.Squares[row][col] != nil && b.Squares[row][col].Color != color {
				fig := *b.Squares[row][col]

				fromPos := string(rune('a'+col)) + string(rune('1'+row))
				mov := Move{
					From: fromPos,
					To:   sussyPos,
				}
				if b.IsLegalMove(fig, mov) { // if any enemy pieces can attack square so this cell is attacked
					return true
				}
			}
		}
	}
	return false
}

func (b *Board) IsCheck(currentColor Color) bool {
	var kingRank, kingFile int
	flag := false
	//finding out where king stays
	for rows := 0; rows < 8; rows++ {
		for cols := 0; cols < 8; cols++ {
			if b.Squares[rows][cols] != nil && b.Squares[rows][cols].Type == King && b.Squares[rows][cols].Color == currentColor {
				kingRank = rows
				kingFile = cols
				flag = true
			}
			if flag {
				break
			}
		}
	}

	kingPos := string(rune('a'+kingFile)) + string(rune('1'+kingRank))

	for rows := 0; rows < 8; rows++ {
		for cols := 0; cols < 8; cols++ {
			if b.Squares[rows][cols] != nil && b.Squares[rows][cols].Color != currentColor {
				fig := *b.Squares[rows][cols]
				fromPose := string(rune('a'+cols)) + string(rune('1'+rows))
				mov := Move{
					From: fromPose,
					To:   kingPos,
				}
				if b.IsLegalMove(fig, mov) { // if any enemy piece can attack king so it is a check
					return true
				}
			}
		}
	}
	return false
}

func (b *Board) IsCheckMate(currentColor Color) bool {
	if !b.IsCheck(currentColor) {
		return false
	}
	var fig Piece
	var kingRank, kingFile int
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			if b.Squares[rank][file] != nil && b.Squares[rank][file].Type == King && b.Squares[rank][file].Color == currentColor {
				kingRank = rank
				kingFile = file
				fig = *b.Squares[rank][file]
				break
			}
		}
	}
	kingPos := string(rune('a'+kingFile)) + string(rune('1'+kingRank))
	for r := kingRank - 1; r <= kingRank+1; r++ { // bro tries to hide
		for f := kingFile - 1; f <= kingFile+1; f++ {
			if isValidCoordinate(f, r) && (r != kingRank || f != kingFile) {
				toPose := string(rune('a'+f)) + string(rune('1'+r))
				mov := Move{
					From: kingPos,
					To:   toPose,
				}
				if b.IsLegalMove(fig, mov) && !b.IsCellAttacked(toPose, fig.Color) {
					return false // bro can escape 
				}
			}
		}
	}

	for r := 0; r < 8; r++ { // bro tries to block the check by any of his pieces
		for f := 0; f < 8; f++ {
			piece := b.Squares[r][f]
			if piece != nil && piece.Color == currentColor {
				fromSquare := string(rune('a'+f)) + string(rune('1'+r))
				for r2 := 0; r2 < 8; r2++ {
					for f2 := 0; f2 < 8; f2++ {
						toSquare := string(rune('a'+f2)) + string(rune('1'+r2))
						move := Move{From: fromSquare, To: toSquare}

						if b.IsLegalMove(*piece, move) {
							// making the move temporarily 
							captured := b.Squares[r2][f2]
							b.Squares[r][f] = nil
							b.Squares[r2][f2] = piece

							kingSafe := !b.IsCellAttacked(
								string(rune('a'+kingFile))+string(rune('1'+kingRank)),
								currentColor,
							)

							// undo the move
							b.Squares[r][f] = piece
							b.Squares[r2][f2] = captured

							if kingSafe {
								return false // bro can block the check or capture the attacker
							}
						}
					}
				}
			}
		}
	}
	return true //Sorry bro, you're checkmated :(
}

func (b *Board) IsStaleMate(currentColor Color) bool {
	if b.IsCheck(currentColor) {
		return false
	}

	for fromRank := 0; fromRank < 8;fromRank++ { // checking all the pieces of current color
		for fromFile := 0;fromFile < 8;fromFile++ {
			if b.Squares[fromRank][fromFile] != nil && b.Squares[fromRank][fromFile].Color == currentColor {
				for toRank := 0;toRank < 8;toRank++ { // checking all the possible moves 
					for toFile := 0;toFile < 8;toFile++ {
						fromSquare := string(rune('a'+fromFile)) + string(rune('1'+fromRank))
						toSquare := string(rune('a'+toFile)) + string(rune('1'+toRank))
						move := Move{From: fromSquare, To: toSquare}
						fig := *b.Squares[fromRank][fromFile]
						if b.IsLegalMove(fig,move) { // if there's any possible move so that isn't stalemate
							return false
						}
					}
				}
			}
		}
	}
	return true
}