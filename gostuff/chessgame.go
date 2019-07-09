package gostuff

import (
	"sync"

	"github.com/notnil/chess"
)

//stores chess game information
type ChessGame struct {
	Type         string
	ID           int
	WhitePlayer  string
	BlackPlayer  string
	WhiteRating  int16
	BlackRating  int16
	GameMoves    []GameMove //stores chess move for games
	Status       string     //white to move, black to move, white won, black won, or draw
	Result       int8       //0 means black won, 1 means white won and 2 means draw. 2 is used intead of 0.5 as database type is int
	GameType     string     //bullet, blitz, standard, correspondence
	TimeControl  int
	BlackMinutes int
	BlackSeconds int
	WhiteMinutes int
	WhiteSeconds int
	StartMinutes int    // used to keep track of start time for correspondence
	PendingDraw  bool   //used to keep track if a player has offered a draw
	Rated        string //Yes if the game is rated, No if the game is unrated
	Spectate     bool
	CountryWhite string
	CountryBlack string
	Validator    *chess.Game

	resetWhiteTime chan bool
	resetBlackTime chan bool
	gameOver       chan bool

	observe Observers // list of user names who are observing this table

	whiteTurn bool //keeps track of whose move it is, true means its whites turn and false means its blacks turn
}

//source and destination of piece moves
type GameMove struct {
	Type string
	ID   int
	S    string // Source move
	T    string // Destination move
	P    string // Promotion piece
	Fen  string // FEN string of the board with the move played
}

// used to unmarshall game ID that is being observed by player(Name)
type SpectateGame struct {
	Type     string
	ID       int `json:"ID,string"`
	Name     string
	Spectate string
}

type Nrating struct {
	Type        string
	WhiteRating float64
	BlackRating float64
}

type Fin struct { //used to hold the result when a player is mated
	Type   string
	ID     int
	Fen    string // FEN string of the board with the move played
	Status string
}

// contains an array of player names observing the table
type Observers struct {
	sync.RWMutex
	Names []string
}

//active and running games on the server
var All = struct {
	sync.RWMutex
	Games map[int]*ChessGame
}{Games: make(map[int]*ChessGame)}

//pending matches in the lobby waiting for someone to accept
var Pending = struct {
	sync.RWMutex
	Matches map[int]*SeekMatch
}{Matches: make(map[int]*SeekMatch)}

//used for quick access to identify two people who are private chatting and playing a game against each other
var PrivateChat = make(map[string]string)

//intitalize all pawns to false as they have not moved yet, and also initialize all en passent to false
func InitGame(gameID int, name string, fighter string) {

	//reset times are used for correspondence
	All.Games[gameID].resetWhiteTime = make(chan bool)
	All.Games[gameID].resetBlackTime = make(chan bool)
	All.Games[gameID].gameOver = make(chan bool)
	All.Games[gameID].whiteTurn = true

	// Long AlgebraicNotation Notation
	All.Games[gameID].Validator = chess.NewGame(chess.UseNotation(chess.LongAlgebraicNotation{}))

	// the players playing the game are also observers
	All.Games[gameID].observe.Names = append(All.Games[gameID].observe.Names, name)
	All.Games[gameID].observe.Names = append(All.Games[gameID].observe.Names, fighter)
}
