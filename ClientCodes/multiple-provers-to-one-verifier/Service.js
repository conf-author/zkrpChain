
var co = require('co');
var path = require('path');
var fs = require('fs');
var util = require('util');
var hfc = require('fabric-client');
var Peer = require('fabric-client/lib/Peer.js');
var User = require('fabric-client/lib/User.js');
var crypto = require('crypto');
var fs = require('fs');
var log4js = require('log4js');
var logger = log4js.getLogger('Helper');
logger.level= 'DEBUG';

var tempdir ="/home/hzaucoi/js-mysql/fabric-client-kvs";
var client = new hfc();

var channel = client.newChannel('vegetablefruitchannel');

var order = client.newOrderer('grpc://172.16.174.45:7050');
channel.addOrderer(order);            


var peer0 = client.newPeer('grpc://172.16.174.45:7051'); //peer0.org1
channel.addPeer(peer0);

var peer1 = client.newPeer('grpc://172.16.174.45:8051'); //peer0.org2
channel.addPeer(peer1);

var peer2 = client.newPeer('grpc://172.16.174.45:9051'); //peer0.org3
channel.addPeer(peer2);

var peer3 = client.newPeer('grpc://172.16.174.39:9051'); //peer0.supervision
channel.addPeer(peer3);

var peer4 = client.newPeer('grpc://172.16.174.39:8051'); //peer0.org4
channel.addPeer(peer4);


var queryCc = function(chaincodeid,func,chaincode_args, uname, orgname){

   return getOrgUserLocal(uname, orgname).then( (user)=> {

        tx_id = client.newTransactionID();
        var request = {
                chaincodeId: chaincodeid,
                txId: tx_id,
                fcn: func,
                args:chaincode_args

        };
        return channel.queryByChaincode(request,peer3);

   },(err)=>{
        console.log('error',e);
   }).then( (sendtransresult) => {
	console.log(sendtransresult)
        return sendtransresult;
   },(err)=>{
        console.log('error',e);
   });
}



var sendTransaction_Dea_Ver = function(chaincodeid,func,chaincode_args, uname, orgname){

   var tx_id = null;
   
   return getOrgUserLocal(uname, orgname).then( (user)=> {
        tx_id = client.newTransactionID();
        var request = {
                chaincodeId: chaincodeid,
                fcn: func,
                args:chaincode_args,
                chainId: "vegetablefruitchannel",
                txId: tx_id

        };
        return channel.sendTransactionProposal(request);

   },(err)=>{

        console.log('error',e);

   }).then( (chaincodeinvokeresult) => {

        var proposalResponses = chaincodeinvokeresult[0];
        var proposal = chaincodeinvokeresult[1];
        var header = chaincodeinvokeresult[2];
        var all_good = true;

        for (var i in proposalResponses){
                let one_good = false;
                if(proposalResponses && proposalResponses[0].response && proposalResponses[0].response.status === 200){
                        one_good = true;
                        console.info('transcation proposal was good');
                }else{
                        console.info('transcation proposal was bad');
                }

                all_good = all_good & one_good;

        }

        if(all_good){
                console.log(util.format(
                        'Successfully :Status - %s,message - "%s",metadata - "%s",endorsement signature :%s',
                        proposalResponses[0].response.status,
                        proposalResponses[0].response.message,
                        proposalResponses[0].response.payload,
                        proposalResponses[0].endorsement.signature));
		
                var request = {
                        proposalResponses:proposalResponses,
                        proposal:proposal,
                        header:header
                };

                var transactionID = tx_id.getTransactionID();
                return channel.sendTransaction(request);
        }

   },(err)=>{
        console.log('error',e);
   }).then( (sendtransresult) => {
	console.log(sendtransresult)
        return sendtransresult;
   },(err)=>{
        console.log('error',e);
   });
}


