var co = require('co');
var fabricservice = require('./Service.js');
var express = require('express');

var app = express();

     
var channelid = "vegetablefruitchannel";
var chaincode_name = ["cc_verifier", "cc_arb_verifier"]


app.get('/GetSetupInfo_StandardRange',function(req,res){

    co(function *(){
		
        var currentIndex = yield fabricservice.queryCc(chaincode_name[0],"Get_KeyID_Sess_Num",[]);
	console.log('currentIndex:'+currentIndex)
        
	var setupInfo = yield fabricservice.queryCc(chaincode_name[0],"Get_Setup_Info",[currentIndex]);
	console.log(setupInfo.toString('utf8'))
	res.send( JSON.stringify({result:setupInfo}) )

    }).catch((err) => {
        res.send(err);
    })
});

app.get('/GetSetupInfo_ArbitraryRange',function(req,res){

    co(function *(){
		
        var currentIndex = yield fabricservice.queryCc(chaincode_name[1],"Get_KeyID_Sess_Num",[]);
	console.log('currentIndex:'+currentIndex)
        
	var setupInfo = yield fabricservice.queryCc(chaincode_name[1],"Get_Setup_Info",[currentIndex]);
	console.log(setupInfo.toString('utf8'))
	res.send( JSON.stringify({result:setupInfo}) )

    }).catch((err) => {
        res.send(err);
    })
});

//Veclength xxx arbrange xxx prover xxx  xxx  dealer xxx

app.get('/InitSetup_StandardRange',function(req,res){

    co(function *(){
	var len = 0;	
	for (var para in req.query){
		len++;
	}
	prover_num = len - 2;
	var veclength = req.query.VecLength;
	var prover = [req.query.Prover1, req.query.Prover2];
	if (prover_num == 4) {
		prover.push(req.query.Prover3);
		prover.push(req.query.Prover4);	
	}else if (prover_num == 8) {
		prover.push(req.query.Prover3);
		prover.push(req.query.Prover4);
		prover.push(req.query.Prover5);
		prover.push(req.query.Prover6);
		prover.push(req.query.Prover7);
		prover.push(req.query.Prover8);
	}
	var dealer = req.query.Dealer;
	// args: Veclength XXX prover xxx  xxx  dealer xxx
	var invokeArgs = ["VecLength", "Prover", "Dealer"]
	invokeArgs.splice(1, 0, veclength);
	for (let i=0; i < prover_num; i++){
		invokeArgs.splice(i+3, 0, prover[i]);
	}
	invokeArgs.push(dealer);
	console.log(invokeArgs);
	
        var result = yield fabricservice.sendTransaction_Dea_Ver(chaincode_name[0],"Init_Setup",invokeArgs, "Admin", "supervision");
	    
            for(let i=0; i < result.length; i++){
                res.send( result[i].toString('utf8'));
            }

    }).catch((err) => {
        res.send(err);
    })
});


app.get('/InitSetup_ArbitraryRange',function(req,res){

    co(function *(){
	var len = 0;	
	for (var para in req.query){
		len++;
	}
	prover_num = len - 3;
	var veclength = req.query.VecLength;
	var arbrange = req.query.ArbRange;
	var prover = [req.query.Prover1, req.query.Prover2];
	if (prover_num == 4) {
		prover.push(req.query.Prover3);
		prover.push(req.query.Prover4);	
	}else if (prover_num == 8) {
		prover.push(req.query.Prover3);
		prover.push(req.query.Prover4);
		prover.push(req.query.Prover5);
		prover.push(req.query.Prover6);
		prover.push(req.query.Prover7);
		prover.push(req.query.Prover8);
	}
	var dealer = req.query.Dealer;
	var invokeArgs = ["VecLength", "ArbRange", "Prover", "Dealer"]
	invokeArgs.splice(1, 0, veclength);
	invokeArgs.splice(3, 0, arbrange);
	for (let i=0; i < prover_num; i++){
		invokeArgs.splice(i+5, 0, prover[i]);
	}
	invokeArgs.push(dealer);
	console.log(invokeArgs);

        var result = yield fabricservice.sendTransaction_Dea_Ver(chaincode_name[1],"Init_Setup",invokeArgs, "Admin", "supervision");
	    
            for(let i=0; i < result.length; i++){
                res.send( result[i].toString('utf8'));
            }

    }).catch((err) => {
        res.send(err);
    })
});



app.get('/VerPrf_StandardRange',function(req,res){

    co(function *(){

        var chaincodequeryresult = yield fabricservice.queryCc(chaincode_name[0],"Ver_Prf",[], "Admin", "supervision");

        for(let i=0; i < chaincodequeryresult.length; i++){
            res.send( chaincodequeryresult[i].toString('utf8'));
        }

    }).catch((err) => {
        res.send(err);
    })
});


app.get('/VerPrf_ArbitraryRange',function(req,res){

    co(function *(){

        var chaincodequeryresult = yield fabricservice.queryCc(chaincode_name[1],"Ver_Prf",[], "Admin", "supervision");

        for(let i=0; i < chaincodequeryresult.length; i++){
            res.send( chaincodequeryresult[i].toString('utf8'));
        }

    }).catch((err) => {
        res.send(err);
    })
});

//chaincode的版本信息
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

var server = app.listen(3000,function(){
    var host = server.address().address;
    var port = server.address().port;

    console.log('Example app listening at http://%s:%s',host,port);
})

process.on('unhandledRejection',function(err){
    console.error(err.stack);
});

process.on('uncaughtException',console.error);
 

