package game

import (
	chesslogic "backendChess/pkg/chessLogic"
	"context"
	"math/rand"
	"strings"
	"sync"
)

type Room struct {
	ID      string
	White *Client
	Black *Client

	Game *chesslogic.Game
	broadcast chan []byte 
	mu sync.Mutex
	closed bool 
}

type Matchmaker struct {
	mu sync.Mutex 
	queue []*Client
}

var rooms sync.Map

func NewRoom(a,b *Client) *Room {
	r := &Room{
		ID:  randRoomID(),
		broadcast: make(chan []byte, 64),
	}

	if rand.Intn(2) == 0 {
		r.White, r.Black = a, b
	} else  {
		r.White, r.Black = b, a
	}
	r.White.Color = chesslogic.White
	r.Black.Color = chesslogic.Black

	playerA := chesslogic.Player{Color: a.Color}
	playerB := chesslogic.Player{Color: b.Color}
	r.Game = chesslogic.NewGame(&playerA,&playerB)
	r.White.room = r 
	r.Black.room = r 

	rooms.Store(r.ID,r)
	return r
}

var matchmaker = &Matchmaker{}

func (m *Matchmaker) Enqueue(c *Client) {
	m.mu.Lock()
	m.queue = append(m.queue, c)
	m.mu.Unlock()
	go m.tryMatch()
}

func (m *Matchmaker) RemoveFromQueue(c *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i,x := range m.queue {
		if x == c {
			m.queue = append(m.queue[:i],m.queue[i+1:]...)
			return
		}
	}
}

func (m *Matchmaker) tryMatch() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.queue) < 2 { // if there's less than 2 players it's useless to search
		return 
	}
	used := make(map[int]bool)
	var pairs [][2]*Client

	for i := 0;i < len(m.queue); i++ {
		if used[i] { // skipping the used one guys 
			continue 
		}
		a := m.queue[i]
		bestJ := -1 // searching for best opponent 
		bestDiff := 1 << 30 
		for j := i + 1;j < len(m.queue);j++ {
			if used[j] { // skipping the used one guys 
				continue
			}
			b := m.queue[j]
			diff := abs(a.Rating - b.Rating)
			if diff <= 100 && diff < bestDiff {
				bestDiff = diff 
				bestJ = j
			}
		}
		if bestJ != -1 { // if best opponent succesfully was found
			used[i], used[bestJ] = true, true // ticking them as used one guys
			pairs = append(pairs, [2]*Client{a,m.queue[bestJ]}) // putting them in pairs
		}
	}

	if len(pairs) == 0 { //no one wanna play chess (that's so mean. start playing chess)
		return 
	}

	newQueue := make([]*Client,0,len(m.queue))
	for idx,c := range m.queue {
		if !used[idx] { 
			newQueue = append(newQueue, c)
		}
	}
	m.queue = newQueue

	for _, pr := range pairs { //Creating game for all who love chess (you must love it too)
		CreateRoom(pr[0],pr[1])
	}
}

func CreateRoom(a,b *Client) {
	r := NewRoom(a,b)

	a.writeJSON(Envelope{Type:MsgMatch,Data: mustRaw(MatchData{RoomID: r.ID, Color: a.Color})})
	b.writeJSON(Envelope{Type:MsgMatch,Data: mustRaw(MatchData{RoomID: r.ID, Color: b.Color})})

	r.pushState("ok")

	ctxA, cancelA := context.WithCancel(context.Background())
	ctxB, cancelB := context.WithCancel(context.Background())

	go func() { a.writePump(ctxA); cancelB() }()
	go func() { b.writePump(ctxB); cancelA()}()
}

func (r *Room) pushState(status string) {
	fen := "" //TODO: make BoardFEN() in chesslogic to send FEN 
	state := Envelope{Type: MsgState, Data: mustRaw(StateData{
		BoardFEN: fen,
		Turn: r.Game.Board.Turn,
		Status: status,
	})}
	r.White.writeJSON(state)
	r.Black.writeJSON(state)
}

func (r *Room) broadcastUpdate(m MoveData, status string) {
	fen := ""
	upd := Envelope{Type: MsgUpdate, Data: mustRaw(UpdateData{
		LastMove: m,BoardFEN: fen, Turn: r.Game.Board.Turn, Status: status,
	})}
	r.White.writeJSON(upd)
	r.Black.writeJSON(upd)
}

func (r *Room) onMove(from *Client,m MoveData) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return
	}

	toMoveColor := r.Game.Board.Turn
	if from.Color != toMoveColor {
		from.writeJSON(errrorEnv("not your turn"))
		return
	}

	mv := chesslogic.Move {
		From: strings.ToLower(m.From),
		To: strings.ToLower(m.To),
		Promotion: m.Promotion,
	}
	if err := r.Game.MakeMove(mv); err != nil {
		from.writeJSON(errrorEnv("illegal move: " + err.Error()))
		return
	}

	status := "ok"
	if r.Game.Board.IsCheck(from.Color) {
		status = "checkmate"
		r.gameOver("checkmate", winnerByColor(chesslogic.OppositeColor(from.Color)))
	} else if r.Game.Board.IsStaleMate(from.Color) {
		status = "stalemate"
		r.gameOver("stalemate", "draw")
	}

	if !r.closed {
		r.broadcastUpdate(m,status)
	}
	
}

func (r *Room) gameOver(reason, winner string) {
	if r.closed {
		return
	}
	r.closed = true
	ev := Envelope{Type: MsgGameOver, Data: mustRaw(GameOverData{Reason: reason, Winner: winner})}
	r.White.writeJSON(ev)
	r.Black.writeJSON(ev)
	rooms.Delete(r.ID)
}

func (r *Room) onResign(from *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return
	}
	win := "white"
	if from.Color == chesslogic.White {
		win = "black"
	}
	r.gameOver("resign", win)
}

func (r *Room) onDisconnect(from *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return 
	}
	win := "white"
	if from.Color == chesslogic.White {
		win = "black"
	}
	r.gameOver("disconnect", win)
}
