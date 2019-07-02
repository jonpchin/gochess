package gostuff

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/notnil/chess"
	"golang.org/x/net/websocket"
)

func (c *Connection) LobbyConnect() {

	defer broadCast(c.username) //remove user when they disconnect from socket
	counter := 0
	start := time.Now()

	logFile, _ := os.OpenFile("logs/chat.txt", os.O_APPEND|os.O_WRONLY, 0666)
	defer logFile.Close()
	log := log.New(logFile, "", log.LstdFlags|log.Lshortfile)

	for {
		var reply string

		if err := websocket.Message.Receive(c.websocket, &reply); err != nil {
			//fmt.Println("A user has drop web socket connection ", err)
			break
		}

		var t MessageType
		message := []byte(reply)
		if err := json.Unmarshal(message, &t); err != nil {
			log.Println("Just receieved a message I couldn't decode:", string(reply), err)
			break
		}

		switch t.Type {

		case "chat_all":
			if t.sendLobbyChatToAll(reply, &start, &counter, c.clientIP, log) == false {
				return
			}
		case "fetch_matches":
			//send in array instead of sending individual
			for _, value := range Pending.Matches {

				if value.Opponent == "" {
					if err := websocket.JSON.Send(c.websocket, &value); err != nil {
						log.Println(err)
					}
				} else {
					// Ensure only the recipient or sender sees the private match proposal
					if t.Name == value.Name || t.Name == value.Opponent {
						if err := websocket.JSON.Send(c.websocket, &value); err != nil {
							log.Println(err)
						}
					}
				}
			}

		case "fetch_players":

			//send in array instead of sending individual
			var player MessageType
			player.Type = "fetch_players"
			var uniquePlayers []string

			// show all players in the lobby and those that are playing a game
			for key := range Chat.Lobby {
				player.Name = key
				uniquePlayers = append(uniquePlayers, player.Name)
				if err := websocket.JSON.Send(c.websocket, player); err != nil {
					log.Println(err)
				}
			}
			for key := range Active.Clients {

				player.Name = key
				found := false

				// this will prevent duplicates if player is in lobby and chess room at the same time
				for _, name := range uniquePlayers {
					if player.Name == name {
						found = true
						break
					}
				}
				if found == false {
					if err := websocket.JSON.Send(c.websocket, &player); err != nil {
						log.Println(err)
					}
				}
			}

		case "match_seek":

			var match SeekMatch
			if err := json.Unmarshal(message, &match); err != nil {
				log.Println("Just receieved a message I couldn't decode:", string(reply), err)
				break
			}

			//check if player already has a game started, if there is a game in progress alert player
			if isPlayersInGame(c.username, match.Opponent) {
				t.Type = "alert"
				if err := websocket.JSON.Send(Chat.Lobby[c.username], &t); err != nil {
					// we could not send the message to a peer
					log.Println("Could not send message to ", c.username, err)
				}
				break
			}

			if match.assignMatchRatingType() == false {
				break
			}

			// if a seek matches an existing one do not post another seek
			var started = false

			// If this seek matches an already existing seek in pending matches then start the seek right away
			// only verify if target's criteria matches seeker's crieria as startPendingMatch checks if
			// the seeker's criteria matches the target's criteria
			for matchID, targetMatch := range Pending.Matches {
				if match.TimeControl == targetMatch.TimeControl &&
					targetMatch.Rating >= match.MinRating && targetMatch.Rating <= match.MaxRating &&
					match.Rated == targetMatch.Rated &&
					match.Name != targetMatch.Name { // a player should not be able to play himself

					started = startPendingMatch(match.Name, matchID)
				}
			}

			// Do not send another seek if it was started already
			if started == false && match.isDuplicateMatch() == false {

				//check to make sure player only has a max of three matches seeks pending, used to prevent flood match seeking
				if countMatches(c.username) >= 3 {
					t.Type = "maxThree"
					if err := websocket.JSON.Send(Chat.Lobby[c.username], &t); err != nil {
						// we could not send the message to a peer
						log.Println("Could not send message to ", c.username, err)
					}
					break //notify user that only three matches pending max are allowed
				}

				start := 0
				for {
					if _, ok := Pending.Matches[start]; ok {
						start++
					} else {
						break
					}
				}

				match.MatchID = start
				//used in backend to keep track of all pending games waiting for a player to accept
				Pending.Matches[start] = &match

				go func() {
					for name, cs := range Chat.Lobby {
						if err := websocket.JSON.Send(cs, &match); err != nil {
							// we could not send the message to a peer
							log.Println("Could not send message to ", name, err)
						}
					}
				}()
			}

		case "cancel_match":

			var match SeekMatch
			if err := json.Unmarshal(message, &match); err != nil {
				log.Println("Just receieved a message I couldn't decode:", string(reply), err)
				break
			}

			delete(Pending.Matches, match.MatchID)

			//check if its a private match, if so then delete it and break out
			if match.Opponent != "" {
				fmt.Println("Private match deleted")
				break //no need to continue as this is a private match
			}

			go func() {
				for _, cs := range Chat.Lobby {
					websocket.Message.Send(cs, reply)
				}
			}()

		case "accept_match":

			var match SeekMatch
			if err := json.Unmarshal(message, &match); err != nil {
				log.Println("Just receieved a message I couldn't decode:", string(reply), err)
				break
			}

			//check if player already has a game started, if there is a game in progress alert player
			if isPlayersInGame(match.Name, match.Opponent) {
				log.Println("Player already has a game. ")
				//alerting player
				t.Type = "alert"
				if err := websocket.JSON.Send(Chat.Lobby[c.username], &t); err != nil {
					// we could not send the message to a peer
					log.Println("Could not send message to ", c.username, err)
				}
				break
			}

			startPendingMatch(match.Name, match.MatchID)

		case "private_match":

			var match SeekMatch
			if err := json.Unmarshal(message, &match); err != nil {
				log.Println("Just receieved a message I couldn't decode:", string(reply), err)
				break
			}
			//check if player already has a game started, if there is a game in progress alert player
			if isPlayersInGame(match.Name, match.Opponent) {
				fmt.Println("Player already has a game.")
				//alerting player
				t.Type = "alert"
				if err := websocket.JSON.Send(Chat.Lobby[c.username], &t); err != nil {
					// we could not send the message to a peer
					log.Println("Could not send message to ", c.username, err)
				}
				break
			}

			//check length of name to make sure its 3-12 characters long
			if len(match.Opponent) < 3 || len(match.Opponent) > 12 {
				fmt.Println("Username is too long or too short")
				break
			}
			//a player should not be able to match himself
			if c.username == match.Opponent {
				fmt.Println("You can't match yourself!")
				break
			}

			//check if opponent is in the lobby or not
			if _, ok := Chat.Lobby[match.Opponent]; !ok {
				//alerting player
				t.Type = "absent"
				if err := websocket.JSON.Send(Chat.Lobby[c.username], &t); err != nil {
					// we could not send the message to a peer
					log.Println("Could not send message to ", c.username, err)
				}
				break
			}

			if match.assignMatchRatingType() == false || match.isDuplicateMatch() {
				break
			}

			//check to make sure player only has a max of three matches seeks pending, used to prevent flood match seeking
			if countMatches(c.username) >= 3 {
				t.Type = "maxThree"
				if err := websocket.JSON.Send(Chat.Lobby[c.username], &t); err != nil {
					// we could not send the message to a peer
					log.Println("Could not send message to ", c.username, err)
				}
				break //notify user that only three matches pending max are allowed
			}

			var start int = 0
			for {
				if _, ok := Pending.Matches[start]; ok {
					start++
				} else {
					break
				}
			}

			match.MatchID = start
			//used in backend to keep track of all pending seeks waiting for a player to accept
			Pending.Matches[start] = &match

			go func() {
				for name := range Chat.Lobby {
					if name == match.Opponent || name == match.Name { //send to self and opponent
						if err := websocket.JSON.Send(Chat.Lobby[name], &match); err != nil {
							// we could not send the message to a peer
							log.Println("Could not send message to ", name, err)
						}
					}
				}
			}()

		default:
			log.Println("I'm not familiar with type ", t.Type, " sent by ", c.username)
		}
	}
}

