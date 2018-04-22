var defaultDepth = 3;

function analyzeGameById(gameID){
    board.position('rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR');
    $.ajax({
        url: 'gameAnalysisById',
        type: 'post',
        dataType: 'html',
        data : { 'id': gameID, 'depth': defaultDepth},
        success : function(data) {			

            var analysisTable = `Analysis:<br> <table class="table3" id="engineTable">
                    <tr>
                    <th>Played Move</th><th>Best Move</th>
                    </tr>
            </table>`;
            document.getElementById('analysis').innerHTML = analysisTable;
            var engineTable = document.getElementById("engineTable");

            var analyzedGame = JSON.parse(data);
            console.log(analyzedGame);
            var moves = analyzedGame.Moves;
            for(var i=1; i<moves.length; ++i){
                var move = game.move({
                    from: moves[i].PlayedMoveSrc,
                    to: moves[i].PlayedMoveTar,
                    promotion: moves[i].PlayedMovePromotion 
                });

                $('#engineTable').html(function() { //when a player chats or connects they will update online players list for everyone
                    return  $(this).html() + '<tr  onclick="goToMove(' + i + ');"><td>' + moves[i].PlayedMoveSrc + 
                        moves[i].PlayedMoveTar + moves[i].PlayedMovePromotion +'</td><td>' + 
                        moves[i].BestMoveSrc + moves[i].BestMoveTar + moves[i].BestMovePromotion
                        '</td></tr>';
                });
			
                if (move !== null){
                    //used to store players own move, moves array is stored in memberchess.js
                    var gameFen = moves[i].PlayedMoveFen;
                    board.position(gameFen);
                    totalFEN.push(gameFen);
                    var pgn = game.pgn();
                    totalPGN.push(pgn);
                    ++moveCounter;
                }else{
                    console.log("Null move for:");
                    console.log(i, moves.length);
                }
            }
        }	
    });
}

// When a row is clicked on the analysis engine table the board will update according to the row clicked
function goToMove(row){
    board.position(totalFEN[row]);	
    setPGN(totalPGN[row]);
    moveCounter = row;
}