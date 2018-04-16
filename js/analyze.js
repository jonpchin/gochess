var defaultDepth = 3;

function analyzeGameById(gameID){
    $.ajax({
        url: 'gameAnalysisById',
        type: 'post',
        dataType: 'html',
        data : { 'id': gameID, 'depth': defaultDepth},
        success : function(data) {			
            // error messages will be less then 100 characters, games always more then 100 characters
            console.log(data);
        }	
    });
}