function checkCookie(){
        var cookieEnabled = (navigator.cookieEnabled) ? true : false;
        
		if (typeof navigator.cookieEnabled == 'undefined' && !cookieEnabled) {
            document.cookie = 'testcookie';
            cookieEnabled = (document.cookie.indexOf('testcookie') != -1) ? true : false;
        }
		
		if(!cookieEnabled){
			$('#submit-result').html("You have cookies disabled. Please enable cookies to login.");	
		}
}

checkCookie();

function parseUrl(){ //checks url to see if should welcome user
	var url = window.location.href; 
	var name = url.split('?user=')[1];
	if(name){
		$('#welcome').html("<img src='img/ajax/available.png'/>Welcome " + name + "! Your account is now activated. You may login.");
	}
}

parseUrl();

function setSrcQuery(e, q) {
	var src  = e.src;
	var p = src.indexOf('?');
	if (p >= 0) {
		src = src.substr(0, p);
	}
	e.src = src + "?" + q
}
document.getElementById('playAudio').onclick = function() {
	var le = document.getElementById("lang");
	var lang = le.options[le.selectedIndex].value;
	var e = document.getElementById('audio')
	setSrcQuery(e, "lang=" + lang)
	e.style.display = 'block';
	e.autoplay = 'true';
	return false;
}

document.getElementById("lang").addEventListener("change", function(){
	var e = document.getElementById('audio')
	if (e.style.display == 'block') {
		playAudio();
	}
})
document.getElementById('reload').onclick = function()  { //used to reload captcha
	setSrcQuery(document.getElementById('image'), "reload=" + (new Date()).getTime());
	setSrcQuery(document.getElementById('audio'), (new Date()).getTime());
	return false;
}

document.getElementById('loginAsGuest').onclick = function() {
	document.getElementById('login').disabled = true;
	document.getElementById('loginAsGuest').disabled = true;

	$('#submit-result').html('<img src="img/ajax/loading.gif" />Login in progress...');
    $.ajax({
  		url: 'enterGuest',
   		type: 'post',
   		dataType: 'html',
   		success : function(data) {			
      		$('#submit-result').html(data);	
   		}	
    });
}

document.getElementById('login').onclick = function(){
	document.getElementById('login').disabled = true;
	document.getElementById('loginAsGuest').disabled = true;
	setTimeout(function() {
		// enable click after 1 second
		document.getElementById('login').disabled = false;
	}, 1000); //  second delay
	
	var user = document.getElementById('user').value;
	var password = document.getElementById('password').value;
	var captchaId = document.getElementById('captchaId').value;
	var captchaSolution = document.getElementById('captchaSolution').value;

	if(user === "" || password === ""){
		$('#submit-result').html("<img src='img/ajax/not-available.png' /> Please fill out all fields.");
		return
	}
	
	$('#submit-result').html('<img src="img/ajax/loading.gif" />Login in progress...');
    $.ajax({
  		url: 'processLogin',
   		type: 'post',
   		dataType: 'html',
   		data : { 'user': user, 'password': password, 'captchaId': captchaId, 'captchaSolution': captchaSolution},
   		success : function(data) {			
      		$('#submit-result').html(data);	
   		}	
    });
	
	$.ajax({
  		url: 'updateCaptcha',
   		type: 'post',
   		success : function(data) {			
			document.getElementById('captchaId').value = data;
			document.getElementById('image').src = "/captcha/" + data + ".png"
   		},
	});

}
//when user presses enter on username input it will jump to password input
$('#user').keypress(function(event) {
    if (event.which === 13) {  
	   $('#password').focus();
    }
});

//when user presses enter on password input box it will trigger the login button
$('#password').keypress(function(event) {
    if (event.which === 13) {  
	   $('#login').click();	
    }
});