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

function copyParams(){ //checks the url parameters if there any and copies them to input box
	
	var url = parseUrl();
	var user = url.user;

	if (typeof user !== "undefined"){
		var token = url.token;
		document.getElementById('user').value = user;
		document.getElementById('token').value = token;
		
		//focus the the cursor onto the captcha box
		$('#captchaSolution').focus();
		   
	}
}

function parseUrl() { //fetches all variables in url and returns them in a json struct
  var query = location.search.substr(1);
  var result = {};
  query.split("&").forEach(function(part) {
    var item = part.split("=");
    result[item[0]] = decodeURIComponent(item[1]);
  });
  return result;
}

copyParams();

document.getElementById('register').onclick = function(){
	document.getElementById('register').disabled = true;
	setTimeout(function() {
		// enable click after 1 second
		document.getElementById('register').disabled = false;
	}, 1000); //  second delay
	
	var user = document.getElementById('user').value;
	var token = document.getElementById('token').value;
	var captchaId = document.getElementById('captchaId').value;
	var captchaSolution = document.getElementById('captchaSolution').value;
	
	if(user === "" || token === "" || captchaSolution === ""){
		$('#submit-result').html("<img src='img/ajax/not-available.png' /> Please fill out all fields.");
		return
	}
	
	$('#submit-result').html('<img src="img/ajax/loading.gif" />Submission in progress...');
    $.ajax({
  		url: 'processActivate',
   		type: 'post',
   		dataType: 'html',
   		data : { 'user': user, 'token': token, 'captchaId': captchaId, 'captchaSolution': captchaSolution},
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

//when user presses enter on captcha it will submit the form
$('#captchaSolution').keypress(function(event) {
    if (event.which === 13) {  
	   $('#register').click();
    }
});