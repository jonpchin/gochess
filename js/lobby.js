if (!window.WebSocket){
	console.log("Your browser doesn't support websockets. Please use the latest version of Firefox, Chrome, IE, Opera or Edge");
	$('#checkwebsocket').html("Your browser doesn't support websockets. Please use the latest version of Firefox, Chrome, IE, Opera or Microsoft Edge.");
}
var wsuri = "wss://goplaychess.com:443/server";
var sock = new WebSocket(wsuri);

var user = document.getElementById('user').value;
var bullet = document.getElementById('bullet').value;
var blitz = document.getElementById('blitz').value;
var standard = document.getElementById('standard').value;

window.onload = function() {	
	//setting up user preferences
	var challenge = new Audio('../sound/challenge.mp3');
	var toggleSound = getCookie("sound");
	var toggleChat = getCookie("chat");
	
    sock.onopen = function() {
		
		var message = {
			Type: "fetch_matches",
			Name: user
		}
	    sock.send(JSON.stringify(message));
		
		var message = {
			Type: "fetch_players",
			Name: user
		}
	    sock.send(JSON.stringify(message));
		
		document.getElementById('textbox').innerHTML = "";
		var message = {
			Type: "chat_all",
			Name: user,
			Text: "has entered the lobby."
		}
	    sock.send(JSON.stringify(message));
    }
	
    sock.onclose = function(e) {
		var datetime =  timeStamp();
		document.getElementById('textbox').innerHTML += (datetime + " " + "Connection lost. Please refresh to reconnect." + '\n');
    }
	
    sock.onmessage = function(e) {
		
		json = JSON.parse(e.data);
		
		if (json.Type === "chat_all"){
			
			var show = true;
			var table = document.getElementById("online");
			var row;
			for (var i = 0; i<table.rows.length; i++) {
				row=table.rows[i];
		
				if(row.cells[0].innerHTML === json.Name){
					show=false;
					break;
				} 
			}
		
			if(show === true && json.Name !== user){
				$('#online').html(function() { //when a player chats or connects they will update online players list for everyone
					return  $(this).html() +'<tr><td  onclick="reviewProfile(\''+ json.Name +'\');">' + json.Name + '</td></tr>';
				});
			}
					
			if(toggleChat !== "false"){
				var datetime =  timeStamp();
				document.getElementById('textbox').innerHTML += (datetime + " " + json.Name +": " + json.Text + '\n');
			}
		}
		else if(json.Type === "match_seek" || json.Type === "fetch_matches"){
			if(json.Name === user){
				
				if(json.Opponent !== ""){ //then its a private match
					$('#seekmatch').html(function() {
					
						return  $(this).html() + '<tr onclick="cancelMatch('+ json.MatchID + ');" id=MatchID' +json.MatchID + ' class=NameID' + json.Opponent + '><td>' + json.Name +  '</td><td>' + json.Rating + '</td><td>' + json.GameType + '</td><td>' + json.TimeControl + " Minutes" + '</td><td>' + "Private-Match" + '</td></tr>';
					});
					if(toggleSound !== "false"){
						challenge.play();//plays sound to notify player they got a private match sent to them
					}
				}
				else{
					
					$('#seekmatch').html(function() {
						
						return  $(this).html() + '<tr onclick="cancelMatch('+ json.MatchID + ');" id=MatchID' +json.MatchID + ' class=NameID' + json.Name + '><td>' + json.Name +  '</td><td>' + json.Rating + '</td><td>' + json.GameType + '</td><td>' + json.TimeControl + " Minutes" + '</td><td>' + json.MinRating + "-" + json.MaxRating + '</td></tr>';
					});
				}
				
			}
			//rejecting matches that do not fit the user criteria
			
			else if(json.GameType === "bullet" && (bullet < json.MinRating || bullet > json.MaxRating)){
					return;
			}	
			else if(json.GameType === "blitz" && (blitz < json.MinRating || blitz > json.MaxRating)){
					return;
			}
			else if (json.GameType === "standard" && (standard < json.MinRating || standard > json.MaxRati)){
					return;		
			}
			
			else{
				
				if(json.Opponent !== ""){ //then its a private match
					
					$('#seekmatch').html(function() {
	   					return  $(this).html() + '<tr onclick="acceptMatch('+ json.MatchID + ');" id=MatchID' +json.MatchID + ' class=NameID' + json.Name + '><td>' + json.Opponent +  '</td><td>' + json.Rating + '</td><td>' + json.GameType + '</td><td>' + json.TimeControl + " Minutes" + '</td><td>' + "Private-Match" + '</td></tr>';
					});
					
				}
				else{
					$('#seekmatch').html(function() {
   						return  $(this).html() + '<tr onclick="acceptMatch('+ json.MatchID + ');" id=MatchID' +json.MatchID + ' class=NameID' + json.Name + '><td>' + json.Name +  '</td><td>' + json.Rating + '</td><td>' + json.GameType + '</td><td>' + json.TimeControl + " Minutes" + '</td><td>' + json.MinRating + "-" + json.MaxRating + '</td></tr>';
					});
				}
				
				
			}
		
		}
		else if(json.Type === "fetch_players"){
			
			$('#online').html(function() {
					return  $(this).html() + '<tr><td onclick="reviewProfile(\''+ json.Name +'\')" onmouseover="getPlayerInfo(\'' + json.Name + '\');">' + json.Name +  '</td></tr>';
			});
			
		}
		else if(json.Type === "broadcast"){
			var table = document.getElementById("online");
			var rows;
			for (var i = 0; i<table.rows.length; i++) {
				row=table.rows[i];
				if(row.cells[0].innerHTML === json.Name){
					table.deleteRow(i);
					console.log(json.Name + " has left the lobby.");
					break;
				} 
			}
		}
		else if(json.Type === "accept_match"){
//			console.log("accepting row")
				
			$(".NameID"+json.Name).remove();
			$(".NameID"+json.TargetPlayer).remove();
			if(json.Name === user || json.TargetPlayer === user){
				//sock.close();
				window.location = "/chess/memberChess";
			}
			
		}
		else if(json.Type === "private_match"){
			if(json.Name === user){
				
				$('#seekmatch').html(function() {
					
					return  $(this).html() + '<tr onclick="cancelMatch('+ json.MatchID + ');" id=MatchID' +json.MatchID + ' class=NameID' + json.Name + '><td>' + json.Opponent +  '</td><td>' + json.Rating + '</td><td>' + json.GameType + '</td><td>' + json.TimeControl + " Minutes" + '</td><td>' + "Your-Match" + '</td></tr>';
				});
			}
			else{
				$('#seekmatch').html(function() {
   					return  $(this).html() + '<tr onclick="acceptMatch('+ json.MatchID + ');" id=MatchID' +json.MatchID + ' class=NameID' + json.Name + '><td>' + json.Name +  '</td><td>' + json.Rating + '</td><td>' + json.GameType + '</td><td>' + json.TimeControl + " Minutes" + '</td><td>' + "Private-Match" + '</td></tr>';
				});
				if(toggleSound !== "false"){
					challenge.play();//plays sound to notify player they got a private match sent to them
				}
				
			}
		}
		else if(json.Type === "cancel_match"){
//			console.log("deleting row")
//			console.log("MatchID is " + json.MatchID)
			$("#MatchID"+json.MatchID).remove();
		}
		else if (json.Type === "alert"){
			alert("You or the opponent are already in a game.");
		}
		else if(json.Type === "maxThree"){
			alert("You are only allowed to have a max of three pending seeks at one time.");
		}
		else if(json.Type === "absent"){
			var datetime =  timeStamp();
			document.getElementById('textbox').innerHTML += (datetime + " That player is not logged in."  + '\n');
		}
		else if(json.Type === "chess_game"){

			window.location = "/chess/memberChess";
		}
		else if(json.Type === "massMessage"){
			
			var datetime = timeStamp();
			document.getElementById('textbox').innerHTML += (datetime + " " + json.Text + '\n');
			sock.close();
			document.getElementById("sendButton").disabled = true;
			document.getElementById("sendSeek").disabled = true;
			document.getElementById("sendPrivateMatch").disabled = true;
		}
		else{
			console.log("Unknown API type");	
		}

    }
	document.getElementById('sendButton').onclick = function(){
	    var msg = document.getElementById('message').value;
		
		var message = {
			Type: "chat_all",
			Name: user,
			Text: msg
		}
	    sock.send(JSON.stringify(message));
		document.getElementById('message').value = "";
		$('#message').focus();
	}
	
	document.getElementById('sendSeek').onclick = function(){
	    var dropdown = document.getElementById('timecontrol');
		var time = dropdown.options[dropdown.selectedIndex].value;
		var dropdown = document.getElementById('minrating');
		var min = dropdown.options[dropdown.selectedIndex].value;
		var dropdown = document.getElementById('maxrating'); 
		var max = dropdown.options[dropdown.selectedIndex].value;
			
		var message = {
			Type: "match_seek",
			Name: user,
			TimeControl: parseInt(time, 10),
			MinRating: parseInt(min, 10),
			MaxRating: parseInt(max, 10)
		}
	    sock.send(JSON.stringify(message));
	}
	
	document.getElementById('sendPrivateMatch').onclick = function(){
		
		var opponent = document.getElementById('privateName').value;
		var dropdown = document.getElementById('timecontrol');
		var time = dropdown.options[dropdown.selectedIndex].value;
		var length = opponent.length;
		
		if(length < 3 || length > 12){ //checking if name meets length requirements
			var datetime =  timeStamp();
			document.getElementById('textbox').innerHTML += (datetime + " " + "Please enter a username with 3-12 characters." + '\n');
			return;
		}
		if(opponent === user){ //you can't send private matches to yourself!
			var datetime =  timeStamp();
			document.getElementById('textbox').innerHTML += (datetime + " " + "You can't match yourself!" + '\n');
			return;
		}

		var message = {
			Type: "private_match",
			Name: user,
			Opponent: opponent,
			TimeControl: parseInt(time, 10)
		}
	    sock.send(JSON.stringify(message));
	}
};

function cancelMatch(matchID){
	var message = {
		Type: "cancel_match",
		Name: user,
		MatchID: matchID
	}
    sock.send(JSON.stringify(message));
}

function acceptMatch(matchID){
	var message = {
		Type: "accept_match",
		Name: user,
		MatchID: matchID
	}
    sock.send(JSON.stringify(message));
}

$('#message').keypress(function(event) {
    if (event.which === 13) {
       var msg = document.getElementById('message').value;
	  
	   $('#sendButton').click();	
    }
});

function timeStamp(){ //returns timestamp
	var currentdate = new Date(); 
	var datetime =  + currentdate.getHours() + ":"  
            		+ currentdate.getMinutes() + ":" 
            		+ currentdate.getSeconds();
	return datetime;
}


function reviewProfile(name){ //when a player clicks a name under online players it will redirect to players profile
	window.location = "/profile?name=" + name;
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

function getPlayerInfo(name){ //returns all three ratings of players and if there are in a game

	$('#playerData').html('<img src="../img/ajax/loading.gif" />');
    $.ajax({
  		url: 'getPlayerData',
   		type: 'post',
   		dataType: 'html',
   		data : { 'username': name},
   		success : function(data) {			
      		$('#playerData').html(data);	
			
   		}	
    });
}