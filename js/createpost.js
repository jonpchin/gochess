if (document.getElementById('createpost') != null){
	document.getElementById('createpost').onclick = function(){
		
		var user = document.getElementById('user').value;
		var forumname = document.getElementById('forumname').value;
		var title = document.getElementById('title').value;
		var message = document.getElementById('message').value;
		var date = timeStamp();
		var firstPost = document.getElementById('firstPost').value;
		var totalPosts = document.getElementById('totalPosts').value;
		var threadId = document.getElementById('threadId').value;

		if(title !== "" && message !== ""){
			$.ajax({
				url: 'sendForumPost',
				type: 'post',
				dataType: 'html',
				data : {'forumname': forumname, 'title': title, 'message': message,
					'firstPost': firstPost, 'totalPosts': totalPosts, 'threadId': threadId},
				success : function(data) {
					if(data !== "createPost"){
						$('#postCenterTable').html(data);	
					}
					else if(firstPost === "Yes"){ // Creates new thread
						$('#thread-title').html(title);	
						var table = '<table class="table1"><tr><th>'+ user + " " + 
							date + '</th></tr><tr><td>'+ message +'</td></tr></table>';	
						$('#postCenterTable').html(table);
						document.getElementById("createpost").disabled = true;	
					}else{ // Add a post
						$('#postCenterTable').html(function() {
							return  $(this).html() + '<tr><th>' + user + ' <span class="right">' + date +
								'</span></th></tr><tr><td>' + message + '</td></tr>';
						});
						document.getElementById("createpost").disabled = true;	
					}			
				}	
			});
		}
	}
}

function getPlayerData(lookUpName){
	$('#playerData').html('<img src="../img/ajax/loading.gif" />');
    $.ajax({
  		url: 'server/getPlayerData',
   		type: 'post',
		dataType: 'html',
   		data : { 'user': lookUpName},
   		success : function(data) {			
      		$('#playerData').html(data);	
   		}	
    });
}

function timeStamp(){ //returns timestamp
	var currentdate = new Date(); 
	var datetime =  + currentdate.getHours() + ":"  
            		+ currentdate.getMinutes() + ":" 
            		+ currentdate.getSeconds();
	return datetime;
}