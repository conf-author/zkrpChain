package main

import (

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"fmt"
	"strconv"
	"time"
	"encoding/json"


)

type txinfo struct {
    
	Txid string
	Txvalue string
	Txstatus string 
	Datestr string
}

type RangeManagement struct{}

func (t *RangeManagement) Init(stub shim.ChaincodeStubInterface) pb.Response {

	return shim.Success([]byte("Success invoke and not opter!!"))

}

func (t *RangeManagement) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	_, args := stub.GetFunctionAndParameters()

	var opttype = args[0]   // upload_range & get_range & get_range_history
	var rangeid = args[1]   // keyid
	
	
	if opttype == "upload_range" {
	
		if len(args) != 4 {
            		return shim.Error("Incorrect number of arguments. Expecting 4")
       		}
		
		Aval, err := strconv.Atoi(args[2])
		if err != nil {
		        return shim.Error("Expecting integer value for asset holding")
		}

		Bval, err := strconv.Atoi(args[3])
		if err != nil {
		        return shim.Error("Expecting integer value for asset holding")
		}

		err = stub.PutState(rangeid, []byte(strconv.Itoa(Aval)+"~"+strconv.Itoa(Bval)))
		if err != nil {
		        return shim.Error(err.Error())
		}

		return shim.Success([]byte("success put:" + strconv.Itoa(Aval)+"~"+strconv.Itoa(Bval)))

	}else if opttype == "get_range"{
		                
		var keyrange []byte
		var err error 

		keyrange,err = stub.GetState(rangeid)
		if(err != nil) {
			return shim.Error(err.Error())
		}

		return shim.Success(keyrange)
		
	}else if opttype == "get_range_history"{
		                
       		keyIters, err := stub.GetHistoryForKey(rangeid)
		if err != nil {
			return shim.Error(err.Error())
		}
		
		defer keyIters.Close()
		
		var history []txinfo
		for keyIters.HasNext() {
		
			response, iterErr := keyIters.Next()
			if iterErr != nil {
				return shim.Error(err.Error())
			}
			
			tx := txinfo{}
			// Transaction id
			tx.Txid = response.TxId
			// Transaction value
			tx.Txvalue = string(response.Value)
			// Current transaction status
			tx.Txstatus = strconv.FormatBool(response.IsDelete)
			// Transaction timestamp
			txtimesamp :=response.Timestamp
			tm := time.Unix(txtimesamp.Seconds, 0)
			tx.Datestr = tm.Format("2006-01-02 03:04:05 PM")
			
			history = append( history , tx) 
			
		}
			
		jsonHistory, err := json.Marshal(history)
		if err != nil {
			return shim.Error("Marshal operation is failed")
		}

		return shim.Success(jsonHistory)
		
	}else{
		return shim.Success( []byte("Success invoke and no operation !!") )
	}

}

func main() {

	err := shim.Start(new(RangeManagement))
	if err != nil{
		fmt.Printf("error starting simple chaincode:%s \n",err)
	}
        
}



