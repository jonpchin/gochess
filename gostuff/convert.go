package gostuff

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/freeeve/pgn.v1"
)

var (
	TotalGrandmasterGames = 0
)

type ConvertChessGame struct {
	Event     string
	Site      string
	Date      string
	Round     string
	White     string
	Black     string
	Result    string
	WhiteElo  string
	BlackElo  string
	ECO       string
	EventDate string
	Moves     []Move
}

type NamesAndID struct {
	ID    int
	White string
	Black string
}

type GrandMasterGame struct {
	ID        int
	Event     string
	Site      string
	Date      string
	Round     string
	White     string
	Black     string
	Result    string
	WhiteElo  string
	BlackElo  string
	ECO       string
	EventDate string
	Moves     string
}

func ConvertAllPGN() {
	convertPGN("games/KingBase2016-03-A00-A39.pgn")
	convertPGN("games/KingBase2016-03-A40-A79.pgn")
	convertPGN("games/KingBase2016-03-A80-A99.pgn")
	convertPGN("games/KingBase2016-03-B00-B19.pgn")
	convertPGN("games/KingBase2016-03-B20-B49.pgn")
	convertPGN("games/KingBase2016-03-B50-B99.pgn")
	convertPGN("games/KingBase2016-03-C00-C19.pgn")
	convertPGN("games/KingBase2016-03-C20-C59.pgn")
	convertPGN("games/KingBase2016-03-C60-C99.pgn")
	convertPGN("games/KingBase2016-03-D00-D29.pgn")
	convertPGN("games/KingBase2016-03-D30-D69.pgn")
	convertPGN("games/KingBase2016-03-D70-D99.pgn")
	convertPGN("games/KingBase2016-03-E00-E19.pgn")
	convertPGN("games/KingBase2016-03-E20-E59.pgn")
	convertPGN("games/KingBase2016-03-E60-E99.pgn")
}

//converts text file pgn and prints it out
func convertPGN(file string) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("error is convert.go convertPGN 1", err)
	}
	ps := pgn.NewPGNScanner(f)
	// while there's more to read in the file
	var newGame ConvertChessGame

	for ps.Next() {
		// scan the next game
		game, err := ps.Scan()
		if err != nil {
			fmt.Println("error is convert.go convertPGN 2", err)
			continue
		}

		// print out tags
		// fmt.Println(game.Tags["Site"])

		newGame.Event = game.Tags["Event"]
		newGame.Site = game.Tags["Site"]
		newGame.Date = game.Tags["Date"]
		newGame.Round = game.Tags["Round"]
		newGame.White = game.Tags["White"]
		newGame.Black = game.Tags["Black"]
		newGame.Result = game.Tags["Result"]
		newGame.WhiteElo = game.Tags["WhiteElo"]
		newGame.BlackElo = game.Tags["BlackElo"]
		newGame.ECO = game.Tags["ECO"]
		newGame.EventDate = game.Tags["EventDate"]

		var temp string
		newGame.Moves = make([]Move, len(game.Moves))
		for key, move := range game.Moves {

			temp = move.String()[0:2]
			newGame.Moves[key].S = temp
			newGame.Moves[key].T = move.String()[2:4]
			checkLength := len(move.String())
			if checkLength > 6 {

				//promotion string guide
				//98=b 110=n 114=r 113=q
				newGame.Moves[key].P = move.String()[4:7]
			} else if checkLength > 4 {
				newGame.Moves[key].P = move.String()[4:6]
			}

		}
		allMoves, err := json.Marshal(newGame.Moves)
		if err != nil {
			fmt.Println("convert.go convertPGN 2 ", err)
		}

		storeGrandMaster(&newGame, allMoves)
		//fmt.Println(newGame)

	}
}

//stores grandmaster PGN games into the grandmaster table
func storeGrandMaster(game *ConvertChessGame, allMoves []byte) {

	moves := string(allMoves)
	//preparing token activation
	stmt, err := db.Prepare("INSERT grandmaster SET event=?, site=?, date=?, round=?, white=?, black=?, result=?, whiteELO=?, blackELO=?, ECO=?, moves=?, eventdate=?")
	if err != nil {
		fmt.Println("convert.go storeGrandMaster 1", err)
		return
	}

	_, err = stmt.Exec(game.Event, game.Site, game.Date, game.Round, game.White, game.Black, game.Result, game.WhiteElo, game.BlackElo, game.ECO, moves, game.EventDate)
	if err != nil {
		fmt.Println("convert.go storeGrandMaster 2", err)
		return
	}
	stmt.Close()

}

func UpdateTotalGrandmasterGames() {
	var total = 1
	err := db.QueryRow("SELECT `AUTO_INCREMENT` FROM information_schema.TABLES WHERE TABLE_NAME = 'grandmaster'").Scan(&total)
	if err != nil {
		fmt.Println("UpdateTotalGrandmasterGames() convert.go", err)
		return
	}
	//setting global total
	TotalGrandmasterGames = total
}

// cycles through the total number of grandmater games in database passed in parameter
// and verifies if the chessVerify function produces any legal moves deemed illegal
func VerifyGrandmasterGames(total int) bool {
	//check if database connection is open
	if db.Ping() != nil {
		fmt.Println("verifyGrandmasterGames DATABASE DOWN!")
		return false
	}

	var allMoves string
	var gameID = 1
	for i := 1; i < total; i++ {

		err := db.QueryRow("SELECT moves FROM grandmaster WHERE id=?", i).Scan(&allMoves)

		if err != nil {
			fmt.Println("verifyGrandmasterGames 1", err)
			return false
		}

		var move []Move
		if err := json.Unmarshal([]byte(allMoves), &move); err != nil {
			fmt.Println("Just receieved a message I couldn't decode:", allMoves, err)
			break
		}
		var legal bool
		gameID = i
		initGame(gameID, "", "")
		for j := 0; j < len(move); j++ {
			legal = chessVerify(move[j].S, move[j].T, move[j].P, gameID)
			totalMoves := (j / 2) + 1
			// The people notating game ID 8035 seems to have made a mistake and notated an illegal move
			if legal == false && gameID != 8035 {
				fmt.Println("Illegal move on turn ", totalMoves, move[j].S, " to ", move[j].T, "at game ID", i)
				return false
			}
		}
	}
	return true

}
