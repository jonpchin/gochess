package gostuff

// import "fmt"
//makes all possible black pawn moves to check for stalemate
func blackPawnStaleMate(x int8, y int8, gameID int16) bool {

	if x == 1 { //then pawn can move two squares
		if Verify.AllTables[gameID].ChessBoard[x+2][y] == "-" { //make sure square is not blocked
			capturedPiece := makeMove(x, y, x+2, y, "bP", gameID)
			if isBlackInCheck(gameID) == true { //if black is in check then undo move as its not stalemate
				undoMove(x, y, x+2, y, "bP", capturedPiece, gameID)		
			}else{
				undoMove(x, y, x+2, y, "bP", capturedPiece, gameID)
				return false //otherwise if there is a legal black pawn move and black is not in check its not stalemate
			}
			
		}
	}
	if x+1 <= 7 && Verify.AllTables[gameID].ChessBoard[x+1][y] == "-" { //checking pawn movement of one square
		capturedPiece := makeMove(x, y, x+1, y, "bP", gameID)
		if isBlackInCheck(gameID) == true { 
			undoMove(x, y, x+1, y, "bP", capturedPiece, gameID)
			
		}else{
			undoMove(x, y, x+1, y, "bP", capturedPiece, gameID)
			return false
		}
		
	}

	if y+1 <= 7 && x+1 <= 7 && (Verify.AllTables[gameID].ChessBoard[x+1][y+1])[0:1] == "w" { //right capture

		capturedPiece := makeMove(x, y, x+1, y+1, "bP", gameID)
		if isBlackInCheck(gameID) == true {
			undoMove(x, y, x+1, y+1, "bP", capturedPiece, gameID)
			
		}else{
			undoMove(x, y, x+1, y+1, "bP", capturedPiece, gameID)
			return false
		}
		
	}

	if y-1 >= 0 && x+1 <= 7 && (Verify.AllTables[gameID].ChessBoard[x+1][y-1])[0:1] == "w" { //left capture

		capturedPiece := makeMove(x, y, x+1, y-1, "bP", gameID)
		if isBlackInCheck(gameID) == true {
			undoMove(x, y, x+1, y-1, "bP", capturedPiece, gameID)
			
		}else{
			undoMove(x, y, x+1, y-1, "bP", capturedPiece, gameID)
			return false
		}
		
	}
	if Verify.AllTables[gameID].blackPass[y] == true { //checking if black can enpassent with this pawn
		if y+1 <= 7 && x+1 <= 7 && Verify.AllTables[gameID].ChessBoard[x][y+1] == "wP" { //checking if white pawn exist on the right

			capturedPiece := makeMove(x, y, x+1, y+1, "bP", gameID)
			if isBlackInCheck(gameID) == true {
				undoMove(x, y, x+1, y+1, "bP", capturedPiece, gameID)
				
			}else{
				undoMove(x, y, x+1, y+1, "bP", capturedPiece, gameID)
				return false
			}
			
			
		}
		if y-1 >= 0 && x+1 <= 7 && Verify.AllTables[gameID].ChessBoard[x][y-1] == "wP" { //checking left

			capturedPiece := makeMove(x, y, x+1, y-1, "bP", gameID)
			if isBlackInCheck(gameID) == true {
				undoMove(x, y, x+1, y-1, "bP", capturedPiece, gameID)
				
			}else{
				undoMove(x, y, x+1, y-1, "bP", capturedPiece, gameID)
				return false
			}
			
		}

	}
	//returns true if all possible move of THIS pawn lead to white king still being in check
	return true
}
func whitePawnStaleMate(x int8, y int8, gameID int16) bool { //x is row y is col

	if x == 6 { //then pawn can move two squares
		if Verify.AllTables[gameID].ChessBoard[x-2][y] == "-" { //make sure square is not blocked
			capturedPiece := makeMove(x, y, x-2, y, "wP", gameID)
			if isWhiteInCheck(gameID) == true { //if white is no longer in check then its not mate
				undoMove(x, y, x-2, y, "wP", capturedPiece, gameID)
				
			}else{
				undoMove(x, y, x-2, y, "wP", capturedPiece, gameID)
				return false
			}
			
		}
	}
	if x-1 >= 0 && Verify.AllTables[gameID].ChessBoard[x-1][y] == "-" { //checking pawn movement of one square
		capturedPiece := makeMove(x, y, x-1, y, "wP", gameID)
		if isWhiteInCheck(gameID) == true { //if white is no longer in check then its not mate
			undoMove(x, y, x-1, y, "wP", capturedPiece, gameID)
			
		}else{
			undoMove(x, y, x-1, y, "wP", capturedPiece, gameID) 
			return false
		}
		
	}

	if x-1 >= 0 && y-1 >= 0 && (Verify.AllTables[gameID].ChessBoard[x-1][y-1])[0:1] == "b" { //left capture

		capturedPiece := makeMove(x, y, x-1, y-1, "wP", gameID)
		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, x-1, y-1, "wP", capturedPiece, gameID)
		}else{
			undoMove(x, y, x-1, y-1, "wP", capturedPiece, gameID)
			return false
		}
		
	}

	if x-1 >= 0 && y+1 <= 7 && (Verify.AllTables[gameID].ChessBoard[x-1][y+1])[0:1] == "b" { //right capture

		capturedPiece := makeMove(x, y, x-1, y+1, "wP", gameID)
		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, x-1, y+1, "wP", capturedPiece, gameID)
		}else{
			undoMove(x, y, x-1, y+1, "wP", capturedPiece, gameID)
			return false
		}
		
	}

	if Verify.AllTables[gameID].whitePass[y] == true { //checking if white can enpassent with this pawn

		if y+1 <= 7 && x-1 >= 0 && Verify.AllTables[gameID].ChessBoard[x][y+1] == "bP" { //checking if black pawn exist on the right

			capturedPiece := makeMove(x, y, x-1, y+1, "wP", gameID)
			if isWhiteInCheck(gameID) == true {
				undoMove(x, y, x-1, y+1, "wP", capturedPiece, gameID)
				
			}else{
				undoMove(x, y, x-1, y+1, "wP", capturedPiece, gameID)
				return false
			}
			
		}
		if y-1 >= 0 && x-1 >= 0 && Verify.AllTables[gameID].ChessBoard[x][y-1] == "bP" { //checking left

			capturedPiece := makeMove(x, y, x-1, y-1, "wP", gameID)
			if isWhiteInCheck(gameID) == true {
				undoMove(x, y, x-1, y-1, "wP", capturedPiece, gameID)
				
			}else{
				undoMove(x, y, x-1, y-1, "wP", capturedPiece, gameID)
				return false
			}
			
		}

	}

	//returns true if all possible move of THIS pawn lead to white king still being in check
	return true
}

