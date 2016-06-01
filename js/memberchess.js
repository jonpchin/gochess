if (!window.WebSocket){
	console.log("Your browser doesn't support websockets. Please use the latest version of Firefox, Chrome, IE, Opera or Edge");
	$('#checkwebsocket').html("Your browser doesn't support websockets. Please use the latest version of Firefox, Chrome, IE, Opera or Microsoft Edge.");
}
var wsuri = "wss://localhost:443/chess";
var sock; 
var matchID;
var moveSound = new Audio('../sound/chessmove.mp3');
var gameSound = new Audio('../sound/startgame.mp3');

//global array of FEN strings, used when pressing back button
var totalFEN = []

var whiteClock = new Tock({
	countdown: true,
	interval: 1000,
	callback: function () {
        if(user === WhiteSide){
	        $('#bottomtime').val(whiteClock.msToTime(whiteClock.lap()));	
	    }
		else{
			$('#toptime').val(whiteClock.msToTime(whiteClock.lap()));	
		}  
    }

});

var blackClock = new Tock({
	countdown: true,
	interval: 1000,
	callback: function () {
    	if(user === BlackSide ){
	        $('#bottomtime').val(blackClock.msToTime(blackClock.lap()));	
	    }
		else{
			$('#toptime').val(blackClock.msToTime(blackClock.lap()));	
		}
    }
 
});
//getting user preferences
var toggleSound = getCookie("sound");
var toggleChat = getCookie("chat");

//store initial game time and use for rematch
var timeGet;