func (match *SeekMatch) isDuplicateMatch() bool {

	for _, targetMatch := range Pending.Matches {
		if match.TimeControl == targetMatch.TimeControl &&
			targetMatch.MinRating == match.MinRating && targetMatch.MaxRating == match.MaxRating &&
			match.Rated == targetMatch.Rated &&
			match.Name == targetMatch.Name {
			return true
		}
	}
	return false
}

// Send chat to all users in the lobby
// Return false if word count and message limit is exceeded
func (t *MessageType) sendLobbyChatToAll(reply string, start *time.Time, counter *int,
	ip string, log *log.Logger) bool {
	if len(reply) > 225 {
		log.Printf("User: %s IP %s has exeeded the 225 character limit by %d byte units.\n",
			t.Name, ip, len(reply))
		return false
	}
	//keeps track of messages are sent in a given interval
	*counter++

	if *counter > 4 {
		elapsed := time.Since(*start)
		if elapsed < time.Second*10 {
			log.Printf("User: %s IP: %s was spamming chat.\n", t.Name, ip)
			return false
		}
		*start = time.Now()
		*counter = 0
	}
	go func() {
		for name, cs := range Chat.Lobby {
			if err := websocket.Message.Send(cs, reply); err != nil {
				// we could not send the message to a peer
				log.Println("Could not send message to ", name, err)
			}
		}
	}()
	return true
}

