//used to pass in the json string of the chess moves so player can review them
function reviewGame(moves, white, black, whiteRating, blackRating, time, result, date){
	window.location = "/chess/memberChess?moves=" + moves + "&white=" + white + "&black=" + black + "&whiteRating=" + whiteRating + "&blackRating=" + blackRating + "&time=" + time + "&result=" + result + "&date=" + date;
}
//used for highscore.html to redirect when a name on highscores is clicked
function reviewProfile(user){
	window.location = "profile?name=" + user; 
}

document.getElementById('searchPlayer').onclick = function(){
	    var name = document.getElementById('playerName').value;
		window.location = "profile?name=" + name; 
}