//brute forces all possible knight moves for white
func whiteKnightStaleMate(x int8, y int8, gameID int16) bool {

	//starting at top left going clockwise
	if x-2 >= 0 && y-1 >= 0 && (Verify.AllTables[gameID].ChessBoard[x-2][y-1])[0:1] != "w" {

		capturedPiece := makeMove(x, y, x-2, y-1, "wN", gameID)
		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, x-2, y-1, "wN", capturedPiece, gameID)
		}else{
			undoMove(x, y, x-2, y-1, "wN", capturedPiece, gameID)
			return false
		}
		
	} //white can't capture own piece
	if x-2 >= 0 && y+1 <= 7 && (Verify.AllTables[gameID].ChessBoard[x-2][y+1])[0:1] != "w" {

		capturedPiece := makeMove(x, y, x-2, y+1, "wN", gameID)
		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, x-2, y+1, "wN", capturedPiece, gameID)
		}else{
			undoMove(x, y, x-2, y+1, "wN", capturedPiece, gameID)
			return false
		}
		
	}
	if x-1 >= 0 && y+2 <= 7 && (Verify.AllTables[gameID].ChessBoard[x-1][y+2])[0:1] != "w" {

		capturedPiece := makeMove(x, y, x-1, y+2, "wN", gameID)
		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, x-1, y+2, "wN", capturedPiece, gameID)
		}else{
			undoMove(x, y, x-1, y+2, "wN", capturedPiece, gameID)
			return false
		}
		
	}
	if x+1 <= 7 && y+2 <= 7 && (Verify.AllTables[gameID].ChessBoard[x+1][y+2])[0:1] != "w" {

		capturedPiece := makeMove(x, y, x+1, y+2, "wN", gameID)
		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, x+1, y+2, "wN", capturedPiece, gameID)
		}else{
			undoMove(x, y, x+1, y+2, "wN", capturedPiece, gameID)
			return false
		}
		
	}
	if x+2 <= 7 && y+1 <= 7 && (Verify.AllTables[gameID].ChessBoard[x+2][y+1])[0:1] != "w" {

		capturedPiece := makeMove(x, y, x+2, y+1, "wN", gameID)
		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, x+2, y+1, "wN", capturedPiece, gameID)
		}else{
			undoMove(x, y, x+2, y+1, "wN", capturedPiece, gameID)
			return false
		}
		
	}
	if x+2 <= 7 && y-1 >= 0 && (Verify.AllTables[gameID].ChessBoard[x+2][y-1])[0:1] != "w" {

		capturedPiece := makeMove(x, y, x+2, y-1, "wN", gameID)
		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, x+2, y-1, "wN", capturedPiece, gameID)
		}else{
			undoMove(x, y, x+2, y-1, "wN", capturedPiece, gameID)
			return false
		}
		
	}
	if x+1 <= 7 && y-2 >= 0 && (Verify.AllTables[gameID].ChessBoard[x+1][y-2])[0:1] != "w" {

		capturedPiece := makeMove(x, y, x+1, y-2, "wN", gameID)
		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, x+1, y-2, "wN", capturedPiece, gameID)
		}else{
			undoMove(x, y, x+1, y-2, "wN", capturedPiece, gameID)
			return false
		}
		
	}
	if x-1 >= 0 && y-2 >= 0 && (Verify.AllTables[gameID].ChessBoard[x-1][y-2])[0:1] != "w" {

		capturedPiece := makeMove(x, y, x-1, y-2, "wN", gameID)
		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, x-1, y-2, "wN", capturedPiece, gameID)
		}else{
			undoMove(x, y, x-1, y-2, "wN", capturedPiece, gameID)
			return false
		}
		
	}
	return true
}

