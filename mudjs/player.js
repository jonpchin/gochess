Mud.Classes = getClasses();
Mud.Races = getRaces();
Mud.Player = {
	Type: "none",
	Username: document.getElementById('user').value, // Go Play Chess account
	Name: "",     // Mud account name
	SessionID: document.getElementById('sessionID').value, // Session Token
	Class: "",
	Race: "",
	Gender: "",
	Inventory: [],
	Equipment: Mud.Equipment,
	Stats:  {
		Health: 50,
		Mana: 50,
		Energy: 50,
		Strength: 10,
		Speed: 10,
		Dexterity: 10,
		Intelligence: 10,
		Wisdom: 10
	}, 
	Status: ["healthy"], 
	Map: "", // Contains 2D array of strings, odd number is always square so MapVision 5 is 11x11
	Bleed: 0,
	Level: 1, 
	Experience: 0,
	Location: Mud.Coordinate,
	Area: Mud.Area,
	Command: ""
}

// Map of mud name to Mud.Player
Mud.Server = [];
Mud.Server[Mud.Player.Name] = Mud.Player;

Mud.Server.Stats = {
	HiddenPlayers: 0
};

// Since this function is async if one wants to call this on load a sleep will need to be made
// or a non async getClasses needs to be made
function getClasses(){

    var classes = [];

    $.getJSON('../data/mud/classes.json', function(data) {        
        
        for (var key in data) {
            // skip loop if the property is from prototype
            if (!data.hasOwnProperty(key)){
                continue;
            }
            classes.push(key);
        }
    });
    return classes;
}

function getRaces(){
	var races = [];

    $.getJSON('../data/mud/races.json', function(data) {        
        
        for (var key in data) {
            // skip loop if the property is from prototype
            if (!data.hasOwnProperty(key)){
                continue;
            }
            races.push(key);
        }
    });
    return races;
}

// Update the player data in memory
function updatePlayer(player){

	Mud.Player.Name      = player.Name;
	Mud.Player.Class     = player.Class;
	Mud.Player.Race      = player.Race;
	Mud.Player.Gender    = player.Gender;
	Mud.Player.Inventory = player.Inventory;	
	Mud.Player.Equipment = player.Equipment;
	Mud.Player.Stats   = player.Stats;
	Mud.Player.Map     = player.Map;
	Mud.Player.Bleed = player.Bleed;
	Mud.Player.Location.Row = player.Location.Row;
	Mud.Player.Location.Col = player.Location.Col;
	Mud.Player.Location.Level = player.Location.Level;

}