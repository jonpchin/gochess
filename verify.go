package gostuff

import (
	"fmt"
	"strings"
)

//returns true if the move is valid otherwise it returns false
func chessVerify(source string, target string, promotion string, gameID int) bool {
	if len(source) != 2 {
		fmt.Println("Invalid input length")
		return false
	}
	var sourceCol = source[0]
	var sourceRow = source[1]

	//if a white piece is picked and its blacks turn or if a black piece is picked and its whites turn or if no piece is picked up return false

	//checking length to ensure no index out of range error
	if len(target) != 2 {
		fmt.Println("Invalid input length")
		return false
	}
	var targetCol = target[0]
	var targetRow = target[1]

	//converting to proper format
	var newSourceCol int8
	var newSourceRow int8
	newSourceCol = convertLetter(sourceCol)
	newSourceRow = 8 - int8(sourceRow-'0')

	var newTargetCol int8
	var newTargetRow int8
	newTargetCol = convertLetter(targetCol)
	newTargetRow = 8 - int8(targetRow-'0')

	//ensuring a digit is entered into the ChessBoard array to prevent index out of range
	if newSourceRow < 0 || newSourceRow > 7 || newSourceCol < 0 || newSourceCol > 7 || newTargetRow < 0 || newTargetRow > 7 || newTargetCol < 0 || newTargetCol > 7 {
		fmt.Println("Invalid input")
		return false
	}

	// make sure Verify.AllTables is inititalized before proceeding
	if _, ok := Verify.AllTables[gameID]; !ok {
		fmt.Println("Please call initGame() function before proceeding")
		return false
	}
	//identifying the piece that was selected
	piece := Verify.AllTables[gameID].ChessBoard[newSourceRow][newSourceCol]

	//no piece was selected
	if piece[0:1] == "-" {
		return false
	}
	//piece without color specification
	noColorPiece := piece[1:2]
	colorOnly := piece[0:1]

	if (Verify.AllTables[gameID].whiteTurn == true && colorOnly == "b") || (Verify.AllTables[gameID].whiteTurn == false && colorOnly == "w") || source == "-" {
		return false
	}
	//checking to make sure player doesn't capture his own pieces
	targetSquare := Verify.AllTables[gameID].ChessBoard[newTargetRow][newTargetCol]
	targetColor := targetSquare[0:1]

	if colorOnly == targetColor {
		fmt.Println("You can't capture your own piece.")
		return false
	}

	//verifying the piece move
	switch noColorPiece {
	case "P":
		if piece == "wP" {
			result := whitePawnMove(newSourceRow, newSourceCol, newTargetRow, newTargetCol, gameID)
			if result == false {
				return false
			}
			//then a succesful pawn move is made
			Verify.AllTables[gameID].pawnMove = (Verify.AllTables[gameID].moveCount + 1) / 2

		} else {

			result := blackPawnMove(newSourceRow, newSourceCol, newTargetRow, newTargetCol, gameID)
			if result == false {
				return false
			}
			Verify.AllTables[gameID].pawnMove = (Verify.AllTables[gameID].moveCount + 1) / 2
		}

	case "N":
		result := knightMove(newSourceRow, newSourceCol, newTargetRow, newTargetCol)
		if result == false {
			return false
		}

	case "B":
		result := bishopMove(newSourceRow, newSourceCol, newTargetRow, newTargetCol, gameID)
		if result == false {
			return false
		}

	case "Q":
		result := queenMove(newSourceRow, newSourceCol, newTargetRow, newTargetCol, gameID)
		if result == false {
			return false
		}

	case "R":
		result := rookMove(newSourceRow, newSourceCol, newTargetRow, newTargetCol, gameID)
		if result == false {
			return false
		} else {
			Verify.AllTables[gameID].rookUpdate = true //used to indicate if a rook has moved, used for castling rights
		}

	case "K":
		result := kingMove(newSourceRow, newSourceCol, newTargetRow, newTargetCol, gameID)
		if result == false {
			return false
		} else { //if its valid king move then update coordinates in the global variables which keeps track of kings location
			Verify.AllTables[gameID].kingUpdate = true
		}

	default:
		fmt.Println(noColorPiece)
		fmt.Println("Error not valid piece")
		return false

	}
	// this promotion check is only used when checking with grand master games
	/*
		if promotion != "" {
			// if promotion is in ASII such as 113 is q then perform ascii conversion
			convertedPromote, err := strconv.Atoi(promotion)
			if err != nil {
				fmt.Println("verify promotion conversion", err)
			}
			switch convertedPromote {
			case 113:
				Verify.AllTables[gameID].promotion = "q"
			case 114:
				Verify.AllTables[gameID].promotion = "r"
			case 110:
				Verify.AllTables[gameID].promotion = "n"
			case 98:
				Verify.AllTables[gameID].promotion = "b"
			default:
				Verify.AllTables[gameID].promotion = promotion
			}
		}
	*/
	capturedPiece := makeMove(newSourceRow, newSourceCol, newTargetRow, newTargetCol, piece, gameID)

	//if a piece is captured within 50 moves then 50 move rule effect is canceled
	if capturedPiece != "-" {
		Verify.AllTables[gameID].lastCapture = (Verify.AllTables[gameID].pawnMove + 1) / 2
	}
	//used to update king position if they are in check
	if Verify.AllTables[gameID].kingUpdate == true {
		if colorOnly == "b" {
			//storing old coordinates
			Verify.AllTables[gameID].blackOldX = newSourceRow
			Verify.AllTables[gameID].blackOldY = newSourceCol

			Verify.AllTables[gameID].blackKingX = newTargetRow
			Verify.AllTables[gameID].blackKingY = newTargetCol

		} else if colorOnly == "w" {
			//storing old coordinates
			Verify.AllTables[gameID].whiteOldX = newSourceRow
			Verify.AllTables[gameID].whiteOldY = newSourceCol

			Verify.AllTables[gameID].whiteKingX = newTargetRow
			Verify.AllTables[gameID].whiteKingY = newTargetCol

		} else {
			fmt.Println("Invalid king color")
		}
	}
	//if the player made a move and his king can be captured that move has to be undone and return false as he didn't stop the check
	if Verify.AllTables[gameID].whiteTurn == true && isWhiteInCheck(gameID) == true {
		undoMove(newSourceRow, newSourceCol, newTargetRow, newTargetCol, piece, capturedPiece, gameID)
		fmt.Println("White cannot make that move as they are in check")
		return false
	} else if Verify.AllTables[gameID].whiteTurn == false && isBlackInCheck(gameID) == true {
		undoMove(newSourceRow, newSourceCol, newTargetRow, newTargetCol, piece, capturedPiece, gameID)
		fmt.Println("Black cannot make that move as they are in check")
		return false
	}
	if Verify.AllTables[gameID].kingUpdate == true { //updating new location of king for quick access
		if colorOnly == "b" {
			Verify.AllTables[gameID].bkMoved = true //can no longer castle with black king
		} else if colorOnly == "w" {

			Verify.AllTables[gameID].wkMoved = true //can no longer castle with white king
		} else {
			fmt.Println("Invalid king color")
		}
		Verify.AllTables[gameID].kingUpdate = false
	}

	if Verify.AllTables[gameID].rookUpdate == true {
		if piece == "bR" && newSourceRow == 0 && newSourceCol == 0 { //black queen rook
			Verify.AllTables[gameID].bqrMoved = true
		} else if piece == "bR" && newSourceRow == 0 && newSourceCol == 7 { //black king rook
			Verify.AllTables[gameID].bkrMoved = true
		} else if piece == "wR" && newSourceRow == 7 && newSourceCol == 0 { //white queen rook
			Verify.AllTables[gameID].wqrMoved = true
		} else if piece == "wR" && newSourceRow == 7 && newSourceCol == 7 { //white king rook move
			Verify.AllTables[gameID].wkrMoved = true
		}
		Verify.AllTables[gameID].rookUpdate = false
	}
	Verify.AllTables[gameID].undoWPass = false //no longer need to watch out for undo en passent
	Verify.AllTables[gameID].undoBPass = false
	Verify.AllTables[gameID].pawnMove++
	switchTurns(gameID)
	return true
}

