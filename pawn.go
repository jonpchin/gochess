package gostuff

import "fmt"

// import "fmt"

//checks if white pawn move is legal, returns true if legal and false if iillegal
func (table *Table) whitePawnMove(sourceRow int8, sourceCol int8, targetRow int8, targetCol int8) bool {
	//moving pawn two squares, pawn should be moving on same file
	if sourceRow-targetRow == 2 && sourceCol == targetCol {

		if (sourceRow-1 >= 0 && table.ChessBoard[sourceRow-1][sourceCol] != "-") || (sourceRow-2 >= 0 && table.ChessBoard[sourceRow-2][sourceCol] != "-") {
			fmt.Println("There is a piece blocking the white pawn move.")
			return false
			//enabling enpassent for the other player if there is pawn on either side
		}
		if targetCol-1 >= 0 && table.ChessBoard[targetRow][targetCol-1] == "bP" {
			table.blackPass[targetCol-1] = true
		} else if targetCol+1 <= 7 && table.ChessBoard[targetRow][targetCol+1] == "bP" {
			table.blackPass[targetCol+1] = true
		}

		//moving pawn one square or a pawn capture
	} else if sourceRow-targetRow == 1 {

		//determine if its a pawn capture or not, if this is a one square pawn move check if the destination is empty
		if sourceRow-1 >= 0 && sourceCol == targetCol && table.ChessBoard[sourceRow-1][sourceCol] == "-" {
			//			fmt.Println("White Pawn moves one square forward.")
			//mark the pawn has moved and can't be moved two squares
			table.whitePawns[sourceCol] = true
			//then its a diagonal pawn capture
		} else if (sourceCol-targetCol == 1 || sourceCol-targetCol == -1) && table.ChessBoard[targetRow][targetCol] != "-" {
			//			fmt.Println("White pawn captures.")
			//check for enpassent
		} else if table.ChessBoard[targetRow][targetCol] == "-" && (sourceCol-targetCol == 1 || targetCol-sourceCol == 1) {
			//remove black pawn left of white pawn
			table.ChessBoard[sourceRow][targetCol] = "-"
			//			fmt.Println("removed black pawn via enpassent")
			table.undoWPass = true //now this can be undone in undo moves if its an illegal move

			//check enpassent the other side now
		} else {
			//			fmt.Println("Invalid pawn move white 1")
			return false
		}
	} else {
		//		fmt.Println("Invalid pawn move white 2")
		return false
	}
	//player can only enpassent on the first oppurtunity
	table.passExpireWhite()
	return true
}

func (table *Table) blackPawnMove(sourceRow int8, sourceCol int8, targetRow int8, targetCol int8) bool {
	//moving pawn two squares, pawn should be moving on same file
	if targetRow-sourceRow == 2 && sourceCol == targetCol {

		//checking if any piece blocks the path of the pawn trying to advance two squares
		if (sourceRow+1 <= 7 && table.ChessBoard[sourceRow+1][sourceCol] != "-") || (sourceRow+2 <= 7 && table.ChessBoard[sourceRow+2][sourceCol] != "-") {
			fmt.Println("There is a piece blocking the black pawn move.")
			return false
		}
		//enabling en passent for other player
		if targetCol-1 >= 0 && table.ChessBoard[targetRow][targetCol-1] == "wP" {
			table.whitePass[targetCol-1] = true
		} else if targetCol+1 <= 7 && table.ChessBoard[targetRow][targetCol+1] == "wP" {
			table.whitePass[targetCol+1] = true
		}

		//moving pawn one square or a pawn capture
	} else if targetRow-sourceRow == 1 {

		//determine if its a pawn capture or not, if this is a one square pawn move check if the destination is empty
		if sourceRow+1 <= 7 && sourceCol == targetCol && table.ChessBoard[sourceRow+1][sourceCol] == "-" {
			//			fmt.Println("Black Pawn moves one square forward.")

			//then its a diagonal pawn capture
		} else if (targetCol-sourceCol == 1 || targetCol-sourceCol == -1) && table.ChessBoard[targetRow][targetCol] != "-" {
			//			fmt.Println("Black pawn captures.")

		} else if table.ChessBoard[targetRow][targetCol] == "-" && (sourceCol-targetCol == 1 || targetCol-sourceCol == 1) {
			//remove black pawn left of white pawn
			table.ChessBoard[sourceRow][targetCol] = "-"
			//			fmt.Println("removed white pawn via enpassent")
			table.undoBPass = true

		} else {
			//			fmt.Println("Invalid pawn move black 1")
			return false
		}
	} else {
		//		fmt.Println("Invalid pawn move black 2")
		return false
	}

	//player can only enpassent on the first oppurtunity
	table.passExpireBlack()
	return true
}

//enPassent expires for the color if they don't make a move
func (table *Table) passExpireWhite() {
	//setting all the values in the map
	for index, _ := range table.whitePass {
		table.whitePass[index] = false
	}
}

func (table *Table) passExpireBlack() {
	//setting all the values to false
	for index, _ := range table.blackPass {
		table.blackPass[index] = false
	}
}
