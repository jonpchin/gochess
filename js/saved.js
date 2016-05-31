//checks if the other player is online in the lobby or game room, if not pops up an alert saying
//the user needs to either be in the lobby or game room
function resumeGame(id, white, black){
	$.ajax({
  		url: 'resumeGame',
   		type: 'post',
   		dataType: 'html',
   		data : {'id': id, 'white': white, 'black': black},
   		success : function(data) {			
//      	$('#submit-result').html(data);	
			if(data === "true"){
				window.location = "/chess/memberChess";
			}
			else{
				alert("Your opponent is not in the lobby.");
			}
   		}	
    });
}