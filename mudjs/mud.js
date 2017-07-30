if (!window.WebSocket){
	$('#checkwebsocket').html("Your browser doesn't support websockets." + 
		"Please use the latest version of Firefox, Chrome, IE, Opera or Microsoft Edge.");
}
var wsuri = "wss://"+ window.location.host +"/mudserver";
var sock = new WebSocket(wsuri);
var codeMirror = CodeMirror.fromTextArea(document.getElementById('textbox'), {
    lineNumbers: true,
    lineWrapping: true,
    readOnly: true
});

// Game state encapsulates all information relevant to the client
var GameState = {
    status: "connect", // Determines what kind of message gets sent over the websocket
    name: "",          // Player's adventurer name
    ingame: false      // If true then player is in the MUD world
};

window.onload = function() {
    
    // Removes blank space at the start
    codeMirror.setValue(""); 
    // Focuses mouse cursor on load
    $("#message").focus();
    
    sock.onopen = function() {
	    displayToTextBox(`                                                                                                        dddddddd
        GGGGGGGGGGGGG                      MMMMMMMM               MMMMMMMM                              d::::::d
     GGG::::::::::::G                      M:::::::M             M:::::::M                              d::::::d
   GG:::::::::::::::G                      M::::::::M           M::::::::M                              d::::::d
  G:::::GGGGGGGG::::G                      M:::::::::M         M:::::::::M                              d:::::d 
 G:::::G       GGGGGG   ooooooooooo        M::::::::::M       M::::::::::Muuuuuu    uuuuuu      ddddddddd:::::d 
G:::::G               oo:::::::::::oo      M:::::::::::M     M:::::::::::Mu::::u    u::::u    dd::::::::::::::d 
G:::::G              o:::::::::::::::o     M:::::::M::::M   M::::M:::::::Mu::::u    u::::u   d::::::::::::::::d 
G:::::G    GGGGGGGGGGo:::::ooooo:::::o     M::::::M M::::M M::::M M::::::Mu::::u    u::::u  d:::::::ddddd:::::d 
G:::::G    G::::::::Go::::o     o::::o     M::::::M  M::::M::::M  M::::::Mu::::u    u::::u  d::::::d    d:::::d 
G:::::G    GGGGG::::Go::::o     o::::o     M::::::M   M:::::::M   M::::::Mu::::u    u::::u  d:::::d     d:::::d 
G:::::G        G::::Go::::o     o::::o     M::::::M    M:::::M    M::::::Mu::::u    u::::u  d:::::d     d:::::d 
 G:::::G       G::::Go::::o     o::::o     M::::::M     MMMMM     M::::::Mu:::::uuuu:::::u  d:::::d     d:::::d 
  G:::::GGGGGGGG::::Go:::::ooooo:::::o     M::::::M               M::::::Mu:::::::::::::::uud::::::ddddd::::::dd
   GG:::::::::::::::Go:::::::::::::::o     M::::::M               M::::::M u:::::::::::::::u d:::::::::::::::::d
     GGG::::::GGG:::G oo:::::::::::oo      M::::::M               M::::::M  uu::::::::uu:::u  d:::::::::ddd::::d
        GGGGGG   GGGG   ooooooooooo        MMMMMMMM               MMMMMMMM    uuuuuuuu  uuuu   ddddddddd   ddddd `)
        displayToTextBox("");
        displayToTextBox("");
        displayToTextBox("Welcome to Go MUD!");
		var message = {
			Type: "connect_mud"
		}
	    sock.send(JSON.stringify(message));
    }

    sock.onclose = function(e) {
        console.log("Socket closing")
	}

     sock.onmessage = function(e) {
        json = JSON.parse(e.data);

        switch(json.Type){
            case "ask_name":
                // There seems to be a default tab spacing
                displayToTextBox("What is your name?", "forestgreen");
                GameState.status="check_name";
                break;
            case "name_taken":
                displayToTextBox("That name is already taken.", "red");
                displayToTextBox("What is your name?", "forestgreen");
                GameState.status="check_name";
                break;
            case "name_available":
                GameState.status="save_name";
                determineMessageType(json.Name); // Saving player name and asking for class
                break;
            case "update_player":
                // TODO: Unmarshal player data that was sent from server into client memory
                console.log(json);
                updatePlayer(json);
                break;
            case "enter_world":
                break;
            case "look":
                displayMap(json.Map);
                break;
            default:
                console.log("No such socket type", json.Type);
        }
     }
};

function displayToTextBox(message, textColor){
    
    // Default text black as the default color unless the background is white,
    // then the text would default to white
    if(textColor){
       //console.log("Font color is set!");
    }else{
        // TODO: Check the background color and set the default text color appropriately
        textColor="#000000";
    }
    var lastLine = codeMirror.lastLine();
    codeMirror.replaceRange(message + "\n", CodeMirror.Pos(lastLine));
    codeMirror.markText({line:lastLine-1,ch:0},{line:lastLine-1,ch:lastLine.length},{css:"color: " + textColor});
    codeMirror.scrollTo(0, codeMirror.getScrollInfo().height);
}

