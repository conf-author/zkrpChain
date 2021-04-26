package main

import (
	"bytes"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"mbp"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Dealer struct{}

type IdentityStore struct {
	Name      string
	Keyid     string
	Numid     int
	VecLength int
	M         int
}

func GetCreator(stub shim.ChaincodeStubInterface) string {
	creatorByte, _ := stub.GetCreator()
	value := strings.Split(string(creatorByte), "BEG")[1]
	creatorStr := "-----BEG" + value
	certText := []byte(creatorStr)
	bl, _ := pem.Decode(certText)
	if bl == nil {
		fmt.Println("Could not decode the PEM structure")
		return ""
	}

	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		fmt.Println("ParseCertificate failed")
		return ""
	}
	uname := cert.Subject.CommonName
	return uname
}

func getPandDInfo(stub shim.ChaincodeStubInterface, arg string) ([]byte, error) {
	var verify_keyid string
	if arg == "" {
		invoke_parm := []string{"Get_KeyID_Sess_Num"}
		queryArgs := make([][]byte, len(invoke_parm))
		for i, arg := range invoke_parm {
			queryArgs[i] = []byte(arg)
		}
		response := stub.InvokeChaincode("cc_verifier", queryArgs, "vegetablefruitchannel")
		if response.Status != shim.OK {
			errStr := fmt.Errorf("failed to query chaincode.got error :%s", response.Payload)
			return []byte(""), errStr
		}
		verify_keyid = string(response.Payload)
	} else {
		verify_keyid = arg
	}

	invoke_parm := []string{"Get_Setup_Info", verify_keyid}
	queryArgs := make([][]byte, len(invoke_parm))
	for i, arg := range invoke_parm {
		queryArgs[i] = []byte(arg)
	}
	response := stub.InvokeChaincode("cc_verifier", queryArgs, "vegetablefruitchannel")
	if response.Status != shim.OK {
		errStr := fmt.Errorf("failed to query chaincode.got error :%s", response.Payload)
		return []byte(""), errStr
	}

	return response.Payload, nil
}

func (t *Dealer) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success([]byte("Success invoke and not opter!!"))
}

func (t *Dealer) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fn, args := stub.GetFunctionAndParameters()
	if fn == "Fait_y_z" {
		return t.Fait_y_z(stub, args)
	} else if fn == "Fait_x" {
		return t.Fait_x(stub, args)
	} else if fn == "Gen_Prf" {
		return t.Gen_Prf(stub, args)
	} else if fn == "Get_MPC_Range_Prf" {
		return t.Get_MPC_Range_Prf(stub, args)
	} else if fn == "Get_Range_History" {
		return t.Get_Range_History(stub, args)
	}
	return shim.Error("Recevied unkown function invocation")
}


