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

function getBulletHistory(lookupName){
	// NOTE: This function must return the value 
    // from calling the $.ajax() method.
	return $.ajax({
  		url: 'fetchBulletHistory',
   		type: 'post',
   		dataType: 'html',
		data : {'user': lookupName},
   		success : function(data) {			
			//console.log(data);		
   		}	
    });
}

function getBlitzHistory(lookupName){
	return $.ajax({
  		url: 'fetchBlitzHistory',
   		type: 'post',
   		dataType: 'html',
		data : {'user': lookupName},
   		success : function(data) {			
			//console.log(data);	
   		}	
    });
}

function getStandardHistory(lookupName){
	return $.ajax({
  		url: 'fetchStandardHistory',
   		type: 'post',
   		dataType: 'html',
		data : {'user': lookupName},
   		success : function(data) {			
			//console.log(data);		
   		}	
    });
}

function getCorrespondenceHistory(lookupName){
	return $.ajax({
  		url: 'fetchCorrespondenceHistory',
   		type: 'post',
   		dataType: 'html',
		data : {'user': lookupName},
   		success : function(data) {			
			//console.log(data);		
   		}	
    });
}

function setupRatingChart(){
	
	var url = parseUrl();
	var lookupName = url.name;
	var name = lookupName;
	if(typeof lookupName !== "undefined"){
		// then do nothing here
	}else{ // this means the player is looking at his/her own profile
		name  = document.getElementById('user').value;
	}

	// the code here will be executed when all three ajax requests resolve.
	// bullet blitz, standard, correspondence are lists of length 3 containing the response text,
	// status, and jqXHR object for each of the three ajax calls respectively.
	$.when(getBulletHistory(name), getBlitzHistory(name), getStandardHistory(name), 
		getCorrespondenceHistory(name)).done(function(bullet, blitz, standard, correspondence){
		
		var ratingHistory = [];
		var showChart = false;

		if (bullet[0] !== ""){
			showChart = true;
			var bulletHistory = JSON.parse(bullet[0]);	

			for(var i=0; i<bulletHistory.length; ++i){
				var oneGame = [];
				var dateString = bulletHistory[i].DateTime;
				var year = dateString.substring(0, 4);
				var month = dateString.substring(4, 6);
				var day = dateString.substring(6, 8);
				var hour = dateString.substring(8, 10);
				var minute = dateString.substring(10, 12);
				var second = dateString.substring(12, 14);
				
				// 00 is Jan, 01 is Feb, 02 is March so month needs to be subtracted by 1 for zero indexing
				oneGame.push(new Date(year, month-1, day, hour, minute, second), 
					bulletHistory[i].Rating, null, null, null);
				ratingHistory.push(oneGame);
			}
		}
		
		if(blitz[0] !== ""){
			showChart = true;
			//console.log("Blitz is");
			//console.log(blitz[0]);
			var blitzHistory = JSON.parse(blitz[0]);	

			for(var i=0; i<blitzHistory.length; ++i){		
				var oneGame = [];
				var dateString = blitzHistory[i].DateTime;
				var year = dateString.substring(0, 4);	
				var month = dateString.substring(4, 6);
				var day = dateString.substring(6, 8);
				var hour = dateString.substring(8, 10);		
				var minute = dateString.substring(10, 12);
				var second = dateString.substring(12, 14);

				oneGame.push(new Date(year, month-1, day, hour, minute, second), null,
					blitzHistory[i].Rating, null, null);
				ratingHistory.push(oneGame);
			}
		}

		if(standard[0] !== ""){
			showChart = true;
			var standardHistory = JSON.parse(standard[0]);

			for(var i=0; i<standardHistory.length; ++i){
				var oneGame = [];
				var dateString = standardHistory[i].DateTime;
				var year = dateString.substring(0, 4);
				var month = dateString.substring(4, 6);
				var day = dateString.substring(6, 8);
				var hour = dateString.substring(8, 10);
				var minute = dateString.substring(10, 12);
				var second = dateString.substring(12, 14);
				oneGame.push(new Date(year, month-1, day, hour, minute, second),null, 
					null, standardHistory[i].Rating, null);
				ratingHistory.push(oneGame);
			}
		}

		if(correspondence[0] !== ""){
			showChart = true;
			var correspondenceHistory = JSON.parse(standard[0]);

			for(var i=0; i<correspondenceHistory.length; ++i){
				var oneGame = [];
				var dateString = correspondenceHistory[i].DateTime;
				var year = dateString.substring(0, 4);
				var month = dateString.substring(4, 6);
				var day = dateString.substring(6, 8);
				var hour = dateString.substring(8, 10);
				var minute = dateString.substring(10, 12);
				var second = dateString.substring(12, 14);
				oneGame.push(new Date(year, month-1, day, hour, minute, second),null, null,
					null, correspondenceHistory[i].Rating);
				ratingHistory.push(oneGame);
			}
		}

		if(showChart){
			
			google.charts.load('current', {'packages':['line']});
			google.charts.setOnLoadCallback(function() { drawChart(ratingHistory) });

			function drawChart(ratingHistory) {

				var data = new google.visualization.DataTable();

				data.addColumn('date', 'Day');
				data.addColumn('number', 'Bullet');
				data.addColumn('number', 'Blitz');
				data.addColumn('number', 'Standard');
				data.addColumn('number', 'Correspondence');

				data.addRows(ratingHistory);

				var formatter = new google.visualization.NumberFormat({
					formatType: 'long'
				});
				formatter.format(data, 1); // Apply formatter to second column

				var formatter_medium = new google.visualization.DateFormat({formatType: 'long'});
				formatter_medium.format(data, 0);

				var options = {
					chart: {
						title: 'Rating History',
					},
					vAxis: { Title: 'Rating' },
					hAxis: { Title: 'Day'},
					width: 1050,
					height: 300
				};
				var chart = new google.charts.Line(document.getElementById('googleLinechart'));
				chart.draw(data, options);
			}
		}
	});
}
setupRatingChart();

function parseUrl() { //fetches all variables in url and returns them in a json struct
	var query = location.search.substr(1);
	var result = {};
	query.split("&").forEach(function(part) {
    	var item = part.split("=");
    	result[item[0]] = decodeURIComponent(item[1]);
	});
  return result;
}