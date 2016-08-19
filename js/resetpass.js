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

document.getElementById('confirm').onkeyup = function(){ //checks if password and confirm password match
	var password = document.getElementById('pass').value;
    var confirm = document.getElementById('confirm').value;

    if (password != confirm){
		$("#checkPass").html('<img src="img/ajax/not-available.png" />');
	}
    else{
		$("#checkPass").html('<img src="img/ajax/available.png" />' +  " Passwords match.");
	}
        
}

document.getElementById('pass').onkeyup = function(){ //checks if password is correct length
	var length = document.getElementById('pass').value.length;
    
    if (length < 5 || length >90){
		$("#checkLen").html('<img src="img/ajax/not-available.png" />');
	} 
    else{
		$("#checkLen").html('<img src="img/ajax/available.png" />');
	}
        
}

function copyParams(){ //checks the url parameters if there any and copies them to input box
	
	var token = parseUrl();
	if (token){
		document.getElementById('token').value = token;
		$('#user').focus();
	}
}

function parseUrl() { 
 //checks url to see if should welcome user
	var url = window.location.href; 
	var token = url.split('?token=')[1];
	return token;
}
copyParams();

document.getElementById('resetpass').onclick = function(){
	document.getElementById('resetpass').disabled = true;
	setTimeout(function() {
		// enable click after 1 second
		document.getElementById('resetpass').disabled = false;
	}, 1000); // 1 second delay
	
	var user = document.getElementById('user').value;
	var token = document.getElementById('token').value;
	var pass = document.getElementById('pass').value;
	var confirm = document.getElementById('confirm').value;
	var captchaId = document.getElementById('captchaId').value;
	var captchaSolution = document.getElementById('captchaSolution').value;
	
	if(user === "" || pass === "" || confirm === "" || token === ""  || captchaSolution === ""){
		$('#submit-result').html("<img src='img/ajax/not-available.png' /> Please fill out all fields.");
		return
	}
	
	$('#submit-result').html('<img src="img/ajax/loading.gif" />Submission in progress...');
    $.ajax({
  		url: 'processResetPass',
   		type: 'post',
   		dataType: 'html',
   		data : { 'user': user, 'token':token, 'pass': pass, 'confirm': confirm, 
		'captchaId': captchaId, 'captchaSolution': captchaSolution},
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

$('#user').keypress(function(event) {
    if (event.which === 13) {  
	   $('#pass').focus();
    }
});

$('#pass').keypress(function(event) {
    if (event.which === 13) {  
	   $('#confirm').focus();
    }
});

$('#confirm').keypress(function(event) {
    if (event.which === 13) {  
	   $('#captchaSolution').focus();
    }
});

$('#captchaSolution').keypress(function(event) {
    if (event.which === 13) {  
	   $('#resetpass').click();
    }
});