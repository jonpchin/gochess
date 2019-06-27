package gostuff

import "fmt"

//runs through the entire ChessBoard array and searches for black pieces and brute all their possible moves
//to see if it can capture the white king in one move
func (table *Table) isWhiteInCheck() bool {
	result := table.canBlackKillSquare(table.whiteKingX, table.whiteKingY)
	if result == true { //then white's king is in check
		return true
	}
	return false
}

func (table *Table) isBlackInCheck() bool {
	result := table.canWhiteKillSquare(table.blackKingX, table.blackKingY)

	if result == true { //then black's king is in check
		return true
	}
	return false
}

//checks to see if white is in checkmate by bruteforcing all possible white moves and seeing if white is still in check
func (table *Table) isWhiteInMate() bool {

	if table.isWhiteInCheck() == false {
		return false
	}

	var i int8
	var j int8
	for i = 0; i < 8; i++ {
		for j = 0; j < 8; j++ {
			color := table.ChessBoard[i][j]
			if color[0:1] == "w" {

				switch color[1:2] {
				case "P":

					result := table.whitePawn(i, j)
					if result == false { //if there is a pawn move in which king is not in check, then its not mate

						return false
					}
				case "N":

					result := table.whiteKnight(i, j)
					if result == false {

						return false
					}
				case "B":

					result := table.whiteBishop(i, j)
					if result == false {

						return false
					}
				case "R":

					result := table.whiteRook(i, j)
					if result == false {

						return false
					}
				case "Q":

					result := table.whiteQueen(i, j)
					if result == false {

						return false
					}
				case "K":

					result := table.whiteKing(i, j)
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

func (table *Table) isBlackInMate() bool {
	if table.isBlackInCheck() == false {
		return false
	}

	var i int8
	var j int8

	for i = 0; i < 8; i++ {
		for j = 0; j < 8; j++ {
			color := table.ChessBoard[i][j]
			if color[0:1] == "b" {

				switch color[1:2] {
				case "P":
					result := table.blackPawn(i, j)
					if result == false { //if there is a pawn move in which king is not in check, then its not mate

						return false
					}
				case "N":
					result := table.blackKnight(i, j)
					if result == false {

						return false
					}
				case "B":
					result := table.blackBishop(i, j)
					if result == false {

						return false
					}
				case "R":
					result := table.blackRook(i, j)
					if result == false {

						return false
					}
				case "Q":
					result := table.blackQueen(i, j)
					if result == false {

						return false
					}
				case "K":
					result := table.blackKing(i, j)
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

func (table *Table) isWhiteStaleMate() bool {
	if table.isWhiteInCheck() || table.isWhiteInMate() {
		return false
	}
	var i int8
	var j int8
	for i = 0; i < 8; i++ {
		for j = 0; j < 8; j++ {
			piece := table.ChessBoard[i][j]
			if piece[0:1] == "w" {
				switch piece[1:2] {

				case "P":
					result := table.whitePawnStaleMate(i, j)
					if result == false {
						return false
					}

				case "N":
					result := table.whiteKnightStaleMate(i, j)
					if result == false {
						return false
					}
				case "B":
					result := table.whiteBishopStaleMate(i, j)
					if result == false {
						return false
					}
				case "R":
					result := table.whiteRookStaleMate(i, j)
					if result == false {
						return false
					}
				case "Q":
					result := table.whiteQueenStaleMate(i, j)
					if result == false {
						return false
					}
				case "K":
					result := table.whiteKingStaleMate(i, j)
					if result == false {
						return false
					}

				default:
					fmt.Println("Mate.go whiteStalemate not valid piece")
				}
			}

		}
	}
	return true
}

func (table *Table) isBlackStaleMate() bool {
	if table.isBlackInCheck() || table.isBlackInMate() {
		return false
	}
	var i int8
	var j int8
	for i = 0; i < 8; i++ {
		for j = 0; j < 8; j++ {
			piece := table.ChessBoard[i][j]
			if piece[0:1] == "b" {
				switch piece[1:2] {

				case "P":
					result := table.blackPawnStaleMate(i, j)
					if result == false {
						return false
					}

				case "N":
					result := table.blackKnightStaleMate(i, j)
					if result == false {
						return false
					}
				case "B":
					result := table.blackBishopStaleMate(i, j)
					if result == false {
						return false
					}
				case "R":
					result := table.blackRookStaleMate(i, j)
					if result == false {
						return false
					}
				case "Q":
					result := table.blackQueenStaleMate(i, j)
					if result == false {
						return false
					}
				case "K":
					result := table.blackKingStaleMate(i, j)
					if result == false {
						return false
					}

				default:
					fmt.Println("Mate.go blackStalemate not valid piece")
				}
			}

		}
	}
	return true
}

//checks if no material for mating, KvK, K+B vs K, K+B vs K+B, K+N vs K, K+N vs K+N.
func (table *Table) noMaterial() bool {

	//used to store piece count pawn=0 knight=1 bishop=2 rook=3 queen=4 king=5
	var white [6]int
	var black [6]int

	var i int8
	var j int8
	for i = 0; i < 8; i++ {
		for j = 0; j < 8; j++ {
			color := table.ChessBoard[i][j]
			if color[0:1] == "w" {

				switch color[1:2] {
				case "P":
					white[0]++
				case "N":
					white[1]++
				case "B":
					white[2]++
				case "R":
					white[3]++
				case "Q":
					white[4]++
				case "K":
					white[5]++
				default:
					fmt.Println("Incorrect piece mate.go no material 1")
				}
			} else if color[0:1] == "b" {
				switch color[1:2] {
				case "P":
					black[0]++
				case "N":
					black[1]++
				case "B":
					black[2]++
				case "R":
					black[3]++
				case "Q":
					black[4]++
				case "K":
					black[5]++
				default:
					fmt.Println("Incorrect piece mate.go no material 2")
				}
			}
		}
	}

	//KvK, K+B vs K, K+B vs K+B, K+N vs K, K+N vs K+N.
	// pawn=0 knight=1 bishop=2 rook=3 queen=4 king=5
	//KvK
	if white[0] == 0 && white[1] == 0 && white[2] == 0 && white[3] == 0 && white[4] == 0 && white[5] == 1 && black[0] == 0 && black[1] == 0 && black[2] == 0 && black[3] == 0 && black[4] == 0 && black[5] == 1 {
		return true
		//K+B vs K
	} else if white[0] == 0 && white[1] == 0 && white[2] == 1 && white[3] == 0 && white[4] == 0 && white[5] == 1 && black[0] == 0 && black[1] == 0 && black[2] == 0 && black[3] == 0 && black[4] == 0 && black[5] == 1 {
		return true
		//K vs K+B
	} else if white[0] == 0 && white[1] == 0 && white[2] == 0 && white[3] == 0 && white[4] == 0 && white[5] == 1 && black[0] == 0 && black[1] == 0 && black[2] == 1 && black[3] == 0 && black[4] == 0 && black[5] == 1 {
		return true
		//K+B vs K+B
	} else if white[0] == 0 && white[1] == 0 && white[2] == 1 && white[3] == 0 && white[4] == 0 && white[5] == 1 && black[0] == 0 && black[1] == 0 && black[2] == 1 && black[3] == 0 && black[4] == 0 && black[5] == 1 {
		return true
		//K+N vs K
	} else if white[0] == 0 && white[1] == 1 && white[2] == 0 && white[3] == 0 && white[4] == 0 && white[5] == 1 && black[0] == 0 && black[1] == 0 && black[2] == 0 && black[3] == 0 && black[4] == 0 && black[5] == 1 {
		return true
		//K vs K+N
	} else if white[0] == 0 && white[1] == 0 && white[2] == 0 && white[3] == 0 && white[4] == 0 && white[5] == 1 && black[0] == 0 && black[1] == 1 && black[2] == 0 && black[3] == 0 && black[4] == 0 && black[5] == 1 {
		return true
		//K+N vs K+N
	} else if white[0] == 0 && white[1] == 1 && white[2] == 0 && white[3] == 0 && white[4] == 0 && white[5] == 1 && black[0] == 0 && black[1] == 1 && black[2] == 0 && black[3] == 0 && black[4] == 0 && black[5] == 1 {
		return true
	}
	//otherwise insufficient mating material
	return false
}

//checks if three repetition rule which leads to a draw. returns false if no three repetition is found
func (game *ChessGame) threeRep() bool {

	var eightSrc string
	var eightTar string
	var sevenSrc string
	var sevenTar string
	var sixSrc string
	var sixTar string
	var fiveSrc string
	var fiveTar string
	var fourSrc string
	var fourTar string
	var threeSrc string
	var threeTar string
	var twoSrc string
	var twoTar string
	var oneSrc string
	var oneTar string

	var length = len(game.GameMoves)
	if length >= 8 {
		eightSrc = game.GameMoves[length-1].S
		eightTar = game.GameMoves[length-1].T

		sevenSrc = game.GameMoves[length-2].S
		sevenTar = game.GameMoves[length-2].T

		sixSrc = game.GameMoves[length-3].S
		sixTar = game.GameMoves[length-3].T

		fiveSrc = game.GameMoves[length-4].S
		fiveTar = game.GameMoves[length-4].T

		fourSrc = game.GameMoves[length-5].S
		fourTar = game.GameMoves[length-5].T

		threeSrc = game.GameMoves[length-6].S
		threeTar = game.GameMoves[length-6].T

		twoSrc = game.GameMoves[length-7].S
		twoTar = game.GameMoves[length-7].T

		oneSrc = game.GameMoves[length-8].S
		oneTar = game.GameMoves[length-8].T

		if eightSrc == fourSrc && eightTar == fourTar && sevenSrc == threeSrc &&
			sevenTar == threeTar && sixSrc == twoSrc && sixTar == twoTar &&
			fiveSrc == oneSrc && fiveTar == oneTar {
			return true
		}
	}

	return false
}

//checks if fifty moves have been made without a pawn push or capture
func (table *Table) fiftyMoves(gameID int) bool {
	var thisMove int
	thisMove = (len(All.Games[gameID].GameMoves) + 1) / 2
	//no capture within 50 moves
	if (thisMove - table.lastCapture) >= 50 {
		return true
		//no pawn move within 50 moves
	} else if (thisMove - table.pawnMove) >= 50 {
		return true
	}
	return false
}

//checks if a square is attacked by white in one turn, used to identify check and checkmates
func (table *Table) canWhiteKillSquare(targetRow int8, targetCol int8) bool {
	var i int8
	var j int8

	for i = 0; i < 8; i++ {
		for j = 0; j < 8; j++ {
			color := table.ChessBoard[i][j]
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
					result := table.bishopAttack(i, j, targetRow, targetCol)
					if result == true {
						return true
					}
				case "R":
					result := table.rookAttack(i, j, targetRow, targetCol)
					if result == true {
						return true
					}
				case "Q":
					result := table.queenAttack(i, j, targetRow, targetCol)
					if result == true {
						return true
					}
				case "K":
					result := kingAttack(i, j, targetRow, targetCol)
					if result == true {
						return true
					}

				default:
					//fmt.Println("Invalid piece type")
				}
			}
		}
	}

	return false
}

func (table *Table) canBlackKillSquare(targetRow int8, targetCol int8) bool {
	var i int8
	var j int8

	for i = 0; i < 8; i++ {
		for j = 0; j < 8; j++ {
			color := table.ChessBoard[i][j]
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
					result := table.bishopAttack(i, j, targetRow, targetCol)
					if result == true {
						return true
					}
				case "R":
					result := table.rookAttack(i, j, targetRow, targetCol)
					if result == true {
						return true
					}
				case "Q":
					result := table.queenAttack(i, j, targetRow, targetCol)
					if result == true {
						return true
					}
				case "K":
					result := kingAttack(i, j, targetRow, targetCol)
					if result == true {
						return true
					}
				default:
					//fmt.Println("Invalid piece type")
				}
			}
		}
	}

	return false
}