func (t *Dealer) Fait_y_z(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	t0 := time.Now()
	length := len(args)
	if (length != 0) && (length != 1) {
		return shim.Error("Incorrect number of arguments. Expecting 0 or 1")
	}

	creator := GetCreator(stub)
	if creator == "" {
		return shim.Error("The operator of Getcreator is failed!")
	}

	// args[0] is verifier_keyid, that is index
	var infoParam string
	if length == 1 {
		infoParam = args[0]
	} else if length == 0 {
		infoParam = ""
	}
	result, err := getPandDInfo(stub, infoParam)
	if err != nil {
		return shim.Error(err.Error())
	}
	var info []IdentityStore
	err = json.Unmarshal(result, &info)
	if err != nil {
		return shim.Error(err.Error())
	}

	deal_name := "deal_" + creator
	length = len(info)
	if info[length-1].Name == deal_name {
		dealInfo := info[length-1]
		dealkeyid := dealInfo.Keyid
		//id := dealInfo.Numid
		VecLength := dealInfo.VecLength
		m := dealInfo.M
		EC := mbp.NewECPrimeGroupKey(VecLength * m)

		mpcrp := mbp.MPCRangeProof{}
		Comms := make([]mbp.ECPoint, m)
		countA := EC.Zero()
		countS := EC.Zero()
		// get prover info
		for j := 0; j < m; j++ {
			proInfo := info[j]
			//proname := proInfo.Name
			prokeyid := proInfo.Keyid
			invoke_parm := []string{"Get_Cur_State", prokeyid}
			queryArgs := make([][]byte, len(invoke_parm))
			for i, arg := range invoke_parm {
				queryArgs[i] = []byte(arg)
			}
			response := stub.InvokeChaincode("cc_prover", queryArgs, "vegetablefruitchannel") 
			if response.Status != shim.OK {
				errStr := fmt.Sprintf("failed to query chaincode.got error :%s", response.Payload)
				return shim.Error(errStr)
			}

			resstr := string(response.Payload)
			values := strings.Split(resstr, "---")

			var as mbp.ASCommitment
			err = json.Unmarshal([]byte(values[1]), &as)
			if err != nil {
				return shim.Error(err.Error())
			}
			countA = countA.Add(as.A, EC)
			countS = countS.Add(as.S, EC)
			Comms[j] = as.Comm
		}
		mpcrp.Comms = Comms
		mpcrp.A = countA
		mpcrp.S = countS

		chal1s256 := sha256.Sum256([]byte(countA.X.String() + countA.Y.String()))
		cy := new(big.Int).SetBytes(chal1s256[:])
		mpcrp.Cy = cy
		chal2s256 := sha256.Sum256([]byte(countS.X.String() + countS.Y.String()))
		cz := new(big.Int).SetBytes(chal2s256[:])
		mpcrp.Cz = cz
		mpcrpJSONasBytes, err := json.Marshal(mpcrp)
		if err != nil {
			return shim.Error(err.Error())
		}

		bytes0 := [][]byte{[]byte("Fait_y_z"), mpcrpJSONasBytes}
		err = stub.PutState(dealkeyid, bytes.Join(bytes0, []byte("---")))
		if err != nil {
			return shim.Error(err.Error())
		}
		elapsed := time.Since(t0)
		runtime := strconv.FormatFloat(elapsed.Seconds(), 'E', -1, 64)
		return shim.Success([]byte("Successful generate cy and cz!!!---" + "runtime:" + runtime))
	}
	return shim.Error("The dealer is not exists or correct.---creator:" + deal_name)
}

func (t *Dealer) Fait_x(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	t0 := time.Now()
	length := len(args)
	if (length != 0) && (length != 1) {
		return shim.Error("Incorrect number of arguments. Expecting 0 or 1")
	}

	creator := GetCreator(stub)
	if creator == "" {
		return shim.Error("The operator of Getcreator is failed!")
	}

	// args[0] is verifier_keyid, that is index
	var infoParam string
	if length == 1 {
		infoParam = args[0]
	} else if length == 0 {
		infoParam = ""
	}
	result, err := getPandDInfo(stub, infoParam)
	if err != nil {
		return shim.Error(err.Error())
	}
	var info []IdentityStore
	err = json.Unmarshal(result, &info)
	if err != nil {
		return shim.Error(err.Error())
	}

	deal_name := "deal_" + creator
	length = len(info)
	if info[length-1].Name == deal_name {
		dealInfo := info[length-1]
		dealkeyid := dealInfo.Keyid
		//id := dealInfo.Numid
		VecLength := dealInfo.VecLength
		m := dealInfo.M
		EC := mbp.NewECPrimeGroupKey(VecLength * m)

		mpcrpBytes, err := stub.GetState(dealkeyid)
		if err != nil {
			return shim.Error(err.Error())
		}
		resstr := string(mpcrpBytes)
		values := strings.Split(resstr, "---")
		var mpcrp mbp.MPCRangeProof
		err = json.Unmarshal([]byte(values[1]), &mpcrp)
		if err != nil {
			return shim.Error(err.Error() + "mpcrp is not unmarshal")
		}

		countT1 := EC.Zero()
		countT2 := EC.Zero()
		// get prover info
		for j := 0; j < m; j++ {
			proInfo := info[j]
			//proname := proInfo.Name
			prokeyid := proInfo.Keyid
			invoke_parm := []string{"Get_Cur_State", prokeyid}
			queryArgs := make([][]byte, len(invoke_parm))
			for i, arg := range invoke_parm {
				queryArgs[i] = []byte(arg)
			}
			response := stub.InvokeChaincode("cc_prover", queryArgs, "vegetablefruitchannel") 
			if response.Status != shim.OK {
				errStr := fmt.Sprintf("failed to query chaincode.got error :%s", response.Payload)
				return shim.Error(errStr)
			}

			result := string(response.Payload)
			values := strings.Split(result, "---")
			var t1t2 mbp.T1T2Commitment
			err = json.Unmarshal([]byte(values[1]), &t1t2)
			if err != nil {
				return shim.Error(err.Error())
			}
			countT1 = countT1.Add(t1t2.T1, EC)
			countT2 = countT2.Add(t1t2.T2, EC)
		}
		mpcrp.T1 = countT1
		mpcrp.T2 = countT2

		chal3s256 := sha256.Sum256([]byte(countT1.X.String() + countT1.Y.String() + countT2.X.String() + countT2.Y.String()))
		cx := new(big.Int).SetBytes(chal3s256[:])
		mpcrp.Cx = cx
		mpcrpJSONasBytes, err := json.Marshal(mpcrp)
		if err != nil {
			return shim.Error(err.Error())
		}

		bytes0 := [][]byte{[]byte("Fait_x"), mpcrpJSONasBytes}
		err = stub.PutState(dealkeyid, bytes.Join(bytes0, []byte("---")))
		if err != nil {
			return shim.Error(err.Error())
		}
		elapsed := time.Since(t0)
		runtime := strconv.FormatFloat(elapsed.Seconds(), 'E', -1, 64)
		return shim.Success([]byte("Successful generate cx!!!---" + "runtime:" + runtime))
	}
	return shim.Error("The dealer is not exists or correct.---creator:" + deal_name)
}

