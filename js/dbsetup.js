var board,
game = new Chess(),
boardEl = $('#board'),
statusEl = $('#status'),
// fenEl = $('#fen'), FEN string is not being used
pgnEl = $('#pgn');

//global array of FEN strings, PGN, and statuus used when reviewing game
var totalFEN = [];
var totalPGN = [];
var totalStatus = []
//used to store what move the game is on, their will be double moves in total one for black and one for white
var moveCounter = 0; 

//user preferences
var pieceTheme = getCookie("pieceTheme");

function defaultTheme(){
	if (pieceTheme === ""){
		pieceTheme = "wikipedia"
	}
}
defaultTheme();

// returns status of current game
var updateStatus = function() {
	var status = '';
	
	var moveColor = 'White';
	if (game.turn() === 'b') {
	  moveColor = 'Black';
	}

	// checkmate?
	if (game.in_checkmate() === true) {
		status = 'Game over, ' + moveColor + ' is in checkmate.';
	}
	else if (game.in_draw() === true) { // draw, todo: need to message server when this is triggered
		status = 'Game over, drawn position';
	}
  	else {   // game still on
	    status = moveColor + ' to move';
	
	    // check?
	    if (game.in_check() === true) {
	      status += ', ' + moveColor + ' is in check';
	    }
	}
	return status;
};

var setStatusAndPGN = function(status, pgn){
	statusEl.html(status);
	//	fenEl.html(game.fen()); FEN string is not being used
	pgnEl.html(pgn);
}

var cfg = {
	draggable: false,
	position: 'start',
	pieceTheme: '../img/chesspieces/'+ pieceTheme +'/{piece}.png'
};

board = ChessBoard('board', cfg);

// defaults the status of the game and pgn
setStatusAndPGN("White to move", "")

$('#flipOrientationBtn').on('click', board.flip);

document.getElementById('goStart').onclick = function(){
	
	board.position('rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR');
	moveCounter = 0;
	setStatusAndPGN("White to move", "")
}

function getCookie(cname) { //gets cookies value
    var name = cname + "=";
    var ca = document.cookie.split(';');
    for(var i=0; i<ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0)==' ') c = c.substring(1);
        if (c.indexOf(name) == 0) return c.substring(name.length,c.length);
    }
    return "";
}

function setBoardColor(){
    var boardColor = getCookie("boardColor");
	switch(boardColor){
		case "blueColor":
			$('.black-3c85d').css({"background-color":"#4682B4"});
			$('.black-3c85d').css({"color":"#B0E0E6"});
			$('.white-1e1d7').css({"background-color":"#B0E0E6"});
			$('.white-1e1d7').css({"color":"#4682B4"});
			break;

		case "greenColor":
			$('.black-3c85d').css({"background-color":"#008000"});
			$('.black-3c85d').css({"color":"#90EE90"});
			$('.white-1e1d7').css({"background-color":"#90EE90"});
			$('.white-1e1d7').css({"color":"#008000"});
			break;

		case "greyColor":
			$('.black-3c85d').css({"background-color":"#696969"});
			$('.black-3c85d').css({"color":"#D3D3D3"});
			$('.white-1e1d7').css({"background-color":"#D3D3D3"});
			$('.white-1e1d7').css({"color":"#696969"});
			break;

		case "pinkColor":
			$('.black-3c85d').css({"background-color":"#FF69B4"});
			$('.black-3c85d').css({"color":"#FFC0CB"});
			$('.white-1e1d7').css({"background-color":"#FFC0CB"});
			$('.white-1e1d7').css({"color":"#FF69B4"});
			break;

		default: //default color
			$('.black-3c85d').css({"background-color":"#b58863"});
			$('.black-3c85d').css({"color":"#f0d9b5"});
	
			$('.white-1e1d7').css({"background-color":"#f0d9b5"});
			$('.white-1e1d7').css({"color":"#b58863"});
	}
}
setBoardColor();