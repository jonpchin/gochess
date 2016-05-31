package gostuff

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"encoding/json"
	"log"
	"os"
	"time"
)

//stores information about players games extracted from database when player clicks there profile
type ProfileGames struct {
	User       string
	SessionID  string
	Bullet     float64
	Blitz      float64
	Standard   float64
	BulletRD   float64
	BlitzRD    float64
	StandardRD float64
	Games      []GoGame
}

//an individual game
type GoGame struct {
	ID          int
	White       string
	Black       string
	GameType    string
	WhiteRating int
	BlackRating int
	WhiteMinutes int
	WhiteSeconds int
	BlackMinutes int
	BlackSeconds int
	TimeControl int
	Moves       string //json this back to a struct of Moves
	Total       int
	Result      int //2 means draw
	Status      string
	Date        string
	Time        string
}

var db *sql.DB
//returns false if database setup failed
func DbSetup() bool{
	
	dbString := ReadFile()
	var err error
	//connecting to database
	db, err = sql.Open("mysql", dbString)
	//	db.SetMaxIdleConns(20)
	if err != nil {
		fmt.Println("Error opening Database DBSetup", err)
		return false
	}

	if db.Ping() != nil {
		fmt.Println("MySQL is down!!!")
		return false
	}
	fmt.Println("MySQL is now connected.")
	return true
}

//fetches players bullet, blitz and standard rating
func GetRating(name string) (errMessage string, bullet, blitz, standard int16) {

	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log.SetOutput(problems)

	//check if database connection is open
	if db.Ping() != nil {
		log.Println("DATABASE DOWN! @GetRating() ping")
		return "Database down", 0, 0, 0
	}

	//looking up players rating
	err2 := db.QueryRow("SELECT bullet, blitz, standard FROM rating WHERE username=?", name).Scan(&bullet, &blitz, &standard)

	if err2 != nil {
		log.Println(err2)
		return "No such player", 0, 0, 0
	}
	return "", bullet, blitz, standard

}

//fetches players bullet, blitz and standard rating and RD
func GetRatingAndRD(name string) (errRate string, bullet, blitz, standard, bulletRD, blitzRD, standardRD float64) {
	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log.SetOutput(problems)
	
	//check if database connection is open
	if db.Ping() != nil {
		log.Println("DATABASE DOWN! @GetRatingAndRD()")
		return "DB down @GetRatingAndRD()", 0, 0, 0, 0, 0, 0
	}

	//looking up players rating
	err2 := db.QueryRow("SELECT bullet, blitz, standard, bulletRD, blitzRD, standardRD FROM rating WHERE username=?", name).Scan(&bullet, &blitz, &standard, &bulletRD, &blitzRD, &standardRD)

	if err2 != nil {
		log.Println(err2)
		return "No such player", 0, 0, 0, 0, 0, 0
	}
	return "", bullet, blitz, standard, bulletRD, blitzRD, standardRD
}

