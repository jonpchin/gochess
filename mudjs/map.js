// Displays a portion of the map that surrounds the player
function displayMap(mapString){

    var rowLength = sqrt(mapstring.length);

    for(var i=0; i<rowLength; ++i){
        
        for(var j=0; j<rowLength; ++j){
            var char = mapString[(i*rowLength)+j];
            displayToTextBoxNoNewLine(char, getMapCharColor(char));
        }
        displayToTextBox("");
    }
}

// Gets the color of the map char that its mapped too
function getMapCharColor(char){
    
    var color = "black"; // Default char color

    switch(char){
        case " ": 
            //unused
            color = "white";
            break;
        case ".": 
            //floor
            color = "yellow";
            break;
        case ",":
            //corridor 
            color = "grey";
            break;
        case "#":
            //wall 
            color = "orange";
            break;
        case "+":
            //closed door 
            color = "red";
            break;
        case "-":
            //open door 
            color = "green";
            break;
        case "<":
            //upstairs 
            color = "firebrick";
            break;
        case ">":
            //down stairs 
            color = "honeydew";
            break;
        case "$":
            //forest
            color = "darkgreen";
            break;
        case "%":
            //water 
            color = "teal";
            break;
        case "@":
            //cloud 
            color = "antiquewhite";
            break;
        case "^":
            //mountain 
            color = "brown";
            break;
        case "!":
            //whirlpool 
            color = "blue";
            break;
        default:
            console.log("No such map char exists");
    }

    return color;
}



Mud.Room = {
    Tiles: [], // Mud.Tile
    Walls: []  // Mud.Tile
}

Mud.Tile = {
	Coordinate: Mud.Coordinate,                         // Contains X (row) and Y (col) coordinates of tile
	Name:        "Default Room",                        // Name of the tile an adventurer will see when they enter the room
	Description: "This is a default room description",  // Description of tile adventurer will see when they enter the room
	Floor:       5,                                     // The floor the tile is located
	Area:        Mud.Area,                              // The area the tile is located
	Room:        Mud.Room,                              // The room the tile is located
	TileType:    ".",                                   // The type of tile such as floor, wall, openDoor, etc, default is floor "."
	Items:       [],                                    // List of items or objects in the tile, Mud.Object
	Players:     [],                                    // Adventurers in the tile
	Monsters:    []                                     // NPC in the tile
}

Mud.Coordinate = {
    Row: 5,
    Col: 5,
    Level: 5
}

Mud.Area = {
	Name: "Cain's Hideout"
}