func blackKnightStaleMate(x int8, y int8, gameID int16) bool {
	if x-2 >= 0 && y-1 >= 0 && (Verify.AllTables[gameID].ChessBoard[x-2][y-1])[0:1] != "b" { //starting at top left going clockwise

		capturedPiece := makeMove(x, y, x-2, y-1, "bN", gameID)
		if isBlackInCheck(gameID) == true {
			undoMove(x, y, x-2, y-1, "bN", capturedPiece, gameID)
		}else{
			undoMove(x, y, x-2, y-1, "bN", capturedPiece, gameID)
			return false
		}
		
	}
	if x-2 >= 0 && y+1 <= 7 && (Verify.AllTables[gameID].ChessBoard[x-2][y+1])[0:1] != "b" {

		capturedPiece := makeMove(x, y, x-2, y+1, "bN", gameID)
		if isBlackInCheck(gameID) == true {
			undoMove(x, y, x-2, y+1, "bN", capturedPiece, gameID)
		}else{
			undoMove(x, y, x-2, y+1, "bN", capturedPiece, gameID)
			return false
		}
		
	}
	if x-1 >= 0 && y+2 <= 7 && (Verify.AllTables[gameID].ChessBoard[x-1][y+2])[0:1] != "b" {

		capturedPiece := makeMove(x, y, x-1, y+2, "bN", gameID)
		if isBlackInCheck(gameID) == true {
			undoMove(x, y, x-1, y+2, "bN", capturedPiece, gameID)
		}else{
			undoMove(x, y, x-1, y+2, "bN", capturedPiece, gameID)
			return false
		}
		
	}
	if x+1 <= 7 && y+2 <= 7 && (Verify.AllTables[gameID].ChessBoard[x+1][y+2])[0:1] != "b" {

		capturedPiece := makeMove(x, y, x+1, y+2, "bN", gameID)
		if isBlackInCheck(gameID) == true {
			undoMove(x, y, x+1, y+2, "bN", capturedPiece, gameID)
		}else{
			undoMove(x, y, x+1, y+2, "bN", capturedPiece, gameID)
			return false
		}
		
	}
	if x+2 <= 7 && y+1 <= 7 && (Verify.AllTables[gameID].ChessBoard[x+2][y+1])[0:1] != "b" {

		capturedPiece := makeMove(x, y, x+2, y+1, "bN", gameID)
		if isBlackInCheck(gameID) == true {
			undoMove(x, y, x+2, y+1, "bN", capturedPiece, gameID)
		}else{
			undoMove(x, y, x+2, y+1, "bN", capturedPiece, gameID)
			return false
		}
		
	}
	if x+2 <= 7 && y-1 >= 0 && (Verify.AllTables[gameID].ChessBoard[x+2][y-1])[0:1] != "b" {

		capturedPiece := makeMove(x, y, x+2, y-1, "bN", gameID)
		if isBlackInCheck(gameID) == true {
			undoMove(x, y, x+2, y-1, "bN", capturedPiece, gameID)
		}else{
			undoMove(x, y, x+2, y-1, "bN", capturedPiece, gameID)
			return false
		}
		
	}
	if x+1 <= 7 && y-2 >= 0 && (Verify.AllTables[gameID].ChessBoard[x+1][y-2])[0:1] != "b" {

		capturedPiece := makeMove(x, y, x+1, y-2, "bN", gameID)
		if isBlackInCheck(gameID) == true {
			undoMove(x, y, x+1, y-2, "bN", capturedPiece, gameID)
		}else{
			undoMove(x, y, x+1, y-2, "bN", capturedPiece, gameID)
			return false
		}
		
	}
	if x-1 >= 0 && y-2 >= 0 && (Verify.AllTables[gameID].ChessBoard[x-1][y-2])[0:1] != "b" {

		capturedPiece := makeMove(x, y, x-1, y-2, "bN", gameID)
		if isBlackInCheck(gameID) == true {
			undoMove(x, y, x-1, y-2, "bN", capturedPiece, gameID)
		}else{
			undoMove(x, y, x-1, y-2, "bN", capturedPiece, gameID)
			return false
		}
		
	}
	return true
}

func whiteBishopStaleMate(x int8, y int8, gameID int16) bool { //moves all possible moves with this one wB

	var i int8
	var j int8
	for i, j = x-1, y-1; i >= 0 && j >= 0; i, j = i-1, j-1 { //left upper

		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "w" {
			break //no need to check rest of squares if the square is already occupied by same piece
		}
		capturedPiece := makeMove(x, y, i, j, "wB", gameID)
		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, i, j, "wB", capturedPiece, gameID)
		}else{
			undoMove(x, y, i, j, "wB", capturedPiece, gameID)
			return false
		}
		

		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "b" {
			break //no need to check rest of squares if the square is already occupied by enemy piece
		}
	}
	for i, j = x-1, y+1; i >= 0 && j <= 7; i, j = i-1, j+1 { //right upper

		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "w" {
			break
		}
		capturedPiece := makeMove(x, y, i, j, "wB", gameID)
		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, i, j, "wB", capturedPiece, gameID)
		}else{
			undoMove(x, y, i, j, "wB", capturedPiece, gameID)
			return false
		}
		

		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "b" {
			break
		}
	}
	for i, j = x+1, y+1; i <= 7 && j <= 7; i, j = i+1, j+1 { //right lower

		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "w" {
			break
		}
		capturedPiece := makeMove(x, y, i, j, "wB", gameID)
		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, i, j, "wB", capturedPiece, gameID)
		}else{
			undoMove(x, y, i, j, "wB", capturedPiece, gameID)
			return false
		}
		
		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "b" {
			break
		}
	}
	for i, j = x+1, y-1; i <= 7 && j >= 0; i, j = i+1, j-1 { //left lower

		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "w" {
			break
		}
		capturedPiece := makeMove(x, y, i, j, "wB", gameID)
		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, i, j, "wB", capturedPiece, gameID)
		}else{
			undoMove(x, y, i, j, "wB", capturedPiece, gameID)
			return false
		}
		
		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "b" {
			break
		}
	}

	return true
}

