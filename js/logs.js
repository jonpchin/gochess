function fetchLogs(){
    $.ajax({
        url: 'fetchLogs',
        type: 'post',
        dataType: 'html',
        data : {'logType': 'chat'},
        success : function(data) {
            document.getElementById('chatLog').innerHTML = data;		
        }	
    });

    $.ajax({
        url: 'fetchLogs',
        type: 'post',
        dataType: 'html',
        data : {'logType': 'errors'},
        success : function(data) {
            document.getElementById('errorLog').innerHTML = data;			
        }	
    });

    $.ajax({
        url: 'fetchLogs',
        type: 'post',
        dataType: 'html',
        data : {'logType': 'main'},
        success : function(data) {
            document.getElementById('mainLog').innerHTML = data;		
        }	
    });
}

fetchLogs();