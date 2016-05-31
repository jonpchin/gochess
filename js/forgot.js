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

document.getElementById('submit').onclick = function(){
	document.getElementById('submit').disabled = true;
	setTimeout(function() {
		// enable click after 1 second
		document.getElementById('submit').disabled = false;
	}, 1000); //  second delay
	
	var user = document.getElementById('user').value;
	var email = document.getElementById('email').value;
	var captchaId = document.getElementById('captchaId').value;
	var captchaSolution = document.getElementById('captchaSolution').value;

	if(user === "" || email === "" || captchaSolution === ""){
		$('#submit-result').html("<img src='img/ajax/not-available.png' /> Please fill out all fields.");
		return
	}
	
	$('#submit-result').html('<img src="img/ajax/loading.gif" />Validating information...');
    $.ajax({
  		url: 'processForgot',
   		type: 'post',
   		dataType: 'html',
   		data : { 'user': user, 'email': email, 'captchaId': captchaId, 'captchaSolution': captchaSolution},
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