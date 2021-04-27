var co = require('co');
var fabricservice = require('./Service.js');
var express = require('express');

var app = express();
var fs = require('fs');

var channelid = "vegetablefruitchannel";
var chaincode_name = ["cc_prover", "cc_arb_prover"]

app.get('/GenV_AandS_StandardRange', function(req,res){
	var data = req.query.data;
	var invokeArgs = [data];
	var uname = req.query.uname;
	var orgname = req.query.orgname;
	var id = req.query.id;

	co(function *(){

		var result = yield fabricservice.sendTransaction_Prover(chaincode_name[0],"V_A_and_S",invokeArgs, uname, orgname, id);

		for(let i=0; i < result.length; i++){
			res.send( result[i].toString('utf8'));
		}
	}).catch((err) => {
		res.send(err);
	})

});


app.get('/GenV_AandS_ArbitraryRange', function(req,res){
  
	var data = req.query.data;
    var invokeArgs = [data];
	var uname = req.query.uname;
	var orgname = req.query.orgname;
	var id = req.query.id;

	co(function *(){
		
		var result = yield fabricservice.sendTransaction_Prover(chaincode_name[1],"V_A_and_S",invokeArgs, uname, orgname, id);
	
		for(let i=0; i < result.length; i++){
			res.send( result[i].toString('utf8'));
		}
	}).catch((err) => {
		res.send(err);
	})

});


app.get('/GenT1andT2_StandardRange', function(req,res){

	var uname = req.query.uname;
	var orgname = req.query.orgname;
    var id = req.query.id;
	var fsname = "GEN_PROVER_" + id + ".txt";
	var data = fs.readFileSync(fsname);
	var invokeArgs = [data.toString()];
	//console.log(data.toString())
	
	co(function *(){
		
		var result = yield fabricservice.sendTransaction_Prover(chaincode_name[0],"T1_and_T2",invokeArgs, uname, orgname, id);
	
		for(let i=0; i < result.length; i++){
			res.send( result[i].toString('utf8'));
		}
	}).catch((err) => {
		res.send(err);
	})

});


app.get('/GenT1andT2_ArbitraryRange', function(req,res){

	var uname = req.query.uname;
	var orgname = req.query.orgname;
    var id = req.query.id;
	var fsname = "GEN_PROVER_" + id + ".txt";
	var data = fs.readFileSync(fsname);
	var invokeArgs = [data.toString()];
	
	co(function *(){
		
		var result = yield fabricservice.sendTransaction_Prover(chaincode_name[1],"T1_and_T2",invokeArgs, uname, orgname, id);
	
		for(let i=0; i < result.length; i++){
			res.send( result[i].toString('utf8'));
		}
	}).catch((err) => {
		res.send(err);
	})

});


app.get('/GenOtherShare_StandardRange', function(req,res){

    var uname = req.query.uname;
	var orgname = req.query.orgname;
    var id = req.query.id;
	var fsname = "GEN_PROVER_" + id + ".txt";
	var data = fs.readFileSync(fsname);
	var invokeArgs = [data.toString()];
	
	co(function *(){
		
		var result = yield fabricservice.sendTransaction_Prover(chaincode_name[0],"OtherShare",invokeArgs, uname, orgname, id);
	
		for(let i=0; i < result.length; i++){
			res.send( result[i].toString('utf8'));
		}
	}).catch((err) => {
		res.send(err);
	})

});


app.get('/GenOtherShare_ArbitraryRange', function(req,res){
	var uname = req.query.uname;
	var orgname = req.query.orgname;
    var id = req.query.id;
	var fsname = "GEN_PROVER_" + id + ".txt";
	var data = fs.readFileSync(fsname);
	var invokeArgs = [data.toString()];
	
	co(function *(){
		
		var result = yield fabricservice.sendTransaction_Prover(chaincode_name[1],"OtherShare",invokeArgs, uname, orgname, id);
	
		for(let i=0; i < result.length; i++){
			res.send( result[i].toString('utf8'));
		}
	}).catch((err) => {
		res.send(err);
	})

});


app.get('/GetCurState_StandardRange', function(req,res){
  
	var keyid = req.query.keyid;
	var uname = req.query.uname;
	var orgname = req.query.orgname;
	
	co(function *(){
        var chaincodequeryresult = yield fabricservice.queryCc(chaincode_name[0], "Get_Cur_State", [keyid], uname, orgname);
		var result = ''
		for(let i=0; i < chaincodequeryresult.length; i++){
		    result += chaincodequeryresult[i].toString('utf8')
		}

		res.send(result)

    }).catch((err) => {
		res.send(err);
    })
});


app.get('/GetCurState_ArbitraryRange', function(req,res){
	
    var keyid = req.query.keyid;
	var uname = req.query.uname;
	var orgname = req.query.orgname;
	
	co(function *(){
        var chaincodequeryresult = yield fabricservice.queryCc(chaincode_name[1], "Get_Cur_State", [keyid], uname, orgname);
		var result = ''

		for(let i=0; i < chaincodequeryresult.length; i++){
		    result += chaincodequeryresult[i].toString('utf8')
		}

		res.send(result)

    }).catch((err) => {
        res.send(err);
    })
});


app.get('GetStateHistory_StandardRange', function(req,res){

	var keyid = req.query.keyid;
	var uname = req.query.uname;
	var orgname = req.query.orgname;
	
	co(function *(){
        var chaincodequeryresult = yield fabricservice.queryCc(chaincode_name[0], "Get_State_History", [keyid], uname, orgname);
		var result = ''
		for(let i=0; i < chaincodequeryresult.length; i++){
		    result += chaincodequeryresult[i].toString('utf8')
		}

		res.send(result)

    }).catch((err) => {
        res.send(err);
    })
});

app.get('/GetStateHistory_ArbitraryRange', function(req,res){

	var keyid = req.query.keyid;
	var uname = req.query.uname;
	var orgname = req.query.orgname;
	
	co(function *(){
        var chaincodequeryresult = yield fabricservice.queryCc(chaincode_name[1], "Get_State_History", [keyid], uname, orgname);
		var result = ''
		for(let i=0; i < chaincodequeryresult.length; i++){
		    result += chaincodequeryresult[i].toString('utf8')
		}
		res.send(result)

    }).catch((err) => {
        res.send(err);
    })
});

//The version information of chaincode
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


var server = app.listen(3001,function(){
	var host = server.address().address;
	var port = server.address().port;
	console.log('Example app listening at http://%s:%s',host,port);
})

process.on('unhandledRejection',function(err){
    console.error(err.stack);
});

process.on('uncaughtException',console.error);
 

