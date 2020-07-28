var co = require('co');
var fabricservice = require('./Service.js');
var express = require('express');
var app = express();
     
var channelid = "vegetableschannel";

/*
Perform the operation which invokes RangeManagement chaincode to upload range to blockchain
*/
app.get('/upload_range',function(req,res){

	co(function *(){

		var rangeid = req.query.rangeid;
		var value1 = req.query.value1;
		var value2 = req.query.value2;

		var result = yield fabricservice.sendTransaction("RangeManagement","invoke",["upload_range",rangeid,value1,value2]);
		res.send(result);

	}).catch((err) => {
		res.send(err);
	})
});

/*
Invoke RangeManagement chaincode to get range from the blockchain
*/
app.get('/get_range',function(req,res){

	co(function *(){

		var rangeid = req.query.rangeid;

		var chaincodequeryresult = yield fabricservice.queryCc("RangeManagement","invoke",["get_range",rangeid]);

		for(let i=0; i < chaincodequeryresult.length; i++){
		    res.send( chaincodequeryresult[i].toString('utf8'));
		}

	}).catch((err) => {
		res.send(err);
	})
});

var server = app.listen(3001,function(){
	
	var host = server.address().address;
	var port = server.address().port;

	console.log('Example app listening at http://%s:%s',host,port);
})

process.on('unhandledRejection',function(err){
	console.error(err.stack);
});

process.on('uncaughtException',console.error);
 
