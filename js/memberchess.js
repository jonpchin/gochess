if (!window.WebSocket){
	$('#checkwebsocket').html("Your browser doesn't support websockets." +
		"Please use the latest version of Firefox, Chrome, IE, Opera or Microsoft Edge.");
}
var wsuri = "wss://" + window.location.host + "/chess";
var sock; 
var matchID;
var moveSound = new Audio('../sound/chessmove.mp3');
var gameSound = new Audio('../sound/startgame.mp3');

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
var toggleSpectate = getCookie("spectate");

//store initial game time and whether game is rated, used for rematch
var timeGet;
var isRatedRematch;

window.onload = function() {
	var whiteRating;
	var blackRating;
	//timeGet is a global variable and is not located here
	var gameDate;
	var gameResult;
	
	//hide button on initial load
    $('#rematchButton').hide();
	document.getElementById("abortButton").disabled = true;
	document.getElementById("drawButton").disabled = true;
	document.getElementById("resignButton").disabled = true;
	
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
		game.header('Site', "Go Play Chess", 'Date', gameDate, 'White', WhiteSide, 'Black', BlackSide, 
			'Result', gameResult, 'WhiteElo', whiteRating, 'BlackElo', blackRating, 'TimeControl', timeGet);

		// second parameter is file name
		download(game.pgn(), WhiteSide + " vs. " + BlackSide + ".pgn", "application/x-chess-pgn");
	}	
	
	user = document.getElementById('user').value;
	
	//always push the default starting position
	totalFEN.push("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1");
	totalStatus.push("White to move");
	totalPGN.push("");

	var token = parseUrl();
	var reviewMoves = token.moves;
	
	if (typeof reviewMoves !== "undefined"){
		
		reviewGame = true;
		
		var whiteName = token.white;
		var blackName = token.black;
		var whiteR = token.whiteRating;
		var blackR = token.blackRating;
		var countryWhite = token.countryWhite;
		var countryBlack = token.countryBlack;
		if(countryWhite === ""){
			countryWhite = "globe";
		}
		if(countryBlack === ""){
			countryBlack = "globe";
		}
		
		//no var needed here as its referencing a var timeGet earlier on
		timeGet = token.time;
		gameDate = token.date;		
		gameResult = token.result;
		
		//assigning global variables which will be used for pgn export
		WhiteSide = whiteName;
		BlackSide = blackName;
		whiteRating = whiteR;
		blackRating = blackR;

		document.getElementById("bottomtime").value = timeGet + ":00"; //setting up name and time of player when they are going over game
		document.getElementById("toptime").value =	timeGet + ":00";
		document.getElementById("bottom").innerHTML = "W: <img src='../img/flags/" + 
			countryWhite + ".png'><a href='/profile?name=" + whiteName + "'>" +
			whiteName + "</a>(" + whiteR +")";
		document.getElementById("top").innerHTML = "B: <img src='../img/flags/" + 
			countryBlack + ".png'><a href='/profile?name=" + blackName  + "'>" +
			blackName + "</a>(" + blackR +")";			

		var review = JSON.parse(reviewMoves);
		var length = 0;
		if(review !== null){
			length = review.length;
		}

		for(var i=0; i<length; i++){
			
			var move = game.move({
			    from: review[i].S,
			    to: review[i].T,
			    promotion: review[i].P
		    });	
    
			totalFEN.push(game.fen());
			totalPGN.push(game.pgn());
			totalStatus.push(updateStatus());
		}
		return; //prevents game from loading if game is being reviewed	
	}
	//hide export PGN button and add favorites button as player is not reviewing a game
	$('#exportPGN').hide();
	
	sock = new WebSocket(wsuri);	
	$(window).on('beforeunload', function(){
		sock.close();
	});  

    sock.onopen = function() {

		// If a game is being spectated then do not load chess_game
		if(typeof token.spectate !== "undefined"){

				// spectators should not be able to do anything but watch the game
				document.getElementById("message").disabled = true;
				document.getElementById("sendMessage").disabled = true;
				document.getElementById("rematchButton").disabled = true;

			var message = {
				Type:  "spectate_game",
				Name:  user,
				ID:    token.id
			}
			sock.send(JSON.stringify(message));

		}else{
			document.getElementById('textbox').innerHTML = "";
			var message = {
				Type: "chat_private",
				Name: user,
				Text: "has joined the chess room."
			}
			sock.send(JSON.stringify(message));
			
			var chess_game = {
				Type: "chess_game",
				Name: user
			}
			sock.send(JSON.stringify(chess_game));
		}		
    }
	
    sock.onclose = function(e) {
		var currentdate = new Date(); 
		var datetime =  + currentdate.getHours() + ":"  
             			+ currentdate.getMinutes() + ":" 
             			+ currentdate.getSeconds();
		document.getElementById('textbox').innerHTML += (datetime + " " +
			 "Connection lost. Please refresh to reconnect if you are in a game." + '\n');
	}

    sock.onmessage = function(e) {
		var json = JSON.parse(e.data);

		switch (json.Type) {
			case "send_move":
			
				// see if the move is legal
			    var move = game.move({
				    from: json.Source,
				    to: json.Target,
				    promotion: json.Promotion
			    });
			
			    // illegal move
			    if (move === null){
					return;
				}
	
			    var gameStatus = updateStatus();
				
				//adding moves to array so player can go forward and backward through chess game
				var fen = game.fen();
				var pgn = game.pgn();
				totalFEN.push(fen);
				totalPGN.push(pgn);
				totalStatus.push(gameStatus);
				setStatusAndPGN(gameStatus, pgn);

				board.position(fen);
				moveCounter++;

				//if draw button is Accept draw then make it say offer draw
				if (document.getElementById("drawButton").innerHTML === "Accept Draw"){
					document.getElementById("drawButton").innerHTML = "Offer Draw";
				}			
	
				//disables abort button
				if($('#abortButton').is(':disabled') === false && totalFEN.length >= 2){
					document.getElementById("abortButton").disabled = true; 
				}
				
				//make the premove that was stored if its stored
				if(preMoveYes){
					var result = onDrop(srcPreMove, targetPreMove)
					if(result === 'snapback'){
						//return
					}
					else{
						// this game.fen() is different then the one stored in fen variable
						var fenPreMove = game.fen();
						var pgnPreMove = game.pgn();
						var preGameStatus = updateStatus();
						
						board.position(fenPreMove);	
						totalFEN.push(fenPreMove);
						totalPGN.push(pgnPreMove);
						totalStatus.push(preGameStatus);
						setStatusAndPGN(preGameStatus, pgnPreMove);
					}
					preMoveYes = false;
					removeHighlights('color');		
				}
					
				whiteClock.pause();
	            blackClock.pause();
				if(toggleSound !== "false"){
					moveSound.play();
				}
				break;
				
			case "sync_clock":
				if(user === WhiteSide){
					if(json.UpdateWhite){
						document.getElementById("bottomtime").value = json.WhiteMinutes + ":" + 
							json.WhiteSeconds + "." + json.WhiteMilli;
					//	console.log("White: " + json.WhiteMinutes + ":" + json.WhiteSeconds + "." + json.WhiteMilli);
					}
					else{
						document.getElementById("toptime").value = json.BlackMinutes + ":" + 
						json.BlackSeconds + "." + json.BlackMilli;	
					//	console.log("Black: " + json.BlackMinutes + ":" + json.BlackSeconds + "." + json.BlackMilli);
					}
				}
				else if(user === BlackSide){
					if(json.UpdateWhite){
						document.getElementById("toptime").value = json.WhiteMinutes + ":" + json.WhiteSeconds + "." +
							json.WhiteMilli;
					//	console.log("White: " + json.WhiteMinutes + ":" + json.WhiteSeconds  + "." + json.WhiteMilli);
					}
					else{
						document.getElementById("bottomtime").value = json.BlackMinutes + ":" + json.BlackSeconds + "." + 
							json.BlackMilli;
					//	console.log("Black: " + json.BlackMinutes + ":" + json.BlackSeconds + "." + json.BlackMilli);
					}			
				}	
				break;
			
			case "chat_private":
				if(toggleChat !== "false"){
					var datetime = timeStamp();
					document.getElementById('textbox').innerHTML += (datetime + " " + json.Name +": " + 
						json.Text + '\n');
					//scrolls chat to the bottom when it goes below the chat window
					document.getElementById("textbox").scrollTop = document.getElementById("textbox").scrollHeight;
				}
				break;
				
			case "chess_game":

				matchID = json.ID;
				WhiteSide = json.WhitePlayer;
				BlackSide = json.BlackPlayer;

				// if spectate is turned off then other people cannot view this game
				var spectateResult = "No";
				if(toggleSpectate !== "false"){
					spectateResult = "Yes";
				}
				var spectateMessage = {
					Type: "update_spectate",
					Name: user,
					ID: json.ID.toString(),
					Spectate: spectateResult
				}
				sock.send(JSON.stringify(spectateMessage));	

				//if its move zero store time control so it can be used in rematch
				timeGet = json.TimeControl;
				isRatedRematch = json.Rated;
				
				//if the game moves are not null then restore the moves back
				if(json.GameMoves !== null){

					var length = json.GameMoves.length;				
					
					for(var i=0 ; i<length; i++){
						
						var move = game.move({
						    from: json.GameMoves[i].S,
						    to: json.GameMoves[i].T,
						    promotion: json.GameMoves[i].P
					     });

						totalFEN.push(game.fen());
						totalPGN.push(game.pgn());
						totalStatus.push(updateStatus());
						moveCounter++;
					}
					if(length-1 >= 0){
						setStatusAndPGN(totalStatus[length-1], totalPGN[length-1]);
					}
					// now go to the last move
					$('#goEnd').click();
				}
				else{
					if(toggleSound !== "false"){
						//play sound when game starts
						gameSound.play();
					}
					chessGameOver = false; 		
					
					//reset game position to brand new game, used for rematch
					totalFEN = [];
					totalFEN.push("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1");
					totalPGN = [];
					totalPGN.push("");
					totalStatus = [];
					moveCounter = 0;
					totalStatus.push("White to move");
					board.position('rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1');
					board.orientation('white');
					game.reset();							
				}

				// enable abort button if less then 3 moves
				if(totalFEN.length <= 2){
					document.getElementById("abortButton").disabled = false; 
				}
				
				//formating time for clock 
				json.WhiteMinutes = json.WhiteMinutes < 10 ? "0" + json.WhiteMinutes : json.WhiteMinutes;
				json.WhiteSeconds = json.WhiteSeconds < 10 ? "0" + json.WhiteSeconds : json.WhiteSeconds;
	            
				json.BlackMinutes = json.BlackMinutes < 10 ? "0" + json.BlackMinutes : json.BlackMinutes;
				json.BlackSeconds = json.BlackSeconds < 10 ? "0" + json.BlackSeconds : json.BlackSeconds;
				
				if(typeof token.spectate !== "undefined"){
					user = json.WhitePlayer; // default to showing white for spectator		
				}else{
					document.getElementById("drawButton").disabled = false;
					document.getElementById("resignButton").disabled = false;
				}

				if(json.CountryWhite === ""){
					json.CountryWhite = "globe";
				}
				if(json.CountryBlack === ""){
					json.CountryBlack = "globe";
				}

				if (user === json.WhitePlayer){			
					document.getElementById("bottom").innerHTML = "W: <img src='../img/flags/" + 
						json.CountryWhite + ".png'><a href='/profile?name=" + json.WhitePlayer + 
						"'>" + json.WhitePlayer + "</a>"  +	json.WhiteRating +")";
					document.getElementById("top").innerHTML = "B: <img src='../img/flags/" + 
						json.CountryBlack + ".png'><a href='/profile?name=" + json.BlackPlayer  + 
						"'>" + json.BlackPlayer + "</a>(" + json.BlackRating +")";		
				}
				else{
					//flips board white on top black on bottom
					$('#flipOrientationBtn').click();
					document.getElementById("bottom").innerHTML = "B: <img src='../img/flags/" + 
						json.CountryBlack + ".png'><a href='/profile?name=" + json.BlackPlayer  + "'>" + 
						json.BlackPlayer + "</a>(" + json.BlackRating +")";
					document.getElementById("top").innerHTML = "W: <img src='../img/flags/" + 
						json.CountryWhite + ".png'><a href='/profile?name=" + json.WhitePlayer + "'>" + 
						json.WhitePlayer + "</a>(" + json.WhiteRating +")";
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
			case "spectate_game":
				document.getElementById('textbox').innerHTML += (timeStamp() + " " + 
					json.Name + " is now spectating this game." + '\n');
				break;
			case "offer_draw":
				document.getElementById('textbox').innerHTML += (timeStamp() + " " + json.Name +" offers you a draw." + '\n');
				document.getElementById("drawButton").innerHTML = "Accept Draw";
				break;
			
			case "accept_draw":
				document.getElementById('textbox').innerHTML += (timeStamp() + " Game drawn by mututal agreement." + '\n');
				//prvents players from moving pieces now that game is over
				gameOver();
				break;
				
			case "draw_game":
				document.getElementById('textbox').innerHTML += (timeStamp() + " Game is a draw." + '\n');
				//prvents players from moving pieces now that game is over
				gameOver();
				break;
				
			case "cancel_draw":
				document.getElementById('textbox').innerHTML += (timeStamp() + " Draw offer declined." + '\n');
				document.getElementById("drawButton").disabled = false; 
				document.getElementById("drawButton").innerHTML = "Offer Draw";
				break;
			
			case "resign":
				document.getElementById('textbox').innerHTML += (timeStamp() + " " + json.Name + " resigned." + '\n');
				gameOver();
				break;
				
			case "leave":
				document.getElementById('textbox').innerHTML += (timeStamp() + " " + json.Text + '\n');
				break;
				
			case "abort_game":
				document.getElementById('textbox').innerHTML += (timeStamp() + " " + json.Name +" has aborted the game. "+ '\n');
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
					console.log("Invalid game status");
				}
				gameOver();
				break;
				
			case "rating":
				var datetime = timeStamp();
				if(user === WhiteSide){
					document.getElementById('textbox').innerHTML += (datetime + " Your new rating is " + 
						json.WhiteRating + '\n');
				}else{
					document.getElementById('textbox').innerHTML += (datetime + " Your new rating is " + 
						json.BlackRating + '\n');
				}
				break;
				
			case "rematch":
				document.getElementById('textbox').innerHTML += (timeStamp() + " Your opponent offers you a rematch." + '\n');
				document.getElementById('rematchButton').innerHTML = "Accept Rematch";
				break;
				
			case "match_three":
				document.getElementById('textbox').innerHTML += (timeStamp() + 
					" Your can only have a max of three pending matches." + '\n');
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

		var message = {
			Type: "chat_private",
			Name: user,
			Text: document.getElementById('message').value
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
		
		if (document.getElementById("drawButton").innerHTML === "Accept Draw"){
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
		
		document.getElementById('textbox').innerHTML += (timeStamp() + " You offer your opponent a draw." + '\n');
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
 	//prevents sending null sock.send() when going over a game
	if(reviewGame){
		return;
	}
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
function sendMove(src, dest, pawnPromotion){
		
	var message = {
		Type: "send_move",
		Name: user,	
		ID: matchID,
		Source: src,
		Target: dest,
		Promotion: pawnPromotion
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

	if(moveCounter < totalFEN.length-1){	
		moveCounter++;
	}
	//make a global array and iterate forwards through the global array when going forward
	board.position(totalFEN[moveCounter]);	
	setStatusAndPGN(totalStatus[moveCounter], totalPGN[moveCounter]);
} 
$('#goBack').on('click', function() {
	if(moveCounter > 0){
		
		moveCounter--;
	}
	//make a global array and iterate backwards through the global array when going back
	board.position(totalFEN[moveCounter]);	
	setStatusAndPGN(totalStatus[moveCounter], totalPGN[moveCounter]);
});
//move forward to last move
document.getElementById('goEnd').onclick = function(){

	for(var i=moveCounter; i<totalFEN.length; i++){
		board.position(totalFEN[i]);
	}
	moveCounter = totalFEN.length-1;
	if(moveCounter>=0){
		setStatusAndPGN(totalStatus[moveCounter], totalPGN[moveCounter]);
	}
} 

//offers player a rematch or accepts it if the other player offers
document.getElementById('rematchButton').onclick = function(){

	//the opponent
	var fighter;
	
	if(user === WhiteSide){
		fighter = BlackSide;
	}
	else{
		fighter = WhiteSide;
	}
	
	if(document.getElementById('rematchButton').innerHTML === "Rematch"){
		var message = {
			Type: "rematch",
			Name: user,	
			Opponent: fighter,
			Rated: isRatedRematch,
			TimeControl: timeGet
		}
		document.getElementById('textbox').innerHTML += (timeStamp() + " You offer your opponent a rematch." + '\n');  
	}
	else{ //else value === "Accept Rematch"
		var message = {
			Type: "accept_rematch",
			Name: user,	
			Opponent: fighter,
			TimeControl: timeGet
		}
		document.getElementById('rematchButton').innerHTML = "Rematch";
	}
	//hiding button after click to prevent rematch abuse
	$('#rematchButton').hide();
	sock.send(JSON.stringify(message));
} 
	
$('#message').keypress(function(event) {
    if (event.which === 13) {
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

function detectMobile(){ //tries to detect if user is using a mobile device
		
	if(screen.width <= 900){
		console.log("mobile device detected...adjusting board size and layout");
		document.getElementById("chatleft").style.display = "none";
		document.getElementById("notation").style.display = "none";	
	}
}

detectMobile(); //calls function