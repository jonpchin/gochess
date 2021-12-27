package gostuff

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jonpchin/chess/engine/uci"
	"github.com/notnil/chess"
)

func StartStockfishBot() {

	log.Println("starting stockfish bot")
	const addr = "goplaychess.com"

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	lobbyUrl := url.URL{Scheme: "wss", Host: addr, Path: "/server"}
	//log.Printf("connecting to %s", lobbyUrl.String())
	username, password := enterGuestForBot()

	lobbyConnection, _, err := websocket.DefaultDialer.Dial(lobbyUrl.String(), http.Header{"Cookie": []string{"username=" + username + "; sessionID=" + password}})
	if err != nil {
		log.Fatal("lobby dial:", err)
	}
	defer lobbyConnection.Close()

	chessRoomUrl := url.URL{Scheme: "wss", Host: addr, Path: "/chess"}
	//log.Printf("connecting to %s", chessRoomUrl.String())

	chessRoomConnection, _, err := websocket.DefaultDialer.Dial(chessRoomUrl.String(), http.Header{"Cookie": []string{"username=" + username + "; sessionID=" + password}})
	if err != nil {
		log.Fatal("chess room dial:", err)
	}
	defer chessRoomConnection.Close()

	doneLobby := make(chan struct{})
	doneChessroom := make(chan struct{})

	chessroomHandler := make(chan []byte)
	lobbyHandler := make(chan []byte)
	var engine *uci.Engine

	timeControl := 5

	go func() {
		defer close(doneLobby)

		sendDefaultMatch(username, lobbyHandler)

		for {
			_, message, err := lobbyConnection.ReadMessage()
			if err != nil {
				log.Println("read lobby :", err)
				return
			}

			//log.Printf("stockfish lobby recv: %s", message)

			result := string(message)
			var t MessageType

			if err := json.Unmarshal(message, &t); err != nil {
				log.Println("Just receieved a message I couldn't decode 1:", result, err)
				break
			}

			switch t.Type {
			case "private_match":
				var match SeekMatch
				if err := json.Unmarshal(message, &match); err != nil {
					log.Println("Just receieved a message I couldn't decode 2:", result, err)
					break
				}
				if match.Opponent == username {
					acceptMatch := struct {
						Type    string
						Name    string
						MatchID int
					}{
						"accept_match",
						username,
						match.MatchID,
					}

					timeControl = match.TimeControl

					acceptMatchResult, err := json.Marshal(acceptMatch)
					if err != nil {
						fmt.Println("Could not unmarshal accept match")
					} else {
						lobbyHandler <- acceptMatchResult
					}
				}
			case "chat_all":
				sendDefaultMatch(username, lobbyHandler)
			default:

			}
		}
	}()

	go func() {
		defer close(doneChessroom)
		for {
			_, message, err := chessRoomConnection.ReadMessage()
			if err != nil {
				log.Println("read chess room:", err)
				return
			}
			//log.Printf("stockfish chessroom recv: %s", message)

			var t MessageType
			result := string(message)

			if err := json.Unmarshal(message, &t); err != nil {
				log.Println("Just receieved a message I couldn't decode 1:", result, err)
				break
			}

			switch t.Type {
			case "chess_game":

				var chessGame ChessGame

				if err := json.Unmarshal(message, &chessGame); err != nil {
					log.Println("Just receieved a message I couldn't decode chessgame:", result, err)
					break
				}

				engine = StartEngine(nil)
				startPosition := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

				// Get latest FEN position
				if len(chessGame.GameMoves) > 0 {
					startPosition = chessGame.GameMoves[len(chessGame.GameMoves)-1].Fen
				}

				fen, err := chess.FEN(startPosition)
				if err != nil {
					fmt.Println("can't start chess position", err)
				}

				game := chess.NewGame(chess.UseNotation(chess.UCINotation{}), fen)

				timeControl = chessGame.TimeControl

				if timeControl > 45 {
					timeControl = 45
				}

				if chessGame.WhitePlayer == username {

					currentFen := game.FEN()
					seconds := rand.Intn(timeControl)
					isOk, bestMove := EngineSearchTimeRaw(currentFen, engine, time.Duration(time.Second*time.Duration(seconds)))

					if isOk {
						err = game.MoveStr(bestMove)
						if err != nil {
							fmt.Println("error with movestr 1", err)
						}

						promotion := ""
						if len(bestMove) > 4 {
							promotion = bestMove[4:5]
						}

						sendMove := struct {
							Type string
							Name string
							ID   int
							S    string
							T    string
							Fen  string
							P    string
						}{
							"send_move",
							username,
							chessGame.ID,
							bestMove[0:2],
							bestMove[2:4],
							game.FEN(),
							promotion,
						}

						sendMoveResult, err := json.Marshal(sendMove)

						if err != nil {
							fmt.Println("chess_game error in marshal", err)
						} else {
							chessroomHandler <- sendMoveResult

						}
					} else {
						fmt.Println("Isokay is false", bestMove)
					}
				}

			case "send_move":

				var gameMove GameMove

				if err := json.Unmarshal(message, &gameMove); err != nil {
					log.Println("Just receieved a message I couldn't decode chessgame:", result, err)
					break
				}

				fen, err := chess.FEN(gameMove.Fen)
				if err != nil {
					fmt.Println("can't start chess position", err)
				}

				game := chess.NewGame(chess.UseNotation(chess.UCINotation{}), fen)
				currentFen := game.FEN()

				if timeControl > 45 {
					timeControl = 45
				}
				seconds := rand.Intn(timeControl)
				isOk, bestMove := EngineSearchTimeRaw(currentFen, engine, time.Duration(time.Second*time.Duration(seconds)))

				if isOk {

					err = game.MoveStr(bestMove)
					if err != nil {
						fmt.Println("error with movestr 2", err)
					}

					promotion := ""
					if len(bestMove) > 4 {
						promotion = bestMove[4:5]
					}

					sendMove := struct {
						Type string
						Name string
						ID   int
						S    string
						T    string
						Fen  string
						P    string
					}{
						"send_move",
						username,
						gameMove.ID,
						bestMove[0:2],
						bestMove[2:4],
						game.FEN(),
						promotion,
					}

					sendMoveResult, err := json.Marshal(sendMove)
					if err != nil {
						fmt.Println("Can't marsal sendMove", sendMove)
					} else {
						chessroomHandler <- sendMoveResult
					}
				} else {
					fmt.Println("Isokay is false", bestMove)
				}

			case "rematch":

				var match SeekMatch

				if err := json.Unmarshal(message, &match); err != nil {
					log.Println("Just receieved a message I couldn't decode rematch:", result, err)
					break
				}

				acceptRematch := struct {
					Type    string
					Name    string
					MatchID int
				}{
					"accept_rematch",
					username,
					match.MatchID,
				}

				timeControl = match.TimeControl

				acceptRematchResult, err := json.Marshal(acceptRematch)
				if err != nil {
					fmt.Println("Can't marsal acceptRematch", acceptRematch)
				} else {
					chessroomHandler <- acceptRematchResult
				}
			case "abort_game":
				engine.Stop()
				sendDefaultMatch(username, lobbyHandler)
			case "resign":
				engine.Stop()
				sendDefaultMatch(username, lobbyHandler)
			case "game_over":
				engine.Stop()
				sendDefaultMatch(username, lobbyHandler)
			default:
			}
		}
	}()
	go func() {
		lobbyGreeting := struct {
			Type string
			Name string
			Text string
		}{
			"chat_all",
			username,
			"has entered the lobby.",
		}

		lobbyGreetingResult, err := json.Marshal(lobbyGreeting)
		if err != nil {
			fmt.Println("Could not unmarshal lobby greeting")
		} else {
			lobbyHandler <- lobbyGreetingResult
		}
	}()

	for {
		select {
		case <-doneLobby:
			return
		case t := <-chessroomHandler:
			err := chessRoomConnection.WriteMessage(websocket.TextMessage, t)
			if err != nil {
				log.Println("chess room write:", err)
				return
			}
		case t := <-lobbyHandler:

			err := lobbyConnection.WriteMessage(websocket.TextMessage, t)
			if err != nil {
				log.Println("lobby write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")
			engine.Quit()

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := lobbyConnection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-doneLobby:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func enterGuestForBot() (string, string) {
	response, err := http.PostForm("https://goplaychess.com/enterGuest", url.Values{})

	//okay, moving on...
	if err != nil {
		fmt.Println("enterGuest 1", err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println("enterGuest 2", err)
	}

	creds := strings.Split(string(body), ",")

	if len(creds) == 2 {
		return creds[0], creds[1]
	} else {
		fmt.Println("creds length is not 2")
	}
	return "", ""
}

func sendDefaultMatch(username string, lobbyHandler chan []byte) {

	var defaultMatch SeekMatch

	defaultMatch.Type = "match_seek"
	defaultMatch.Name = username
	defaultMatch.TimeControl = 5 // bot will send out a 5 minute seek by default
	defaultMatch.MinRating = 500
	defaultMatch.MaxRating = 2700
	defaultMatch.Rated = "No"

	acceptMatchResult, err := json.Marshal(defaultMatch)
	if err != nil {
		fmt.Println("Could not unmarshal accept match")
	} else {
		lobbyHandler <- acceptMatchResult
	}
}
