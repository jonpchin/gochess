package gostuff

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"time"
	"bufio"
	"encoding/base64"
	"encoding/hex"
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
	ID           int
	White        string
	Black        string
	GameType     string
	WhiteRating  int
	BlackRating  int
	WhiteMinutes int
	WhiteSeconds int
	BlackMinutes int
	BlackSeconds int
	TimeControl  int
	Moves        string //json this back to a struct of Moves
	Total        int
	Result       int //2 means draw
	Status       string
	Date         string
	Time         string
	Rated		 string
}

var db *sql.DB

//returns false if database setup failed
func DbSetup(backLocation string) bool {

	//Checks if backup folder for database export exists
	exists, err := isDirOrFileExists(backLocation)
	if err != nil{
		fmt.Println("database.go DbSetup 0, error checking if directory exists", err)
	}
	if exists == false{
		err := os.Mkdir(backLocation, 0777)
		if err != nil{
			fmt.Println("database.go DbSetup 1, error creating backup directory", err)
		}
	}	
	
	dbString, database := ReadFile()
	//connecting to database
	db, err = sql.Open("mysql", dbString)
	//	db.SetMaxIdleConns(20)
	if err != nil {
		fmt.Println("Error opening Database DBSetup 2", err)
		return false
	}
	
	if db.Ping() != nil {
		
		var result string
		//checking if database exist
		db.QueryRow("SHOW DATABASES LIKE '" + database + "'").Scan(&result)
		if result == ""{
			fmt.Println("database.go DbSetup 3 Database", database, "does not exist")
			
			// TODO: Attempt to import an existing backup database, if
			// that is not available then import template database
			// if that is also not available then return false
			// if database is imported attempt to reconnect to database
			// using same credentials
			
			return false
			
		}else{
			fmt.Println("MySQL is down!!!")
			return false
		}
		
	}
	fmt.Println("MySQL is now connected.")
	return true
}

func ReadFile() (string, string) {
	config, err := os.Open("secret/config.txt")
	defer config.Close()
	if err != nil {
		log.Println("database.go ReadFile 1 ", err)
	}

	scanner := bufio.NewScanner(config)
	//creating new string to append database info
	dbString := ""
	scanner.Scan()
	//user
	dbData := scanner.Text()
	dbString = dbString + dbData + ":"

	//pass
	scanner.Scan()
	dbData = scanner.Text()
	//decode
	ans, _ := hex.DecodeString(dbData)

	result, _ := base64.StdEncoding.DecodeString(string(ans))
	answer := string(result)

	dbString = dbString + answer + "@tcp("
	//host
	scanner.Scan()
	dbData = scanner.Text()
	dbString = dbString + dbData + ":"
	//port
	scanner.Scan()
	dbData = scanner.Text()
	dbString = dbString + dbData + ")/"
	//database name
	scanner.Scan()
	dbData = scanner.Text()
	db := dbData
	dbString = dbString + dbData

	return dbString, db
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
	err2 := db.QueryRow("SELECT bullet, blitz, standard FROM rating WHERE username=?",
		name).Scan(&bullet, &blitz, &standard)

	if err2 != nil {
		log.Println("database.go GetRating 1 ", err2)
		return "No such player", 0, 0, 0
	}
	return "", bullet, blitz, standard
}

//fetches players bullet, blitz and standard rating and RD
func GetRatingAndRD(name string) (errRate string, bullet, blitz, standard, bulletRD,
	blitzRD, standardRD float64) {
		
	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log.SetOutput(problems)

	//check if database connection is open
	if db.Ping() != nil {
		log.Println("DATABASE DOWN! @GetRatingAndRD()")
		return "DB down @GetRatingAndRD()", 0, 0, 0, 0, 0, 0
	}

	//looking up players rating
	err2 := db.QueryRow("SELECT bullet, blitz, standard, bulletRD, blitzRD, standardRD " +
		"FROM rating WHERE username=?", name).Scan(&bullet, &blitz, &standard,
		 &bulletRD, &blitzRD, &standardRD)

	if err2 != nil {
		log.Println("database.go GetRating 2 ", err2)
		return "No such player", 0, 0, 0, 0, 0, 0
	}
	return "", bullet, blitz, standard, bulletRD, blitzRD, standardRD
}

