var express = require('express'),
  	app = express(),
  	request = require('http'),
	port = process.env.PORT || 8080,
	bodyParser = require('body-parser'),
	proxy = require('http-proxy-middleware');

if(process.env.MAINIP != undefined) {
	app.use('**', proxy({target: 'http://' + process.env.MAINIP, changeOrigin: true, onError: function (err, req, res) {
		res.status(501);
		res.json({
			"result": "Error",
			"msg": "Server unavailable"
		});
	}}));
}


app.use('/keyValue-store/:key', bodyParser.json({ limit: '1mb' }))
//app.use(bodyParser.json());

// Error handling
app.use(function (err, req, res, next) {
	if(err.statusCode == 413) // payload too large
		res.json({
			'result': 'Error',
			'msg': 'Object too large. Size limit is 1MB'
		});
});

require('./api/') (app);

app.listen(port);


console.log('Server started on: ' + port);
