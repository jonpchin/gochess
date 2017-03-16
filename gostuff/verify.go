package gostuff

import (
	"fmt"
	"strings"
)

//returns true if the move is valid otherwise it returns false
func ChessVerify(source string, target string, promotion string, gameID int) bool {
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
		fmt.Println("Please call InitGame() function before proceeding")
		return false
	}
	table := Verify.AllTables[gameID]
	//identifying the piece that was selected
	piece := table.ChessBoard[newSourceRow][newSourceCol]

	//no piece was selected
	if piece[0:1] == "-" {
		return false
	}
	//piece without color specification
	noColorPiece := piece[1:2]
	colorOnly := piece[0:1]

	if (table.whiteTurn && colorOnly == "b") || (table.whiteTurn == false && colorOnly == "w") || source == "-" {
		return false
	}
	//checking to make sure player doesn't capture his own pieces
	targetSquare := table.ChessBoard[newTargetRow][newTargetCol]
	targetColor := targetSquare[0:1]

	if colorOnly == targetColor {
		fmt.Println("You can't capture your own piece.")
		return false
	}

	//verifying the piece move
	switch noColorPiece {
	case "P":
		if piece == "wP" {
			result := table.whitePawnMove(newSourceRow, newSourceCol, newTargetRow, newTargetCol)
			if result == false {
				return false
			}
			//then a succesful pawn move is made
			table.pawnMove = (table.moveCount + 1) / 2

		} else {

			result := table.blackPawnMove(newSourceRow, newSourceCol, newTargetRow, newTargetCol)
			if result == false {
				return false
			}
			table.pawnMove = (table.moveCount + 1) / 2
		}

	case "N":
		result := knightMove(newSourceRow, newSourceCol, newTargetRow, newTargetCol)
		if result == false {
			return false
		}

	case "B":
		result := table.bishopMove(newSourceRow, newSourceCol, newTargetRow, newTargetCol)
		if result == false {
			return false
		}

	case "Q":
		result := table.queenMove(newSourceRow, newSourceCol, newTargetRow, newTargetCol)
		if result == false {
			return false
		}

	case "R":
		result := table.rookMove(newSourceRow, newSourceCol, newTargetRow, newTargetCol)
		if result == false {
			return false
		} else {
			table.rookUpdate = true //used to indicate if a rook has moved, used for castling rights
		}

	case "K":
		result := table.kingMove(newSourceRow, newSourceCol, newTargetRow, newTargetCol)
		if result == false {
			return false
		} else { //if its valid king move then update coordinates in the global variables which keeps track of kings location
			table.kingUpdate = true
		}

	default:
		fmt.Println(noColorPiece)
		fmt.Println("Error not valid piece")
		return false

	}
	// Do NOT MOVE, this allows selected promotion piece to be updated in the back end
	table.promotion = promotion

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
				table.promotion = "q"
			case 114:
				table.promotion = "r"
			case 110:
				table.promotion = "n"
			case 98:
				table.promotion = "b"
			default:
				// then do nothing as promotion piece is already set
			}
		}
	*/
	capturedPiece := table.makeMove(newSourceRow, newSourceCol, newTargetRow, newTargetCol, piece)

	//if a piece is captured within 50 moves then 50 move rule effect is canceled
	if capturedPiece != "-" {
		table.lastCapture = (table.pawnMove + 1) / 2
	}
	//used to update king position if they are in check
	if table.kingUpdate {
		if colorOnly == "b" {
			//storing old coordinates
			table.blackOldX = newSourceRow
			table.blackOldY = newSourceCol

			table.blackKingX = newTargetRow
			table.blackKingY = newTargetCol

		} else if colorOnly == "w" {
			//storing old coordinates
			table.whiteOldX = newSourceRow
			table.whiteOldY = newSourceCol

			table.whiteKingX = newTargetRow
			table.whiteKingY = newTargetCol

		} else {
			fmt.Println("Invalid king color")
		}
	}
	//if the player made a move and his king can be captured that move has to be undone and return false as he didn't stop the check
	if table.whiteTurn && table.isWhiteInCheck() {
		table.undoMove(newSourceRow, newSourceCol, newTargetRow, newTargetCol, piece, capturedPiece)
		fmt.Println("White cannot make that move as they are in check")
		return false
	} else if table.whiteTurn == false && table.isBlackInCheck() {
		table.undoMove(newSourceRow, newSourceCol, newTargetRow, newTargetCol, piece, capturedPiece)
		fmt.Println("Black cannot make that move as they are in check")
		return false
	}
	if table.kingUpdate { //updating new location of king for quick access
		if colorOnly == "b" {
			table.bkMoved = true //can no longer castle with black king
		} else if colorOnly == "w" {

			table.wkMoved = true //can no longer castle with white king
		} else {
			fmt.Println("Invalid king color")
		}
		table.kingUpdate = false
	}

	if table.rookUpdate {
		if piece == "bR" && newSourceRow == 0 && newSourceCol == 0 { //black queen rook
			table.bqrMoved = true
		} else if piece == "bR" && newSourceRow == 0 && newSourceCol == 7 { //black king rook
			table.bkrMoved = true
		} else if piece == "wR" && newSourceRow == 7 && newSourceCol == 0 { //white queen rook
			table.wqrMoved = true
		} else if piece == "wR" && newSourceRow == 7 && newSourceCol == 7 { //white king rook move
			table.wkrMoved = true
		}
		table.rookUpdate = false
	}
	table.undoWPass = false //no longer need to watch out for undo en passent
	table.undoBPass = false
	table.pawnMove++
	table.switchTurns()
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