window.onload = function() {
	
	var whiteRating;
	var blackRating;
	//timeGet is a global variable and is not located here
	var gameDate;
	var gameResult;
	
	//hide button on initial load
    $('#rematchButton').hide();
	
	//action listener for exporting game to PGN file
	document.getElementById('exportPGN').onclick = function(){
		if(gameResult === "0"){ //black won
			gameResult = "0-1";
		}
		else if(gameResult === "1"){ //white won
			gameResult = "1-0";
		}
		else{ //game is a draw
			gameResult = "1/2-1/2";
		}
		game.header('Site', "Go Play Chess", 'Date', gameDate, 'White', WhiteSide, 'Black', BlackSide, 'Result', gameResult, 'WhiteElo', whiteRating, 'BlackElo', blackRating, 'TimeControl', timeGet);
		var pgn = game.pgn();
		var fileName = WhiteSide + "vs" + BlackSide + ".pgn";
		download(pgn, fileName, "application/x-chess-pgn");
	}	
	
	user = document.getElementById('user').value;
	
	var token = parseUrl();
	var reviewMoves = token.moves;
	
	if (typeof reviewMoves !== "undefined"){
		
		var whiteName = token.white;
		var blackName = token.black;
		var whiteR = token.whiteRating;
		var blackR = token.blackRating;
		
		//no var needed here as its referencing a var timeGet earlier on
		timeGet = token.time;
		gameDate = token.date;		
		gameResult = token.result;
		
		//assigning global variables which will be used for pgn export
		WhiteSide = whiteName;
		BlackSide = blackName;
		whiteRating = whiteR;
		blackRating = blackR;
			
		document.getElementById("abortButton").disabled = true;
		document.getElementById("drawButton").disabled = true;
		document.getElementById("resignButton").disabled = true;
		document.getElementById("bottomtime").value = timeGet + ":00"; //setting up name and time of player when they are going over game
		document.getElementById("toptime").value =	timeGet + ":00";
		document.getElementById("bottom").innerHTML = "W: " + whiteName + "(" + whiteR + ")";
		document.getElementById("top").innerHTML = "B: " + blackName + "(" + blackR + ")";
		
		var review = JSON.parse(reviewMoves);
		var length = 0;
		if(review !== null){
			length = review.length;
		}
		//always push the default starting position
		totalFEN.push("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1");

		for(var i=0; i<length; i++){
			
			var move = game.move({
			    from: review[i].S,
			    to: review[i].T,
			    promotion: 'q' // NOTE: always promote to a queen for example simplicity
		    });	

		    updateStatus();
			
			board.move(review[i].S + "-" + review[i].T);
			moves.push([ review[i].S, review[i].T ]);
			//pushing on next FEN string
			var fenString = game.fen();
			totalFEN.push(fenString);

			moveCounter++;
		}
		onSnapEnd();
		reviewGame = true;
		return; //prevents game from loading if game is being reviewed	
	}
	//hide export PGN button and add favorites button as player is not reviewing a game
	$('#exportPGN').hide();
	sock = new WebSocket(wsuri);

    sock.onopen = function() {

		document.getElementById('textbox').innerHTML = "";
		var message = {
			Type: "chat_private",
			Name: user,
			Text: "has joined the chess room."
		}
	    sock.send(JSON.stringify(message));
		
		
		//checks to see if a game exist for the player in order to resume the game
		var checkGame = {
			Type: "chess_game",
			Name: user
		}
		sock.send(JSON.stringify(checkGame));
		
    }
	
    sock.onclose = function(e) {
		var currentdate = new Date(); 
		var datetime =  + currentdate.getHours() + ":"  
             			+ currentdate.getMinutes() + ":" 
             			+ currentdate.getSeconds();
		document.getElementById('textbox').innerHTML += (datetime + " " + "Connection lost. Please refresh to reconnect if you are in a game." + '\n');
    }

    sock.onmessage = function(e) {
		var json = JSON.parse(e.data);

		switch (json.Type) {
			case "send_move":
				// see if the move is legal
			    var move = game.move({
				    from: json.Source,
				    to: json.Target,
				    promotion: 'q' // NOTE: always promote to a queen for example simplicity
			    });
			
			    // illegal move
			    if (move === null){
					return;
				}
	
			    updateStatus();
				
				//adding moves to array so player can go forward and backward through chess game
				moves.push([json.Source, json.Target])
				//if draw button is Accept draw then make it say offer draw
				if (document.getElementById("drawButton").value === "Accept Draw"){
					document.getElementById("drawButton").value= "Offer Draw";
				}			
	
				//disables abort button
				if($('#abortButton').is(':disabled') === false && moves.length >= 1){
					document.getElementById("abortButton").disabled = true; 
					
				}
				board.move(json.Source + "-" + json.Target);
				moveCounter++;
				//make the premove that was stored if its stored
				if(preMoveYes === true){
					var result = onDrop(srcPreMove, targetPreMove)
					if(result === 'snapback'){
						//return
					}
					else{
						board.move(srcPreMove + "-" + targetPreMove);	
					}
					preMoveYes = false;
					removeHighlights('color');		
				}
	
				if(user === WhiteSide){
//					document.getElementById("bottomtime").value = json.WhiteMinutes + ":" + json.WhiteSeconds;
				//	document.getElementById("toptime").value = json.BlackMinutes + ":" + json.BlackSeconds;	
				}
				else if(user === BlackSide){
//					document.getElementById("bottomtime").value = json.BlackMinutes + ":" + json.BlackSeconds;
				//	document.getElementById("toptime").value = json.WhiteMinutes + ":" + json.WhiteSeconds;
				}		
				onSnapEnd();
				whiteClock.pause();
	            blackClock.pause();
				if(toggleSound !== "false"){
					moveSound.play();
				}
				break;
			
			case "chat_private":
				if(toggleChat !== "false"){
					var datetime = timeStamp();
					document.getElementById('textbox').innerHTML += (datetime + " " + json.Name +": " + json.Text + '\n');
				}
				break;
				
			case "chess_game":
				//storing matchID in global variable used in sending moves for verification
				matchID = json.ID;
				WhiteSide = json.WhitePlayer;
				BlackSide = json.BlackPlayer;

				//if its move zero store time control so it can be used in rematch
				timeGet = json.TimeControl;
				
				//if the game moves are not null then restore the moves back
				if(json.GameMoves != null){
					//disables abort button
					document.getElementById("abortButton").disabled = true;
					
					var length = json.GameMoves.length				
					
					for(var i=0 ; i<length; i++){
						
						var move = game.move({
						    from: json.GameMoves[i].S,
						    to: json.GameMoves[i].T,
						    promotion: 'q' // NOTE: always promote to a queen for example simplicity
					     });
					
			
					    updateStatus();
						
						board.move(json.GameMoves[i].S + "-" + json.GameMoves[i].T);
						moves.push([ json.GameMoves[i].S, json.GameMoves[i].T ]);
						moveCounter++;
					}
					onSnapEnd(); //used to update castling and en passent positions
			
				}
				else{
					if(toggleSound !== "false"){
						//play sound when game starts
						gameSound.play();
					}
					chessGameOver = false; 
					
					
					//reset game position to brand new game, used for rematch
					board.position('rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1');
					board.orientation('white');
					game.reset();		
					
					document.getElementById("abortButton").disabled = false;
					document.getElementById("drawButton").disabled = false;
					document.getElementById("resignButton").disabled = false;
					
				}
				
				//formating time for clock
				json.WhiteMinutes = json.WhiteMinutes < 10 ? "0" + json.WhiteMinutes : json.WhiteMinutes;
				json.WhiteSeconds = json.WhiteSeconds < 10 ? "0" + json.WhiteSeconds : json.WhiteSeconds;
	            
				json.BlackMinutes = json.BlackMinutes < 10 ? "0" + json.BlackMinutes : json.BlackMinutes;
				json.BlackSeconds = json.BlackSeconds < 10 ? "0" + json.BlackSeconds : json.BlackSeconds;
				
				if (user === json.WhitePlayer){
					
					document.getElementById("bottom").innerHTML = "W: " + json.WhitePlayer + "(" + json.WhiteRating +")";
					document.getElementById("top").innerHTML = "B: " + json.BlackPlayer + "(" + json.BlackRating +")";
					
				}
				else{
					//flips board white on top black on bottom
					$('#flipOrientationBtn').click();
					document.getElementById("bottom").innerHTML = "B: " + json.BlackPlayer + "(" + json.BlackRating +")";
					document.getElementById("top").innerHTML = "W: " + json.WhitePlayer + "(" + json.WhiteRating +")";
	
				}
				
				if(json.Status === "White"){
					if(user === json.WhitePlayer){
						document.getElementById("bottomtime").value = json.WhiteMinutes + ":" + json.WhiteSeconds;
						whiteClock.start($('#bottomtime').val());
					    document.getElementById("toptime").value = json.BlackMinutes + ":" + json.BlackSeconds;
						blackClock.start($('#toptime').val());
						blackClock.pause();
					}
					else{
						document.getElementById("toptime").value = json.WhiteMinutes + ":" + json.WhiteSeconds;
						whiteClock.start($('#toptime').val());
						document.getElementById("bottomtime").value = json.BlackMinutes + ":" + json.BlackSeconds;
						blackClock.start($('#bottomtime').val());
						blackClock.pause();
					}
					
				}
				//else if (json.Status === "Black)
				else {
					if(user === json.WhitePlayer){
						document.getElementById("toptime").value = json.BlackMinutes + ":" + json.BlackSeconds;
						blackClock.start($('#toptime').val());
						document.getElementById("bottomtime").value = json.WhiteMinutes + ":" + json.WhiteSeconds;
						whiteClock.start($('#bottomtime').val());
						whiteClock.pause();
					}
					else{
						document.getElementById("bottomtime").value = json.BlackMinutes + ":" + json.BlackSeconds;
						blackClock.start($('#bottomtime').val());
						document.getElementById("toptime").value = json.WhiteMinutes + ":" + json.WhiteSeconds;
						whiteClock.start($('#toptime').val());
						whiteClock.pause();
					}				
				}
				break;

			case "offer_draw":
				var datetime = timeStamp();
				document.getElementById('textbox').innerHTML += (datetime + " " + json.Name +" offers you a draw." + '\n');
				document.getElementById("drawButton").value= "Accept Draw";
				break;
			
			case "accept_draw":
				var datetime = timeStamp();
				document.getElementById('textbox').innerHTML += (datetime + " Game drawn by mututal agreement." + '\n');
				//prvents players from moving pieces now that game is over
				gameOver();
				break;
				
			case "draw_game":
				var datetime = timeStamp();
				document.getElementById('textbox').innerHTML += (datetime + " Game is a draw." + '\n');
				//prvents players from moving pieces now that game is over
				gameOver();
				break;
				
			case "cancel_draw":
				var datetime = timeStamp();
				document.getElementById('textbox').innerHTML += (datetime + " Draw offer declined." + '\n');
				document.getElementById("drawButton").disabled = false; 
				document.getElementById("drawButton").value= "Offer Draw";
				break;
			
			case "resign":
				var datetime = timeStamp();
				document.getElementById('textbox').innerHTML += (datetime + " " + json.Name + " resigned." + '\n');
				gameOver();
				break;
			case "leave":
				var datetime = timeStamp();
				document.getElementById('textbox').innerHTML += (datetime + " " + json.Text + '\n');
				break;
			case "abort_game":
				var datetime = timeStamp();
				document.getElementById('textbox').innerHTML += (datetime + " " + json.Name +" has aborted the game. "+ '\n');
				gameOver();
				break;
			
			case "game_over":
				var datetime = timeStamp();
				
				if(json.Status === "Black won on time"){
					document.getElementById('textbox').innerHTML += (datetime + " Black won on time."+ '\n');
				}
				else if(json.Status === "White won on time"){
					document.getElementById('textbox').innerHTML += (datetime + " White won on time."+ '\n');
				}
				else if(json.Status === "White"){
					document.getElementById('textbox').innerHTML += (datetime + " White is in checkmate. Black won!"+ '\n');
				}
				else if(json.Status === "Black"){
					document.getElementById('textbox').innerHTML += (datetime + " Black is in checkmate. White won!"+ '\n');
				}
				else{
					console.log("Invalid game status line 243");
				}
				gameOver();
				break;
				
			case "rating":
				var datetime = timeStamp();
				if(user === WhiteSide){
					document.getElementById('textbox').innerHTML += (datetime + " Your new rating is " + json.WhiteRating + '\n');
				}else{
					document.getElementById('textbox').innerHTML += (datetime + " Your new rating is " + json.BlackRating + '\n');
				}
				
				break;
			case "rematch":
				var datetime = timeStamp();
				document.getElementById('textbox').innerHTML += (datetime + " Your opponent offers you a rematch." + '\n');
				document.getElementById('rematchButton').value = "Accept Rematch";
				break;
			case "match_three":
				var datetime = timeStamp();
				document.getElementById('textbox').innerHTML += (datetime + " Your can only have a max of three pending matches." + '\n');
				break;
	
			case "massMessage":
				var datetime = timeStamp();
				document.getElementById('textbox').innerHTML += (datetime + " " + json.Text + '\n');
				document.getElementById('textbox').innerHTML += (datetime + " Game adjourned." + '\n');
				document.getElementById("abortButton").disabled = true;
				document.getElementById("drawButton").disabled = true;
				document.getElementById("resignButton").disabled = true;
	
				whiteClock.stop();
				blackClock.stop();
				chessGameOver = true; 
				if(toggleSound !== "false"){
					gameSound.play();
				}
				break;
			default:
				console.log("Unknown API type json.Type is " + json.Type);
		}

    }
	document.getElementById('sendMessage').onclick = function(){
	    var msg = document.getElementById('message').value;
		
		var message = {
			Type: "chat_private",
			Name: user,
			Text: msg
		}
	    sock.send(JSON.stringify(message));
		//clears input box so next message can be typed without having to clear
		document.getElementById('message').value = "";
		$('#message').focus();
	}
	//allows player to quit the game before move two
	document.getElementById('abortButton').onclick = function(){
		
		var message = {
			Type: "abort_game",
			Name: user,
			ID: matchID
		}
	    sock.send(JSON.stringify(message));	
	}
	
	//a player can only offer draw once per turn
	//allows player to quit the game before move two, this button is also used to accept draws
	document.getElementById('drawButton').onclick = function(){
		
		if (document.getElementById("drawButton").value == "Accept Draw"){
			var message = {
				Type: "accept_draw",
				Name: user,
				ID: matchID
			}
		    sock.send(JSON.stringify(message));
		}
		else{
			var message = {
				Type: "offer_draw",
				Name: user,
				ID: matchID
			}
		    sock.send(JSON.stringify(message));
		}
		
		
		var datetime = timeStamp();
		document.getElementById('textbox').innerHTML += (datetime + " You offer your opponent a draw." + '\n');
		document.getElementById("drawButton").disabled = true; 
		
	}
	
	document.getElementById('resignButton').onclick = function(){
		if (window.confirm("Do you really want to resign?")) { 
  			var message = {
				Type: "resign",
				Name: user,
				ID: matchID
			}
		    sock.send(JSON.stringify(message));
			
		}
	}
	

};
function finishGame(color){ //color is the color player that is checkmated
 
	var message = {
			Type: "game_over",
			Name: user,	
			ID: matchID,
			Status: color
	}
	sock.send(JSON.stringify(message));

}
	
