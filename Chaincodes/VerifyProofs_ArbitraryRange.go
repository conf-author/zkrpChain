package main

import (

	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"mbp"
	"strconv"
	"strings"

)

type txinfo struct {
    
	Txid string
	Txvalue string
	Txstatus string 
	Datestr string
}


type VerifyProofs_ArbitraryRange struct{}


func (t *VerifyProofs_ArbitraryRange) Init(stub shim.ChaincodeStubInterface) pb.Response {

    return shim.Success([]byte("Success invoke and not opter!!"))

}


func (t *VerifyProofs_ArbitraryRange) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
        
    _, args := stub.GetFunctionAndParameters()
	keyid := args[0]  

	//Verify that the range are correct or not modified
	min := args[1]
	max := args[2]
	rangestr := min + "~" + max
		
	//Get proof from GenProofs_ArbitraryRang chaincode according to keyid 
	valuerange_parm := []string{"invoke","get_proof",keyid}
	queryArgs := make([][]byte ,len(valuerange_parm))
	for i,arg := range valuerange_parm{
			queryArgs[i] = []byte(arg)
	}

	//invoke GenProofs_ArbitraryRange chaincode
	response := stub.InvokeChaincode("GenProofs_ArbitraryRange",queryArgs,"vegetableschannel") 
	if response.Status != shim.OK {
			errStr := fmt.Sprintf("failed to query chaincode.got error :%s",response.Payload)
			return shim.Error(errStr)
	}

	result := string(response.Payload)
	values := strings.Split(result,"---")
	
	//Get history of range  from RangeManagement chaincode according to rangeid
	valuerange_parm = []string{"invoke","get_range_history",values[0]}
	queryArgs = make([][]byte ,len(valuerange_parm))
	for i,arg := range valuerange_parm{
		queryArgs[i] = []byte(arg)
	}
	//invoke RangeManagement chaincode
	response = stub.InvokeChaincode("RangeManagement",queryArgs,"vegetableschannel") 
	if response.Status != shim.OK {
		errStr := fmt.Sprintf("failed to query chaincode.got error :%s",response.Payload)
		return shim.Error(errStr)
	}

	//procese the returned result
	var info []txinfo
	err := json.Unmarshal(response.Payload, &info)
	if err != nil {
		return shim.Error(err.Error())
	}

	for _, v := range info{
	
		if v.Txvalue != rangestr {
			return shim.Error("Value range is error or modified !!!") 
		}
	}
	
	VecLength,err := strconv.Atoi(values[1])                
	if err != nil {
			return shim.Error("VecLength strconv operation is error")
	}
	
	m,err := strconv.Atoi(values[2])                
	if err != nil {
			return shim.Error("strconv Atoi is error")
	}

	EC := mbp.NewECPrimeGroupKey(VecLength*2*m)

	var proof mbp.MultiRangeProof
	err = json.Unmarshal([]byte(values[3]), &proof)
	if err != nil {
			return shim.Error(err.Error())
	}

	// Verify the arbitrary range proof
	ok := mbp.MRPVerify(proof,EC)
	if ok {

		return shim.Success([]byte("Multi Range Proof Verification is true !"))

	} else {
		return shim.Error("Multi Range Proof Verification is false")
	}
        
}

func main() {

	err1 := shim.Start(new(VerifyProofs_ArbitraryRange))
	if err1 != nil {
		fmt.Printf("error starting simple chaincode:%s\n", err1)
	}

}