func blackBishopStaleMate(x int8, y int8, gameID int16) bool {
	var i int8
	var j int8
	for i, j = x-1, y-1; i >= 0 && j >= 0; i, j = i-1, j-1 { //left upper

		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "b" {
			break //no need to check rest of squares if the square is already occupied by a piece
		}
		capturedPiece := makeMove(x, y, i, j, "bB", gameID)
		if isBlackInCheck(gameID) == true {
			undoMove(x, y, i, j, "bB", capturedPiece, gameID)
		}else{
			undoMove(x, y, i, j, "bB", capturedPiece, gameID)
			return false
		}
		

		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "w" {
			break
		}
	}
	for i, j = x-1, y+1; i >= 0 && j <= 7; i, j = i-1, j+1 { //right upper

		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "b" {
			break //no need to check rest of squares if the square is already occupied by a piece
		}
		capturedPiece := makeMove(x, y, i, j, "bB", gameID)
		if isBlackInCheck(gameID) == true {
			undoMove(x, y, i, j, "bB", capturedPiece, gameID)
		}else{
			undoMove(x, y, i, j, "bB", capturedPiece, gameID)
			return false
		}
		
		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "w" {
			break //no need to check rest of squares if the square is already occupied by a piece
		}
	}
	for i, j = x+1, y+1; i <= 7 && j <= 7; i, j = i+1, j+1 { //right lower

		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "b" {
			break //no need to check rest of squares if the square is already occupied by a piece
		}
		capturedPiece := makeMove(x, y, i, j, "bB", gameID)
		if isBlackInCheck(gameID) == true {
			undoMove(x, y, i, j, "bB", capturedPiece, gameID)
		}else{
			undoMove(x, y, i, j, "bB", capturedPiece, gameID)
			return false
		}
		
		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "w" {
			break
		}
	}
	for i, j = x+1, y-1; i <= 7 && j >= 0; i, j = i+1, j-1 { //left lower

		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "b" {
			break
		}
		capturedPiece := makeMove(x, y, i, j, "bB", gameID)
		if isBlackInCheck(gameID) == true {
			undoMove(x, y, i, j, "bB", capturedPiece, gameID)
		}else{
			undoMove(x, y, i, j, "bB", capturedPiece, gameID)
			return false
		}
		
		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "w" {
			break
		}
	}

	return true
}

func whiteRookStaleMate(x int8, y int8, gameID int16) bool {

	var i int8
	var j int8
	for i = x - 1; i >= 0; i-- { //up

		if (Verify.AllTables[gameID].ChessBoard[i][y])[0:1] == "w" {
			break
		}
		capturedPiece := makeMove(x, y, i, y, "wR", gameID)
		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, i, y, "wR", capturedPiece, gameID)
		}else{
			undoMove(x, y, i, y, "wR", capturedPiece, gameID)
			return false
		}
		
		if (Verify.AllTables[gameID].ChessBoard[i][y])[0:1] == "b" {
			break
		}
	}
	for i = x + 1; i <= 7; i++ { //down

		if (Verify.AllTables[gameID].ChessBoard[i][y])[0:1] == "w" {
			break //no need to check rest of squares if the square is already occupied by a piece
		}
		capturedPiece := makeMove(x, y, i, y, "wR", gameID)
		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, i, y, "wR", capturedPiece, gameID)
		}else{
			undoMove(x, y, i, y, "wR", capturedPiece, gameID)
			return false
		}
		
		if (Verify.AllTables[gameID].ChessBoard[i][y])[0:1] == "b" {
			break
		}
	}
	for j = y - 1; j >= 0; j-- { //left

		if (Verify.AllTables[gameID].ChessBoard[x][j])[0:1] == "w" {
			break //no need to check rest of squares if the square is already occupied by a piece
		}
		capturedPiece := makeMove(x, y, x, j, "wR", gameID)
		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, x, j, "wR", capturedPiece, gameID)
		}else{
			undoMove(x, y, x, j, "wR", capturedPiece, gameID)
			return false
		}
		
		if (Verify.AllTables[gameID].ChessBoard[x][j])[0:1] == "b" {
			break //no need to check rest of squares if the square is already occupied by a piece
		}
	}
	for j = y + 1; j <= 7; j++ { //right

		if (Verify.AllTables[gameID].ChessBoard[x][j])[0:1] == "w" {
			break //no need to check rest of squares if the square is already occupied by a piece
		}
		capturedPiece := makeMove(x, y, x, j, "wR", gameID)
		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, x, j, "wR", capturedPiece, gameID)
		}else{
			undoMove(x, y, x, j, "wR", capturedPiece, gameID)
			return false
		}
		
		if (Verify.AllTables[gameID].ChessBoard[x][j])[0:1] == "b" {
			break //no need to check rest of squares if the square is already occupied by a piece
		}
	}
	return true
}

func blackRookStaleMate(x int8, y int8, gameID int16) bool {

	var i int8
	var j int8
	for i = x - 1; i >= 0; i-- { //up

		if (Verify.AllTables[gameID].ChessBoard[i][y])[0:1] == "b" {
			break //no need to check rest of squares if the square is already occupied by a piece
		}
		capturedPiece := makeMove(x, y, i, y, "bR", gameID)
		if isBlackInCheck(gameID) == true {
			undoMove(x, y, i, y, "bR", capturedPiece, gameID)
		}else{
			undoMove(x, y, i, y, "bR", capturedPiece, gameID)
			return false
		}
		
		if (Verify.AllTables[gameID].ChessBoard[i][y])[0:1] == "w" {
			break
		}
	}
	for i = x + 1; i <= 7; i++ { //down

		if (Verify.AllTables[gameID].ChessBoard[i][y])[0:1] == "b" {
			break
		}
		capturedPiece := makeMove(x, y, i, y, "bR", gameID)
		if isBlackInCheck(gameID) == true {
			undoMove(x, y, i, y, "bR", capturedPiece, gameID)
		}else{
			undoMove(x, y, i, y, "bR", capturedPiece, gameID)
			return false
		}
		
		if (Verify.AllTables[gameID].ChessBoard[i][y])[0:1] == "w" {
			break
		}
	}
	for j = y - 1; j >= 0; j-- { //left

		if (Verify.AllTables[gameID].ChessBoard[x][j])[0:1] == "b" {
			break //no need to check rest of squares if the square is already occupied by a piece
		}
		capturedPiece := makeMove(x, y, x, j, "bR", gameID)
		if isBlackInCheck(gameID) ==  true{
			undoMove(x, y, x, j, "bR", capturedPiece, gameID)
		}else{
			undoMove(x, y, x, j, "bR", capturedPiece, gameID)
			return false
		}
		
		if (Verify.AllTables[gameID].ChessBoard[x][j])[0:1] == "w" {
			break
		}
	}
	for j = y + 1; j <= 7; j++ { //right

		if (Verify.AllTables[gameID].ChessBoard[x][j])[0:1] == "b" {
			break
		}
		capturedPiece := makeMove(x, y, x, j, "bR", gameID)
		if isBlackInCheck(gameID) == true {
			undoMove(x, y, x, j, "bR", capturedPiece, gameID)
		}else{
			undoMove(x, y, x, j, "bR", capturedPiece, gameID)
			return false
		}
		
		if (Verify.AllTables[gameID].ChessBoard[x][j])[0:1] == "w" {
			break
		}
	}
	return true
}

