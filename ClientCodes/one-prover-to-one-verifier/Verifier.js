var co = require('co');
var fabricservice = require('./Service.js');
var express = require('express');
var app = express();

var channelid = "vegetableschannel";
var cc_verify_arbitrary = "VerifyProofs_ArbitraryRange";
var cc_verify_standard = "VerifyProofs_StandardRange";

/*
The verifier perform the operation which invokes VerifyProofs_ArbitraryRange chaincode to verify arbitrary range proof.
*/
app.get('/VerifyProofs_ArbitraryRange',function(req,res){

	co(function *(){

		var keyid = req.query.keyid;
		var minvalue = req.query.minvalue;
		var maxvalue = req.query.maxvalue;

		var chaincodequeryresult = yield fabricservice.queryCc(cc_verify_arbitrary,"invoke",[keyid, minvalue, maxvalue]);

		for(let i=0; i < chaincodequeryresult.length; i++){
			res.send( chaincodequeryresult[i].toString('utf8'));
		}

	}).catch((err) => {
		res.send(err);
	})
});

/*
The verifier perform the operation which invokes VerifyProofs_StandardRange chaincode to verify standard range proof.
*/
app.get('/VerifyProofs_StandardRange',function(req,res){

	co(function *(){

		var keyid = req.query.keyid;

		var chaincodequeryresult = yield fabricservice.queryCc(cc_verify_standard,"invoke",[keyid]);

		for(let i=0; i < chaincodequeryresult.length; i++){
			res.send( chaincodequeryresult[i].toString('utf8'));
		}

	}).catch((err) => {
		res.send(err);
	})
});

/*
Accoding the version information of instantiated chaincode to verify chaincode is modified or not.
*/
app.get('/chaincodes',function(req,res){

	co(function *(){

		var ccname = req.query.ccname;
		var ccversion = req.query.ccversion;
		console.info(ccname);
		console.info(ccversion);

		var info = yield fabricservice.getInstantiatedChaincodes();
		for (let i = 0; i < info.chaincodes.length; i++) {
			//console.info('name: ' + info.chaincodes[i].name + ', version: ' +info.chaincodes[i].version + ', path: ' + info.chaincodes[i].path);
			if ((info.chaincodes[i].name == ccname) && (info.chaincodes[i].version == ccversion)){
				res.send("chaincode is true")
			}
		}
		res.send("chaincode is false")

	}).catch((err) => {
		res.send(err);
	})
});

var server = app.listen(3003,function(){
	
	var host = server.address().address;
	var port = server.address().port;

	console.log('Example app listening at http://%s:%s',host,port);
})

process.on('unhandledRejection',function(err){
    console.error(err.stack);
});

process.on('uncaughtException',console.error);
 
