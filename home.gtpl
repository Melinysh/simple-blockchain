<html>
<head>
<title>Your Blockchain</title>
</head>
<body>
	<form action="/create" method="post">
    	Add to Blockchain:<input type="text" name="BlockData">
    	<input type="submit" value="Add block">
	</form>
	<form action="/peer" method="post">
    	Add peer:<input type="text" name="PeerURL">
    	<input type="submit" value="Add peer">
	</form>
	<h1>Your Blockchain</h1>
	<table border="1px solid black">
		{{range .Blocks}}
		<tr>
			<td>{{.}}</td>
		</tr>
		{{end}}
	</table>
	<h1>Your Peers</h1>
	<table border="1px solid black">
		{{range .Peers}}
		<tr>
			<td>{{.}}</td>
		</tr>
		{{end}}
	</table>

</body>
