user = document.getElementById('user').value;

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
        //action listener for exporting game to PGN file
        document.getElementById('exportPGN').onclick = function(){
            game.header('Site', json.Site, 'Event', json.Event, 'Date', json.Date, 'White',
                json.White, 'Black', json.Black, 'Result', json.Result, 'WhiteELO',
                json.WhiteElo, 'BlackElo', json.BlackElo);

            // second parameter is file name
            download(game.pgn(), json.White + "vs" + json.Black + ".pgn", "application/x-chess-pgn");
        }	

        document.getElementById("bottom").innerHTML = "W: " + json.White + "(" +
            json.WhiteElo + ")" ;
        document.getElementById("top").innerHTML = "B: " + json.Black  + "(" +
            json.BlackElo +")";
        document.getElementById("gameID").innerHTML = json.ID;
        document.getElementById("event").innerHTML = json.Event;
        document.getElementById("site").innerHTML = json.Site;
        document.getElementById("eco").innerHTML = json.ECO;
        document.getElementById("result").innerHTML = json.Result;

        //always push the default starting position
        totalFEN.push("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1");
        totalStatus.push("White to move");
        totalPGN.push("");

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

//validates input in seach box is a number that was entered
function isNumber(evt){
    
    var charCode = (evt.which) ? evt.which : evt.keyCode;
    
	// Allow backspace, left and right arrow keys
	if(charCode === 8 || charCode === 37 || charCode === 39){
		return true;
	}
    if (charCode < 48 || charCode > 58){
		return false;
	}   
    return true;
}