// Displays to textbox without appending the newline, used for printing out the map
function displayToTextBoxNoNewLine(message, textColor){

    // Default text black as the default color unless the background is white,
    // then the text would default to white
    if(textColor){
       //console.log("Font color is set!");
    }else{
        // TODO: Check the background color and set the default text color appropriately
        textColor="#000000";
    }
    var lastLine = codeMirror.lastLine();
    codeMirror.replaceRange(message, CodeMirror.Pos(lastLine));
    codeMirror.markText({line:lastLine-1,ch:lastLine.length-1},{line:lastLine-1,ch:lastLine.length},{css:"color: " + textColor});
}

document.getElementById('sendButton').onclick = function(){
    var message = document.getElementById('message').value
    displayToTextBox(message);
    determineMessageType(message);
}

// Checks if name is available
function checkName(name){
    var message = {
        Type: "check_name",
        Name: name
    }
    sock.send(JSON.stringify(message));
}

function enterWorldFirstTime(){
    Mud.Player.Type = "enter_world_first_time";
    sock.send(JSON.stringify(Mud.Player));
    GameState.ingame = true;
}


function fetchMap(){
    
    var credentials = {
        Username: Mud.Player.Username,
        Name: Mud.Player.Name,
        SessionID: Mud.Player.SessionID
    }
    var message = {
        Type: "fetch_map",
        Creds: credentials
    }
    sock.send(JSON.stringify(message));
}

function searchCommands(command){

    $.getJSON('../data/mud/commands.json', function(data) {        
        
        for (var key in data) {
            // skip loop if the property is from prototype
            if (!data.hasOwnProperty(key)){
                continue;
            }
            if(command.startsWith(key)){
                return key;
            }
        }
    });
    return "";
}

function determineMessageType(message){
    var status = GameState.status;
    
    if (GameState.ingame){
        // Searches command for matching starting substring, all commands are lowercase
        switch(searchCommands(message).toLowerCase()){
            case "":
                displayToTextBox("I do not understand.");
                break;
            case "look":
                fetchMap();
                break;
            case "north":
                displayToTextBox("You walk north.");
                break;
            case "east":
                displayToTextBox("You walk east.");
                break;
            case "south":
                displayToTextBox("You walk south.");
                break;
            case "west":
                displayToTextBox("You walk west.");
                break;
            default:
                // This means its a command in commands.json but its not registered in the switch statement
                displayToTextBox("That does not make sense.");
                break;
        }
    }else{
         switch(status){
            case "connect":
                break;
            case "check_name":
                checkName(message);
                break;
            case "save_name":
                savePlayerData(status, message);
                askClass();
                break;
            case "register_class":
                if(isValidClass(message) === false){
                    displayToTextBox("That is not a valid class. Please try again.");
                    askClass();
                }else{
                    savePlayerData(status, message);
                    askRace();
                }
                break;
            case "save_race":
                if(isValidRace(message) === false){
                    displayToTextBox("That is not a valid race. Please try again.");
                    askRace();
                }else{
                    savePlayerData(status, message);
                    enterWorldFirstTime();
                }
                break;
            default:
                console.log("No matching message type for", status);
        }
    }

   
}

function askClass(){
    displayToTextBox("Select your class. Options are: " + Mud.Classes.join(", "), "forestgreen");
    GameState.status="register_class";
} 

function askRace(){
    displayToTextBox("Select your race. Options are: " + Mud.Races.join(", "), "forestgreen");
    GameState.status="save_race";
}

function isValidClass(mudClass){
    for(var i=0; i<Mud.Classes.length; ++i){
        if(Mud.Classes[i].toLowerCase() === mudClass.toLowerCase()){
            return true;
        }
    }
    return false;
}

function isValidRace(mudRace){
    for(var i=0; i<Mud.Races.length; ++i){
        if(Mud.Races[i].toLowerCase() === mudRace.toLowerCase()){
            return true;
        }
    }
    return false;
}
function savePlayerData(type, message){
    
    switch(type){
        case "save_name":
            Mud.Player.Name = message;
            break;
        case "register_class":
            Mud.Player.Class = message;
            break;
        case "save_race":
            Mud.Player.Race = message;
            break;
        default:
            console.log(type, "is not a valid savePlayerData type");
    }
}

// If enter is pressed auto submit
$('#message').keypress(function(event) {
    if (event.which === 13) {
        $('#sendButton').click();	
        if (document.getElementById('forget').checked) {  
            document.getElementById('message').value = "";
        }
    }
});