//changes chess letter notation to a number a=0, b=1, c=2, etc
func convertLetter(letter byte) int8 {
	switch letter {
	case 'a': //this is a file on chess board
		return 0
	case 'b':
		return 1
	case 'c':
		return 2
	case 'd':
		return 3
	case 'e':
		return 4
	case 'f':
		return 5
	case 'g':
		return 6
	case 'h':
		return 7
	default:
		fmt.Println("Invalid file on chess board")
		return -1
	}

}

//makes the chess move on the board, verify the move first, returns captured piece as a string to be used in case of a move undo
func makeMove(sourceRow int8, sourceCol int8, targetRow int8, targetCol int8, piece string, gameID int) string {

	capturedPiece := Verify.AllTables[gameID].ChessBoard[targetRow][targetCol]
	//make the source square blank as now the piece is no longer there
	Verify.AllTables[gameID].ChessBoard[sourceRow][sourceCol] = "-"

	if targetRow == 0 && piece == "wP" {
		Verify.AllTables[gameID].ChessBoard[targetRow][targetCol] = "w" + strings.ToUpper(Verify.AllTables[gameID].promotion)
	} else if targetRow == 7 && piece == "bP" {
		Verify.AllTables[gameID].ChessBoard[targetRow][targetCol] = "b" + strings.ToUpper(Verify.AllTables[gameID].promotion)
	} else {
		Verify.AllTables[gameID].ChessBoard[targetRow][targetCol] = piece //place the piece to its new target square
	}
	return capturedPiece

}

