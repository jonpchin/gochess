var WhiteELO;
var BlackELO;
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
        'Result', gameResult, 'WhiteELO', WhiteELO, 'BlackElo', BlackELO, 'TimeControl', timeGet);

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

function getGame(gameID){
    $.ajax({
        url: 'fetchgameID',
        type: 'post',
        dataType: 'html',
        data : { 'gameID': gameID},
        success : function(data) {			
           //console.log(data);
            loadGame(data);
        }	
    });
}

getGame(3);

// loads game based on the JSON string in data that is passed in
function loadGame(gameData){
   
    var json = JSON.parse(gameData);
    var moves = JSON.parse(json.Moves);

    if(moves !== null){

        document.getElementById("bottom").innerHTML = "W: " + json.White + "(" +
            json.WhiteElo + ")" ;
        document.getElementById("top").innerHTML = "B: " + json.Black  + "(" +
            json.BlackElo +")";
            
        whiteELO = json.WhiteElo;
        blackELO = json.BlackElo;

        var length = moves.length;				
        
        for(var i=0 ; i<length; i++){
            var move = game.move({
                from: moves[i].S,
                to: moves[i].T,
                promotion: moves[i].P
            });

            totalFEN.push(game.fen());
            totalPGN.push(game.pgn());
            totalStatus.push(updateStatus());
        }
        if(length-1 >= 0){
            setStatusAndPGN(totalStatus[length-1], totalPGN[length-1]);
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