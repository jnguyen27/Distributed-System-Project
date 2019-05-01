/* CMPS 128 Key Value Store Assignment 2 */

// KVS data structure
keyValueStore = {
	store: {},
	// returns: true if value is successfully updated (changed) else false
	set: function (key, value) {
		if(this.hasKey(key) && value == this.store[key]) // Doesn't work for object equality - not sure if problematic
			return false;
		this.store[key] = value;
		return true;
	},
	hasKey: function(key) { // returns boolean
		return key in this.store;
	},
	get: function (key) {
		return this.store[key]; // returns value
	},
	// returns: true if the key-value pair was deleted, else false
	// if the given key does not exist, returns false
	remove: function (key) {
		if(!this.hasKey(key))
			return false;
		return delete this.store[key]
	}
}

module.exports = function (app) {

	/* GET getValue given key method --> returns value for given key */
	app.get('/keyValue-store/:key', (req, res) => {
		console.log('LEADER GET KEYVALUESTORE');
		res.status(200).json({
			'result': 'Success',
			'value': keyValueStore.get(req.params.key)
		});
	});

	/* GET hasKey given key method --> returns true if KVS contains the given key */
	app.get('/keyValue-store/search/:key', (req, res) => {
		res.status(200).json({
			'result': 'Success',
			'value': keyValueStore.hasKey(req.params.key) // Not sure if we should return boolean or string here - specs are not clear
		});
	});

	/* Sets value for given key for KVS */
	app.put('/keyValue-store/:key', (req, res) => {
		if(req.params.key.length < 1 || req.params.key.length > 200)
			res.json({
				'result': 'Error',
				'msg': 'Key not valid'
			});
		else {
			var responseBody = {};
			if(keyValueStore.hasKey(req.params.key)) {
				res.status(200);
				responseBody.msg = "Updated successfully";
				if(keyValueStore.set(req.params.key, req.body.val))
					responseBody.replaced = "True";
				else
					responseBody.replaced = "False";
			} else {
				keyValueStore.set(req.params.key, req.body.val);
				res.status(201);
				responseBody.replaced = "False";
				responseBody.msg = "Added successfully";
			}
			res.json(responseBody);
		}
	});

	/* Deletes given key-value pair from KVS */
	app.delete('/keyValue-store/:key', (req, res) => {
		if(keyValueStore.remove(req.params.key)) {
			res.status(200).json({
				'result': 'Success'
			});
		} else {
			res.status(404).json({
				'result': 'Error',
				'msg': 'Status code 404'
			})
		}
	});

   /* post test method assignment 1 */
   app.post('/test', (req, res) => {
   	res.send('POST message received: ' + req.query.msg);
   });
	 /* post hello mehthod assignment 1 */
   app.post("/hello", function(req, res) {
   	res.status(405).end();
   });
   /* get methods test assignment 1 */
   app.get('/test', (req, res) => {
   	res.send('GET request received');
   });
	 /* get methods hello assignment 1 */
   app.get("/hello", (req, res, next) => { // Why is there next here
   	res.send("Hello world!");
   });

}