//updates both players chess rating
func updateRating(gameType string, white string, whiteRating float64, whiteRD float64,
	black string, blackRating float64, blackRD float64) {

	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log.SetOutput(problems)

	//check if database connection is open
	if db.Ping() != nil {
		log.Println("database.go updateRating 1 DATABASE DOWN!")
		return
	}

	//setting verify to yes and deleting row from activate table
	stmt, err := db.Prepare("UPDATE rating SET " + gameType + "=?," + gameType +
		"RD=?" + " where username=?")
	if err != nil {
		log.Println("database.go updateRating 2 ",err)
		return
	}

	res, err := stmt.Exec(whiteRating, whiteRD, white)
	if err != nil {
		log.Println("database.go updateRating 3 ", err)
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		log.Println("database.go updateRating 4 ", err)
		return
	}

	log.Printf("%s rating has changed and %d row was updated.\n", white, affect)

	//setting verify to yes and deleting row from activate table
	stmt, err = db.Prepare("UPDATE rating SET " + gameType + "=?," + gameType +
		"RD=?" + " where username=?")
	if err != nil {
		log.Println("database.go updateRating 5 ", err)
		return
	}

	res, err = stmt.Exec(blackRating, blackRD, black)
	if err != nil {
		log.Println("database.go updateRating 6 ", err)
		return
	}
	affect, err = res.RowsAffected()
	if err != nil {
		log.Println("database.go updateRating 7 ", err)
		return
	}

	log.Printf("%s rating has changed and %d row was updated.\n", black, affect)
}

//stores game into MySQL database as a string
func storeGame(totalMoves int, allMoves []byte, game *ChessGame) {
	moves := string(allMoves)

	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log.SetOutput(problems)

	//check if database connection is open
	if db.Ping() != nil {
		log.Println("database.go storeGame 1 DATABASE DOWN!")
		return
	}

	//preparing token activation
	stmt, err := db.Prepare("INSERT games SET white=?, black=?, gametype=?, rated=?, " +
		"whiterating=?, blackrating=?, timecontrol=?, moves=?, totalmoves=?, " +
		"result=?, status=?, date=?, time=?")
	if err != nil {
		log.Println("database.go storeGame 2 ", err)
		return
	}
	date := time.Now()
	res, err := stmt.Exec(game.WhitePlayer, game.BlackPlayer, game.GameType, game.Rated,
		game.WhiteRating, game.BlackRating, game.TimeControl, moves, totalMoves,
		game.Result, game.Status, date, date)
	if err != nil {
		log.Println("database.go storeGame 3 ", err)
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Println("database.go storeGame 4 ", err)
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
		log.Println("database.go GetGames 1 ", err)
	}
	defer rows.Close()

	var all GoGame

	for rows.Next() {

		err = rows.Scan(&all.ID, &all.White, &all.Black, &all.GameType, &all.Rated, 
			&all.WhiteRating, &all.BlackRating, &all.TimeControl, &all.Moves, 
			&all.Total, &all.Result, &all.Status, &all.Date, &all.Time)
		if err != nil {

			log.Println("database.go GetGames 2 ", err)
		}
		storage = append(storage, all)

	}
	return storage

}

func GetSaved(name string) (storage []GoGame) {
	problems, err := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log.SetOutput(problems)

	rows, err := db.Query("SELECT * FROM saved WHERE white=? or black=?", name, name)
	if err != nil {
		log.Println("database.go GetSaved 1 ", err)
	}
	defer rows.Close()

	var all GoGame

	for rows.Next() {

		err = rows.Scan(&all.ID, &all.White, &all.Black, &all.GameType, &all.Rated, 
			&all.WhiteRating, &all.BlackRating, &all.BlackMinutes, &all.BlackSeconds, 
			&all.WhiteMinutes, &all.WhiteSeconds, &all.TimeControl, &all.Moves, &all.Total, 
			&all.Status, &all.Date, &all.Time)
		if err != nil {
			log.Println("database.go GetSaved 2 ", err)
		}
		storage = append(storage, all)
	}
	return storage
}

//fetches saved game from database
func fetchSavedGame(id string, user string) bool {

	problems, err := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log.SetOutput(problems)

	var white string
	var black string
	var gametype string
	var rated string
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

	err = db.QueryRow("SELECT white, black, gametype, rated, whiterating, " +
		"blackrating, blackminutes, blackseconds, whiteminutes, whiteseconds, " +
		"timecontrol, moves, totalmoves, status FROM saved WHERE id=?", id).Scan(&white,
		&black, &gametype, &rated, &whiterating, &blackrating, &blackminutes, 
		&blackseconds, &whiteminutes, &whiteseconds, &timecontrol, &moves,
		&totalmoves, &status)
	if err != nil {
		log.Println("database.go fetchSavedGame 1 ", err)
	}

	var game ChessGame
	game.Type = "chess_game"
	var holder []Move
	//White for white to move or Black for black to move, white won, black won, stalemate or draw.
	game.Status = status

	storage := []byte(moves)
	err = json.Unmarshal(storage, &holder)
	if err != nil {
		log.Println("database.go fetchSavedGame 2", err)
	}
	game.GameMoves = holder
	game.WhitePlayer = white
	game.BlackPlayer = black
	game.WhiteRating = whiterating
	game.BlackRating = blackrating
	game.TimeControl = timecontrol
	game.GameType = gametype
	game.Rated = rated
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

	//intitalizes all the variables of the game
	initGame(game.ID)

	var result bool
	total := len(game.GameMoves)

	for i := 0; i < total; i++ {
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
