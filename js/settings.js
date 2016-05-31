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
//name of cookie, value and days until expire
function setCookie(cname, cvalue, exdays) {
    var d = new Date();
    d.setTime(d.getTime() + (exdays*24*60*60*1000));
    var expires = "expires="+d.toUTCString();
    document.cookie = cname + "=" + cvalue + "; " + expires;
} 
//sets up cookies
function setupCookies() {

    var toggleSound = getCookie("sound");
	var togglePremove = getCookie("premove");
	var toggleChat = getCookie("chat");
	var toggleSpectate = getCookie("spectate");
	var togglePromote = getCookie("promote");
	var pieceTheme = getCookie("pieceTheme");

    if (toggleSound.value === "") { //if no cookie exist then set it to true
        setCookie("sound", "true", 30);
		console.log("yes");
    }
	else if(toggleSound === "false"){ //if its true then leave as default otherwise set it as false
		document.getElementById("toggleSound").checked = false;
	
	}
	
	if (togglePremove === "") { //if no cookie exist then set it to true
        setCookie("premove", "true", 30);
    }
	else if(togglePremove === "false"){ //if its true then leave as default otherwise set it as false
		document.getElementById("togglePremove").checked = false;
	}
	
	if (toggleChat === "") { //if no cookie exist then set it to true
        setCookie("chat", "true", 30);
    }
	else if(toggleChat === "false"){ //if its true then leave as default otherwise set it as false
		document.getElementById("toggleChat").checked = false;
	}
	
	if (toggleSpectate === "") { //if no cookie exist then set it to true
        setCookie("spectate", "true", 30);
    }
	else if(toggleSpectate === "false"){ //if its true then leave as default otherwise set it as false
		document.getElementById("toggleSpectate").checked = false;
	}
	
	if (togglePromote === "") { //if no cookie exist then set it to true
        setCookie("promote", "true", 30);
    }
	else if(togglePromote === "false"){ //if its true then leave as default otherwise set it as false
		document.getElementById("togglePromote").checked = false;
	}
	
	if (pieceTheme === "") { //if no cookie exist then set it to true
        setCookie("pieceTheme", "wikipedia", 30); //Wikipedia is default piece theme
    }
	else if(pieceTheme === "uscf"){ 
		document.getElementById("uscf").checked = true;
	}
	else if(pieceTheme === "alpha"){
		document.getElementById("alpha").checked = true;
	}
}

//setCookie("Age", 14, 3)
setupCookies();

document.getElementById('toggleSound').onclick = function(){
	if(document.getElementById("toggleSound").checked === false){
		setCookie("sound", "false", 30);	
	}
	else{	
		setCookie("sound",  "true", 30);	
	}		
}
document.getElementById('togglePremove').onclick = function(){
	if(document.getElementById("togglePremove").checked === false){
		setCookie("premove", "false", 30);	
	}
	else{	
		setCookie("premove",  "true", 30);	
	}
}
document.getElementById('toggleChat').onclick = function(){
	if(document.getElementById("toggleChat").checked === false){
		setCookie("chat", "false", 30);		
	}
	else{	
		setCookie("chat",  "true", 30);	
	}
}
document.getElementById('toggleSpectate').onclick = function(){
	if(document.getElementById("toggleSpectate").checked === false){
		setCookie("spectate", "false", 30);
	}
	else{	
		setCookie("spectate",  "true", 30);	
	}
}
document.getElementById('togglePromote').onclick = function(){
	if(document.getElementById("togglePromote").checked === false){
		setCookie("promote", "false", 30);
	}
	else{	
		setCookie("promote",  "true", 30);	
	}
}

document.getElementById('wiki').onclick = function(){
	setCookie("pieceTheme", "wikipedia", 30);
}

document.getElementById('uscf').onclick = function(){
	setCookie("pieceTheme", "uscf", 30);	
}

document.getElementById('alpha').onclick = function(){
	setCookie("pieceTheme", "alpha", 30);	
}