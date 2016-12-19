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

function getBulletHistory(){
	// NOTE: This function must return the value 
    // from calling the $.ajax() method.
	return $.ajax({
  		url: 'fetchBulletHistory',
   		type: 'post',
   		dataType: 'html',
   		success : function(data) {			
      		console.log("Bullet rating history is: ");
			if(data !== ""){
				console.log(data);		
			}
   		}	
    });
}

function getBlitzHistory(){
	return $.ajax({
  		url: 'fetchBlitzHistory',
   		type: 'post',
   		dataType: 'html',
   		success : function(data) {			
      		console.log("Blitz rating history is: ");
            if(data !== ""){
				console.log(data);	
			}
   		}	
    });
}

function getStandardHistory(){
	return $.ajax({
  		url: 'fetchStandardHistory',
   		type: 'post',
   		dataType: 'html',
   		success : function(data) {			
      		console.log("Standard rating history is: ");
            if(data !== ""){
				console.log(data);		
			}
   		}	
    });
}
getBulletHistory();
getBlitzHistory();
getStandardHistory();

// the code here will be executed when all three ajax requests resolve.
// bullet blitz, standard are lists of length 3 containing the response text,
// status, and jqXHR object for each of the three ajax calls respectively.
$.when(getBulletHistory(), getBlitzHistory(), getStandardHistory()).done(function(bullet, blitz, standard){
	
	var ratingHistory = [];
	
	if (bullet[0] !== ""){

		console.log("Bullet is");
		console.log(bullet[0]);
		var bulletHistory = JSON.parse(bullet[0]);	

		for(var i=0; i<bulletHistory.length; ++i){
			console.log("Date: ")
			console.log(bulletHistory[i].DateTime);
			console.log("Rating: ")
			console.log(bulletHistory[i].Rating);
		}
	}
	
	if(blitz[0] !== ""){

		console.log("Blitz is");
		console.log(blitz[0]);
		var blitzHistory = JSON.parse(blitz[0]);	
		console.log(blitzHistory);

		for(var i=0; i<blitzHistory.length; ++i){
			console.log("Date: ")
			console.log(blitzHistory[i].DateTime);
			console.log("Rating: ")
			console.log(blitzHistory[i].Rating);
		}
	}
	if(standard[0] !== ""){

		console.log("Standard is");
		console.log(standard[0]);
		var standardHistory = JSON.parse(standard[0]);

		for(var i=0; i<standardHistory.length; ++i){
			console.log("Date: ")
			console.log(standardHistory[i].DateTime);
			console.log("Rating: ")
			console.log(standardHistory[i].Rating);
		}
	}
	/* 
	google.charts.load('current', {'packages':['line']});
    google.charts.setOnLoadCallback(drawChart);

    function drawChart() {

      var data = new google.visualization.DataTable();
      data.addColumn('date', 'Day');
      data.addColumn('number', 'Bullet');
      data.addColumn('number', 'Blitz');
      data.addColumn('number', 'Standard');

      data.addRows([
        [1,  37.8, 80.8, 41.8],
        [2,  30.9, 69.5, 32.4],
        [3,  25.4,   57, 25.7],
        [4,  11.7, 18.8, 10.5],
        [5,  11.9, 17.6, 10.4],
        [6,   8.8, 13.6,  7.7],
        [7,   7.6, 12.3,  9.6],
        [8,  12.3, 29.2, 10.6],
        [9,  16.9, 42.9, 14.8],
        [10, 12.8, 30.9, 11.6],
        [11,  5.3,  7.9,  4.7],
        [12,  6.6,  8.4,  5.2],
        [13,  4.8,  6.3,  3.6],
        [14,  4.2,  6.2,  3.4]
      ]);

      var options = {
        chart: {
          title: 'Box Office Earnings in First Two Weeks of Opening',
          subtitle: 'in millions of dollars (USD)'
        },
        width: 900,
        height: 500
      };

      var chart = new google.charts.Line(document.getElementById('linechart_material'));

      chart.draw(data, options);
    }
	*/
});

// parses JSON rating history string and returns 
function parseRatingstring(ratingHistory){
	var bulletHistory = JSON.parse(bullet);	
	for (var key in data) {
		// skip loop if the property is from prototype
		if (!data.hasOwnProperty(key)){
			continue;
		}

		var option = document.createElement('option');
		option.text = key + ": " + data[key].name;
		option.value = key;
		openingDropDown.add(option);
	}
}