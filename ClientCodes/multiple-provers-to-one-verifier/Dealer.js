var co = require('co');
var fabricservice = require('./Service.js');
var express = require('express');

var app = express();
var fs = require('fs');

var channelid = "vegetablefruitchannel";
var chaincode_name = ["cc_dealer", "cc_arb_dealer"]


app.get('/Fait_y_z_StandardRange',function(req,res){


	var uname = req.query.uname;
	var orgname = req.query.orgname;

        co(function *(){
            
            var result = yield fabricservice.sendTransaction_Dea_Ver(chaincode_name[0],"Fait_y_z",[], uname, orgname);
	    
            for(let i=0; i < result.length; i++){
                res.send( result[i].toString('utf8'));
            }
        }).catch((err) => {
            res.send(err);
        })

});


app.get('/Fait_y_z_ArbitraryRange',function(req,res){


	var uname = req.query.uname;
	var orgname = req.query.orgname;

        co(function *(){
            
            var result = yield fabricservice.sendTransaction_Dea_Ver(chaincode_name[1],"Fait_y_z",[], uname, orgname);
	    
            for(let i=0; i < result.length; i++){
                res.send( result[i].toString('utf8'));
            }
        }).catch((err) => {
            res.send(err);
        })

});


app.get('/Fait_x_StandardRange',function(req,res){


	var uname = req.query.uname;
	var orgname = req.query.orgname;

        co(function *(){
            
            var result = yield fabricservice.sendTransaction_Dea_Ver(chaincode_name[0],"Fait_x",[], uname, orgname);
	    
            for(let i=0; i < result.length; i++){
                res.send( result[i].toString('utf8'));
            }
        }).catch((err) => {
            res.send(err);
        })

});


app.get('/Fait_x_ArbitraryRange',function(req,res){


	var uname = req.query.uname;
	var orgname = req.query.orgname;

        co(function *(){
            
            var result = yield fabricservice.sendTransaction_Dea_Ver(chaincode_name[1],"Fait_x",[], uname, orgname);
	    
            for(let i=0; i < result.length; i++){
                res.send( result[i].toString('utf8'));
            }
        }).catch((err) => {
            res.send(err);
        })

});


app.get('/GenPrf_StandardRange',function(req,res){


	var uname = req.query.uname;
	var orgname = req.query.orgname;

        co(function *(){
            
            var result = yield fabricservice.sendTransaction_Dea_Ver(chaincode_name[0],"Gen_Prf",[], uname, orgname);
	    
            for(let i=0; i < result.length; i++){
                res.send( result[i].toString('utf8'));
            }
        }).catch((err) => {
            res.send(err);
        })

});


app.get('/GenPrf_ArbitraryRange',function(req,res){


	var uname = req.query.uname;
	var orgname = req.query.orgname;

        co(function *(){
            
            var result = yield fabricservice.sendTransaction_Dea_Ver(chaincode_name[1],"Gen_Prf",[], uname, orgname);
	    
            for(let i=0; i < result.length; i++){
                res.send( result[i].toString('utf8'));
            }
        }).catch((err) => {
            res.send(err);
        })

});

app.get('/GetMPCRangePrf_StandardRange',function(req,res){

  
	var keyid = req.query.keyid;
        var chaincodequeryresult = yield fabricservice.queryCc(chaincode_name[0],"Get_MPC_Range_Prf",[keyid]);

	var result = ''

        for(let i=0; i < chaincodequeryresult.length; i++){
            result += chaincodequeryresult[i].toString('utf8')
        }

	res.send(result)

    }).catch((err) => {
        res.send(err);
    })
});


app.get('/GetMPCRangePrf_ArbitraryRange',function(req,res){

  
	var keyid = req.query.keyid;
        var chaincodequeryresult = yield fabricservice.queryCc(chaincode_name[1],"Get_MPC_Range_Prf",[keyid]);

	var result = ''

        for(let i=0; i < chaincodequeryresult.length; i++){
            result += chaincodequeryresult[i].toString('utf8')
        }

	res.send(result)

    }).catch((err) => {
        res.send(err);
    })
});


app.get('/GetRangeHistory_StandardRange',function(req,res){

  
	var keyid = req.query.keyid;
        var chaincodequeryresult = yield fabricservice.queryCc(chaincode_name[0],"Get_Range_History",[keyid]);

	var result = ''

        for(let i=0; i < chaincodequeryresult.length; i++){
            result += chaincodequeryresult[i].toString('utf8')
        }

	res.send(result)

    }).catch((err) => {
        res.send(err);
    })
});

app.get('/GetRangeHistory_ArbitraryRange',function(req,res){

  
	var keyid = req.query.keyid;
        var chaincodequeryresult = yield fabricservice.queryCc(chaincode_name[1],"Get_Range_History",[keyid]);

	var result = ''

        for(let i=0; i < chaincodequeryresult.length; i++){
            result += chaincodequeryresult[i].toString('utf8')
        }

	res.send(result)

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


var server = app.listen(3002,function(){
    var host = server.address().address;
    var port = server.address().port;

    console.log('Example app listening at http://%s:%s',host,port);
})

process.on('unhandledRejection',function(err){
    console.error(err.stack);
});

process.on('uncaughtException',console.error);
 

