 var matchID;
var moveSound = new Audio('../sound/chessmove.mp3');
var gameSound = new Audio('../sound/startgame.mp3');

//getting user preferences
var toggleSound = getCookie("sound");

var whiteRating;
var blackRating;
//timeGet is a global variable and is not located here
var gameDate;
var gameResult;

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
    download(game.pgn(), WhiteSide + "vs" + BlackSide + ".pgn", "application/x-chess-pgn");
}	

user = document.getElementById('user').value;

//always push the default starting position
totalFEN.push("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1");
totalStatus.push("White to move");
totalPGN.push("");


function loadDatabase(){
     $.ajax({
        url: 'load_database',
        type: 'post',
        dataType: 'html',
        data : { 'user': user, 'password': password, 'captchaId': captchaId, 'captchaSolution': captchaSolution},
        success : function(data) {			
            $('#submit-result').html(data);	
        }	
    });
}

function loadGame(){
    $.ajax({
        url: 'load_game',
        type: 'post',
        dataType: 'html',
        data : { 'user': user, 'password': password, 'captchaId': captchaId, 'captchaSolution': captchaSolution},
        success : function(data) {			
            $('#submit-result').html(data);	
        }	
    });

    WhiteSide = json.WhitePlayer;
    BlackSide = json.BlackPlayer;
    
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
    }
    
    //formating time for clock 
    json.WhiteMinutes = json.WhiteMinutes < 10 ? "0" + json.WhiteMinutes : json.WhiteMinutes;
    json.WhiteSeconds = json.WhiteSeconds < 10 ? "0" + json.WhiteSeconds : json.WhiteSeconds;
    
    json.BlackMinutes = json.BlackMinutes < 10 ? "0" + json.BlackMinutes : json.BlackMinutes;
    json.BlackSeconds = json.BlackSeconds < 10 ? "0" + json.BlackSeconds : json.BlackSeconds;

    if (user === json.WhitePlayer){
                                            
        document.getElementById("bottom").innerHTML = "W: <a href='/profile?name=" + json.WhitePlayer + 
            "'>" + json.WhitePlayer + "</a>(" + json.WhiteRating +")";
        document.getElementById("top").innerHTML = "B: <a href='/profile?name=" + json.BlackPlayer  + 
            "'>" + json.BlackPlayer + "</a>(" + json.BlackRating +")";			
    }
    else{
        //flips board white on top black on bottom
        $('#flipOrientationBtn').click();
        document.getElementById("bottom").innerHTML = "B: <a href='/profile?name=" + json.BlackPlayer  + "'>" + 
            json.BlackPlayer + "</a>(" + json.BlackRating +")";
        document.getElementById("top").innerHTML = "W: <a href='/profile?name=" + json.WhitePlayer + "'>" + 
            json.WhitePlayer + "</a>(" + json.WhiteRating +")";
    }
    
    if(json.Status === "White"){
        if(user === json.WhitePlayer){
            document.getElementById("bottomtime").value = json.WhiteMinutes + ":" + json.WhiteSeconds;
            document.getElementById("toptime").value = json.BlackMinutes + ":" + json.BlackSeconds;
        }
        else{
            document.getElementById("toptime").value = json.WhiteMinutes + ":" + json.WhiteSeconds;
            document.getElementById("bottomtime").value = json.BlackMinutes + ":" + json.BlackSeconds;
        }		
    }
    //else if (json.Status === "Black)
    else {
        if(user === json.WhitePlayer){
            document.getElementById("toptime").value = json.BlackMinutes + ":" + json.BlackSeconds;
            blackClock.start($('#toptime').val());
            document.getElementById("bottomtime").value = json.WhiteMinutes + ":" + json.WhiteSeconds;

        }
        else{
            document.getElementById("bottomtime").value = json.BlackMinutes + ":" + json.BlackSeconds;
            blackClock.start($('#bottomtime').val());
            document.getElementById("toptime").value = json.WhiteMinutes + ":" + json.WhiteSeconds;
        }				
    }
}

function searchGame(){
    $.ajax({
        url: 'search_game',
        type: 'post',
        dataType: 'html',
        data : { 'user': user, 'password': password, 'captchaId': captchaId, 'captchaSolution': captchaSolution},
        success : function(data) {			
            $('#submit-result').html(data);	
        }	
    });
}

//go forward one move
document.getElementById('goForward').onclick = function(){

	if(moveCounter < totalFEN.length){	
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
	moveCounter = totalFEN.length;
	if(moveCounter>=0){
		setStatusAndPGN(totalStatus[moveCounter-1], totalPGN[moveCounter-1]);
	}
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