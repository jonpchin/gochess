var board,
game = new Chess(),
statusEl = $('#status'),
fenEl = $('#fen'),
pgnEl = $('#pgn');

// the color the computer is playing, switch this to switch sides for computer
var computer = 'b';

//global array of FEN strings, PGN, and statuus used when reviewing game
var totalFEN = [];
var totalPGN = [];
//used to store what move the game is on, their will be double moves in total one for black and one for white
var moveCounter = 0;
var pieceTheme = getCookie("pieceTheme");

var onDragStart = function(source, piece) {
	// do not pick up pieces if the game is over
	// or if it's not that side's turn
	if (game.game_over() === true ||
		(game.turn() === 'w' && piece.search(/^b/) !== -1) ||
		(game.turn() === 'b' && piece.search(/^w/) !== -1)) {
		return false;
	}
};

var removeGreySquares = function() {
	$('#board .square-55d63').css('background', '');
};

var greySquare = function(square) {
	var squareEl = $('#board .square-' + square);

	var background = '#a9a9a9';
	if (squareEl.hasClass('black-3c85d') === true) {
		background = '#696969';
	}

	squareEl.css('background', background);
};

var updateStatus = function(engineType) {

	if (game.in_checkmate()){
		alert('Game over, ' + moveColor + ' is in checkmate.');
		return;
	}

	if (game.turn() === computer && engineType === "cinnamon") {
		engineGo()
	}
	var status = '';

	var moveColor = 'White';
	if (game.turn() === 'b') {
		moveColor = 'Black';
	}

	// checkmate?
	if (game.in_checkmate() === true) {
		status = 'Game over, ' + moveColor + ' is in checkmate.';
	} else if (game.in_draw() === true) { // draw?
		status = 'Game over, drawn position';
	} else { // game still on
		status = moveColor + ' to move';

		// check?
		if (game.in_check() === true) {
			status += ', ' + moveColor + ' is in check';
		}
	}
	statusEl.html(status);
	fenEl.html(game.fen());
	pgnEl.html(game.pgn());

	return status;
};

var onSnapEnd = function() {
	board.position(game.fen());
};

function engineGo(){
	
	cinnamonCommand("setMaxTimeMillsec", document.getElementById('thinkTime').value);
	cinnamonCommand("position",game.fen());
	var move=cinnamonCommand("go","");

	var from=move.substring(0,2);
	var to=move.substring(2,4);
	var move = game.move({
		from: from,
		to: to,
		promotion: 'q' // NOTE: always promote to a queen for example simplicity
	});

	totalFEN.push(game.fen());
	totalPGN.push(game.pgn());
	updateStatus();
	++moveCounter;
}
var onMouseoverSquare = function(square, piece) {
	// get list of possible moves for this square
	var moves = game.moves({
		square: square,
		verbose: true
	});

	// exit if there are no moves available for this square
	if (moves.length === 0){
		return;
	} 

	// highlight the square they moused over
	greySquare(square);

	// highlight the possible squares for this piece
	for (var i = 0; i < moves.length; i++) {
		greySquare(moves[i].to);
	}
};

var onMouseoutSquare = function(square, piece) {
	removeGreySquares();
};

