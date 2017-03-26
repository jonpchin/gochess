document.getElementById('createfirstpost').onclick = function(){
    
	var user = document.getElementById('user').value;
	var forumname = document.getElementById('forumname').value;
	var title = document.getElementById('title').value;
	var message = document.getElementById('message').value;
	var date = timeStamp();

    $.ajax({
  		url: 'sendFirstForumPost',
   		type: 'post',
   		dataType: 'html',
   		data : {'forumname': forumname, 'title': title, 'message': message},
   		success : function(data) {
			if(data === "<img src='img/ajax/not-available.png' /> Invalid credentials"){
				$('#submit-result').html(data);	
			}
			else{
				$('#thread-title').html(title);	
				var table = '<table class="table1"><tr><th></th></tr> <tr><td>' + user + " " + 
					date +'</td><td></td>' + message + '</tr></table>'	
				$('#submit-result').html(table);	
			}			
   		}	
    });
    //window.location = "/thread?forumname=" + document.getElementById('createThread').value; 
}

function timeStamp(){ //returns timestamp
	var currentdate = new Date(); 
	var datetime =  + currentdate.getHours() + ":"  
            		+ currentdate.getMinutes() + ":" 
            		+ currentdate.getSeconds();
	return datetime;
}