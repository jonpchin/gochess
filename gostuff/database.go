package gostuff

import (
	"bufio"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//stores information about players games extracted from database when player clicks there profile
type ProfileGames struct {
	User             string
	SessionID        string
	Bullet           float64
	Blitz            float64
	Standard         float64
	Correspondence   float64
	BulletRD         float64
	BlitzRD          float64
	StandardRD       float64
	CorrespondenceRD float64
	Games            []GoGame
	GameID           int
	Opponent         string
	Days             string
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
	Rated        string
	CountryWhite string
	CountryBlack string
}

var db *sql.DB

//returns false if database setup failed, backup directory is passed in
func DbSetup(backup string) bool {

	//Checks if backup folder for database export exists
	exists, err := isDirOrFileExists(backup)
	if err != nil {
		fmt.Println("database.go DbSetup 0, error checking if directory exists", err)
	}
	if exists == false {
		err := os.Mkdir(backup, 0777)
		if err != nil {
			fmt.Println("database.go DbSetup 1, error creating backup directory", err)
		}
	}

	var checkDBConnectFile = "secret/checkdb.txt"
	var sqlOpenFile = "secret/config.txt"

	if IsEnvironmentTravis() {
		checkDBConnectFile = "data/dbtravis.txt"
		sqlOpenFile = "data/dbtravis.txt"
	}

	// make sure MySQL connection is alive before proceeding
	if CheckDBConnection(checkDBConnectFile) == false {
		return false
	}

	dbString, database := ReadFile(sqlOpenFile)
	db, err = sql.Open("mysql", dbString)

	if err != nil {
		fmt.Println("Error opening Database DBSetup 2", err)
		return false
	}

	//if database ping fails here that means connection is alive but database is missing
	if db.Ping() != nil {
		fmt.Println("Database", database, "does not exist")
		fmt.Println("Please wait while database is imported...")

		result := importDatabase()
		if result == false {
			result = importTemplateDatabase()
			if result == false {
				fmt.Println("database.go Dbsetup FAILED to import both databases!")
				return false
			} else {
				fmt.Println("Empty template database successfully imported!")
			}
		} else {
			fmt.Println(database, "database successfully imported!")
		}

		// Pinging database again to see if newly database exists
		if db.Ping() != nil {
			fmt.Println("database.go Dbsetup 5 ", database, " is still missing after import!!!")
			return false
		}
	}
	fmt.Println("MySQL is now connected.")
	return true
}

// Returns global database handler
func DbConnect() *sql.DB {
	return db
}

// Sets the global database handler
func SetDb(dataDb *sql.DB) {
	db = dataDb
}

// Returns true if the environment is in Travis
func IsEnvironmentTravis() bool {
	if os.Getenv("GOCHESSENV") == "travis" {
		return true
	}
	return false
}

// Returns true if the environment is in App Veyor
func IsEnvironmentAppVeyor() bool {
	if os.Getenv("APPVEYOR") == "True" {
		return true
	}
	return false
}

// checks if database connection is open, returns true if MySQL is running
// the parameter path is where the text file is located
func CheckDBConnection(path string) bool {

	config, err := os.Open(path)
	defer config.Close()
	if err != nil {
		fmt.Println("database.go checkDBConnection 1", err)
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
	dbString = dbString + dbData

	var testDB *sql.DB
	testDB, err = sql.Open("mysql", dbString)

	if err != nil {
		fmt.Println("Error opening Database checkDBConnection 2", err)
		return false
	}
	defer testDB.Close()

	if testDB.Ping() != nil {
		fmt.Println("Database ping failed. Please check if the database server is running.")
		return false
	}
	return true
}

// the parameter path is where the text file is located containing the database connection info
// if password is blank when encoded it will be blank when decoded
func ReadFile(path string) (string, string) {
	config, err := os.Open(path)
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
func GetRating(name string) (errMessage string, bullet, blitz, standard int16, correspondence int16) {

	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

	//check if database connection is open
	if db.Ping() != nil {
		log.Println("DATABASE DOWN!")
		return "Database down", 0, 0, 0, 0
	}

	//looking up players rating
	err2 := db.QueryRow("SELECT bullet, blitz, standard, correspondence FROM rating WHERE username=?",
		name).Scan(&bullet, &blitz, &standard, &correspondence)

	if err2 != nil {
		log.Println(err2)
		return "No such player", 0, 0, 0, 0
	}
	return "", bullet, blitz, standard, correspondence
}

//fetches players bullet, blitz and standard rating and RD
func GetRatingAndRD(name string) (errRate string, bullet, blitz, standard, correspondence, bulletRD,
	blitzRD, standardRD float64, correspondenceRD float64) {

	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

	//check if database connection is open
	if db.Ping() != nil {
		log.Println("DATABASE DOWN!")
		return "DB down @GetRatingAndRD()", 0, 0, 0, 0, 0, 0, 0, 0
	}

	//looking up players rating
	err2 := db.QueryRow("SELECT bullet, blitz, standard, correspondence, bulletRD,"+
		" blitzRD, standardRD, correspondenceRD FROM rating WHERE username=?",
		name).Scan(&bullet, &blitz, &standard, &correspondence,
		&bulletRD, &blitzRD, &standardRD, &correspondenceRD)

	if err2 != nil {
		log.Println(err2)
		return "No such player", 0, 0, 0, 0, 0, 0, 0, 0
	}
	return "", bullet, blitz, standard, correspondence, bulletRD, blitzRD, standardRD, correspondenceRD
}

//updates both players chess rating
func updateRating(gameType string, white string, whiteRating float64, whiteRD float64,
	black string, blackRating float64, blackRD float64) {

	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

	//check if database connection is open
	if db.Ping() != nil {
		log.Println("DATABASE DOWN!")
		return
	}

	//setting verify to yes and deleting row from activate table
	stmt, err := db.Prepare("UPDATE rating SET " + gameType + "=?," + gameType +
		"RD=?" + " where username=?")
	defer stmt.Close()
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
	stmt, err = db.Prepare("UPDATE rating SET " + gameType + "=?," + gameType +
		"RD=?" + " where username=?")
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

	problems, _ := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

	//check if database connection is open
	if db.Ping() != nil {
		log.Println("DATABASE DOWN!")
		return
	}

	//preparing token activation
	stmt, err := db.Prepare("INSERT games SET white=?, black=?, gametype=?, rated=?, " +
		"whiterating=?, blackrating=?, timecontrol=?, moves=?, totalmoves=?, " +
		"result=?, status=?, date=?, time=?, countrywhite=?, countryblack=?")
	defer stmt.Close()
	if err != nil {
		log.Println(err)
		return
	}
	date := time.Now()
	res, err := stmt.Exec(game.WhitePlayer, game.BlackPlayer, game.GameType, game.Rated,
		game.WhiteRating, game.BlackRating, game.TimeControl, moves, totalMoves,
		game.Result, game.Status, date, date, game.CountryWhite, game.CountryBlack)
	if err != nil {
		log.Println(err)
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
	log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

	rows, err := db.Query("SELECT * FROM games WHERE white=? or black=?", name, name)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var all GoGame

	for rows.Next() {

		err = rows.Scan(&all.ID, &all.White, &all.Black, &all.GameType, &all.Rated,
			&all.WhiteRating, &all.BlackRating, &all.TimeControl, &all.Moves,
			&all.Total, &all.Result, &all.Status, &all.Date, &all.Time, &all.CountryWhite, &all.CountryBlack)
		if err != nil {
			log.Println(err)
		}
		storage = append(storage, all)
	}
	return storage
}

func GetSaved(name string) (storage []GoGame) {
	problems, err := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

	rows, err := db.Query("SELECT * FROM saved WHERE white=? or black=?", name, name)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var all GoGame

	for rows.Next() {

		err = rows.Scan(&all.ID, &all.White, &all.Black, &all.GameType, &all.Rated,
			&all.WhiteRating, &all.BlackRating, &all.BlackMinutes, &all.BlackSeconds,
			&all.WhiteMinutes, &all.WhiteSeconds, &all.TimeControl, &all.Moves, &all.Total,
			&all.Status, &all.Date, &all.Time, &all.CountryWhite, &all.CountryBlack)
		if err != nil {
			log.Println(err)
		}
		storage = append(storage, all)
	}
	return storage
}

//fetches saved/adjourned game from database
func (game *ChessGame) fetchSavedGame(id string, user string) bool {

	problems, err := os.OpenFile("logs/errors.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer problems.Close()
	log := log.New(problems, "", log.LstdFlags|log.Lshortfile)

	//these are in the database but not part of the ChessGame struct
	var moves string
	var totalmoves int

	//game status options are:
	//White for white to move or Black for black to move, white won, black won, stalemate or draw.

	err = db.QueryRow("SELECT white, black, gametype, rated, whiterating, "+
		"blackrating, blackminutes, blackseconds, whiteminutes, whiteseconds, "+
		"timecontrol, moves, totalmoves, status, countrywhite, countryblack FROM saved WHERE id=?", id).Scan(&game.WhitePlayer,
		&game.BlackPlayer, &game.GameType, &game.Rated, &game.WhiteRating, &game.BlackRating, &game.BlackMinutes,
		&game.BlackSeconds, &game.WhiteMinutes, &game.WhiteSeconds, &game.TimeControl, &moves,
		&totalmoves, &game.Status, &game.CountryWhite, &game.CountryBlack)
	if err != nil {
		log.Println(err)
	}

	game.Type = "chess_game"
	var holder []Move

	storage := []byte(moves)
	err = json.Unmarshal(storage, &holder)
	if err != nil {
		log.Println(err)
	}
	game.GameMoves = holder
	game.PendingDraw = false

	var start int = 0
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
	All.Games[start] = game

	//intitalizes all the variables of the game
	initGame(game.ID, game.WhitePlayer, game.BlackPlayer)

	var result bool

	for i := 0; i < len(game.GameMoves); i++ {
		result = chessVerify(game.GameMoves[i].S, game.GameMoves[i].T, game.GameMoves[i].P, game.ID)
		if result == false {
			log.Println("something went wrong in move validation in fetchSavedGame of saved game id ", game.ID)
			//undo all game setup and break out
			delete(Verify.AllTables, game.ID)
			delete(All.Games, game.ID)
			return false
		}
	}
	PrivateChat[game.WhitePlayer] = game.BlackPlayer
	PrivateChat[game.BlackPlayer] = game.WhitePlayer

	//starting white's clock first, this goroutine will keep track of both players clock for this game
	table := Verify.AllTables[game.ID]
	go table.startClock(game.ID, game.WhiteMinutes, game.WhiteSeconds, user)

	//delete saved game from database now that its in memory
	stmt, err := db.Prepare("DELETE FROM saved where id=?")
	defer stmt.Close()
	if err != nil {
		log.Println(err)
		return false
	}

	res, err := stmt.Exec(id)
	if err != nil {
		log.Println(err)
		return false
	}
	affect, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return false
	}
	log.Printf("%d row was deleted from the saved table by user %s\n", affect, user)
	return true
}

// gets rating history of player based on type, returns JSON string of ratings with their date time
// returns false if there was an error
func getRatingHistory(name string, gametype string) (string, bool) {

	var ratingHistory string

	err := db.QueryRow("SELECT "+gametype+" FROM ratinghistory WHERE username=?", name).Scan(&ratingHistory)
	if err == sql.ErrNoRows { // this will occur if there is no name exist
		log.Println("No name found in ratinghistory for ", name)
		return "", false
	} else if ratingHistory == "" { // Then there is no history but there is a name
		return "", false
	} else if err != nil { //
		log.Println(err)
		return "", false
	}
	return ratingHistory, true
}

// returns true if username already exists, this function assumes database is already pinged
func CheckUserNameInDb(username string) bool {

	var name string
	//checking if name exists
	_ = db.QueryRow("SELECT username FROM userinfo WHERE username=?", username).Scan(&name)
	if username == name { // already exists
		return true
	} else {
		return false
	}
}

/*
// execute line by line in sql script, does not work for stored procedures
func executeSqlScript(filepath string) {
	file, err := ioutil.ReadAll("/path/to/file")

	if err != nil {
		// handle error
	}

	requests := strings.Split(string(file), ";")

	for _, request := range requests {
		result, err := db.Exec(request)
		// do whatever you need with result and error
	}
}
*/