var init = function() {

	//always push the default starting position
	totalFEN.push("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1");
	totalPGN.push("");
	
	var onDrop = function(source, target) {
		//removeGreySquares();

		// see if the move is legal
		var move = game.move({
			from: source,
			to: target,
			promotion: 'q' // NOTE: always promote to a queen for example simplicity
		});

		// illegal move
		if (move === null){
			return 'snapback';
		} 
		totalFEN.push(game.fen());
		totalPGN.push(game.pgn());
		updateStatus("cinnamon");
		++moveCounter;
	};

	// update the board position after the piece snap 
	// for castling, en passant, pawn promotion
	

	var cfg = {
		draggable: true,
		position: 'start',
		onDragStart: onDragStart,
		onDrop: onDrop,
		moveSpeed: 'slow',
		onMouseoutSquare: onMouseoutSquare,
		onMouseoverSquare: onMouseoverSquare,
		onSnapEnd: onSnapEnd,
		pieceTheme: '../img/chesspieces/'+ pieceTheme +'/{piece}.png'
	};
	board = new ChessBoard('board', cfg);

	updateStatus("cinnamon");

	// action listener for board flip
	$('#flipOrientationBtn').on('click', board.flip);

	// force the computer to make a move thus switching sides of the game
	document.getElementById('forceMoveButton').onclick = function(){
		if(computer === 'b'){
			computer = 'w';
		}else{
			computer = 'b';
		}
		// need to make sure its the most current move before forcing engine to move
		document.getElementById('goEnd').click();
		updateStatus("cinnamon");
		board.position(game.fen());
	}

	document.getElementById('goStart').onclick = function(){
		board.position('rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR');
		moveCounter = 0;
		setPGN(totalPGN[moveCounter]);
	}

	//go forward one move
	document.getElementById('goForward').onclick = function(){

		if(moveCounter < totalFEN.length-1){	
			moveCounter++;
		}
		//make a global array and iterate forwards through the global array when going forward
		board.position(totalFEN[moveCounter]);	
		setPGN(totalPGN[moveCounter]);
	} 

	$('#goBack').on('click', function() {
		if(moveCounter > 0){
			
			moveCounter--;
		}
		//make a global array and iterate backwards through the global array when going back
		board.position(totalFEN[moveCounter]);	
		setPGN(totalPGN[moveCounter]);
	});

	//move forward to last move
	document.getElementById('goEnd').onclick = function(){

		for(var i=moveCounter; i<totalFEN.length; i++){
			board.position(totalFEN[i]);
		}
		moveCounter = totalFEN.length-1;
		if(moveCounter>=0){
			setPGN(totalPGN[moveCounter]);
		}
	}

	var setPGN = function(pgn){
		pgnEl.html(pgn);
	}

	document.getElementById('newGameButton').onclick = function(){
		totalFEN = [];
		totalFEN.push("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1");
		totalPGN = [];
		totalPGN.push("");
		moveCounter = 0;
		board.position('rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1');
		board.orientation('white');
		game.reset();
	} 

	//action listener for exporting game to PGN file
	document.getElementById('exportPGN').onclick = function(){
		
		// TODO Fill out parameters
		var gameResult = "???";
		var whiteRating = "????";
		var blackRating = "????"
		var timeGet = "???";
		var gameDate = Date();
		var whitePlayer = document.getElementById('user').value;;
		var blackPlayer = "Cinnamon Computuer";
		if(computer === 'w'){
			whitePlayer =  "Cinnamon Computuer";
			blackPlayer = document.getElementById('user').value;
		}

		game.header('Site', "Go Play Chess", 'Date', gameDate, 'White', whitePlayer, 'Black', blackPlayer, 
			'Result', gameResult, 'WhiteElo', whiteRating, 'BlackElo', blackRating, 'TimeControl', timeGet);

		// second parameter is file name
		download(game.pgn(), whitePlayer + " vs. " + blackPlayer + ".pgn", "application/x-chess-pgn");
	}
};

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

function defaultTheme(){
	if (pieceTheme === ""){
		pieceTheme = "wikipedia"
	}
}
defaultTheme();

function detectMobile(){ //tries to detect if user is using a mobile device
		
	if(screen.width <= 900){
		console.log("mobile device detected...adjusting board size and layout");
		document.getElementById("chatleft").style.display = "none";
		document.getElementById("notation").style.display = "none";	
	}
}

function startCinnamon(){
	//$(document).ready(init);
	newGameButton.click();
	cinnamonCommand = Module.cwrap('command', 'string', ['string','string']);
	// milliseconds for engine to think before making a move
	cinnamonCommand("setMaxTimeMillsec", 1000);
	init();
}


function startStockFish(){
	newGameButton.click();
	var stockfish = new Worker("../third-party-js/stockfish.js");
	stockfish.postMessage("position startpos moves e2e4");

	stockfish.onmessage = function(event) {
		//NOTE: Web Workers wrap the response in an object.
		console.log(event.data ? event.data : event);
	};
	stockfish.postMessage("go movetime 5000");

	var onDrop = function(source, target) {
		//removeGreySquares();

		// see if the move is legal
		var move = game.move({
			from: source,
			to: target,
			promotion: 'q' // NOTE: always promote to a queen for example simplicity
		});
		console.log("result is:");
		console.log(move);
		// illegal move
		if (move === null){
			return 'snapback';
		} 
		totalFEN.push(game.fen());
		totalPGN.push(game.pgn());
		updateStatus("stockfish");
		++moveCounter;
	};

	var cfg = {
		draggable: true,
		position: 'start',
		onDragStart: onDragStart,
		onDrop: onDrop,
		moveSpeed: 'slow',
		onMouseoutSquare: onMouseoutSquare,
		onMouseoverSquare: onMouseoverSquare,
		onSnapEnd: onSnapEnd,
		pieceTheme: '../img/chesspieces/'+ pieceTheme +'/{piece}.png'
	};
	board = new ChessBoard('board', cfg);

	

	//always push the default starting position
	totalFEN.push("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1");
	totalPGN.push("");
}

// starts the engine that radio button is checked for
function startCheckedEngine(){
	if(document.getElementById("cinnamonRadioButton").checked === true){
		startCinnamon();
	}else{
		startStockFish();
	}
}

startCheckedEngine();
detectMobile(); //calls function