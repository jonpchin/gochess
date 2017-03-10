
function testGetPlayerData(){
    $('#get-player-data').html('<img src="../img/ajax/loading.gif" />');
    $.ajax({
  		url: 'server/getPlayerData',
   		type: 'post',
		dataType: 'html',
   		data : { 'user': "Jon"},
   		success : function(data) {
            if(data !== "Service is down."){
                $('#get-player-data').html('<img src="../img/ajax/available.png" />');
            }else{
                $('#get-player-data').html('<img src="../img/ajax/not-available.png" />');
            }	
   		}	
    });
}

function testResumeGame(){
    $('#resume-game').html('<img src="../img/ajax/loading.gif" />');
    $.ajax({
  		url: 'resumeGame',
   		type: 'post',
		dataType: 'html',
   		data : { 'id': "1", 'white': "Jon", 'black': "Carl"},
   		success : function(data) {
            if(data === "false"){
                $('#resume-game').html('<img src="../img/ajax/available.png" />');
            }else{
                $('#resume-game').html('<img src="../img/ajax/not-available.png" />');
            }
   		}	
    });
}

function testFetchGameById1(){
    $('#fetch-game-by-id-1').html('<img src="../img/ajax/loading.gif" />');
    $.ajax({
  		url: 'fetchgameID',
   		type: 'post',
		dataType: 'html',
   		data : { 'gameID': "24"},
   		success : function(data) {
            var result = JSON.parse(data);
            if(result.Black != "Socko, M" || result.BlackElo != "2448" || result.ECO != "A07" ||
                result.Event != "IMSA Blitz w 2016" || result.Site != "Huai'an CHN" || 
                result.White != "Khotenashvili, Bela"){
                console.log("Data is:")
                console.log(result);
                $('#fetch-game-by-id-1').html('<img src="../img/ajax/not-available.png" />');
            }else{
                $('#fetch-game-by-id-1').html('<img src="../img/ajax/available.png" />');
            }
   		}	
    });
}

function testFetchGameById2(){
    $('#fetch-game-by-id-2').html('<img src="../img/ajax/loading.gif" />');
    $.ajax({
  		url: 'fetchgameID',
   		type: 'post',
		dataType: 'html',
   		data : { 'gameID': "1856326"},
   		success : function(data) {
            var result = JSON.parse(data);
            if(result.Black != "Georgiev, K" || result.BlackElo != "2455" || result.ECO != "E81" ||
                result.Event != "Stara Zagora II zt" || result.Site != "Stara Zagora II zt" || 
                result.White != "Polgar, Z"){
                console.log("Data is:")
                console.log(result);
                $('#fetch-game-by-id-2').html('<img src="../img/ajax/not-available.png" />');
            }else{
                $('#fetch-game-by-id-2').html('<img src="../img/ajax/available.png" />');
            }
   		}	
    });
}

function testFetchGameByECO(){
    $('#fetch-game-by-eco').html('<img src="../img/ajax/loading.gif" />');
    $.ajax({
  		url: 'fetchgameByECO',
   		type: 'post',
		dataType: 'html',
   		data : { 'ECO': "B20", 'ECOIndex': "23"},
   		success : function(data) {
            var result = JSON.parse(data);
            if(result.Black != "Lagashin, Pavel" || result.BlackElo != "2139" || result.ECO != "B20" ||
                result.Event != "Moscow Open A 2016" || result.Site != "Moscow RUS" || 
                result.White != "Afanasiev, Nikita"){
                console.log("Data is:")
                console.log(result);
                $('#fetch-game-by-eco').html('<img src="../img/ajax/not-available.png" />');
            }else{
                $('#fetch-game-by-eco').html('<img src="../img/ajax/available.png" />');
            }
   		}	
    });
}

function testCheckUserName(){
    $('#check-username').html('<img src="../img/ajax/loading.gif" />');
    $.ajax({
  		url: 'checkname',
   		type: 'post',
		dataType: 'html',
   		data : { 'username': "Jon"},
   		success : function(data) {
            if(data === "<img src='img/ajax/not-available.png' /> Username taken"){
                $('#check-username').html('<img src="../img/ajax/available.png" />');         
            }else{
                $('#check-username').html('<img src="../img/ajax/not-available.png" />');
            }
   		}	
    });
}

function testFetchBulletHistory(){
    $('#fetch-bullet-history').html('<img src="../img/ajax/loading.gif" />');
    $.ajax({
  		url: 'fetchBulletHistory',
   		type: 'post',
		dataType: 'html',
   		data : { 'user': "Jon"},
   		success : function(data) {
            if(data !== ""){
                $('#fetch-bullet-history').html('<img src="../img/ajax/available.png" />');     
            }else{
                $('#fetch-bullet-history').html('<img src="../img/ajax/not-available.png" />');
            }
   		}	
    });
}

function testFetchBlitzHistory(){
    $('#fetch-blitz-history').html('<img src="../img/ajax/loading.gif" />');
    $.ajax({
  		url: 'fetchBlitzHistory',
   		type: 'post',
		dataType: 'html',
   		data : { 'user': "Jon"},
   		success : function(data) {
            if(data !== ""){
                $('#fetch-blitz-history').html('<img src="../img/ajax/available.png" />');     
            }else{
                $('#fetch-blitz-history').html('<img src="../img/ajax/not-available.png" />');
            }
   		}	
    });
}

function testFetchStandardHistory(){
    $('#fetch-standard-history').html('<img src="../img/ajax/loading.gif" />');
    $.ajax({
  		url: 'fetchStandardHistory',
   		type: 'post',
		dataType: 'html',
   		data : { 'user': "Jon"},
   		success : function(data) {
            if(data !== ""){
                $('#fetch-standard-history').html('<img src="../img/ajax/available.png" />');     
            }else{
                $('#fetch-standard-history').html('<img src="../img/ajax/not-available.png" />');
            }
   		}	
    });
}

function testFetchCorrespondenceHistory(){
    $('#fetch-correspondence-history').html('<img src="../img/ajax/loading.gif" />');
    $.ajax({
  		url: 'fetchCorrespondenceHistory',
   		type: 'post',
		dataType: 'html',
   		data : { 'user': "Jon"},
   		success : function(data) {
            if(data !== ""){
                $('#fetch-correspondence-history').html('<img src="../img/ajax/available.png" />');     
            }else{
                $('#fetch-correspondence-history').html('<img src="../img/ajax/not-available.png" />');
            }
   		}	
    });
}

function testCheckInGame(){
    $('#check-in-game').html('<img src="../img/ajax/loading.gif" />');
    $.ajax({
  		url: 'checkInGame',
   		type: 'post',
		dataType: 'html',
   		data : { 'user': "Jon"},
   		success : function(data) {
            if(data === "Safe"){
                $('#check-in-game').html('<img src="../img/ajax/available.png" />');     
            }else{
                $('#check-in-game').html('<img src="../img/ajax/not-available.png" />');
            }
   		}	
    });
}

function runAjaxUnitTests(){
    testGetPlayerData();
    testResumeGame();
    testFetchGameById1();
    testFetchGameById2();
    testFetchGameByECO();
    testCheckUserName();
    testFetchBulletHistory();
    testFetchBlitzHistory();
    testFetchStandardHistory();
    testFetchCorrespondenceHistory();
    testCheckInGame();
}

runAjaxUnitTests();