func whiteQueenStaleMate(x int8, y int8, gameID int16) bool {

	var i int8
	var j int8

	for i, j = x-1, y-1; i >= 0 && j >= 0; i, j = i-1, j-1 { //left upper

		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "w" {
			break
		}
		capturedPiece := makeMove(x, y, i, j, "wQ", gameID)
		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, i, j, "wQ", capturedPiece, gameID)
		}else{
			undoMove(x, y, i, j, "wQ", capturedPiece, gameID)
			return false
		}
		
		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "b" {
			break
		}
	}
	for i, j = x-1, y+1; i >= 0 && j <= 7; i, j = i-1, j+1 { //right upper

		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "w" {
			break
		}
		capturedPiece := makeMove(x, y, i, j, "wQ", gameID)
		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, i, j, "wQ", capturedPiece, gameID)
		}else{
			undoMove(x, y, i, j, "wQ", capturedPiece, gameID)
			return false
		}
		
		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "b" {
			break
		}
	}
	for i, j = x+1, y+1; i <= 7 && j <= 7; i, j = i+1, j+1 { //right lower

		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "w" {
			break
		}
		capturedPiece := makeMove(x, y, i, j, "wQ", gameID)
		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, i, j, "wQ", capturedPiece, gameID)
		}else{
			undoMove(x, y, i, j, "wQ", capturedPiece, gameID)
			return false
		}
		
		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "b" {
			break
		}
	}
	for i, j = x+1, y-1; i <= 7 && j >= 0; i, j = i+1, j-1 { //left lower

		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "w" {
			break
		}
		capturedPiece := makeMove(x, y, i, j, "wQ", gameID)
		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, i, j, "wQ", capturedPiece, gameID)
		}else{
			undoMove(x, y, i, j, "wQ", capturedPiece, gameID)
			return false
		}
		
		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "b" {
			break
		}
	}

	//rook moves are now checked
	for i = x - 1; i >= 0; i-- { //up

		if (Verify.AllTables[gameID].ChessBoard[i][y])[0:1] == "w" {
			break
		}
		capturedPiece := makeMove(x, y, i, y, "wQ", gameID)
		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, i, y, "wQ", capturedPiece, gameID)
		}else{
			undoMove(x, y, i, y, "wQ", capturedPiece, gameID)
			return false
		}
		
		if (Verify.AllTables[gameID].ChessBoard[i][y])[0:1] == "b" {
			break
		}
	}
	for i = x + 1; i <= 7; i++ { //down

		if (Verify.AllTables[gameID].ChessBoard[i][y])[0:1] == "w" {
			break
		}
		capturedPiece := makeMove(x, y, i, y, "wQ", gameID)
		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, i, y, "wQ", capturedPiece, gameID)
		}else{
			undoMove(x, y, i, y, "wQ", capturedPiece, gameID)
			return false
		}
		
		if (Verify.AllTables[gameID].ChessBoard[i][y])[0:1] == "b" {
			break
		}
	}
	for j = y - 1; j >= 0; j-- { //left

		if (Verify.AllTables[gameID].ChessBoard[x][j])[0:1] == "w" {
			break
		}
		capturedPiece := makeMove(x, y, x, j, "wQ", gameID)
		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, x, j, "wQ", capturedPiece, gameID)
		}else{
			undoMove(x, y, x, j, "wQ", capturedPiece, gameID)
			return false
		}
		
		if (Verify.AllTables[gameID].ChessBoard[x][j])[0:1] == "b" {
			break
		}
	}
	for j = y + 1; j <= 7; j++ { //right

		if (Verify.AllTables[gameID].ChessBoard[x][j])[0:1] == "w" {
			break
		}
		capturedPiece := makeMove(x, y, x, j, "wQ", gameID)
		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, x, j, "wQ", capturedPiece, gameID)
		}else{
			undoMove(x, y, x, j, "wQ", capturedPiece, gameID)
			return false
		}
		
		if (Verify.AllTables[gameID].ChessBoard[x][j])[0:1] == "b" {
			break
		}
	}

	return true
}

