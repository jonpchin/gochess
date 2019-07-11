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
	engine = nil

	matchID := 0
	timeControl := 15
	rand.Seed(time.Now().Unix())

	go func() {
		defer close(doneLobby)
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
			case "match_seek":
				var match SeekMatch
				if err := json.Unmarshal(message, &match); err != nil {
					log.Println("Just receieved a message I couldn't decode 2:", result, err)
					break
				}

				acceptMatch := struct {
					Type    string
					Name    string
					MatchID int
				}{
					"accept_match",
					username,
					match.MatchID,
				}

				matchID = match.MatchID
				timeControl = match.TimeControl

				acceptMatchResult, err := json.Marshal(acceptMatch)
				if err != nil {
					fmt.Println("Could not unmarshal accept match")
				} else {
					lobbyHandler <- acceptMatchResult
				}

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

				if engine == nil {
					engine = startEngine(nil)
				}

				if chessGame.WhitePlayer == username {

					startPosition := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

					// Get latest FEN position
					if len(chessGame.GameMoves) > 0 {
						startPosition = chessGame.GameMoves[len(chessGame.GameMoves)-1].Fen
					}

					fen, err := chess.FEN(startPosition)
					if err != nil {
						fmt.Println("can't start chess position", err)
					}

					game := chess.NewGame(chess.UseNotation(chess.LongAlgebraicNotation{}), fen)
					currentFen := game.FEN()

					seconds := rand.Intn(timeControl)
					isOk, bestMove := engineSearchTimeRaw(currentFen, engine, time.Duration(time.Second*time.Duration(seconds)))

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
						matchID,
						bestMove[0:2],
						bestMove[2:4],
						game.FEN(),
						promotion,
					}

					sendMoveResult, err := json.Marshal(sendMove)

					if err != nil {
						fmt.Println("chess_game error in marshal", err)
					} else {
						if isOk {
							chessroomHandler <- sendMoveResult

						} else {
							fmt.Println("move is not okay, exiting")
						}
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

				game := chess.NewGame(chess.UseNotation(chess.LongAlgebraicNotation{}), fen)
				currentFen := game.FEN()

				seconds := rand.Intn(timeControl)
				isOk, bestMove := engineSearchTimeRaw(currentFen, engine, time.Duration(time.Second*time.Duration(seconds)))

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
					matchID,
					bestMove[0:2],
					bestMove[2:4],
					game.FEN(),
					promotion,
				}

				sendMoveResult, err := json.Marshal(sendMove)
				if err != nil {
					fmt.Println("Can't marsal sendMove", sendMove)
				} else {
					if isOk {
						chessroomHandler <- sendMoveResult

					} else {
						fmt.Println("send move is not okay, exiting")
					}
				}
			case "rematch":

				var match SeekMatch

				if err := json.Unmarshal(message, &match); err != nil {
					log.Println("Just receieved a message I couldn't decode rematch:", result, err)
					break
				}

				acceptRematch := struct {
					Type        string
					Name        string
					Opponent    string
					TimeControl int
				}{
					"accept_rematch",
					username,
					match.Opponent,
					match.TimeControl,
				}

				acceptRematchResult, err := json.Marshal(acceptRematch)
				if err != nil {
					fmt.Println("Can't marsal acceptRematch", acceptRematch)
				} else {
					chessroomHandler <- acceptRematchResult
				}
			case "abort_game":
				if engine != nil {
					engine.Quit()
					engine = nil
				}
			case "resign":
				if engine != nil {
					engine.Quit()
					engine = nil
				}
			case "game_over":
				if engine != nil {
					engine.Quit()
					engine = nil
				}
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
			if engine != nil {
				engine.Quit()
			}

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
