// global object which stores opening ECO, opening names and their moves
var allOpenings;
var ECOIndex = 0;

var openingDropDown = document.getElementById('openingDropDown');

function getGame(gameID){
    $.ajax({
        url: 'fetchgameID',
        type: 'post',
        dataType: 'html',
        data : { 'gameID': gameID},
        success : function(data) {			
            // error messages will be less then 100 characters, games always more then 100 characters
            if(data.length <= 100){
                document.getElementById('textbox').innerHTML += (timeStamp() + " " +
			        data + "\n"); 
            }else{
                loadGame(data);
            } 
            document.getElementById("textbox").scrollTop = document.getElementById("textbox").scrollHeight;
        }	
    });
}

getGame(1);

function getGameByECO(ECO, localECOIndex){
    $.ajax({
        url: 'fetchgameByECO',
        type: 'post',
        dataType: 'html',
        data : { 'ECO':ECO, 'ECOIndex': localECOIndex},
        success : function(data) {			
            // error messages will be less then 100 characters, games always more then 100 characters
            if(data.length <= 100){
                document.getElementById('textbox').innerHTML += (timeStamp() + " " +
			        data + "\n"); 
            }else{
                loadGame(data);
            } 
            document.getElementById("textbox").scrollTop = document.getElementById("textbox").scrollHeight;
        }	
    });
}
// Fills opening drop down with ECO and opening name
setupOpening();

function setupOpening(){
    
    $.getJSON('/data/openings.json', function(data) {         
       
        // setting global variable with the opening object
        allOpenings = data;

        for (var key in data) {
            // skip loop if the property is from prototype
            if (!data.hasOwnProperty(key)){
                continue;
            }

            var option = document.createElement('option');
            option.text = key + ": " + data[key].name;
            option.value = key;
            openingDropDown.add(option);
        }
    });
}

// loads game based on the JSON string in data that is passed in
function loadGame(gameData){
   
    var json = JSON.parse(gameData);
    var moves = JSON.parse(json.Moves);
    document.getElementById('textbox').innerHTML += (timeStamp() + " " +
			"Game ID " +  json.ID + " has loaded.\n");
    if(moves !== null){
        //action listener for exporting game to PGN file
        document.getElementById('exportPGN').onclick = function(){
            game.header('Site', json.Site, 'Event', json.Event, 'Date', json.Date, 'White',
                json.White, 'Black', json.Black, 'Result', json.Result, 'WhiteELO',
                json.WhiteElo, 'BlackElo', json.BlackElo);

            // second parameter is file name
            download(game.pgn(), json.White + " vs. " + json.Black + ".pgn", "application/x-chess-pgn");
        }	
        // updates the ID in the Search ID input box
        document.getElementById('searchID').value = json.ID;

        document.getElementById("bottom").innerHTML = "W: " + json.White + "(" +
            json.WhiteElo + ")" ;
        document.getElementById("top").innerHTML = "B: " + json.Black  + "(" +
            json.BlackElo +")";
        document.getElementById("gameID").innerHTML = json.ID;
        document.getElementById("event").innerHTML = json.Event;
        document.getElementById("site").innerHTML = json.Site;
        document.getElementById("eco").innerHTML = json.ECO;
        document.getElementById("result").innerHTML = json.Result;

        //resetting array as its a brand new game
        totalFEN = [];
        totalStatus = [];
        totalPGN = [];
        $('#goStart').click();
        game.reset();
        //always push the default starting position
        totalFEN.push("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1");
        totalStatus.push("White to move");
        totalPGN.push("");

        var length = moves.length;				
        
        for(var i=0 ; i<length; i++){
            // If promotion check to see if ASCII conversion needs to take place
            if (moves[i].P !== ""){
                switch(parseInt(moves[i].P)){
                    case 113:
                        moves[i].P = "q";
                        break;
                    case 114:
                        moves[i].P = "r";
                        break;
                    case 110:
                        moves[i].P = "n";
                        break;
                    case 98:
                        moves[i].P = "b";
                        break;
                    default:
                        // then no conversion is needed
                }
            }
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

document.getElementById('searchGameButton').onclick = function(){
    getGame(document.getElementById('searchID').value);
    ECOIndex = 0;
}

document.getElementById('searchOpeningButton').onclick = function(){
    getGameByECO(document.getElementById('openingDropDown').value, 0);
    ECOIndex = 0;
}

document.getElementById('previousOpeningButton').onclick = function(){
    if(ECOIndex-1 >= 0){
        ECOIndex--;
        getGameByECO(document.getElementById('openingDropDown').value, ECOIndex);
    }
}

document.getElementById('nextOpeningButton').onclick = function(){
    if(ECOIndex+1 < 100){
        ECOIndex++;
        getGameByECO(document.getElementById('openingDropDown').value, ECOIndex);
    }
}

document.getElementById('goForwardGame').onclick = function(){
    var value = parseInt(document.getElementById('searchID').value) + 1
    getGame(value);
} 

document.getElementById('goBackGame').onclick = function(){
    var value = parseInt(document.getElementById('searchID').value) - 1;
    if(value > 0){
        getGame(value);
    }
} 

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

function timeStamp(){ //returns timestamp
	var currentdate = new Date(); 
	var datetime =  + currentdate.getHours() + ":"  
            		+ currentdate.getMinutes() + ":" 
            		+ currentdate.getSeconds();
	return datetime;
}

$('#searchID').keypress(function(event) {
    if (event.which === 13) {  
	   $('#searchGameButton').click();	
    }
});
