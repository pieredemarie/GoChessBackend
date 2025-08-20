package chesslogic

type PieceType int 
const (
	Pawn PieceType = iota //0
	Knight // 1
	Bishop // 2
	Rook // 3
	Queen // 4
	King  // 5
)

type Color int 
const (
	White Color = iota // 0
	Black 
)

type Piece struct {
	Type PieceType
	Color Color
}

func NewPiece(Type PieceType, color Color) *Piece{
	return &Piece{
		Type: Type,
		Color: color,
	}
}

func (p *Piece) IsWhite() bool {
	return p != nil && p.Color == White
}

func (p *Piece) IsBlack() bool {
	return p != nil && p.Color == Black
}

func squareToIndex(square string) (row,col int,ok bool) {
	if len(square) != 2 {
		return 0,0, false
	}

	file := square[0]
	rank := square[1]

	col = int(file - 'a')
	row = int(rank - '1')

	if col < 0 || col > 7 || row < 0 || row > 7 {
		return 0,0, false
	}
	return row,col, true
} 