var sendTransaction_Prover = function(chaincodeid,func,chaincode_args, uname, orgname,id){

   var tx_id = null;
   
   return getOrgUserLocal(uname, orgname).then( (user)=> {
	
        tx_id = client.newTransactionID();
        var request = {
                chaincodeId: chaincodeid,
                fcn: func,
                args:chaincode_args,
                chainId: "vegetablefruitchannel",
                txId: tx_id

        };
        return channel.sendTransactionProposal(request);

   },(err)=>{

        console.log('error',e);

   }).then( (chaincodeinvokeresult) => {

        var proposalResponses = chaincodeinvokeresult[0];
        var proposal = chaincodeinvokeresult[1];
        var header = chaincodeinvokeresult[2];
        var all_good = true;

        for (var i in proposalResponses){
                let one_good = false;
                if(proposalResponses && proposalResponses[0].response && proposalResponses[0].response.status === 200){
                        one_good = true;
                        console.info('transcation proposal was good');
                }else{
                        console.info('transcation proposal was bad');
                }

                all_good = all_good & one_good;

        }

        if(all_good){
                console.log(util.format(
                        'Successfully :Status - %s,message - "%s",metadata - "%s",endorsement signature :%s',
                        proposalResponses[0].response.status,
                        proposalResponses[0].response.message,
                        proposalResponses[0].response.payload,
                        proposalResponses[0].endorsement.signature));

		var prip = proposalResponses[0].response.payload.toString()
		//console.log(prip)
		
		var fsname = "GEN_PROVER_" + id + ".txt"
		var writeStream = fs.createWriteStream(fsname);
		writeStream.write(prip.split("---")[0], 'utf-8');

		// 标记写入完成
		writeStream.end();
		writeStream.on('finish', function() {
		    console.log('写入完成');
		})
		// 失败
		writeStream.on('error', function() {
		    console.log('写入失败');
		})		
		
                var request = {
                        proposalResponses:proposalResponses,
                        proposal:proposal,
                        header:header
                };

                var transactionID = tx_id.getTransactionID();
                return channel.sendTransaction(request);
        }

   },(err)=>{
        console.log('error',e);
   }).then( (sendtransresult) => {
	console.log(sendtransresult)
        return sendtransresult;
   },(err)=>{
        console.log('error',e);
   });
}


var getInstantiatedChaincodes = function() {

        return getOrgAdmin4Local().then( (user)=> {

                return channel.queryInstantiatedChaincodes(peer);

        },(err)=>{
        console.log('error',e);
        }).then( (instantiatedresult) => {
        return instantiatedresult;
   },(err)=>{
        console.log('error',e);
   });
}


function getOrgUserLocal(uname,orgname) { 
	if (orgname == "supervision") {
		mspID = "SupervisionMSP"
	}else if (orgname == "vegetablefruit1"){
		mspID = "Vegetablefruit1MSP"
	}else if (orgname == "vegetablefruit2"){
		mspID = "Vegetablefruit2MSP"
	}else if (orgname == "vegetablefruit3"){
		mspID = "Vegetablefruit3MSP"
	}else if (orgname == "vegetablefruit4"){
		mspID = "Vegetablefruit4MSP"
	}
        var keyPath="/home/hzaucoi/aberic/crypto-config/peerOrganizations/"+orgname+'.'+"qklfood.com"+'/users/'+uname+'@'+orgname+'.'+"qklfood.com"+"/msp/keystore";
        var keyPEM = Buffer.from(readAllFiles(keyPath)[0]).toString();
        var certPath="/home/hzaucoi/aberic/crypto-config/peerOrganizations/"+orgname+'.'+"qklfood.com"+'/users/'+uname+'@'+orgname+'.'+"qklfood.com"+"/msp/signcerts";
        var certPEM = readAllFiles(certPath)[0].toString();
      

        return hfc.newDefaultKeyValueStore({
                path:tempdir
        }).then((store) => {
                client.setStateStore(store);

                return client.createUser({
                        username: uname,
			mspid: mspID,
                        cryptoContent: {
                                privateKeyPEM: keyPEM,
                                signedCertPEM: certPEM
                        }
                });
        });
};


function readAllFiles(dir) {
        var files = fs.readdirSync(dir);
        var certs = [];
        files.forEach((file_name) => {
                let file_path = path.join(dir, file_name);
                let data = fs.readFileSync(file_path);
                certs.push(data);
        });
        return certs;
}


exports.sendTransaction_Dea_Ver = sendTransaction_Dea_Ver;
exports.queryCc = queryCc;
exports.getInstantiatedChaincodes = getInstantiatedChaincodes;
exports.sendTransaction_Prover = sendTransaction_Prover;


