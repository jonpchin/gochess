package gostuff

//can't castle in check or through check, can't castle when king has already moved
func (table *Table) canWhiteCastleKing() bool {

	if table.wkMoved == false && table.wkrMoved == false && table.canBlackKillSquare(7, 4) == false &&
		table.canBlackKillSquare(7, 5) == false && table.canBlackKillSquare(7, 6) == false {
		return true
	}
	return false
}

func (table *Table) canBlackCastleKing() bool {
	if table.bkMoved == false && table.bkrMoved == false && table.canWhiteKillSquare(0, 4) == false &&
		table.canWhiteKillSquare(0, 5) == false && table.canWhiteKillSquare(0, 6) == false {
		return true
	}
	return false
}

func (table *Table) canWhiteCastleQueen() bool {
	if table.wkMoved == false && table.wqrMoved == false && table.canBlackKillSquare(7, 3) == false &&
		table.canBlackKillSquare(7, 2) == false {
		return true
	}
	return false
}

func (table *Table) canBlackCastleQueen() bool {
	if table.bkMoved == false && table.bqrMoved == false && table.canWhiteKillSquare(0, 3) == false &&
		table.canWhiteKillSquare(0, 2) == false {
		return true
	}
	return false
}
