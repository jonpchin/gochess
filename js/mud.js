if (!window.WebSocket){
	$('#checkwebsocket').html("Your browser doesn't support websockets." + 
		"Please use the latest version of Firefox, Chrome, IE, Opera or Microsoft Edge.");
}
var wsuri = "wss://"+ window.location.host +"/mudserver";
var sock = new WebSocket(wsuri);
$(window).on('beforeunload', function(){
	sock.close();
}); 

window.onload = function() {
    sock.onopen = function() {
		
		var message = {
			Type: "connect_mud"
		}
	    sock.send(JSON.stringify(message));

    }

    sock.onclose = function(e) {

	}

     sock.onmessage = function(e) {
		
		json = JSON.parse(e.data);

        if (json.Type === "chat_all"){
        }
        else if(json.Type === "test"){
            
        }
     }

     

};