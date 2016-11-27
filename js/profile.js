//used to pass in the json string of the chess moves so player can review them
function reviewGame(moves, white, black, whiteRating, blackRating, time, result, date, countryWhite, countryBlack){
	window.location = "/chess/memberChess?moves=" + moves + "&white=" + white + "&black=" + black + 
    "&whiteRating=" + whiteRating + "&blackRating=" + blackRating + "&time=" + time + 
    "&result=" + result + "&date=" + date + "&countryWhite=" + countryWhite + "&countryBlack=" + countryBlack;
}
//used for highscore.html to redirect when a name on highscores is clicked
function reviewProfile(user){
	window.location = "profile?name=" + user; 
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

// gets country flag stored in cookie and updates the img source
function setFlag(){
	var countryFlag = getCookie("country");
	document.getElementById('countryFlag').src = "img/flags/" + countryFlag + ".png"
}
setFlag();