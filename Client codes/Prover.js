var co = require('co');
var fabricservice = require('./Service.js');
var express = require('express');
var app = express();
     
var channelid = "vegetableschannel";
var cc_gen_arbitrary = "GenProofs_ArbitraryRange";
var cc_gen_standard = "GenProofs_StandardRange";

/*             
Get the datas which is used to generate range proof from the mysql database 
*/
var data = function getv(){
	
	// connect database
	var promise = new Promise(function (resolve, reject) {
		var mysql = require('mysql');
		var connection = mysql.createConnection({
			host: '',
			user: '',
			password: '',
			database: ''
		});
		connection.connect();
		connection.query(
			"SELECT Vedata FROM VegetablesInfo",
			function selectCb(err, results) {
				if (results) {
					//console.log(results);
					resolve(results);
				}
				if (err) {
					console.log(err);
					reject(results[0].Vebacteria);
				}
				connection.end();
			}
		);
	});
	promise.then(function (value) {
    
		return value;
   
	}, function (value) {

	});

	return promise;
};

/*
The prover perform the operation which invokes GenProofs_ArbitraryRange chaincode to 
generate arbitrary range proof and upload proof to blockchain.
*/
app.get('/GenProofs_ArbitraryRange',function(req,res){
	
    data().then(data => {

        var keyid = req.query.keyid;
	var rangeid = req.query.rangeid;
	var veclength = req.query.veclength;
        var m = req.query.countvalues;

        var invokeArgs = []
        invokeArgs.push("generate_upload_proof");
        invokeArgs.push(keyid);
        invokeArgs.push(rangeid);
	invokeArgs.push(veclength);

        for(var i=0; i<m; i++){

         	invokeArgs.push(String(data[i].Vedata));

        }
        //console.log(invokeArgs);

        co(function *(){
            
		var result = yield fabricservice.sendTransaction(cc_gen_arbitrary,"invoke",invokeArgs);
		res.send(result);
		
        }).catch((err) => {
		
        	res.send(err);
		
        })

    })

});

/*
The prover perform the operation which invokes  GenProofs_StandardRange chaincode to 
generate standard range proof and upload proof to blockchain.
*/
app.get('/GenProofs_StandardRange',function(req,res){
	
    data().then(data => {

        var keyid = req.query.keyid;
        var veclength = req.query.veclength;
        var m = req.query.countvalues;

        var invokeArgs = []
        invokeArgs.push("generate_upload_proof");
        invokeArgs.push(keyid);
        invokeArgs.push(veclength)

        for(var i=0; i<m; i++){

                invokeArgs.push(String(data[i].Vedata));

        }

        co(function *(){
            
		var result = yield fabricservice.sendTransaction(cc_gen_standard,"invoke",invokeArgs);
		res.send(result);
            
        }).catch((err) => {
         	res.send(err);
        })

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
				res.send("chaincode is true");
			}
		}
		res.send("chaincode is false");
	}).catch((err) => {
		res.send(err);
	})
});

var server = app.listen(3002,function(){
	
	var host = server.address().address;
	var port = server.address().port;

	console.log('Example app listening at http://%s:%s',host,port);
})

process.on('unhandledRejection',function(err){
	console.error(err.stack);
});

process.on('uncaughtException',console.error);
 
