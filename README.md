# GoChessBackend

Simple online chess backend written in Go.  
Supports user accounts, matchmaking, WebSocket play, and game storage in PostgreSQL.

## Database schema (PostgreSQL)

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    rating INT DEFAULT 900
);

CREATE TABLE games (
    id TEXT PRIMARY KEY,
    white_id INT REFERENCES users(id),
    black_id INT REFERENCES users(id),
    status TEXT,
    created_at TIMESTAMP,
    finished_at TIMESTAMP,
    result TEXT
);

CREATE TABLE moves (
    id SERIAL PRIMARY KEY,
    game_id TEXT REFERENCES games(id),
    move_number INT,
    from_square TEXT,
    to_square TEXT,
    piece TEXT,
    created_at TIMESTAMP
);
```
API Endpoints
POST /auth/register - register user

POST /auth/login - login user

GET /user/me - get current user (JWT required)

POST /game/find - find a game

GET /game/:id - get game state

GET /game/ws - Websocket for moves and updates 

TODO: 
PGN export
