package gostuff

//can't castle in check or through check, can't castle when king has already moved
func canWhiteCastleKing(gameID int) bool {

	if Verify.AllTables[gameID].wkMoved == false && Verify.AllTables[gameID].wkrMoved == false && canBlackKillSquare(7, 4, gameID) == false && canBlackKillSquare(7, 5, gameID) == false && canBlackKillSquare(7, 6, gameID) == false {
		return true
	}
	return false
}

func canBlackCastleKing(gameID int) bool {
	if Verify.AllTables[gameID].bkMoved == false && Verify.AllTables[gameID].bkrMoved == false && canWhiteKillSquare(0, 4, gameID) == false && canWhiteKillSquare(0, 5, gameID) == false && canWhiteKillSquare(0, 6, gameID) == false {
		return true
	}
	return false
}

func canWhiteCastleQueen(gameID int) bool {
	if Verify.AllTables[gameID].wkMoved == false && Verify.AllTables[gameID].wqrMoved == false && canBlackKillSquare(7, 3, gameID) == false && canBlackKillSquare(7, 2, gameID) == false {
		return true
	}
	return false
}

func canBlackCastleQueen(gameID int) bool {
	if Verify.AllTables[gameID].bkMoved == false && Verify.AllTables[gameID].bqrMoved == false && canWhiteKillSquare(0, 3, gameID) == false && canWhiteKillSquare(0, 2, gameID) == false {
		return true
	}
	return false
}