func switchTurns(gameID int) {
	if Verify.AllTables[gameID].whiteTurn == true {
		Verify.AllTables[gameID].whiteTurn = false
	} else {
		Verify.AllTables[gameID].whiteTurn = true
	}
}

func printBoard(gameID int) {

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			fmt.Printf("%s ", Verify.AllTables[gameID].ChessBoard[i][j])
		}
		fmt.Printf("\n")
	}
}

//undo a move if a player makes a move and doesn't stop check
func undoMove(sourceRow int8, sourceCol int8, targetRow int8, targetCol int8, piece string, capturedPiece string, gameID int) {
	Verify.AllTables[gameID].ChessBoard[sourceRow][sourceCol] = piece
	Verify.AllTables[gameID].ChessBoard[targetRow][targetCol] = capturedPiece

	//revert back to original location coordinates for the king
	if Verify.AllTables[gameID].kingUpdate == true {

		if Verify.AllTables[gameID].whiteTurn == true {
			Verify.AllTables[gameID].whiteKingX = Verify.AllTables[gameID].whiteOldX
			Verify.AllTables[gameID].whiteKingY = Verify.AllTables[gameID].whiteOldY
		} else {
			Verify.AllTables[gameID].blackKingX = Verify.AllTables[gameID].blackOldX
			Verify.AllTables[gameID].blackKingY = Verify.AllTables[gameID].blackOldY
		}

		Verify.AllTables[gameID].kingUpdate = false
	}
	//checking if there is an enpassent that needs to be undone
	if Verify.AllTables[gameID].undoWPass == true {
		Verify.AllTables[gameID].ChessBoard[targetRow+1][targetCol] = "bP" //placing back the pawn
	}
	if Verify.AllTables[gameID].undoBPass == true {
		Verify.AllTables[gameID].ChessBoard[targetRow-1][targetCol] = "wP"
	}
}

//checks if the time choices are valid
func checkTime(choice int) bool {
	timeChoices := []int{1, 2, 3, 4, 5, 10, 15, 20, 30, 45}

	var v int

	for _, v = range timeChoices {
		if choice == v {
			return true
		}
	}
	return false
}
