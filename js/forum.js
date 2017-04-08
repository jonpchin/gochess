if (document.getElementById('createThread') != null){
    document.getElementById('createThread').onclick = function(){
        window.location = "/createthread?forumname=" + document.getElementById('createThread').value; 
    }
}

// Use data.value to get the id of the thread
function updateThreadLock(data){
    if (document.getElementById(data.value).innerHTML === "Unlock Thread"){
        $.ajax({
            url: 'unlockThread',
            type: 'post',
            dataType: 'html',
            data : {'id': data.value},
            success : function(unused) {
                document.getElementById(data.value).innerHTML = "Lock Thread"		
            }	
        });
        
    }else{
        $.ajax({
            url: 'lockThread',
            type: 'post',
            dataType: 'html',
            data : {'id': data.value},
            success : function(unused) {
                document.getElementById(data.value).innerHTML = "Unlock Thread"		
            }	
        });
    }
}

function lockThread(data){
    
}
