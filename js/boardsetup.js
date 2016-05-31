var board,
game = new Chess(),
boardEl = $('#board'),
statusEl = $('#status'),
fenEl = $('#fen'),
pgnEl = $('#pgn');


//used to store whether or not the game is over due to resignation or mututal draw agreement
var chessGameOver;
//stores the moves of all the games, use length of moves to determine if abort button should be disabled
var moves = [];
//used to store what move the game is on, their will be double moves in total one for black and one for white
var moveCounter = 0;
var user; 
var WhiteSide;
var BlackSide;
//whether or not a prove move is stored
var preMoveYes = false;
//stores premove in string
var srcPreMove;
var targetPreMove;
//user preferences
var togglePremove = getCookie("premove");
var pieceTheme = getCookie("pieceTheme");
//used to check if a player is viewing a game
var reviewGame = false;

function defaultTheme(){
	if (pieceTheme === ""){
		pieceTheme = "wikipedia"
	}
}
defaultTheme();

// do not pick up pieces if the game is over
// only pick up pieces for the side to move
var onDragStart = function(source, piece, position, orientation) {
	//onclick premove should be undone	
	if (preMoveYes === true){
		removeHighlights('color');
  		preMoveYes = false;	
	}

	if (game.game_over() === true ||
		chessGameOver === true ||
    	(WhiteSide === user && piece.search(/^b/) !== -1) ||
    	(BlackSide === user && piece.search(/^w/) !== -1)) {
    	return false;
    }
	
    
};

var onDrop = function(source, target) {
	//only allow premove if user enabled in preferences, by default premove is enabled
	if(togglePremove !== "false"){
		if( (game.turn() === 'w' && BlackSide === user) || (game.turn() === 'b' && WhiteSide === user)   ){
			preMoveYes = true;
			srcPreMove = source;
			targetPreMove = target;
			boardEl.find('.square-' + source).addClass('highlight-color'); //adds premove color
			boardEl.find('.square-' + target).addClass('highlight-color');
			return;
		}
	}
    
	// see if the move is legal
	var move = game.move({
    	from: source,
    	to: target,
    	promotion: 'q' // NOTE: always promote to a queen for example simplicity
  	});
	
  	// illegal move
	if (move === null) return 'snapback';

  
	//used to store players own move, moves array is stored in memberchess.js
	moves.push([source, target])
 	moveCounter++;
	//starting player's clock on move 1
  
	sendMove(source, target)

	updateStatus();	
};

// update the board position after the piece snap 
// for castling, en passant, pawn promotion
var onSnapEnd = function() {
	if(reviewGame === false){
		board.position(game.fen());
	}
};

var updateStatus = function() {
	var status = '';
	
	var moveColor = 'White';
	if (game.turn() === 'b') {
	  moveColor = 'Black';
	}

	// checkmate?
	if (game.in_checkmate() === true) {
		status = 'Game over, ' + moveColor + ' is in checkmate.';
		if(WhiteSide === user){ // prevents game over duplication being sent to server
			finishGame(moveColor); //function call located in memberchess.js
		}
	
	}
	else if (game.in_draw() === true) { // draw, todo: need to message server when this is triggered
		status = 'Game over, drawn position';
		if(WhiteSide === user){ // prevents game over duplication being sent to server
			drawGame(); //function call located in memberchess.js
		}
	}
  	else {   // game still on
	    status = moveColor + ' to move';
	
	    // check?
	    if (game.in_check() === true) {
	      status += ', ' + moveColor + ' is in check';
	    }
	}

	statusEl.html(status);
	if(reviewGame === false){
		fenEl.html(game.fen());
		pgnEl.html(game.pgn());
	}
	
};
//removes premove coor
var removeHighlights = function(color) {
  boardEl.find('.square-55d63')
    .removeClass('highlight-' + color);
};

var cfg = {
	draggable: true,
	position: 'start',
	onDragStart: onDragStart,
	onDrop: onDrop,
	onSnapEnd: onSnapEnd,
	pieceTheme: '../img/chesspieces/'+ pieceTheme +'/{piece}.png'
};


board = ChessBoard('board', cfg);

updateStatus();

$('#flipOrientationBtn').on('click', board.flip);

document.getElementById('goStart').onclick = function(){
	
	board.position('rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR');
	moveCounter = 0;
		
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