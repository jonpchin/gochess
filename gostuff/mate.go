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
	
	//used to store piece count pawn=0 knight=1 bishop=2 rook=3 queen=4 king=5
	var white [6]int
	var black [6]int
	
	var i int8
	var j int8
	for i = 0; i < 8; i++ {
		for j = 0; j < 8; j++ {
			color := Verify.AllTables[gameID].ChessBoard[i][j]
			if color[0:1] == "w" {
			
				switch color[1:2] {
				case "P":
					white[0]++;
				case "N":
					white[1]++;
				case "B":
					white[2]++;
				case "R":
					white[3]++;
				case "Q":
					white[4]++;
				case "K":
					white[5]++;
				default:
					fmt.Println("Incorrect piece mate.go no material 1")
				}
			}else if color[0:1] == "b"{
				switch color[1:2] {
				case "P":
					black[0]++;
				case "N":
					black[1]++;
				case "B":
					black[2]++;
				case "R":
					black[3]++;
				case "Q":
					black[4]++;
				case "K":
					black[5]++;
				default:
					fmt.Println("Incorrect piece mate.go no material 2")
				}
			}else{
				fmt.Printf("Invalid color type mate.go no material 3 location of i and j is %d %d", i, j)
			}
		}
	}
	
	//KvK, K+B vs K, K+B vs K+B, K+N vs K, K+N vs K+N.
	// pawn=0 knight=1 bishop=2 rook=3 queen=4 king=5
	//KvK
	if white[0] == 0 && white[1] == 0 && white[2]== 0 && white[3] == 0 && white[4] == 0 && white[5] == 1 && black[0] == 0 && black[1] == 0 && black[2]== 0 && black[3] == 0 && black[4] == 0 && black[5] == 1{ 
		return true
		//K+B vs K
	}else if white[0] == 0 && white[1] == 0 && white[2]== 1 && white[3] == 0 && white[4] == 0 && white[5] == 1 && black[0] == 0 && black[1] == 0 && black[2]== 0 && black[3] == 0 && black[4] == 0 && black[5] == 1{
		return true
		//K vs K+B	
	}else if white[0] == 0 && white[1] == 0 && white[2]== 0 && white[3] == 0 && white[4] == 0 && white[5] == 1 && black[0] == 0 && black[1] == 0 && black[2]== 1 && black[3] == 0 && black[4] == 0 && black[5] == 1{
		return true
		//K+B vs K+B
	}else if white[0] == 0 && white[1] == 0 && white[2]== 1 && white[3] == 0 && white[4] == 0 && white[5] == 1 && black[0] == 0 && black[1] == 0 && black[2]== 1 && black[3] == 0 && black[4] == 0 && black[5] == 1{
		return true
		//K+N vs K
	}else if white[0] == 0 && white[1] == 1 && white[2]== 0 && white[3] == 0 && white[4] == 0 && white[5] == 1 && black[0] == 0 && black[1] == 0 && black[2]== 0 && black[3] == 0 && black[4] == 0 && black[5] == 1{
		return true
		//K vs K+N
	}else if white[0] == 0 && white[1] == 0 && white[2]== 0 && white[3] == 0 && white[4] == 0 && white[5] == 1 && black[0] == 0 && black[1] == 1 && black[2]== 0 && black[3] == 0 && black[4] == 0 && black[5] == 1{
		return true
		//K+N vs K+N
	}else if white[0] == 0 && white[1] == 1 && white[2]== 0 && white[3] == 0 && white[4] == 0 && white[5] == 1 && black[0] == 0 && black[1] == 1 && black[2]== 0 && black[3] == 0 && black[4] == 0 && black[5] == 1{
		return true
	}
	//otherwise insufficient mating material
 	return false
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
func fiftyMoves(gameID int16) bool{
	var thisMove int
	thisMove = (len(All.Games[gameID].GameMoves) + 1) / 2
	//no capture within 50 moves
	if (thisMove - Verify.AllTables[gameID].lastCapture) >= 50{
		return true
		//no pawn move within 50 moves
	}else if (thisMove - Verify.AllTables[gameID].pawnMove) >= 50{
		return true
	}
	return false
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
