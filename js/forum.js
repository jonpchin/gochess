document.getElementById('createThread').onclick = function(){
    window.location = "/createthread?forumname=" + document.getElementById('createThread').value; 
}