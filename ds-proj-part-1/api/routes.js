module.exports = function (app, db) {
	app.post('/test', (req, res) => {
		res.send('POST message received: ' + req.query.msg);
	});

	app.get('/test', (req, res) => {
		res.send('GET request received');
	});

	app.get("/hello", (req, res, next) => {
		res.send("Hello world!");
	});

	app.post("/hello", function(req, res) {
		res.status(405).end();
	});
}