//game is now drawn
function drawGame(){
	var message = {
			Type: "draw_game",
			Name: user,	
			ID: matchID
	}
	sock.send(JSON.stringify(message));
}

//this function is called from boardsetup.js
function sendMove(src, dest){
		
	var message = {
		Type: "send_move",
		Name: user,	
		ID: matchID,
		Source: src,
		Target: dest
	}
    sock.send(JSON.stringify(message));

	whiteClock.pause();
	blackClock.pause();
	if(toggleSound !== "false"){
		moveSound.play();
	}
}
//go forward one move
document.getElementById('goForward').onclick = function(){

	if(moveCounter < moves.length){
		
		moveCounter++;
	}
	//make a global array and iterate forwards through the global array when going forward
	board.position(totalFEN[moveCounter]);	
	
} 
$('#goBack').on('click', function() {
	if(moveCounter > 0){
		
		moveCounter--;
		
	}
	//make a global array and iterate backwards through the global array when going back
	board.position(totalFEN[moveCounter]);
	
});
//move forward to last move
document.getElementById('goEnd').onclick = function(){

	for(var i=moveCounter; i<moves.length; i++){

		board.move(moves[i][0] + "-" + moves[i][1]);
	}
	moveCounter = moves.length;	
	
} 

//offers player a rematch or accepts it if the other player offers
document.getElementById('rematchButton').onclick = function(){
	
	var value = document.getElementById('rematchButton').value;
	//the opponent
	var fighter;
	
	if(user === WhiteSide){
		fighter = BlackSide;
	}
	else{
		fighter = WhiteSide;
	}
	
	if(value === "Rematch"){
		var message = {
			Type: "rematch",
			Name: user,	
			Opponent: fighter,
			TimeControl: timeGet
		}
		var datetime = timeStamp();
		document.getElementById('textbox').innerHTML += (datetime + " You offer your opponent a rematch." + '\n');
	    
	}
	else{ //else value === "Accept Rematch"
		var message = {
			Type: "accept_rematch",
			Name: user,	
			Opponent: fighter,
			TimeControl: timeGet
		}
		document.getElementById('rematchButton').value = "Rematch";

	}
	//hiding button after click to prevent rematch abuse
	$('#rematchButton').hide();
	sock.send(JSON.stringify(message));

} 

	
$('#message').keypress(function(event) {
    if (event.which === 13) {
       var msg = document.getElementById('message').value;
	  
	   $('#sendMessage').click();	
    }
});

