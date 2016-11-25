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

$(document).ready(function() {
    var x_timer;    
    $('#username').keyup(function (e){
        clearTimeout(x_timer);
        var user_name = $(this).val();
        x_timer = setTimeout(function(){
            check_username_ajax(user_name);
        }, 1000);
    });

	function check_username_ajax(username){
		var length = username.length;
		if(length < 3 ){
			$('#user-result').html("<img src='img/ajax/not-available.png' /> Too short.");
			return
		}
	    $('#user-result').html('<img src="img/ajax/loading.gif" />');
	      $.ajax({
	      	url: 'checkname',
	      	type: 'post',
	      	dataType: 'html',
	      	data : { 'username': username},
	      	success : function(data) {
		        $('#user-result').html(data);
	      	},
	    });
	}

});

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

document.getElementById('register').onclick = function(){
	document.getElementById('register').disabled = true;
	setTimeout(function() {
		// enable click after 1 second
		document.getElementById('register').disabled = false;
	}, 1000); // 1 second delay
	
	var user = document.getElementById('username').value;
	var pass = document.getElementById('pass').value;
	var confirm = document.getElementById('confirm').value;
	var email = document.getElementById('email').value;	
	var captchaId = document.getElementById('captchaId').value;
	var captchaSolution = document.getElementById('captchaSolution').value;
	
	if(user === "" || pass === "" || confirm === "" || email === ""  || captchaSolution === ""){
		$('#submit-result').html("<img src='img/ajax/not-available.png' /> Please fill out all fields.");
		return
	}
	
	$('#submit-result').html('<img src="img/ajax/loading.gif" />Submission in progress...');
    $.ajax({
  		url: 'processRegister',
   		type: 'post',
   		dataType: 'html',
   		data : { 'username': user, 'pass': pass, 'confirm': confirm,  'email': email,
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

//when user presses enter on username input it will jump to password input
$('#username').keypress(function(event) {
    if (event.which === 13) {  
	   $('#pass').focus();
    }
});

//when user presses enter on password input it will jump to confirm input
$('#pass').keypress(function(event) {
    if (event.which === 13) {  
	   $('#confirm').focus();
    }
});

//when user presses enter on confirm input it will jump to email input
$('#confirm').keypress(function(event) {
    if (event.which === 13) {  
	   $('#email').focus();
    }
});

//when user press enter on email it will jump to captcha
$('#email').keypress(function(event) {
    if (event.which === 13) {  
	   $('#captchaSolution').focus();
    }
});
//when user presses enter on captcha it will submit the form
$('#captchaSolution').keypress(function(event) {
    if (event.which === 13) {  
	   $('#register').click();
    }
});

function loadFlag(){
	$.ajax({
  		url: 'getCountry',
   		type: 'post',
   		success : function(data) {			
			document.getElementById('countryFlag').src = "img/flags/" + data + ".png"
   		},
	});
}
// Load country flag on page load
loadFlag();