func blackQueenStaleMate(x int8, y int8, gameID int16) bool {
	var i int8
	var j int8

	for i, j = x-1, y-1; i >= 0 && j >= 0; i, j = i-1, j-1 { //left upper

		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "b" {
			break //no need to check rest of squares if the square is already occupied by same color
		}
		capturedPiece := makeMove(x, y, i, j, "bQ", gameID)
		if isBlackInCheck(gameID) == true {
			undoMove(x, y, i, j, "bQ", capturedPiece, gameID)
		}else{
			undoMove(x, y, i, j, "bQ", capturedPiece, gameID)
			return false
		}
		

		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "w" {
			break //no need to check rest of squares if the square is already occupied by enemy piece
		}
	}
	for i, j = x-1, y+1; i >= 0 && j <= 7; i, j = i-1, j+1 { //right upper

		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "b" {
			break //no need to check rest of squares if the square is already occupied by a piece
		}
		capturedPiece := makeMove(x, y, i, j, "bQ", gameID)
		if isBlackInCheck(gameID) == true {
			undoMove(x, y, i, j, "bQ", capturedPiece, gameID)
		}else{
			undoMove(x, y, i, j, "bQ", capturedPiece, gameID)
			return false
		}
		

		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "w" {
			break //no need to check rest of squares if the square is already occupied by enemy piece
		}
	}
	for i, j = x+1, y+1; i <= 7 && j <= 7; i, j = i+1, j+1 { //right lower

		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "b" {
			break //no need to check rest of squares if the square is already occupied by a piece
		}
		capturedPiece := makeMove(x, y, i, j, "bQ", gameID)
		if isBlackInCheck(gameID) == true {
			undoMove(x, y, i, j, "bQ", capturedPiece, gameID)
		}else{
			undoMove(x, y, i, j, "bQ", capturedPiece, gameID)
			return false
		}
		

		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "w" {
			break //no need to check rest of squares if the square is already occupied by enemy piece
		}
	}
	for i, j = x+1, y-1; i <= 7 && j >= 0; i, j = i+1, j-1 { //left lower

		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "b" {
			break //no need to check rest of squares if the square is already occupied by a piece
		}
		capturedPiece := makeMove(x, y, i, j, "bQ", gameID)
		if isBlackInCheck(gameID) == true {
			undoMove(x, y, i, j, "bQ", capturedPiece, gameID)
		}else{
			undoMove(x, y, i, j, "bQ", capturedPiece, gameID)
			return false
		}
		

		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "w" {
			break //no need to check rest of squares if the square is already occupied by enemy piece
		}
	}

	//rook moves are now checked
	for i = x - 1; i >= 0; i-- { //up

		if (Verify.AllTables[gameID].ChessBoard[i][y])[0:1] == "b" {
			break //no need to check rest of squares if the square is already occupied by a piece
		}
		capturedPiece := makeMove(x, y, i, y, "bQ", gameID)
		if isBlackInCheck(gameID) == true {
			undoMove(x, y, i, y, "bQ", capturedPiece, gameID)
		}else{
			undoMove(x, y, i, y, "bQ", capturedPiece, gameID)
			return false
		}
		
		

		if (Verify.AllTables[gameID].ChessBoard[i][y])[0:1] == "w" {
			break //no need to check rest of squares if the square is already occupied by enemy piece
		}
	}
	for i = x + 1; i <= 7; i++ { //down

		if (Verify.AllTables[gameID].ChessBoard[i][y])[0:1] == "b" {
			break //no need to check rest of squares if the square is already occupied by same piece
		}
		capturedPiece := makeMove(x, y, i, y, "bQ", gameID)
		if isBlackInCheck(gameID) == true {
			undoMove(x, y, i, y, "bQ", capturedPiece, gameID)
		}else{
			undoMove(x, y, i, y, "bQ", capturedPiece, gameID)
			return false
		}
		

		if (Verify.AllTables[gameID].ChessBoard[i][y])[0:1] == "w" {
			break //no need to check rest of squares if the square is already occupied by enemy piece
		}
	}
	for j = y - 1; j >= 0; j-- { //left

		if (Verify.AllTables[gameID].ChessBoard[x][j])[0:1] == "b" {
			break //no need to check rest of squares if the square is already occupied by a piece
		}
		capturedPiece := makeMove(x, y, x, j, "bQ", gameID)
		if isBlackInCheck(gameID) == true {
			undoMove(x, y, x, j, "bQ", capturedPiece, gameID)
		}else{
			undoMove(x, y, x, j, "bQ", capturedPiece, gameID)
			return false
		}
		

		if (Verify.AllTables[gameID].ChessBoard[i][j])[0:1] == "w" {
			break //no need to check rest of squares if the square is already occupied by enemy piece
		}
	}
	for j = y + 1; j <= 7; j++ { //right

		if (Verify.AllTables[gameID].ChessBoard[x][j])[0:1] == "b" {
			break //no need to check rest of squares if the square is already occupied by a piece
		}
		capturedPiece := makeMove(x, y, x, j, "bQ", gameID)
		if isBlackInCheck(gameID) == true {
			undoMove(x, y, x, j, "bQ", capturedPiece, gameID)
		}else{
			undoMove(x, y, x, j, "bQ", capturedPiece, gameID)
			return false
		}
		

		if (Verify.AllTables[gameID].ChessBoard[x][j])[0:1] == "w" {
			break //no need to check rest of squares if the square is already occupied by a piece
		}
	}

	return true
}

