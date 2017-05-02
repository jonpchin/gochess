function generateBelts(){
	for(var i=0; i<50; ++i){
		document.write(generator$armour$belts(1));
		document.write("<br>");
	}
	for(var i=0; i<50; ++i){
		document.write(generator$armour$belts(0));
		document.write("<br>");
	}
}

function generateBoots(){
	for(var i=0; i<50; ++i){
		document.write(generator$armour$boots(2));
		document.write("<br>");
	}
	for(var i=0; i<50; ++i){
		document.write(generator$armour$boots(0));
		document.write("<br>");
	}
}

function generateChests(){
	for(var i=0; i<50; ++i){
		document.write(generator$armour$chests(1));
		document.write("<br>");
	}
	for(var i=0; i<50; ++i){
		document.write(generator$armour$chests(0));
		document.write("<br>");
	}
}

function generateHelmets(){
	for(var i=0; i<50; ++i){
		document.write(generator$armour$helmets(2));
		document.write("<br>");
	}
	for(var i=0; i<50; ++i){
		document.write(generator$armour$helmets(0));
		document.write("<br>");
	}
}

function generateLegs(){
	for(var i=0; i<50; ++i){
		document.write(generator$armour$legs(2));
		document.write("<br>");
	}
	for(var i=0; i<50; ++i){
		document.write(generator$armour$legs(0));
		document.write("<br>");
	}
}

function generateShields(){
	for(var i=0; i<100; ++i){
		document.write(generator$armour$shields());
		document.write("<br>");
	}
}

function generateDaggers(){
	for(var i=0; i<100; ++i){
		document.write(generator$weapons$daggers());
		document.write("<br>");
	}
}
//generateBelts();
//generateBoots();
//generateChests();
//generateHelmets();
//generateLegs();
//generateShields();
generateDaggers();