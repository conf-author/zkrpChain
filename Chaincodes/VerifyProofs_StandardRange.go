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


type VerifyProofs_StandardRange struct{}


func (t *VerifyProofs_StandardRange) Init(stub shim.ChaincodeStubInterface) pb.Response {

    return shim.Success([]byte("Success invoke and not opter!!"))

}


func (t *VerifyProofs_StandardRange) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
        
	_, args := stub.GetFunctionAndParameters()

	//Get proof from GenProofs_StandardRange chaincode according to keyid 
	keyid := args[0]  
	valuerange_parm := []string{"invoke","get_proof",keyid}
	queryArgs := make([][]byte ,len(valuerange_parm))
	for i,arg := range valuerange_parm{
		queryArgs[i] = []byte(arg)
	}
	
	//invoke GenProofs_StandardRange chaincode
	response := stub.InvokeChaincode("GenProofs_StandardRange",queryArgs,"vegetableschannel") 
	if response.Status != shim.OK {
			errStr := fmt.Sprintf("failed to query chaincode.got error :%s",response.Payload)
			return shim.Error(errStr)
	}

	result := string(response.Payload)
	values := strings.Split(result,"---")

	VecLength,err := strconv.Atoi(values[0])                
	if err != nil {
			return shim.Error("VecLength strconv operation is error")
	}

	m,err := strconv.Atoi(values[1])                
	if err != nil {
			return shim.Error("strconv Atoi is error")
	}

	EC := mbp.NewECPrimeGroupKey(VecLength*m)

	var proof mbp.MultiRangeProof
	err = json.Unmarshal([]byte(values[2]), &proof)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Verify the standard range proof
	ok := mbp.MRPVerify(proof,EC)
	if ok {

		return shim.Success([]byte("Multi Range Proof Verification is true!!!"))

	} else {

		return shim.Error("Multi Range Proof Verification is false")

	}
        
}

func main() {

	err1 := shim.Start(new(VerifyProofs_StandardRange))
	if err1 != nil {
			fmt.Printf("error starting simple chaincode:%s\n", err1)
	}

}