//returns timestamp
function timeStamp(){
	var currentdate = new Date(); 
	var datetime =  + currentdate.getHours() + ":"  
              			+ currentdate.getMinutes() + ":" 
              			+ currentdate.getSeconds();
	return datetime;
}

//when game is over abort, draw and resign button should be disabled
function gameOver(){
	document.getElementById("abortButton").disabled = true;
	document.getElementById("drawButton").disabled = true;
	document.getElementById("resignButton").disabled = true;
	//show rematch button now game is over
    $('#rematchButton').show();

	whiteClock.stop();
	blackClock.stop();
	chessGameOver = true; 
	if(toggleSound !== "false"){
		gameSound.play();
	}
}

function parseUrl() { //fetches all variables in url and returns them in a json struct
  var query = location.search.substr(1);
  var result = {};
  query.split("&").forEach(function(part) {
    var item = part.split("=");
    result[item[0]] = decodeURIComponent(item[1]);
  });
  return result;
}

function detectMobile(){ //tries to detect if user is using a mobile device
	
	var screenWidth = screen.width;	
	if(screenWidth <= 900){

		console.log("mobile device detected...adjusting board size and layout");
		document.getElementById("chatleft").style.display = "none";
		document.getElementById("notation").style.display = "none";	
	}
}

detectMobile(); //calls function
 