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
var chessPieceTheme = getCookie("pieceTheme");

function defaultTheme(){
	if (chessPieceTheme === ""){
		chessPieceTheme = "wikipedia"
	}
}
defaultTheme();

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

var updateStatus = function(moveString) {

	if(document.getElementById("cinnamonRadioButton").checked === true){
		cinnamonEngineGo()
	}else{ // then send move to stockfish
		stockfish.postMessage("position startpos moves " + moveString);
		stockfish.postMessage("go movetime " + document.getElementById('thinkTime').value);
	}
	var status = '';

	status = checkGameStatus(getMoveColor());

	statusEl.html(status);
	fenEl.html(game.fen());
	pgnEl.html(game.pgn());

	return status;
};

function getMoveColor(){
	var moveColor = 'White';
	if (game.turn() === 'b') {
		moveColor = 'Black';
	}
	return moveColor;
}

function checkGameStatus(moveColor){
	var status = "";
	if (game.in_checkmate() === true) {
		status = 'Game over, ' + moveColor + ' is in checkmate.';
		sendNotification(status);
	} else if (game.in_draw() === true) { // draw?
		status = 'Game over, drawn position';
		sendNotification(status);
	} else { // game still on
		status = moveColor + ' to move';

		// check?
		if (game.in_check() === true) {
			status += ', ' + moveColor + ' is in check';
		}
	}
}

var onSnapEnd = function() {
	board.position(game.fen());
};

function cinnamonEngineGo(){
	
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
	++moveCounter;
	checkGameStatus(getMoveColor());
}

// updates the board with the move stock fish had in mine
// @param src the starting location of the square the piece/pawn is from
// @param tar the target location of the square the piece/pawn is moving too
function stockFishEngineGo(src, tar){
	var move = game.move({
		from: src,
		to: tar,
		promotion: 'q' // NOTE: always promote to a queen for example simplicity
	});

	board.position(game.fen());
	totalFEN.push(game.fen());
	totalPGN.push(game.pgn());
	++moveCounter;	
	checkGameStatus(getMoveColor());
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

var onDrop = function(source, target) {
	
	$.ajax({
  		url: 'checkInGame',
   		type: 'post',
		data : { 'user': document.getElementById('user').value},
   		success : function(data) {		
			if(data === "inGame"){
				alert("Engine use when you are playing a game against a real person is not allowed!");
				window.location = "/memberHome";
			}
   		},
	});

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

	makeEngineMove();

	totalFEN.push(game.fen());
	totalPGN.push(game.pgn());
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
	pieceTheme: '../img/chesspieces/'+ chessPieceTheme +'/{piece}.png'
};
board = new ChessBoard('board', cfg);

// force the computer to make a move thus switching sides of the game
document.getElementById('forceMoveButton').onclick = function(){
	if(computer === 'b'){
		computer = 'w';
	}else{
		computer = 'b';
	}
	// need to make sure its the most current move before forcing engine to move
	document.getElementById('goEnd').click();
	makeEngineMove();
	board.position(game.fen());
}

//checks to see which engine is turned on and makes a move for that engine
function makeEngineMove(){
	if(document.getElementById("cinnamonRadioButton").checked === true){
		updateStatus("");
	}else{
		var gameHistory = game.history({verbose: true});
		var length = gameHistory.length;
		var moveString = "";
		for(var i = 0; i<length; ++i){
			moveString += (gameHistory[i].from + gameHistory[i].to + " ");
		}
		updateStatus(moveString);
	}
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

// action listener for board flip
$('#flipOrientationBtn').on('click', board.flip);
//always push the default starting position
totalFEN.push("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1");
totalPGN.push("");

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

// Ask user to support notifications
function askNotificationPermission(){
	
	if (!Notification) {
		alert('Desktop notifications not available in your browser. Try Chromium.'); 
		return;
	}

	if (Notification.permission !== "granted"){
		Notification.requestPermission();
	}	
}
askNotificationPermission();

function sendNotification(message) {
	
	// check if the browser supports notifications
	if (!("Notification" in window)) {
		alert("This browser does not support system notifications");
	}
  
	// check whether notification permissions have already been granted
	else if (Notification.permission === "granted") {
		// create a notification
		var notification = new Notification(message);
	}
  
	// Otherwise, ask the user for permission
	else if (Notification.permission !== 'denied') {
		Notification.requestPermission(function (permission) {
			// If the user accepts, create a notification
			if (permission === "granted") {
				var notification = new Notification(message);
			}
		});
	}
	else{
		document.getElementById('textbox').innerHTML += (message + '\n');
	}
}

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
}

var stockfish;

function startStockFish(){
	
	newGameButton.click();
	stockfish = new Worker("../third-party/js/stockfish.js");

	stockfish.onmessage = function(event) {
		//NOTE: Web Workers wrap the response in an object.
		var data = event.data;
		if(data){
			if(data.startsWith("bestmove")){
				stockFishEngineGo(data.substring(9, 11), data.substring(11, 13));
			}
		}	
	};
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