func whiteKingStaleMate(x int8, y int8, gameID int16) bool {
	//starting from top left corner
	if x-1 >= 0 && y-1 >= 0 && (Verify.AllTables[gameID].ChessBoard[x-1][y-1])[0:1] != "w" { //top left

		capturedPiece := makeMove(x, y, x-1, y-1, "wK", gameID)

		Verify.AllTables[gameID].whiteOldX = x //storing king information
		Verify.AllTables[gameID].whiteOldY = y

		Verify.AllTables[gameID].whiteKingX = x - 1
		Verify.AllTables[gameID].whiteKingY = y - 1
		Verify.AllTables[gameID].kingUpdate = true

		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, x-1, y-1, "wK", capturedPiece, gameID)
		}else{
			undoMove(x, y, x-1, y-1, "wK", capturedPiece, gameID)
			return false
		}
		

	}

	if x-1 >= 0 && (Verify.AllTables[gameID].ChessBoard[x-1][y])[0:1] != "w" { //top middle

		capturedPiece := makeMove(x, y, x-1, y, "wK", gameID)

		Verify.AllTables[gameID].whiteOldX = x //storing king information
		Verify.AllTables[gameID].whiteOldY = y

		Verify.AllTables[gameID].whiteKingX = x - 1
		Verify.AllTables[gameID].whiteKingY = y
		Verify.AllTables[gameID].kingUpdate = true

		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, x-1, y, "wK", capturedPiece, gameID)
		}else{
			undoMove(x, y, x-1, y, "wK", capturedPiece, gameID)
			return false
		}
		
	}

	if x-1 >= 0 && y+1 <= 7 && (Verify.AllTables[gameID].ChessBoard[x-1][y+1])[0:1] != "w" { //top right

		capturedPiece := makeMove(x, y, x-1, y+1, "wK", gameID)

		Verify.AllTables[gameID].whiteOldX = x //storing king information
		Verify.AllTables[gameID].whiteOldY = y

		Verify.AllTables[gameID].whiteKingX = x - 1
		Verify.AllTables[gameID].whiteKingY = y + 1
		Verify.AllTables[gameID].kingUpdate = true

		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, x-1, y+1, "wK", capturedPiece, gameID)
		}else{
			undoMove(x, y, x-1, y+1, "wK", capturedPiece, gameID)
			return false
		}
		
	}

	if y-1 >= 0 && (Verify.AllTables[gameID].ChessBoard[x][y-1])[0:1] != "w" { //middle left

		capturedPiece := makeMove(x, y, x, y-1, "wK", gameID)

		Verify.AllTables[gameID].whiteOldX = x //storing king information
		Verify.AllTables[gameID].whiteOldY = y

		Verify.AllTables[gameID].whiteKingX = x
		Verify.AllTables[gameID].whiteKingY = y - 1
		Verify.AllTables[gameID].kingUpdate = true

		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, x, y-1, "wK", capturedPiece, gameID)
		}else{
			undoMove(x, y, x, y-1, "wK", capturedPiece, gameID)
			return false
		}
		
	}

	if y+1 <= 7 && (Verify.AllTables[gameID].ChessBoard[x][y+1])[0:1] != "w" { //middle right

		capturedPiece := makeMove(x, y, x, y+1, "wK", gameID)

		Verify.AllTables[gameID].whiteOldX = x //storing king information
		Verify.AllTables[gameID].whiteOldY = y

		Verify.AllTables[gameID].whiteKingX = x
		Verify.AllTables[gameID].whiteKingY = y + 1
		Verify.AllTables[gameID].kingUpdate = true

		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, x, y+1, "wK", capturedPiece, gameID)
		}else{
			undoMove(x, y, x, y+1, "wK", capturedPiece, gameID)
			return false
		}
		
	}

	if x+1 <= 7 && y-1 >= 0 && (Verify.AllTables[gameID].ChessBoard[x+1][y-1])[0:1] != "w" { //bottom left

		capturedPiece := makeMove(x, y, x+1, y-1, "wK", gameID)

		Verify.AllTables[gameID].whiteOldX = x //storing king information
		Verify.AllTables[gameID].whiteOldY = y

		Verify.AllTables[gameID].whiteKingX = x + 1
		Verify.AllTables[gameID].whiteKingY = y - 1
		Verify.AllTables[gameID].kingUpdate = true

		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, x+1, y-1, "wK", capturedPiece, gameID)
		}else{
			undoMove(x, y, x+1, y-1, "wK", capturedPiece, gameID)
			return false
		}
		
	}

	if x+1 <= 7 && (Verify.AllTables[gameID].ChessBoard[x+1][y])[0:1] != "w" { //bottom middle

		capturedPiece := makeMove(x, y, x+1, y, "wK", gameID)

		Verify.AllTables[gameID].whiteOldX = x //storing king information
		Verify.AllTables[gameID].whiteOldY = y

		Verify.AllTables[gameID].whiteKingX = x + 1
		Verify.AllTables[gameID].whiteKingY = y
		Verify.AllTables[gameID].kingUpdate = true

		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, x+1, y, "wK", capturedPiece, gameID)
		}else{
			undoMove(x, y, x+1, y, "wK", capturedPiece, gameID)
			return false
		}
		
	}

	if x+1 <= 7 && y+1 <= 7 && (Verify.AllTables[gameID].ChessBoard[x+1][y+1])[0:1] != "w" { //bottom right

		capturedPiece := makeMove(x, y, x+1, y+1, "wK", gameID)

		Verify.AllTables[gameID].whiteOldX = x //storing king information
		Verify.AllTables[gameID].whiteOldY = y

		Verify.AllTables[gameID].whiteKingX = x + 1
		Verify.AllTables[gameID].whiteKingY = y + 1
		Verify.AllTables[gameID].kingUpdate = true

		if isWhiteInCheck(gameID) == true {
			undoMove(x, y, x+1, y+1, "wK", capturedPiece, gameID)
		}else{
			undoMove(x, y, x+1, y+1, "wK", capturedPiece, gameID)
			return false
		}
		
	}

	return true
}

