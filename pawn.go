package gostuff

// import "fmt"

//checks if white pawn move is legal, returns true if legal and false if iillegal
func whitePawnMove(sourceRow int8, sourceCol int8, targetRow int8, targetCol int8, gameID int16) bool {
	//moving pawn two squares, pawn should be moving on same file
	if sourceRow-targetRow == 2 && sourceCol == targetCol {

		//if the pawns already moved two squares on their first move then they can't move two squares again
		if Verify.AllTables[gameID].whitePawns[sourceCol] == true {
			//			fmt.Println("You already moved the white pawn two squares.")
			return false
			//checking if any piece blocks the path of the pawn trying to advance two squares
		} else if (sourceRow-1 >= 0 && Verify.AllTables[gameID].ChessBoard[sourceRow-1][sourceCol] != "-") || (sourceRow-2 >= 0 && Verify.AllTables[gameID].ChessBoard[sourceRow-2][sourceCol] != "-") {

			//			fmt.Println("There is a piece blocking the white pawn move.")
			return false
			//enabling enpassent for the other player if there is pawn on either side
		}
		if targetCol-1 >= 0 && Verify.AllTables[gameID].ChessBoard[targetRow][targetCol-1] == "bP" {
			Verify.AllTables[gameID].blackPass[targetCol-1] = true
		} else if targetCol+1 <= 7 && Verify.AllTables[gameID].ChessBoard[targetRow][targetCol+1] == "bP" {
			Verify.AllTables[gameID].blackPass[targetCol+1] = true
		}
		//mark the pawn has moved and can't be moved two squares again
		Verify.AllTables[gameID].whitePawns[sourceCol] = true

		//moving pawn one square or a pawn capture
	} else if sourceRow-targetRow == 1 {

		//determine if its a pawn capture or not, if this is a one square pawn move check if the destination is empty
		if sourceRow-1 >= 0 && sourceCol == targetCol && Verify.AllTables[gameID].ChessBoard[sourceRow-1][sourceCol] == "-" {
			//			fmt.Println("White Pawn moves one square forward.")
			//mark the pawn has moved and can't be moved two squares
			Verify.AllTables[gameID].whitePawns[sourceCol] = true
			//then its a diagonal pawn capture
		} else if (sourceCol-targetCol == 1 || sourceCol-targetCol == -1) && Verify.AllTables[gameID].ChessBoard[targetRow][targetCol] != "-" {
			//			fmt.Println("White pawn captures.")
			//check for enpassent
		} else if Verify.AllTables[gameID].whitePass[sourceCol] == true && Verify.AllTables[gameID].ChessBoard[targetRow][targetCol] == "-" && (sourceCol-targetCol == 1 || targetCol-sourceCol == 1) {
			//remove black pawn left of white pawn
			Verify.AllTables[gameID].ChessBoard[sourceRow][targetCol] = "-"
			//			fmt.Println("removed black pawn via enpassent")
			Verify.AllTables[gameID].undoWPass = true //now this can be undone in undo moves if its an illegal move

			//check enpassent the other side now
		} else {
			//			fmt.Println("Invalid pawn move")
			return false
		}
	} else {
		//		fmt.Println("Invalid pawn move")
		return false
	}
	//player can only enpassent on the first oppurtunity
	passExpireWhite(gameID)
	return true
}

func blackPawnMove(sourceRow int8, sourceCol int8, targetRow int8, targetCol int8, gameID int16) bool {
	//moving pawn two squares, pawn should be moving on same file
	if targetRow-sourceRow == 2 && sourceCol == targetCol {

		//if the pawns already moved two squares on their first move then they can't move two squares again
		if Verify.AllTables[gameID].blackPawns[sourceCol] == true {
			//			fmt.Println("You already moved the black pawn two squares.")
			return false
			//checking if any piece blocks the path of the pawn trying to advance two squares
		} else if (sourceRow+1 <= 7 && Verify.AllTables[gameID].ChessBoard[sourceRow+1][sourceCol] != "-") || (sourceRow+2 <= 7 && Verify.AllTables[gameID].ChessBoard[sourceRow+2][sourceCol] != "-") {

			//			fmt.Println("There is a piece blocking the black pawn move.")
			return false
		}
		//enabling en passent for other player
		if targetCol-1 >= 0 && Verify.AllTables[gameID].ChessBoard[targetRow][targetCol-1] == "wP" {
			Verify.AllTables[gameID].whitePass[targetCol-1] = true
		} else if targetCol+1 <= 7 && Verify.AllTables[gameID].ChessBoard[targetRow][targetCol+1] == "wP" {
			Verify.AllTables[gameID].whitePass[targetCol+1] = true
		}
		//mark the pawn has moved two squares and can't be moved two squares again
		Verify.AllTables[gameID].blackPawns[sourceCol] = true

		//moving pawn one square or a pawn capture
	} else if targetRow-sourceRow == 1 {

		//determine if its a pawn capture or not, if this is a one square pawn move check if the destination is empty
		if sourceRow+1 <= 7 && sourceCol == targetCol && Verify.AllTables[gameID].ChessBoard[sourceRow+1][sourceCol] == "-" {
			//			fmt.Println("Black Pawn moves one square forward.")
			//mark the pawn has moved and can't be moved two squares
			Verify.AllTables[gameID].blackPawns[sourceCol] = true

			//then its a diagonal pawn capture
		} else if (targetCol-sourceCol == 1 || targetCol-sourceCol == -1) && Verify.AllTables[gameID].ChessBoard[targetRow][targetCol] != "-" {
			//			fmt.Println("Black pawn captures.")

		} else if Verify.AllTables[gameID].blackPass[sourceCol] == true && Verify.AllTables[gameID].ChessBoard[targetRow][targetCol] == "-" && (sourceCol-targetCol == 1 || targetCol-sourceCol == 1) {
			//remove black pawn left of white pawn
			Verify.AllTables[gameID].ChessBoard[sourceRow][targetCol] = "-"
			//			fmt.Println("removed white pawn via enpassent")
			Verify.AllTables[gameID].undoBPass = true

		} else {

			//			fmt.Println("Invalid pawn move")
			return false
		}
	} else {
		//		fmt.Println("Invalid pawn move")
		return false
	}

	//player can only enpassent on the first oppurtunity
	passExpireBlack(gameID)
	return true
}

//enPassent expires for the color if they don't make a move
func passExpireWhite(gameID int16) {
	//setting all the values in the map
	for index, _ := range Verify.AllTables[gameID].whitePass {
		Verify.AllTables[gameID].whitePass[index] = false
	}
}

func passExpireBlack(gameID int16) {
	//setting all the values to false
	for index, _ := range Verify.AllTables[gameID].blackPass {
		Verify.AllTables[gameID].blackPass[index] = false
	}
}