func (t *Dealer) Gen_Prf(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	t0 := time.Now()
	length := len(args)
	if (length != 0) && (length != 1) {
		return shim.Error("Incorrect number of arguments. Expecting 0 or 1")
	}
	creator := GetCreator(stub)
	if creator == "" {
		return shim.Error("The operator of Getcreator is failed!")
	}

	var verify_keyid string
	if length == 0 {
		invoke_parm := []string{"Get_KeyID_Sess_Num"}
		queryArgs := make([][]byte, len(invoke_parm))
		for i, arg := range invoke_parm {
			queryArgs[i] = []byte(arg)
		}
		response := stub.InvokeChaincode("cc_verifier", queryArgs, "vegetablefruitchannel")
		if response.Status != shim.OK {
			errStr := fmt.Sprintf("failed to query chaincode.got error :%s", response.Payload)
			return shim.Error(errStr)
		}
		verify_keyid = string(response.Payload)
	} else if length == 1 {
		verify_keyid = args[0]
	}
	invoke_parm := []string{"Get_Setup_Info", verify_keyid}
	queryArgs := make([][]byte, len(invoke_parm))
	for i, arg := range invoke_parm {
		queryArgs[i] = []byte(arg)
	}
	response := stub.InvokeChaincode("cc_verifier", queryArgs, "vegetablefruitchannel")
	if response.Status != shim.OK {
		errStr := fmt.Sprintf("failed to query chaincode.got error :%s", response.Payload)
		return shim.Error(errStr)
	}

	var info []IdentityStore
	err := json.Unmarshal(response.Payload, &info)
	if err != nil {
		return shim.Error(err.Error())
	}

	deal_name := "deal_" + creator
	length = len(info)
	if info[length-1].Name == deal_name {
		dealInfo := info[length-1]
		dealkeyid := dealInfo.Keyid
		//id := dealInfo.Numid
		VecLength := dealInfo.VecLength
		m := dealInfo.M
		EC := mbp.NewECPrimeGroupKey(VecLength * m)

		mpcrpBytes, err := stub.GetState(dealkeyid)
		if err != nil {
			return shim.Error(err.Error())
		}
		result := string(mpcrpBytes)
		values := strings.Split(result, "---")
		var mpcrp mbp.MPCRangeProof
		if verify_keyid == values[1] {
			err = json.Unmarshal([]byte(values[2]), &mpcrp)
		} else {
			err = json.Unmarshal([]byte(values[1]), &mpcrp)
		}
		if err != nil {
			return shim.Error(err.Error() + "mpcrp is not unmarshal")
		}

		bitsPerValue := EC.V / m
		countThat := big.NewInt(0)
		countTaux := big.NewInt(0)
		countMu := big.NewInt(0)
		left := make([]*big.Int, VecLength*m)
		right := make([]*big.Int, VecLength*m)

		for j := 0; j < m; j++ {
			proInfo := info[j]
			//proname := proInfo.Name
			prokeyid := proInfo.Keyid
			invoke_parm := []string{"Get_Cur_State", prokeyid}
			queryArgs := make([][]byte, len(invoke_parm))
			for i, arg := range invoke_parm {
				queryArgs[i] = []byte(arg)
			}
			response := stub.InvokeChaincode("cc_prover", queryArgs, "vegetablefruitchannel") 
			if response.Status != shim.OK {
				errStr := fmt.Sprintf("failed to query chaincode.got error :%s", response.Payload)
				return shim.Error(errStr)
			}

			result := string(response.Payload)
			values := strings.Split(result, "---")

			var share mbp.OtherShare
			err = json.Unmarshal([]byte(values[1]), &share)
			if err != nil {
				return shim.Error(err.Error())
			}
			countThat = new(big.Int).Add(countThat, share.That)
			countTaux = new(big.Int).Add(countTaux, share.Taux)
			countMu = new(big.Int).Add(countMu, share.Mu)
			for i := 0; i < bitsPerValue; i++ {
				left[bitsPerValue*j+i] = share.Left[i]
				right[bitsPerValue*j+i] = share.Right[i]
			}
		}
		mpcrp.Tau = countTaux
		mpcrp.Th = countThat
		mpcrp.Mu = countMu
		countLeft := left
		countRight := right
		HPrime := make([]mbp.ECPoint, len(EC.BPH))
		A := mpcrp.A
		chal1s256 := sha256.Sum256([]byte(A.X.String() + A.Y.String()))
		cy := new(big.Int).SetBytes(chal1s256[:])
		PowerOfCY := mbp.PowerVector(EC.V, cy, EC)
		for j := 0; j < m; j++ {
			for i := 0; i < bitsPerValue; i++ {
				HPrime[j*bitsPerValue+i] = EC.BPH[j*bitsPerValue+i].Mult(new(big.Int).ModInverse(PowerOfCY[j*bitsPerValue+i], EC.N), EC)
			}
		}

		P := mbp.TwoVectorPCommitWithGens(EC.BPG, HPrime, countLeft, countRight, EC)
		that := mbp.InnerProduct(countLeft, countRight, EC)
		IPP := mbp.InnerProductProve(countLeft, countRight, that, P, EC.U, EC.BPG, HPrime, EC)
		mpcrp.IPP = IPP

		mpcrpJSONasBytes, err := json.Marshal(mpcrp)
		if err != nil {
			return shim.Error(err.Error())
		}
		bytes0 := [][]byte{[]byte("genProof"), []byte(verify_keyid), mpcrpJSONasBytes}
		err = stub.PutState(dealkeyid, bytes.Join(bytes0, []byte("---")))
		if err != nil {
			return shim.Error(err.Error())
		}
		elapsed := time.Since(t0)
		runtime := strconv.FormatFloat(elapsed.Seconds(), 'E', -1, 64)
		return shim.Success([]byte("Successful generate proof!!!---" + "runtime:" + runtime))
	}
	return shim.Error("The dealer is not exists or correct.---creator:" + deal_name)
}

// get the part and all of MPC range proof
func (t *Dealer) Get_MPC_Range_Prf(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	keyid := args[0]
	state, err := stub.GetState(keyid)
	fmt.Println(string(state))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(state)
}

func (t *Dealer) Get_Range_History(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1 (keyid)")
	}

	keyid := args[0]
	keysIter, err := stub.GetHistoryForKey(keyid)
	if err != nil {
		return shim.Error(fmt.Sprintf("GetHistoryForKey failed. Error accessing state:%s", err))
	}
	defer keysIter.Close()

	var keys []string
	//var buffer bytes.Buffer
	for keysIter.HasNext() {
		response, iterErr := keysIter.Next()
		if iterErr != nil {
			return shim.Error(fmt.Sprintf("GetHistoryForKey operation failed. Error accessing state:%s", err))
		}
		keys = append(keys, string(response.Value))

	}
	keysJSONasBytes, err2 := json.Marshal(keys)
	if err2 != nil {
		return shim.Error(err2.Error())
	}
	return shim.Success(keysJSONasBytes)
}

func main() {
	err1 := shim.Start(new(Dealer))
	if err1 != nil {
		fmt.Printf("error starting simple chaincode:%s \n", err1)
	}
}
