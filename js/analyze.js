var defaultDepth = 3;

function analyzeGameById(gameID){
    $.ajax({
        url: 'gameAnalysisById',
        type: 'post',
        dataType: 'html',
        data : { 'id': gameID, 'depth': defaultDepth},
        success : function(data) {			
            var game = JSON.parse(data);
            console.log(game);
            //var moves = game.Moves;
            //console.log(moves);
            //for(var i=0; i<moves.length; ++i){
             //   var playedMove = moves[i].PlayedMove;
             //   console.log(playedMove);
             //   board.position(playedMove);
            //    totalFEN.push(playedMove);
            //}
        }	
    });
}