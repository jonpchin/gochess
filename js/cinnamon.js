cinnamonCommand = Module.cwrap('command', 'string', ['string','string']);

// milliseconds for engine to think before making a move
cinnamonCommand("setMaxTimeMillsec", 1000);
// the color the computer is playing, switch this to switch sides for computer
var computer = 'b';

var init = function() {

	//--- start example JS ---
	var board,
	game = new Chess(),
	statusEl = $('#status'),
	fenEl = $('#fen'),
	pgnEl = $('#pgn');

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
	var onDrop = function(source, target) {
		removeGreySquares();

		// see if the move is legal
		var move = game.move({
			from: source,
			to: target,
			promotion: 'q' // NOTE: always promote to a queen for example simplicity
		});

		// illegal move
		if (move === null) return 'snapback';
		updateStatus();
	};

	// update the board position after the piece snap 
	// for castling, en passant, pawn promotion
	var onSnapEnd = function() {
		board.position(game.fen());
	};

	function engineGo(){
		
		cinnamonCommand("position",game.fen());
		var move=cinnamonCommand("go","");
;
		var from=move.substring(0,2);
		var to=move.substring(2,4);
		var move = game.move({
			from: from,
			to: to,
			promotion: 'q' // NOTE: always promote to a queen for example simplicity
		});
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

	var updateStatus = function() {

		if (game.turn() === computer) {
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
	};

	var cfg = {
		draggable: true,
		position: 'start',
		onDragStart: onDragStart,
		onDrop: onDrop,
		moveSpeed: 'slow',
		onMouseoutSquare: onMouseoutSquare,
		onMouseoverSquare: onMouseoverSquare,
		onSnapEnd: onSnapEnd
	};
	board = new ChessBoard('board', cfg);

	updateStatus();
	$('#startPositionBtn').on('click', function() {
		board.destroy();
		$(document).ready(init);
	});

	// action listener for board flip
	$('#flipOrientationBtn').on('click', board.flip);

	// force the computer to make a move thus switching sides of the game
	document.getElementById('forceMoveButton').onclick = function(){
		if(computer === 'b'){
			computer = 'w';
		}else{
			computer = 'b';
		}
		updateStatus();
		board.position(game.fen());
	}
};

$(document).ready(init);

document.getElementById('setThinkTime').onclick = function() {
	cinnamonCommand("setMaxTimeMillsec", document.getElementById('thinkTime').value);
}

function detectMobile(){ //tries to detect if user is using a mobile device
		
	if(screen.width <= 900){
		console.log("mobile device detected...adjusting board size and layout");
		document.getElementById("chatleft").style.display = "none";
		document.getElementById("notation").style.display = "none";	
	}
}

detectMobile(); //calls function
