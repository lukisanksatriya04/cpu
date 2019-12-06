package main

func html() string {
	file := `
	<!DOCTYPE html>
	<html>

	<head>
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<link rel="stylesheet" href="bootstrap.min.css">
		<link rel="stylesheet" href="home.css">
		<script type="text/javascript" src="storage.json"></script>
	</head>

	<body>

		<div class="position-fixed homepanel" style=" top:100px; left:31px;">
			<div class="position-fixed text"  style=" top:105px; left:40px;">
				Consumption</div>
			<div class="value" id="invval">0.00</div>
			<div class="satuan">KW</div>
			<button class="button">Click for detail</button>
		</div>

		<div class="position-fixed homepanel" style=" top:100px; left:362px;">
			<div class="position-fixed text"  style=" top:105px; left:436px;">
				Storage</div>
			<div class="value" id="strval">0.00</div>
			<div class="satuan">KWH</div>
			<button class="button">Click for detail</button>
		</div>

		<div class="position-fixed homepanel" style=" top:100px; left:693px;">
			<div class="position-fixed text"  style=" top:105px; left:735px;">
				Production</div>
			<div class="value" id="cnvval">0.00</div>
			<div class="satuan">KW</div>
			<button class="button">Click for detail</button>
		</div>

		<script>

		var i = 1
		setInterval(function() {
			reloadvalue();
		}, 1000);

		function reloadvalue(){
			var obj = JSON.parse(JSON.stringify(str));
			document.getElementById("invval").innerHTML = obj.data.Temperaturepack;
			document.getElementById("strval").innerHTML = obj.data.Current;
			document.getElementById("cnvval").innerHTML = obj.data.Voltagepack;
		}
		</script>

	</body>

	</html>`

	return file
}
