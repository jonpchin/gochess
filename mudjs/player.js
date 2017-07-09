
Mud = {}
Mud.Classes = getClasses();
Mud.Races = getRaces();
Mud.Player = {
	Username: document.getElementById('user').value, // Go Play Chess account
	Name: "",     // Mud account name
	Class: "",
	Race: "",
	Gender: "",
	Inventory: [],
	Equipment: [],
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
	Status: "healthy", 
	Bleed: 0,  
	Level: 1, 
	Experience: 0,
	Location: [],
	Area: "Cain's Hideout",
}

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
	console.log("start updatePlayer");
	console.log(Mud.Player);
	Mud.Player.Name      = player.Name;
	Mud.Player.Class     = player.Class;
	Mud.Player.Race      = player.Race;
	Mud.Player.Gender    = player.Gender;
	Mud.Player.Inventory = player.Inventory.split(',');
	Mud.Player.Equipment = player.Equipment.split(',');

	var stats = JSON.parse(Mud.Player.Stats);

	Mud.Player.Stats.Health    = stats.Health;
	Mud.Player.Stats.Mana      = stats.Mana;
	Mud.Player.Stats.Energy    = stats.Energy;
	Mud.Player.Stats.Strength  = stats.Strength;
	Mud.Player.Stats.Speed     = stats.Speed;
	Mud.Player.Stats.Dexterity = stats.Dexterity;
	Mud.Player.Stats.Intelligence = stats.Intelligence;
	Mud.Player.Stats.Wisdom = stats.Wisdom;
	console.log(Mud.Player);
}