// if a pending match is accepted start game for both players that are waiting
// seekerName is name of seeker and matchID is the ID that belongs to player
// waiting in pending matches
// if a match cannot be started then return false to indicate the match did not start succesfully
func startPendingMatch(seekerName string, matchID int) bool {

	var game ChessGame

	//checking to make sure player's rating is in range, used as a backend rating check
	errMessage, bullet, blitz, standard, correspondence := GetRating(seekerName)
	if errMessage != "" {
		fmt.Println("Cannot get rating lobby.go startPendingMatch")
		return false
	}

	match := Pending.Matches[matchID]
	if match == nil {
		return false
	}
	//isPlayersInGame function is located in socket.go
	if isPlayersInGame(match.Name, match.Opponent) {
		return false
	}

	if match.Opponent == "" { //only use this case for public matches
		if match.GameType == "bullet" && (bullet < match.MinRating || bullet > match.MaxRating) {
			//fmt.Println("Bullet Rating not in range.")
			return false
		} else if match.GameType == "blitz" && (blitz < match.MinRating || blitz > match.MaxRating) {
			//fmt.Println("Blitz Rating not in range.")
			return false
		} else if match.GameType == "standard" && (standard < match.MinRating || standard > match.MaxRating) {
			//fmt.Println("Standard Rating not in range.")
			return false
		} else if match.GameType == "correspondence" && (correspondence < match.MinRating ||
			correspondence > match.MaxRating) {
			//fmt.Println("Correspondence Rating not in range.")
			return false
		}
	}

	//bullet, blitz, standard or correspondence game type
	game.GameType = match.GameType
	game.Type = "chess_game"

	//seting up the game info such as white/black player, time control, etc
	rand.Seed(time.Now().UnixNano())

	//randomly selects both players to be white or black
	if rand.Intn(2) == 0 {
		game.WhitePlayer = seekerName
		if game.GameType == "bullet" {
			game.WhiteRating = bullet

		} else if game.GameType == "blitz" {
			game.WhiteRating = blitz

		} else {
			game.WhiteRating = standard
		}

		game.BlackRating = match.Rating
		game.BlackPlayer = match.Name

	} else {
		game.WhitePlayer = match.Name
		if game.GameType == "bullet" {
			game.BlackRating = bullet

		} else if game.GameType == "blitz" {
			game.BlackRating = blitz
		} else {
			game.BlackRating = standard
		}

		game.WhiteRating = match.Rating
		game.BlackPlayer = seekerName
	}
	//White for white to move or Black for black to move, white won, black won, stalemate or draw.
	game.Status = "White"

	//no moves yet so nill/null
	game.GameMoves = nil
	game.StartMinutes = match.TimeControl

	game.TimeControl = match.TimeControl
	//for simplicity we will only allow minutes
	game.WhiteMinutes = match.TimeControl
	game.WhiteSeconds = 0
	game.BlackMinutes = match.TimeControl
	game.BlackSeconds = 0
	game.PendingDraw = false
	game.Rated = match.Rated
	game.Spectate = false
	game.CountryWhite = GetCountry(game.WhitePlayer)
	game.CountryBlack = GetCountry(game.BlackPlayer)

	// Long AlgebraicNotation Notation
	game.Validator = chess.NewGame(chess.UseNotation(chess.LongAlgebraicNotation{}))

	// Guests should always be unrated games
	if strings.Contains(game.WhitePlayer, "guest") || strings.Contains(game.BlackPlayer, "guest") {
		game.Rated = "No"
	}

	var start int = 0
	for {
		if _, ok := All.Games[start]; ok {
			start++
		} else {
			break
		}
	}

	game.ID = start
	//used in backend to keep track of all pending games waiting for a player to accept
	All.Games[start] = &game

	//no longer need all the pending matches as game will be started
	for key, value := range Pending.Matches {
		//deletes all pending matches for either players
		if value.Name == game.WhitePlayer || value.Name == game.BlackPlayer {
			delete(Pending.Matches, key)
		}
	}

	//sending to front end for url redirection
	var acceptmatch AcceptMatch
	acceptmatch.Type = "accept_match"
	acceptmatch.Name = game.WhitePlayer
	acceptmatch.TargetPlayer = game.BlackPlayer

	//setting up the private chat between two players and send move connection
	PrivateChat[acceptmatch.Name] = acceptmatch.TargetPlayer
	PrivateChat[acceptmatch.TargetPlayer] = acceptmatch.Name

	//intitalizes all the variables of the game
	InitGame(game.ID, acceptmatch.Name, acceptmatch.TargetPlayer)

	// Redirects when players in the lobby
	for _, cs := range Chat.Lobby {
		if err := websocket.JSON.Send(cs, &acceptmatch); err != nil {
			fmt.Println(err)
		}
	}

	// Redirects for players in the game room
	for _, name := range Verify.AllTables[game.ID].observe.Names {
		if client, ok := Active.Clients[name]; ok {
			if err := websocket.JSON.Send(client, &game); err != nil {
				fmt.Println(err)
			}
		}
	}

	//starting white's clock first, this goroutine will keep track of both players clock for this game
	// the name of person passed in does not matter as long as its one of the two players
	table := Verify.AllTables[game.ID]
	go table.StartClock(game.ID, game.WhiteMinutes, game.WhiteSeconds, game.WhitePlayer)

	return true
}

