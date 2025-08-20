package chesslogic

type Board struct {
	Squares [8][8]*Piece
	Turn Color
}

type CapablancaBoard struct { //In the future will be done this type of chess
	Squares [10][8]*Piece
	Turn Color
}

func NewBoard() *Board {
	b := &Board{}
	b.ArrangeFigures()
	return b
}

func (b *Board) ArrangeFigures() {
	//pawns 
	for cols := 0; cols < 8;cols++ {
		b.Squares[1][cols] = NewPiece(Pawn,White)
		b.Squares[6][cols] = NewPiece(Pawn,Black)
	}

	//Rooks
	b.Squares[0][0] = NewPiece(Rook,White)
	b.Squares[0][7] = NewPiece(Rook,Black)
	b.Squares[7][0] = NewPiece(Rook,White)
	b.Squares[7][7] = NewPiece(Rook,Black)

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
	row,col, ok := squareToIndex(pos)
	if !ok {
		return nil
	}

	return b.Squares[row][col]
}