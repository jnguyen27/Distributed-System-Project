var express = require('express'),
  	app = express(),
  	port = process.env.PORT || 8080;
    temp = 'nothings'

require('./api/') (app, {});

app.listen(port);


console.log('todo list RESTful API server started on: ' + port);