func blackKingStaleMate(x int8, y int8, gameID int16) bool {
	//starting from top left corner
	if x-1 >= 0 && y-1 >= 0 && (Verify.AllTables[gameID].ChessBoard[x-1][y-1])[0:1] != "b" { //top left

		capturedPiece := makeMove(x, y, x-1, y-1, "bK", gameID)

		Verify.AllTables[gameID].blackOldX = x //storing king information
		Verify.AllTables[gameID].blackOldY = y

		Verify.AllTables[gameID].blackKingX = x - 1
		Verify.AllTables[gameID].blackKingY = y - 1
		Verify.AllTables[gameID].kingUpdate = true

		if isBlackInCheck(gameID) == true {
			undoMove(x, y, x-1, y-1, "bK", capturedPiece, gameID)
		}else{
			undoMove(x, y, x-1, y-1, "bK", capturedPiece, gameID)
			return false
		}
		
	}

	if x-1 >= 0 && (Verify.AllTables[gameID].ChessBoard[x-1][y])[0:1] != "b" { //top middle

		capturedPiece := makeMove(x, y, x-1, y, "bK", gameID)

		Verify.AllTables[gameID].blackOldX = x //storing king information
		Verify.AllTables[gameID].blackOldY = y

		Verify.AllTables[gameID].blackKingX = x - 1
		Verify.AllTables[gameID].blackKingY = y
		Verify.AllTables[gameID].kingUpdate = true

		if isBlackInCheck(gameID) == true {
			undoMove(x, y, x-1, y, "bK", capturedPiece, gameID)
		}else{
			undoMove(x, y, x-1, y, "bK", capturedPiece, gameID)

			return false
		}
		
	}

	if x-1 >= 0 && y+1 <= 7 && (Verify.AllTables[gameID].ChessBoard[x-1][y+1])[0:1] != "b" { //top right

		capturedPiece := makeMove(x, y, x-1, y+1, "bK", gameID)

		Verify.AllTables[gameID].blackOldX = x //storing king information
		Verify.AllTables[gameID].blackOldY = y

		Verify.AllTables[gameID].blackKingX = x - 1
		Verify.AllTables[gameID].blackKingY = y + 1
		Verify.AllTables[gameID].kingUpdate = true

		if isBlackInCheck(gameID) == true {
			undoMove(x, y, x-1, y+1, "bK", capturedPiece, gameID)
		}else{
			undoMove(x, y, x-1, y+1, "bK", capturedPiece, gameID)

			return false
		}
		
	}

	if y-1 >= 0 && (Verify.AllTables[gameID].ChessBoard[x][y-1])[0:1] != "b" { //middle left

		capturedPiece := makeMove(x, y, x, y-1, "bK", gameID)

		Verify.AllTables[gameID].blackOldX = x //storing king information
		Verify.AllTables[gameID].blackOldY = y

		Verify.AllTables[gameID].blackKingX = x
		Verify.AllTables[gameID].blackKingY = y - 1
		Verify.AllTables[gameID].kingUpdate = true

		if isBlackInCheck(gameID) == true {
			undoMove(x, y, x, y-1, "bK", capturedPiece, gameID)
		}else{
			undoMove(x, y, x, y-1, "bK", capturedPiece, gameID)

			return false
		}
		
	}

	if y+1 <= 7 && (Verify.AllTables[gameID].ChessBoard[x][y+1])[0:1] != "b" { //middle right

		capturedPiece := makeMove(x, y, x, y+1, "bK", gameID)

		Verify.AllTables[gameID].blackOldX = x //storing king information
		Verify.AllTables[gameID].blackOldY = y

		Verify.AllTables[gameID].blackKingX = x
		Verify.AllTables[gameID].blackKingY = y + 1
		Verify.AllTables[gameID].kingUpdate = true

		if isBlackInCheck(gameID) == true {
			undoMove(x, y, x, y+1, "bK", capturedPiece, gameID)
		}else{
			undoMove(x, y, x, y+1, "bK", capturedPiece, gameID)

			return false
		}
		
	}

	if x+1 <= 7 && y-1 >= 0 && (Verify.AllTables[gameID].ChessBoard[x+1][y-1])[0:1] != "b" { //bottom left

		capturedPiece := makeMove(x, y, x+1, y-1, "bK", gameID)

		Verify.AllTables[gameID].blackOldX = x //storing king information
		Verify.AllTables[gameID].blackOldY = y

		Verify.AllTables[gameID].blackKingX = x + 1
		Verify.AllTables[gameID].blackKingY = y - 1
		Verify.AllTables[gameID].kingUpdate = true

		if isBlackInCheck(gameID) == true {
			undoMove(x, y, x+1, y-1, "bK", capturedPiece, gameID)
		}else{
			undoMove(x, y, x+1, y-1, "bK", capturedPiece, gameID)

			return false
		}
		
	}

	if x+1 <= 7 && (Verify.AllTables[gameID].ChessBoard[x+1][y])[0:1] != "b" { //bottom middle

		capturedPiece := makeMove(x, y, x+1, y, "bK", gameID)

		Verify.AllTables[gameID].blackOldX = x //storing king information
		Verify.AllTables[gameID].blackOldY = y

		Verify.AllTables[gameID].blackKingX = x + 1
		Verify.AllTables[gameID].blackKingY = y
		Verify.AllTables[gameID].kingUpdate = true

		if isBlackInCheck(gameID) == true {
			undoMove(x, y, x+1, y, "bK", capturedPiece, gameID)
		}else{
			undoMove(x, y, x+1, y, "bK", capturedPiece, gameID)
			return false
		}
		
	}

	if x+1 <= 7 && y+1 <= 7 && (Verify.AllTables[gameID].ChessBoard[x+1][y+1])[0:1] != "b" { //bottom right

		capturedPiece := makeMove(x, y, x+1, y+1, "bK", gameID)

		Verify.AllTables[gameID].blackOldX = x //storing king information
		Verify.AllTables[gameID].blackOldY = y

		Verify.AllTables[gameID].blackKingX = x + 1
		Verify.AllTables[gameID].blackKingY = y + 1
		Verify.AllTables[gameID].kingUpdate = true

		if isBlackInCheck(gameID) == true {
			undoMove(x, y, x+1, y+1, "bK", capturedPiece, gameID)
		}else{
			undoMove(x, y, x+1, y+1, "bK", capturedPiece, gameID)
			return false
		}
		
	}

	return true
}