//updates both players chess rating
func updateRating(gameType string, white string, whiteRating float64, whiteRD float64, black string, blackRating float64, blackRD float64) {

	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log.SetOutput(problems)

	//check if database connection is open
	if db.Ping() != nil {
		log.Println("DATABASE DOWN!")
		return
	}

	//setting verify to yes and deleting row from activate table
	stmt, err := db.Prepare("UPDATE rating SET " + gameType + "=?," + gameType + "RD=?" + " where username=?")
	if err != nil {
		log.Println(err)
		return
	}

	res, err := stmt.Exec(whiteRating, whiteRD, white)
	if err != nil {
		log.Println(err)
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("%s rating has changed and %d row was updated.\n", white, affect)

	//setting verify to yes and deleting row from activate table
	stmt, err = db.Prepare("UPDATE rating SET " + gameType + "=?," + gameType + "RD=?" + " where username=?")
	if err != nil {
		log.Println(err)
		return
	}

	res, err = stmt.Exec(blackRating, blackRD, black)
	if err != nil {
		log.Println(err)
		return
	}
	affect, err = res.RowsAffected()
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("%s rating has changed and %d row was updated.\n", black, affect)
}

//stores game into MySQL database as a string
func storeGame(totalMoves int, allMoves []byte, game *ChessGame) {
	moves := string(allMoves)
	//	fmt.Println("The game moves are ", moves)

	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log.SetOutput(problems)

	//check if database connection is open
	if db.Ping() != nil {
		log.Println("DATABASE DOWN!")
		return
	}

	//preparing token activation
	stmt, err := db.Prepare("INSERT games SET white=?, black=?, gametype=?, whiterating=?, blackrating=?, timecontrol=?, moves=?, totalmoves=?, result=?, status=?, date=?, time=?")
	if err != nil {
		log.Println(err)
		return
	}
	date := time.Now()
	res, err := stmt.Exec(game.WhitePlayer, game.BlackPlayer, game.GameType, game.WhiteRating, game.BlackRating, game.TimeControl, moves, totalMoves, game.Result, game.Status, date, date)
	if err != nil {
		log.Println(err)
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Game ID %d has been updated in the games table.\n", id)

}

//gets all games by player from database and stores them in array of structs
func GetGames(name string) (storage []GoGame) {

	problems, err := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log.SetOutput(problems)

	rows, err := db.Query("SELECT * FROM games WHERE white=? or black=?", name, name)
	if err != nil {
		fmt.Println("error is", err)
	}
	defer rows.Close()

	var all GoGame

	for rows.Next() {

		err = rows.Scan(&all.ID, &all.White, &all.Black, &all.GameType, &all.WhiteRating, &all.BlackRating, &all.TimeControl, &all.Moves, &all.Total, &all.Result, &all.Status, &all.Date, &all.Time)
		if err != nil {

			fmt.Println("The error  is", err)
		}
		storage = append(storage, all)

	}
	return storage

}

func GetSaved(name string)(storage []GoGame){
	problems, err := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log.SetOutput(problems)

	rows, err := db.Query("SELECT * FROM saved WHERE white=? or black=?", name, name)
	if err != nil {
		fmt.Println("error is", err)
	}
	defer rows.Close()

	var all GoGame

	for rows.Next() {

		err = rows.Scan(&all.ID, &all.White, &all.Black, &all.GameType, &all.WhiteRating, &all.BlackRating, &all.BlackMinutes, &all.BlackSeconds, &all.WhiteMinutes, &all.WhiteSeconds, &all.TimeControl, &all.Moves, &all.Total, &all.Status, &all.Date, &all.Time)
		if err != nil {

			fmt.Println("The error  is", err)
		}
		storage = append(storage, all)

	}
	return storage
}

//fetches saved game from database
func fetchSavedGame(id string, user string) bool{
	
	problems, err := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log.SetOutput(problems)
	
	var white string
	var black string
	var gametype string
	var whiterating int16
	var blackrating int16
	var blackminutes int
	var blackseconds int
	var whiteminutes int
	var whiteseconds int
	var timecontrol int
	var moves string
	var totalmoves int
	var status string
	
	
	err = db.QueryRow("SELECT white, black, gametype, whiterating, blackrating, blackminutes, blackseconds, whiteminutes, whiteseconds, timecontrol, moves, totalmoves, status FROM saved WHERE id=?", id).Scan(&white, &black, &gametype, &whiterating, &blackrating, &blackminutes, &blackseconds, &whiteminutes, &whiteseconds, &timecontrol, &moves, &totalmoves, &status)
	if err != nil{
		fmt.Println("database.go fetchSavedGame 1 ", err)
	}
	
	
	var game ChessGame
	game.Type = "chess_game"
	var holder []Move
	//White for white to move or Black for black to move, white won, black won, stalemate or draw.
	game.Status = status
	
	storage := []byte(moves)
	err = json.Unmarshal(storage, &holder)
	if err != nil{
		fmt.Println("database.go fetchSavedGame 2", err)
	}
	game.GameMoves = holder
    game.WhitePlayer = white
	game.BlackPlayer = black
	game.WhiteRating = whiterating
	game.BlackRating = blackrating
	game.TimeControl = timecontrol
	game.GameType = gametype
	game.Status = status

	game.WhiteMinutes = whiteminutes
	game.WhiteSeconds = whiteseconds
	game.BlackMinutes = blackminutes
	game.BlackSeconds = blackseconds

	game.PendingDraw = false

	var start int16 = 0
	for {
		if _, ok := All.Games[start]; ok {
			start++
		} else {
			break
		}
	}
	//value := fmt.Sprintf("%d", start)
	game.ID = start
	//used in backend to keep track of all pending games waiting for a player to accept
	All.Games[start] = &game

	//setting up back end move verification
	var table Table
	table.ChessBoard = [][]string{
		[]string{"bR", "bN", "bB", "bQ", "bK", "bB", "bN", "bR"},
		[]string{"bP", "bP", "bP", "bP", "bP", "bP", "bP", "bP"},
		[]string{"-", "-", "-", "-", "-", "-", "-", "-"},
		[]string{"-", "-", "-", "-", "-", "-", "-", "-"},
		[]string{"-", "-", "-", "-", "-", "-", "-", "-"},
		[]string{"-", "-", "-", "-", "-", "-", "-", "-"},
		[]string{"wP", "wP", "wP", "wP", "wP", "wP", "wP", "wP"},
		[]string{"wR", "wN", "wB", "wQ", "wK", "wB", "wN", "wR"},
	}
	Verify.AllTables[game.ID] = &table
	//intitalizes all the variables of the game
	initGame(game.ID)


	var result bool
	total := len(game.GameMoves)	
	
	for i:=0; i<total; i++{
		result = chessVerify(game.GameMoves[i].S, game.GameMoves[i].T, game.ID)
		if result == false {
			log.Println("something went wrong in move validation in fetchSavedGame of saved game id ", game.ID)
			//undo all game setup and break out
			delete(Verify.AllTables, game.ID)	
			delete(All.Games, game.ID)
			return false
		}
	}
	PrivateChat[white] = black
	PrivateChat[black] = white

	//starting white's clock first, this goroutine will keep track of both players clock for this game
	go setClocks(game.ID, user)
	
	//delete saved game from database now that its in memory
	stmt, err := db.Prepare("DELETE FROM saved where id=?")
	if err != nil {
		log.Println("database.go fetchSavedGame 3 ", err)
		return false
	}

	res, err := stmt.Exec(id)
	if err != nil {
		log.Println("database.go fetchSavedGame 4 ", err)
		return false
	}
	stmt.Close()
	affect, err := res.RowsAffected()
	if err != nil {
		log.Println("database.go fetchSavedGame 5 ", err)
		return false
	}
	log.Printf("%d row was deleted from the saved table by user %s\n", affect, user)

	return true
}