//makes the chess move on the board, verify the move first, returns captured piece
//as a string to be used in case of a move undo
func (table *Table) makeMove(sourceRow int8, sourceCol int8, targetRow int8,
	targetCol int8, piece string) string {

	capturedPiece := table.ChessBoard[targetRow][targetCol]
	//make the source square blank as now the piece is no longer there
	table.ChessBoard[sourceRow][sourceCol] = "-"

	if targetRow == 0 && piece == "wP" {
		table.ChessBoard[targetRow][targetCol] = "w" + strings.ToUpper(table.promotion)
	} else if targetRow == 7 && piece == "bP" {
		table.ChessBoard[targetRow][targetCol] = "b" + strings.ToUpper(table.promotion)
	} else {
		table.ChessBoard[targetRow][targetCol] = piece //place the piece to its new target square
	}
	return capturedPiece

}

func (table *Table) switchTurns() {
	if table.whiteTurn {
		table.whiteTurn = false
	} else {
		table.whiteTurn = true
	}
}

func (table *Table) printBoard() {

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			fmt.Printf("%s ", table.ChessBoard[i][j])
		}
		fmt.Printf("\n")
	}
}

//undo a move if a player makes a move and doesn't stop check
func (table *Table) undoMove(sourceRow int8, sourceCol int8, targetRow int8, targetCol int8,
	piece string, capturedPiece string) {

	table.ChessBoard[sourceRow][sourceCol] = piece
	table.ChessBoard[targetRow][targetCol] = capturedPiece

	//revert back to original location coordinates for the king
	if table.kingUpdate {

		if table.whiteTurn {
			table.whiteKingX = table.whiteOldX
			table.whiteKingY = table.whiteOldY
		} else {
			table.blackKingX = table.blackOldX
			table.blackKingY = table.blackOldY
		}

		table.kingUpdate = false
	}
	//checking if there is an enpassent that needs to be undone
	if table.undoWPass {
		table.ChessBoard[targetRow+1][targetCol] = "bP" //placing back the pawn
	}
	if table.undoBPass {
		table.ChessBoard[targetRow-1][targetCol] = "wP"
	}
}

// checks if the time choices are valid
// gameType can be used to check correspondence times in minutes as well
func checkTime(choice int) bool {

	// 1440, 2880, 4320, 5760 are minutes for correspondence which are 1, 2, 3, 4 days
	var timeChoices = []int{1, 2, 3, 4, 5, 10, 15, 20, 30, 45, 1440, 2880, 4320, 5760}

	for _, v := range timeChoices {
		if choice == v {
			return true
		}
	}
	return false
}
