document.getElementById('createfirstpost').onclick = function(){
    
	var forumname = document.getElementById('forumname').value;

    $.ajax({
  		url: 'sendFirstForumPost',
   		type: 'post',
   		dataType: 'html',
   		data : {'forumname': forumname, 'title': title, 'message': message},
   		success : function(data) {			
      		$('#submit-result').html(data);	
   		}	
    });
    
    window.location = "/thread?forumname=" + document.getElementById('createThread').value; 
}