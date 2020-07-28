package main

import (

	"fmt"
	"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"math/big"
	"mbp"
	"bytes"
        
)

type GenProofs_StandardRange struct{}

func (t *GenProofs_StandardRange) Init(stub shim.ChaincodeStubInterface) pb.Response {

    return shim.Success([]byte("success invoke and not opter!!"))
}

func (t *GenProofs_StandardRange) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
        
	_, args := stub.GetFunctionAndParameters()
	var opttype = args[0]    // generate_upload_proof & get_proof
	keyid := args[1] 

	if opttype == "generate_upload_proof" {
	
	    // process VecLength: 8bit 16bit 32bit 64bit
		VecLength, err := strconv.Atoi(args[2])
		if err != nil {
				return shim.Error("VecLength strconv operation is error")
		}

		//m represents the number of secret values
		m := len(args)-3
		EC := mbp.NewECPrimeGroupKey(VecLength*m)

		var i int
		var secretvalues uint64
		var values []*big.Int = make([]*big.Int, m)

		for i=0;i<m;i++{

			secretvalues, err = strconv.ParseUint(args[i+3], 10, 64)
			if err != nil {
					return shim.Error("secretvalues strconv ParseInt is error")
			}
			
			v_bigInt := new(big.Int).SetUint64(secretvalues)
			values[i] = v_bigInt

		}

		//Generate standard range proof 
		proof := mbp.MRPProve(values,EC)

		proofJSONasBytes, err1 := json.Marshal(proof)
		if err1 != nil {
				return shim.Error(err1.Error())
		}

		// upload VecLength, the number of secret values, proof together into the blockchain
		bytes0 := [][]byte{[]byte(args[2]),[]byte(strconv.Itoa(m)),proofJSONasBytes}              

		err1 = stub.PutState(keyid, bytes.Join(bytes0,[]byte("---")))
		if err1 != nil {
				 return shim.Error(err1.Error())
		}

		return shim.Success([]byte("Successfully generate multi range proof !"))

	}else if opttype == "get_proof"{

		keyproof,err := stub.GetState(keyid)
		if(err != nil) {
				return shim.Error(err.Error())
		}

		return shim.Success(keyproof)

	}else{

		return shim.Success( []byte("Success invoke and no operation !!") )
	}

}

func main() {

	err1 := shim.Start(new(GenProofs_StandardRange))
	if err1 != nil {
		fmt.Printf("error starting simple chaincode:%s \n", err1)
	}

}
 
  
