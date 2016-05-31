package gostuff

import (
	"sync"
)

//stores chess game information
type ChessGame struct {
	Type         string
	ID           int16
	WhitePlayer  string
	BlackPlayer  string
	WhiteRating  int16
	BlackRating  int16
	GameMoves    []Move //stores chess move for games
	Status       string //white to move, black to move, white won, black won, or draw
	Result       int8   //0 means black won, 1 means white won and 2 means draw. 2 is used intead of 0.5 as database type is int
	GameType     string //bullet, blitz or standard
	TimeControl  int
	BlackMinutes int
	BlackSeconds int
	WhiteMinutes int
	WhiteSeconds int
	PendingDraw  bool //used to keep track if a player has offered a draw
}

//source and destination of piece moves
type GameMove struct {
	Type         string
	ID           int16
	Source       string
	Target       string
	WhiteMinutes int
	WhiteSeconds int
	BlackMinutes int
	BlackSeconds int
}

//only holds source and destination
type Move struct {
	S string
	T string
}

type Nrating struct {
	Type        string
	WhiteRating float64
	BlackRating float64
}

type Fin struct { //used to hold the result when a player is mated
	Type   string
	ID     int16
	Status string
}

//used to verify players games
type Table struct {
	ChessBoard [][]string

	whitePawns [8]bool //stores array of booleans indicating whether or not the pawns have made their first move yet
	blackPawns [8]bool
	whitePass  [8]bool //stores array to indicate whether or not the pawns can perform an enpassent
	blackPass  [8]bool

	whiteTurn bool //keeps track of whose move it is, true means its whites turn and false means its blacks turn

	wkMoved bool //if king moved or not
	bkMoved bool

	wkrMoved bool
	wqrMoved bool
	bkrMoved bool
	bqrMoved bool

	whiteKingX int8 //stores location of king for easy access
	whiteKingY int8
	blackKingX int8
	blackKingY int8

	kingUpdate bool //used to figure out if king position needs to be updated
	rookUpdate bool //castling rights

	blackOldX int8 //used to store old coordinates for king when in check
	blackOldY int8
	whiteOldX int8
	whiteOldY int8

	undoWPass bool //if this is true then white just did an en passent and it is used in undoMove()
	undoBPass bool
	
	threeRepS [6]string //stores last 6 moves to check for three repetition of source move
	threeRepT [6]string //stores last 6 moves to check for three repetition of target move
	repIndex int8 //stores the index of where threeRep should replace the next spot
	
	whiteTimeOut chan bool
	blackTimeOut chan bool
	gameOver 	 chan bool
	Connection   chan bool
}

//active and running games on the server
var All = struct {
	sync.RWMutex
	Games map[int16]*ChessGame
}{Games: make(map[int16]*ChessGame)}

//pending matches in the lobby waiting for someone to accept
var Pending = struct {
	sync.RWMutex
	Matches map[int16]*SeekMatch
}{Matches: make(map[int16]*SeekMatch)}

//used to verify each move on the board
var Verify = struct {
	sync.RWMutex
	AllTables map[int16]*Table
}{AllTables: make(map[int16]*Table)}

//used for quick access to identify two people who are private chatting and playing a game against each other
var PrivateChat = make(map[string]string)

//intitalize all pawns to false as they have not moved yet, and also initialize all en passent to false
func initGame(gameID int16) {
	for i := 0; i < 8; i++ {
		Verify.AllTables[gameID].whitePawns[i] = false
		Verify.AllTables[gameID].blackPawns[i] = false
		Verify.AllTables[gameID].whitePass[i] = false
		Verify.AllTables[gameID].blackPass[i] = false
	}
	//castling rights init for kings
	Verify.AllTables[gameID].wkMoved = false
	Verify.AllTables[gameID].bkMoved = false
	//castling rights int for rooks
	Verify.AllTables[gameID].wkrMoved = false
	Verify.AllTables[gameID].wqrMoved = false
	Verify.AllTables[gameID].bkrMoved = false
	Verify.AllTables[gameID].bqrMoved = false
	//storing coordinates for starting location of both kings, X is row and Y is col
	Verify.AllTables[gameID].whiteKingX = 7
	Verify.AllTables[gameID].whiteKingY = 4
	Verify.AllTables[gameID].blackKingX = 0
	Verify.AllTables[gameID].blackKingY = 4

	Verify.AllTables[gameID].kingUpdate = false
	Verify.AllTables[gameID].rookUpdate = false

	Verify.AllTables[gameID].whiteTurn = true
	Verify.AllTables[gameID].whiteOldX = 7
	Verify.AllTables[gameID].whiteOldY = 4
	Verify.AllTables[gameID].blackOldX = 0
	Verify.AllTables[gameID].blackOldY = 4

	Verify.AllTables[gameID].undoWPass = false
	Verify.AllTables[gameID].undoBPass = false
	
	Verify.AllTables[gameID].repIndex = 0
	
	Verify.AllTables[gameID].whiteTimeOut = make(chan bool)
	Verify.AllTables[gameID].blackTimeOut = make(chan bool)
	Verify.AllTables[gameID].gameOver = make(chan bool)
	//intitalize all starting moves in 
}