// Takes a colored string and returns it hex
// Returns false if there was an error
function colourNameToHex(colour)
{
    var colours = {"aliceblue":"#f0f8ff","antiquewhite":"#faebd7","aqua":"#00ffff","aquamarine":"#7fffd4","azure":"#f0ffff",
    "beige":"#f5f5dc","bisque":"#ffe4c4","black":"#000000","blanchedalmond":"#ffebcd","blue":"#0000ff","blueviolet":"#8a2be2","brown":"#a52a2a","burlywood":"#deb887",
    "cadetblue":"#5f9ea0","chartreuse":"#7fff00","chocolate":"#d2691e","coral":"#ff7f50","cornflowerblue":"#6495ed","cornsilk":"#fff8dc","crimson":"#dc143c","cyan":"#00ffff",
    "darkblue":"#00008b","darkcyan":"#008b8b","darkgoldenrod":"#b8860b","darkgray":"#a9a9a9","darkgreen":"#006400","darkkhaki":"#bdb76b","darkmagenta":"#8b008b","darkolivegreen":"#556b2f",
    "darkorange":"#ff8c00","darkorchid":"#9932cc","darkred":"#8b0000","darksalmon":"#e9967a","darkseagreen":"#8fbc8f","darkslateblue":"#483d8b","darkslategray":"#2f4f4f","darkturquoise":"#00ced1",
    "darkviolet":"#9400d3","deeppink":"#ff1493","deepskyblue":"#00bfff","dimgray":"#696969","dodgerblue":"#1e90ff",
    "firebrick":"#b22222","floralwhite":"#fffaf0","forestgreen":"#228b22","fuchsia":"#ff00ff",
    "gainsboro":"#dcdcdc","ghostwhite":"#f8f8ff","gold":"#ffd700","goldenrod":"#daa520","gray":"#808080","green":"#008000","greenyellow":"#adff2f",
    "honeydew":"#f0fff0","hotpink":"#ff69b4",
    "indianred ":"#cd5c5c","indigo":"#4b0082","ivory":"#fffff0","khaki":"#f0e68c",
    "lavender":"#e6e6fa","lavenderblush":"#fff0f5","lawngreen":"#7cfc00","lemonchiffon":"#fffacd","lightblue":"#add8e6","lightcoral":"#f08080","lightcyan":"#e0ffff","lightgoldenrodyellow":"#fafad2",
    "lightgrey":"#d3d3d3","lightgreen":"#90ee90","lightpink":"#ffb6c1","lightsalmon":"#ffa07a","lightseagreen":"#20b2aa","lightskyblue":"#87cefa","lightslategray":"#778899","lightsteelblue":"#b0c4de",
    "lightyellow":"#ffffe0","lime":"#00ff00","limegreen":"#32cd32","linen":"#faf0e6",
    "magenta":"#ff00ff","maroon":"#800000","mediumaquamarine":"#66cdaa","mediumblue":"#0000cd","mediumorchid":"#ba55d3","mediumpurple":"#9370d8","mediumseagreen":"#3cb371","mediumslateblue":"#7b68ee",
    "mediumspringgreen":"#00fa9a","mediumturquoise":"#48d1cc","mediumvioletred":"#c71585","midnightblue":"#191970","mintcream":"#f5fffa","mistyrose":"#ffe4e1","moccasin":"#ffe4b5",
    "navajowhite":"#ffdead","navy":"#000080",
    "oldlace":"#fdf5e6","olive":"#808000","olivedrab":"#6b8e23","orange":"#ffa500","orangered":"#ff4500","orchid":"#da70d6",
    "palegoldenrod":"#eee8aa","palegreen":"#98fb98","paleturquoise":"#afeeee","palevioletred":"#d87093","papayawhip":"#ffefd5","peachpuff":"#ffdab9","peru":"#cd853f","pink":"#ffc0cb","plum":"#dda0dd","powderblue":"#b0e0e6","purple":"#800080",
    "rebeccapurple":"#663399","red":"#ff0000","rosybrown":"#bc8f8f","royalblue":"#4169e1",
    "saddlebrown":"#8b4513","salmon":"#fa8072","sandybrown":"#f4a460","seagreen":"#2e8b57","seashell":"#fff5ee","sienna":"#a0522d","silver":"#c0c0c0","skyblue":"#87ceeb","slateblue":"#6a5acd","slategray":"#708090","snow":"#fffafa","springgreen":"#00ff7f","steelblue":"#4682b4",
    "tan":"#d2b48c","teal":"#008080","thistle":"#d8bfd8","tomato":"#ff6347","turquoise":"#40e0d0",
    "violet":"#ee82ee",
    "wheat":"#f5deb3","white":"#ffffff","whitesmoke":"#f5f5f5",
    "yellow":"#ffff00","yellowgreen":"#9acd32"};

    if (typeof colours[colour.toLowerCase()] != 'undefined')
        return colours[colour.toLowerCase()];

    return false;
}