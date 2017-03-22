document.getElementById('createThread').onclick = function(){
    window.location = "/createpost?forumid=" + document.getElementById('createThread').value; 
}