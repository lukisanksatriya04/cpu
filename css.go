package main

func css() string {

	file := `
	@font-face {
		font-family: "sylendra";
		src: url("D-DIN.otf");
	}
	body {
		overflow: hidden;
		background-color: #0c090a;
	}
	
	.homepanel {
		width: 300px;
		height: 400px;
		background-color: #f8f8ff;
		border-radius: 15px;
		text-align: center;
	}
	
	.text {
		color: #0c090a;
		font-family: "sylendra";
		font-size: 50px;
		transform: scale(1, .85);
	}
	
	.value {
		margin: 0 auto;
		margin-top: 40px;
		color: #0c090a;
		font-family: "sylendra";
		font-size: 150px;
		transform: scale(.85, 1);
	}
	
	.satuan {
		margin: 0 auto;
		margin-top: -50px;
		color: #0c090a;
		font-family: "sylendra";
		font-size: 75px;
		transform: scale(1, .95);
	}
	
	.button {
		margin: 0 auto;
		margin-top: 10px;
		border-radius: 10px;
		color: #0c090a;
		background-color: #f8f8ff;
		border: 2px solid #0c090a;
		font-family: "sylendra";
		font-size: 25px;
		box-shadow: 0 2px #000000;
	}
	
	.button:active {
		background-color: #f8f8ff;
		box-shadow: 0 3px #0c090a;
		transform: translateY(2px);
		color: #0c090a;
	}`

	return file
}
