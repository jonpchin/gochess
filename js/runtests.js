
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
    
}

function testFetchGameByID(){

}

function testFetchGameByECO(){

}

function testCheckUserName(){

}

function testFetchBulletHistory(){

}

function testFetchBlitzHistory(){

}

function testFetchStandardHistory(){

}

function testFetchCorrespondenceHistory(){

}

function testCheckInGame(){

}

function runAjaxUnitTests(){
    testGetPlayerData();
    testResumeGame();
    testFetchGameByID();
    testFetchGameByECO();
    testCheckUserName();
    testFetchBulletHistory();
    testFetchBlitzHistory();
    testFetchStandardHistory();
    testFetchCorrespondenceHistory();
    testCheckInGame();
}

runAjaxUnitTests();