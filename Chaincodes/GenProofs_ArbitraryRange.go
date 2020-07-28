package main

import (

  "fmt"
  "strings"
  "strconv"
  "encoding/json"
  "github.com/hyperledger/fabric/core/chaincode/shim"
  pb "github.com/hyperledger/fabric/protos/peer"
  "math/big"
  "mbp"
  "bytes"
    
)


type GenProofs_ArbitraryRange struct{}

func (t *GenProofs_ArbitraryRange) Init(stub shim.ChaincodeStubInterface) pb.Response {

    return shim.Success([]byte("Success invoke and not opter!!"))
}

func (t *GenProofs_ArbitraryRange) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
        
	_, args := stub.GetFunctionAndParameters()

	var opttype = args[0]   // generate_upload_proof & get_proof
	keyid := args[1]   // keyid
   
	if opttype == "generate_upload_proof" {

		//Get range [min,max) from RangeManagement chaincode according to rangeid 
		rangeid := args[2]
		valuerange_parm := []string{"invoke","get_range",rangeid}
		queryArgs := make([][]byte ,len(valuerange_parm))
		for i,arg := range valuerange_parm{
			queryArgs[i] = []byte(arg)
		}

		//invoke RangeManagement chaincode
		response := stub.InvokeChaincode("RangeManagement",queryArgs,"vegetableschannel") 
		if response.Status != shim.OK {
			errStr := fmt.Sprintf("failed to query chaincode.got error :%s",response.Payload)
			return shim.Error(errStr)
		}

		//procese the returned result
		result := string(response.Payload)
		rangestr := strings.Split(result,"~")
		
		value1, err := strconv.ParseUint(rangestr[0], 10, 64)
		if err != nil {
			return shim.Error("Value strconv ParseUint is error")
		}

		value2, err := strconv.ParseUint(rangestr[1], 10, 64)
		if err != nil {
			return shim.Error("Value strconv ParseUint is error")
		}
			
		var minvalue, maxvalue uint64
		if value1 < value2 {
			minvalue = value1
			maxvalue = value2
		}else{
			minvalue = value2
			maxvalue = value1
		}
			
		// process VecLength: 8bit 16bit 32bit 64bit
		VecLength,err := strconv.Atoi(args[3]) 
		if err != nil {
				return shim.Error("VecLength strconv operation is error")
		}
		
		//m represents the number of secret values
		m := len(args)-4
		EC := mbp.NewECPrimeGroupKey(VecLength*2*m)

		var i int
		var secretvalues uint64
		var values []*big.Int = make([]*big.Int, 2*m)

		for i=0;i<m;i++{

			secretvalues, err = strconv.ParseUint(args[i+4], 10, 64)
			if err != nil {
					return shim.Error("Secretvalues strconv operation is error")
			}
			
			v_bigInt := new(big.Int).SetUint64(secretvalues)

			// v - maxvalue + 2^n (n= 8 / 16 / 32 /64)
			p2 := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(VecLength)), EC.N)  //2^n
			v_max := new(big.Int).Sub(v_bigInt, new(big.Int).SetUint64(maxvalue))
			v_max.Add(v_max, p2)

			// v - minvalue
			v_min:= new(big.Int).Sub(v_bigInt, new(big.Int).SetUint64(minvalue))

			values[2*i] = v_max
			values[2*i+1] = v_min
		
		}
			
		//Generate arbitrary range proof (2m)
		proof := mbp.MRPProve(values,EC)
			
		proofJSONasBytes, err1 := json.Marshal(proof)
		if err1 != nil {
				return shim.Error(err1.Error())
		}

		//upload the rangeid, VecLength, the number of secret values, proof together into the blockchain
		bytes0 := [][]byte{[]byte(rangeid),[]byte(args[3]),[]byte(strconv.Itoa(m)),proofJSONasBytes}  

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

	err1 := shim.Start(new(GenProofs_ArbitraryRange))
	if err1 != nil {
		fmt.Printf("error starting simple chaincode:%s \n", err1)
	}

}
 

