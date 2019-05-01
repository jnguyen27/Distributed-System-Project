var express = require('express'),
  	app = express(),
  	http = require('http'),
	port = process.env.PORT || 8080,
	bodyParser = require('body-parser');


//app.use('/keyValue-store/:key', bodyParser.json({ limit: '1mb' }))
app.use(bodyParser());

// app.use(bodyParser.raw()); // for parsing application/json
// app.use(bodyParser.text()); // for parsing application/x-www-form-urlencoded


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
