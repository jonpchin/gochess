function setBoardColor(){
    var boardColor = getCookie("boardColor");
	switch(boardColor){
		case "blueColor":
			$('.black-3c85d').css({"background-color":"#4682B4"});
			$('.black-3c85d').css({"color":"#B0E0E6"});
			$('.white-1e1d7').css({"background-color":"#B0E0E6"});
			$('.white-1e1d7').css({"color":"#4682B4"});
			break;

		case "greenColor":
			$('.black-3c85d').css({"background-color":"#008000"});
			$('.black-3c85d').css({"color":"#90EE90"});
			$('.white-1e1d7').css({"background-color":"#90EE90"});
			$('.white-1e1d7').css({"color":"#008000"});
			break;

		case "greyColor":
			$('.black-3c85d').css({"background-color":"#696969"});
			$('.black-3c85d').css({"color":"#D3D3D3"});
			$('.white-1e1d7').css({"background-color":"#D3D3D3"});
			$('.white-1e1d7').css({"color":"#696969"});
			break;

		case "pinkColor":
			$('.black-3c85d').css({"background-color":"#FF69B4"});
			$('.black-3c85d').css({"color":"#FFC0CB"});
			$('.white-1e1d7').css({"background-color":"#FFC0CB"});
			$('.white-1e1d7').css({"color":"#FF69B4"});
			break;
			
		default: //default color
			$('.black-3c85d').css({"background-color":"#b58863"});
			$('.black-3c85d').css({"color":"#f0d9b5"});
	
			$('.white-1e1d7').css({"background-color":"#f0d9b5"});
			$('.white-1e1d7').css({"color":"#b58863"});
	}
}
setBoardColor();