// Assigns a match as bullet, blitz, standard or correspondence type,
// Also checks if time control is a valid option returns false
// if there was an error
func (match *SeekMatch) assignMatchRatingType() bool {
	//verify.go
	if checkTime(match.TimeControl) == false {
		fmt.Println("An invalid time control has been selected.")
		return false
	}

	//fetching rating from back end
	errRate, bullet, blitz, standard, correspondence := GetRating(match.Name)
	if errRate != "" {
		fmt.Println("Cannot get rating lobby.go match_seek")
		return false
	}

	switch match.TimeControl {
	case 1:
		match.Rating = bullet
		match.GameType = "bullet"
	case 2:
		match.Rating = bullet
		match.GameType = "bullet"
	case 3:
		match.Rating = blitz
		match.GameType = "blitz"
	case 4:
		match.Rating = blitz
		match.GameType = "blitz"
	case 5:
		match.Rating = blitz
		match.GameType = "blitz"
	case 10:
		match.Rating = blitz
		match.GameType = "blitz"
	case 15:
		match.Rating = standard
		match.GameType = "standard"
	case 20:
		match.Rating = standard
		match.GameType = "standard"
	case 30:
		match.Rating = standard
		match.GameType = "standard"
	case 45:
		match.Rating = standard
		match.GameType = "standard"
	default: //for 1440, 2880, 4320 or 5760 minute game defaults to correspondence
		match.Rating = correspondence
		match.GameType = "correspondence"
	}
	return true
}
