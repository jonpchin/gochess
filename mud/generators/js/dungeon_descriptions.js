function generator$descriptions$dungeons() {
	var nm1 = ["A grand", "A large", "A massive", "A minor", "A modest", "A narrow", "A short", "A small", "A tall", "A wide"];
	var nm2 = ["overgrown boulder", "granite door", "pair of granite doors", "broken statue", "worn statue", "pair of worn statues", "boulder", "dark cave", "murky cave", "fallen tree", "waterfall", "crypt", "broken temple", "fallen temple", "graveyard", "fallen tower"];
	var nm3 = ["bog", "boulder field", "cliff side", "forest", "grove", "marsh", "morass", "mountain base", "mountain range", "mountain top", "snowland", "swamp", "thicket", "wasteland", "woodlands", "woods"];
	var nm4 = ["large", "small", "massive", "grand", "modest", "scanty", "narrow"];
	var nm5 = ["broken", "clammy", "crumbling", "damp", "dank", "dark", "deteriorated", "dusty", "filthy", "foggy", "grimy", "humid", "putrid", "ragged", "shady", "timeworn", "weary", "worn"];
	var nm6 = ["ash", "bat droppings", "broken pottery", "broken stone", "cobwebs", "crawling insects", "dead insects", "dead vermin", "dirt", "large bones", "moss", "puddles of water", "rat droppings", "remains", "roots", "rubble", "small bones"];
	var nm7 = ["a broken statue part of a fountain", "a broken tomb", "a pillaged treasury", "an altar", "an overgrown underground garden", "broken arrows, rusty swords and skeletal remains", "broken cages and torture devices", "broken mining equipment", "broken vats and flasks", "carved out openings filled with pottery", "cases of explosives and mining equipment", "drawings and symbols on the walls", "empty chests and broken statues", "empty shelves and broken pots", "locked chests, vats, crates and pieces of broken wood", "prison cells", "remnants of a small camp", "remnants of sacks, crates and caskets", "remnants of statues", "remnants of what once must've been a mess hall of sorts", "remnants of what was once a decorated room with a now unknown purpose", "rows of statues", "rows of tombs and several statues", "rows of vertical tombs", "ruins of what seems to be a crude throne room", "the remnants of a pillaged burial chamber", "triggered traps and skeletal remains", "warped and molten metal remnants", "weapons racks and locked crates", "what seems like some form of a sacrificial chamber"];
	var nm8 = ["is a single path", "are two paths, you take the right", "are two paths, but the right is a dead end", "are two paths, you take the left", "are two paths, but the left is a dead end", "are three paths, you take the right", "are three paths, you take the left", "are three paths, you take the middle"];
	var nm9 = ["downwards", "onwards", "passed broken and pillaged tombs", "passed collapsed rooms and pillaged treasuries", "passed countless other pathways", "passed countless rooms", "passed long lost rooms and tombs", "passed lost treasuries, unknown rooms and armories", "passed pillaged rooms", "passed several empty rooms"];
	var nm10 = ["clammy", "crumbled", "damp", "dank", "dark", "dusty", "filthy", "foggy", "ghastly", "ghostly", "grimy", "humid", "putrid", "ragged", "shady", "timeworn", "weary", "worn"];
	var nm11 = ["An altar in the center is covered in runes, some of which are glowing", "An enormous beastly skeleton is chained to the walls", "Countless traps, swinging axes and other devices move all around. They're either still active, or just activated", "It's filled with hanging cages which still hold skeletal remains", "It's filled with strange glowing crystals and countless dead vermin", "It's filled with tombs, but their owners are spread across the floor", "It's filled with tombs, some of which no longer hold their owner", "It's littered with skeletons, but no weaponry in sight", "It's packed with boxes full of runes and magical equipment, as well as skeletons", "Piles and piles of gold lie in the center, several skeletons lie next to it", "Remnants of a makeshift barricade still 'guard' the group of skeletons behind it", "Rows upon rows of shelves are packed with books or remnants of books. In the center sits a single skeleton", "Several cages hold skeletal remains of various animals. Next to the cages are odd machines", "Several stacks of gunpowder barrels are stacked against a wall. A skeleton holding a torch lies before it", "Small holes and carved paths cover the walls, it looks like a community or burrow for small creatures", "Spiderwebs cover everything, large figures seem to be wrapped in the same web", "The floor is riddled with shredded blue prints and a half finished machine sits in a corner", "The room is filled with lifelike statues of terrified people", "There are several braziers scattered around, somehow they're still burning, or burning again", "There's a demolished door with a sign that says \"don't open\"", "There's a huge skeleton in the center, along with dozens of human skeletons", "There's a pile of skeletons in the center, all burned and black", "There's a round stone in the center, around it are a dozen skeletons in a circle", "There's a seemingly endless hole in the center. Around it are what seem like runes", "There's machinery all over the place, probably part of a workshop of sorts"];
	var nm12 = ["advance carefully", "carefully continue", "cautiously proceed", "continue", "move", "press", "proceed", "slowly march", "slowly move", "tread"];
	var nm13 = ["darkness", "depths", "expanse", "mysteries", "secret passages", "secrets", "shadows"];
	var nm14 = ["a few more passages", "a few more rooms and passages", "countless passages", "dozens of similar rooms and passages", "many different passages", "many rooms and passages", "various different rooms and countless passages", "various passages"];
	var nm15 = ["A big", "A grand", "A huge", "A large", "A massive", "A mysterious", "A tall", "A thick", "A vast", "A wide", "An enormous", "An immense", "An ominous"];
	var nm16 = ["wooden", "granite", "metal"];
	var nm17 = ["some are dead ends, others lead to who knows where, or what", "some have collapsed, others are dead ends or too dangerous to try", "most of which are far too ominous looking to try out", "most of which have collapsed or were dead ends to begin with", "some of them have collapsed, others seem to go on forever", "some are dead ends, others seem to have no end at all", "each leading to who knows where, or what", "most of which probably lead to other depths of this dungeon", "most of which look just like the other", "they all look so similar, this whole place is a maze", "each seem to go on forever, leading to who knows what", "some look awfully familiar, others stranger everything else", "each with their own twists, turns and destinations", "most lead to nowhere or back to this same path", "it's one big labyrinth of twists and turns"];
	var nm18 = ["Ash and soot is", "Countless odd symbols are", "Countless runes are", "Dire warning messages are", "Dried blood splatters are", "Intricate carvings are", "Large claw marks are", "Messages in strange languages are", "Ominous symbols are", "Strange writing is", "Various odd symbols are"];
	var nm19 = ["did something just move behind this door?", "you're sure you saw a shadow under the cracks of the door.", "did the door just change its appearance?", "what was that sound?", "you're pretty sure you're being watched.", "was that a growl coming from behind the door?", "did somebody just knock on the door?", "you hear the faint sound of footsteps behind you.", "is that a scratching sound coming from behind the door?", "you think you can hear a whisper coming from behind the door.", "light's coming through the gap below the door.", "you hear a loud bang in the distance from which you came.", "you hear a faint laugh coming from behind the door.", "suddenly the door slowly opens on its own.", "something just grabbed your shoulder."];
	var nm20 = ["bleak", "dark", "dire", "eerie", "foggy", "gloomy", "grim", "misty", "murky", "overcast", "shadowy", "shady", "sinister", "somber"];
	var nm21 = ["aged", "battered", "busted", "decayed", "demolished", "destroyed", "deteriorated", "forgotten", "frayed", "long lost", "pillaged", "tattered", "wasted", "weathered", "worn", "worn down"];
	var nm22 = ["absorbed", "butchered", "claimed", "consumed", "defaced", "desolated", "devoured", "dismantled", "drained", "eaten", "maimed", "mutilated", "ravaged", "ravished", "spoiled", "taken", "wiped out", "wrecked"];

	var rnd1 = Math.floor(Math.random() * nm1.length);
	var rnd2 = Math.floor(Math.random() * nm2.length);
	var rnd3 = Math.floor(Math.random() * nm3.length);
	var rnd4 = Math.floor(Math.random() * nm4.length);
	var rnd5 = Math.floor(Math.random() * nm5.length);
	var rnd6 = Math.floor(Math.random() * nm6.length);
	var rnd6a = Math.floor(Math.random() * nm6.length);
	var rnd6b = Math.floor(Math.random() * nm6.length);
	var rnd7 = Math.floor(Math.random() * nm7.length);
	var rnd8 = Math.floor(Math.random() * nm8.length);
	var rnd9 = Math.floor(Math.random() * nm9.length);
	var rnd10 = Math.floor(Math.random() * nm10.length);
	var rnd11 = Math.floor(Math.random() * nm11.length);
	var rnd12 = Math.floor(Math.random() * nm12.length);
	var rnd13 = Math.floor(Math.random() * nm13.length);
	var rnd14 = Math.floor(Math.random() * nm14.length);
	var rnd15 = Math.floor(Math.random() * nm15.length);
	var rnd16 = Math.floor(Math.random() * nm16.length);
	var rnd17 = Math.floor(Math.random() * nm17.length);
	var rnd18 = Math.floor(Math.random() * nm18.length);
	var rnd19 = Math.floor(Math.random() * nm19.length);
	var rnd20 = Math.floor(Math.random() * nm20.length);
	var rnd21 = Math.floor(Math.random() * nm21.length);
	var rnd22 = Math.floor(Math.random() * nm22.length);

	var name = nm1[rnd1] + " " + nm2[rnd2] + " in a " + nm20[rnd20] + " " + nm3[rnd3] + " marks the entrance to this dungeon. Beyond the " + nm2[rnd2] + " lies a " + nm4[rnd4] + ", " + nm5[rnd5] + " room. It's covered in " + nm6[rnd6] + ", " + nm6[rnd6a] + " and " + nm6[rnd6b] + ".";
	var name2 = "Your torch allows you to see " + nm7[rnd7] + ", " + nm21[rnd21] + " and " + nm22[rnd22] + " by time itself.";

	var name3 = "Further ahead " + nm8[rnd8] + ". Its twisted trail leads " + nm9[rnd9] + " and soon you enter a " + nm10[rnd10] + " area. " + nm11[rnd11] + ". What happened in this place?";

	var name4 = "You " + nm12[rnd12] + " onwards, deeper into the dungeon's " + nm13[rnd13] + ". You pass " + nm14[rnd14] + ", " + nm17[rnd17] + ". You eventually make it to what is likely the final room. " + nm15[rnd15] + " " + nm16[rnd16] + " door blocks your path. " + nm18[rnd18] + " all over it, somehow untouched by time and the elements. You step closer to inspect it and.. wait.. " + nm19[rnd19];

	var result = "";
	result += name;
	result += "\n";
	result += name2;
	result += "\n";
	result += "\n";
	result += name3;
	result += "\n";
	result += "\n";
	result += name4;
	return result;
}