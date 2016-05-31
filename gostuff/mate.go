package gostuff

import "fmt"

//runs through the entire ChessBoard array and searches for black pieces and brute all their possible moves
//to see if it can capture the white king in one move
func isWhiteInCheck(gameID int16) bool {
	result := canBlackKillSquare(Verify.AllTables[gameID].whiteKingX, Verify.AllTables[gameID].whiteKingY, gameID)
	if result == true { //then white's king is in check
		return true
	}
	return false
}

func isBlackInCheck(gameID int16) bool {
	result := canWhiteKillSquare(Verify.AllTables[gameID].blackKingX, Verify.AllTables[gameID].blackKingY, gameID)

	if result == true { //then black's king is in check
		return true
	}
	return false
}

//checks to see if white is in checkmate by bruteforcing all possible white moves and seeing if white is still in check
func isWhiteInMate(gameID int16) bool {

	if isWhiteInCheck(gameID) == false{
		return false
	}

	var i int8
	var j int8
	for i = 0; i < 8; i++ {
		for j = 0; j < 8; j++ {
			color := Verify.AllTables[gameID].ChessBoard[i][j]
			if color[0:1] == "w" {

				switch color[1:2] {
				case "P":

					result := whitePawn(i, j, gameID)
					if result == false { //if there is a pawn move in which king is not in check, then its not mate

						return false
					}
				case "N":

					result := whiteKnight(i, j, gameID)
					if result == false {

						return false
					}
				case "B":

					result := whiteBishop(i, j, gameID)
					if result == false {

						return false
					}
				case "R":

					result := whiteRook(i, j, gameID)
					if result == false {

						return false
					}
				case "Q":

					result := whiteQueen(i, j, gameID)
					if result == false {

						return false
					}
				case "K":

					result := whiteKing(i, j, gameID)
					if result == false {

						return false
					}
				default:
					fmt.Println("Invalid piece type")
				}
			}
		}
	}

	return true
}

func isBlackInMate(gameID int16) bool {
	if isBlackInCheck(gameID) == false{
		return false
	}

	var i int8
	var j int8
	for i = 0; i < 8; i++ {
		for j = 0; j < 8; j++ {
			color := Verify.AllTables[gameID].ChessBoard[i][j]
			if color[0:1] == "b" {

				switch color[1:2] {
				case "P":

					result := blackPawn(i, j, gameID)
					if result == false { //if there is a pawn move in which king is not in check, then its not mate

						return false
					}
				case "N":
					result := blackKnight(i, j, gameID)
					if result == false {

						return false
					}
				case "B":
					result := blackBishop(i, j, gameID)
					if result == false {

						return false
					}
				case "R":
					result := blackRook(i, j, gameID)
					if result == false {

						return false
					}
				case "Q":
					result := blackQueen(i, j, gameID)
					if result == false {

						return false
					}
				case "K":
					result := blackKing(i, j, gameID)
					if result == false {

						return false
					}
				default:
					fmt.Println("Invalid piece type")
				}
			}
		}
	}

	return true
}

func isWhiteStaleMate(gameID int16) bool {
	if isWhiteInCheck(gameID) == true || isWhiteInMate(gameID) == true{
		return false
	}
	return true
}

func isBlackStaleMate(gameID int16) bool {
	if isBlackInCheck(gameID) == true || isBlackInMate(gameID) == true{
		return false
	}
	return true
}

//checks if no material for mating, KvK, K+B vs K, K+B vs K+B, K+N vs K, K+N vs K+N.
func noMaterial(gameID int16) bool{ 
	return true
}
//checks if three reptition rule which leads to a draw. returns false if no three repetition is found
func threeRep(gameID int16) bool{
	if Verify.AllTables[gameID].threeRepS[5] != Verify.AllTables[gameID].threeRepS[1] || Verify.AllTables[gameID].threeRepT[3] != Verify.AllTables[gameID].threeRepS[1]{
		return false
	}
	if Verify.AllTables[gameID].threeRepT[5] != Verify.AllTables[gameID].threeRepT[1] || Verify.AllTables[gameID].threeRepS[3] != Verify.AllTables[gameID].threeRepT[1]{
		return false
	}
	
	if Verify.AllTables[gameID].threeRepS[0] != Verify.AllTables[gameID].threeRepS[4] ||  Verify.AllTables[gameID].threeRepS[4] != Verify.AllTables[gameID].threeRepT[2]{
		return false
	}
	if Verify.AllTables[gameID].threeRepT[0] != Verify.AllTables[gameID].threeRepT[4] ||  Verify.AllTables[gameID].threeRepT[4] != Verify.AllTables[gameID].threeRepS[2]{
		return false
	}
	return true
}
//checks if fifty moves have been made without a pawn push or capture
func fiftyMoves(){
	
}

//checks if a square is attacked by white in one turn, used to identify check and checkmates
func canWhiteKillSquare(targetRow int8, targetCol int8, gameID int16) bool {
	var i int8
	var j int8

	for i = 0; i < 8; i++ {
		for j = 0; j < 8; j++ {
			color := Verify.AllTables[gameID].ChessBoard[i][j]
			if color[0:1] == "w" {
				switch color[1:2] {
				case "P":
					result := whitePawnAttack(i, j, targetRow, targetCol)
					if result == true {
						return true
					}
				case "N":
					result := knightAttack(i, j, targetRow, targetCol)
					if result == true {
						return true
					}
				case "B":
					result := bishopAttack(i, j, targetRow, targetCol, gameID)
					if result == true {
						return true
					}
				case "R":
					result := rookAttack(i, j, targetRow, targetCol, gameID)
					if result == true {
						return true
					}
				case "Q":
					result := queenAttack(i, j, targetRow, targetCol, gameID)
					if result == true {
						return true
					}
				case "K":
					result := kingAttack(i, j, targetRow, targetCol)
					if result == true {
						return true
					}

				default:
					//					fmt.Println("Invalid piece type")
				}
			}
		}
	}

	return false
}

func canBlackKillSquare(targetRow int8, targetCol int8, gameID int16) bool {
	var i int8
	var j int8

	for i = 0; i < 8; i++ {
		for j = 0; j < 8; j++ {
			color := Verify.AllTables[gameID].ChessBoard[i][j]
			if color[0:1] == "b" {
				switch color[1:2] {
				case "P":
					result := blackPawnAttack(i, j, targetRow, targetCol)
					if result == true {
						return true
					}
				case "N":
					result := knightAttack(i, j, targetRow, targetCol)
					if result == true {
						return true
					}
				case "B":
					result := bishopAttack(i, j, targetRow, targetCol, gameID)
					if result == true {
						return true
					}
				case "R":
					result := rookAttack(i, j, targetRow, targetCol, gameID)
					if result == true {
						return true
					}
				case "Q":
					result := queenAttack(i, j, targetRow, targetCol, gameID)
					if result == true {
						return true
					}
				case "K":
					result := kingAttack(i, j, targetRow, targetCol)
					if result == true {
						return true
					}
				default:
					//					fmt.Println("Invalid piece type")
				}
			}
		}
	}

	return false
}
