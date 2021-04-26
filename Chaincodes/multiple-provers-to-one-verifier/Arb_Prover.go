package main

import (
	"bytes"
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

type Arb_Prover struct{}

type IdentityStore struct {
	Name      string
	Keyid     string
	Numid     int
	VecLength int
	M         int
	MinRange  uint64
	MaxRange  uint64
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
		response := stub.InvokeChaincode("cc_arb_verifier", queryArgs, "vegetablefruitchannel")
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
	response := stub.InvokeChaincode("cc_arb_verifier", queryArgs, "vegetablefruitchannel")
	if response.Status != shim.OK {
		errStr := fmt.Errorf("failed to query chaincode.got error :%s", response.Payload)
		return []byte(""), errStr
	}

	return response.Payload, nil
}

func (t *Arb_Prover) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success([]byte("Success invoke and not opter!!"))
}

func (t *Arb_Prover) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fn, args := stub.GetFunctionAndParameters()
	if fn == "V_A_and_S" {
		return t.V_A_and_S(stub, args)
	} else if fn == "T1_and_T2" {
		return t.T1_and_T2(stub, args)
	} else if fn == "OtherShare" {
		return t.OtherShare(stub, args)
	} else if fn == "Get_Cur_State" {
		return t.Get_Cur_State(stub, args)
	} else if fn == "Get_State_History" {
		return t.Get_State_History(stub, args)
	}
	
	return shim.Error("Recevied unkown function invocation")
}

// AandS function can be invoked by everyone, so we will judge the identity of user from verifier.
func (t *Arb_Prover) V_A_and_S(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	t0 := time.Now()
	// if length is equal 2, the second arg is verifier_keyid
	length := len(args)
	if (length != 1) && (length != 2) {
		return shim.Error("Incorrect number of arguments. Expecting 1 or 2")
	}

	creator := GetCreator(stub)
	if creator == "" {
		return shim.Error("The operator of Getcreator is failed!!!")
	}

	// args[1] is verifier_keyid, that is index
	var infoParam string
	if length == 2 {
		infoParam = args[1]
	} else if length == 1 {
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

	pro_name := "pro_" + creator
	for i, _ := range info {
		if info[i].Name == pro_name {
			proInfo := info[i]
			keyid := proInfo.Keyid
			id := proInfo.Numid
			VecLength := proInfo.VecLength
			m := 2 * proInfo.M
			EC := mbp.NewECPrimeGroupKey(VecLength * m)
			MinRange := proInfo.MinRange
			MaxRange := proInfo.MaxRange

			secret, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return shim.Error("The strconv operation of secret value is error")
			}
			value := new(big.Int).SetUint64(secret)

			// v - MaxRange + 2^n (n= 8 / 16 / 32 /64)
			p2 := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(VecLength)), EC.N) //2^n
			vmax := new(big.Int).Sub(value, new(big.Int).SetUint64(MaxRange))
			vmax.Add(vmax, p2)
			// v - MinRange
			vmin := new(big.Int).Sub(value, new(big.Int).SetUint64(MinRange))

			prip := make([]mbp.PrivateParams, 2)
			as := make([]mbp.ASCommitment, 2)
			as[0] = mbp.AandS(vmax, 2*id, m, &prip[0], EC)
			as[1] = mbp.AandS(vmin, 2*id+1, m, &prip[1], EC)
			asJSONasBytes, err := json.Marshal(as)
			if err != nil {
				return shim.Error("The marshal of AandS is error")
			}

			bytes0 := [][]byte{[]byte("AandS"), asJSONasBytes}
			err = stub.PutState(keyid, bytes.Join(bytes0, []byte("---")))
			if err != nil {
				return shim.Error("The operation of putstate is error")
			}

			pripJSONasBytes, err := json.Marshal(prip)
			if err != nil {
				return shim.Error("The marshal of prip is error")
			}

			elapsed := time.Since(t0)
			runtime := strconv.FormatFloat(elapsed.Seconds(), 'E', -1, 64)
			bytes1 := [][]byte{pripJSONasBytes, []byte("runtime:" + runtime)}
			return shim.Success(bytes.Join(bytes1, []byte("---")))
		}
	}
	return shim.Error("---The creator of the prover is not exists or correct!!!---Now the unexpected creator is " + pro_name)
}

func (t *Arb_Prover) T1_and_T2(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	t0 := time.Now()
	length := len(args)
	if (length != 1) && (length != 2) {
		return shim.Error("Incorrect number of arguments. Expecting 1 or 2")
	}

	creator := GetCreator(stub)
	if creator == "" {
		return shim.Error("The operator of Getcreator is failed!")
	}

	var infoParam string
	if length == 2 {
		infoParam = args[1]
	} else if length == 1 {
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

	length = len(info)
	//dealname := info[length-1].Name
	dealkeyid := info[length-1].Keyid
	invoke_parm2 := []string{"Get_MPC_Range_Prf", dealkeyid}
	queryArgs2 := make([][]byte, len(invoke_parm2))
	for i, arg := range invoke_parm2 {
		queryArgs2[i] = []byte(arg)
	}
	response := stub.InvokeChaincode("cc_arb_dealer", queryArgs2, "vegetablefruitchannel")
	if response.Status != shim.OK {
		errStr := fmt.Sprintf("failed to query chaincode.got error :%s", response.Payload)
		return shim.Error(errStr)
	}
	resstr := string(response.Payload)
	values := strings.Split(resstr, "---")
	if len(values) == 3 {
		return shim.Error("Dealer doesn't work normally!")
	}
	if values[0] != "Fait_y_z" {
		return shim.Error("The result of getMPCRP is not unexpected!!!")
	}
	var mpcrp mbp.MPCRangeProof
	err = json.Unmarshal([]byte(values[1]), &mpcrp)
	if err != nil {
		return shim.Error(err.Error() + "mpcrp is not unmarshal")
	}
	cy := mpcrp.Cy
	cz := mpcrp.Cz

	pro_name := "pro_" + creator
	for i, _ := range info {
		if info[i].Name == pro_name {
			proInfo := info[i]
			prokeyid := proInfo.Keyid
			id := proInfo.Numid
			VecLength := proInfo.VecLength
			m := 2 * proInfo.M
			EC := mbp.NewECPrimeGroupKey(VecLength * m)

			str_prip := args[0]
			var prip []mbp.PrivateParams
			err = json.Unmarshal([]byte(str_prip), &prip)
			if err != nil {
				return shim.Error(err.Error())
			}

			t1t2 := make([]mbp.T1T2Commitment, 2)
			t1t2[0] = mbp.T1andT2(prip[0].V, 2*id, m, cy, cz, &prip[0], EC)
			t1t2[1] = mbp.T1andT2(prip[1].V, 2*id+1, m, cy, cz, &prip[1], EC)
			t1t2JSONasBytes, err := json.Marshal(t1t2)
			if err != nil {
				return shim.Error(err.Error())
			}

			bytes0 := [][]byte{[]byte("T1andT2"), t1t2JSONasBytes}
			err = stub.PutState(prokeyid, bytes.Join(bytes0, []byte("---")))
			if err != nil {
				return shim.Error(err.Error())
			}

			pripJSONasBytes, err := json.Marshal(prip)
			if err != nil {
				return shim.Error(err.Error())
			}

			elapsed := time.Since(t0)
			runtime := strconv.FormatFloat(elapsed.Seconds(), 'E', -1, 64)
			bytes1 := [][]byte{pripJSONasBytes, []byte("runtime:" + runtime)}
			return shim.Success(bytes.Join(bytes1, []byte("---")))
		}
	}
	return shim.Error("---The creator of the prover is not exists or correct!!!---Now the unexpected creator is " + pro_name)
}

func (t *Arb_Prover) OtherShare(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	t0 := time.Now()
	length := len(args)
	if (length != 1) && (length != 2) {
		return shim.Error("Incorrect number of arguments. Expecting 1 or 2")
	}
	creator := GetCreator(stub)
	if creator == "" {
		return shim.Error("The operator of Getcreator is failed!")
	}

	var infoParam string
	if length == 2 {
		infoParam = args[1]
	} else if length == 1 {
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

	length = len(info)
	dealkeyid := info[length-1].Keyid
	invoke_parm2 := []string{"Get_MPC_Range_Prf", dealkeyid}
	queryArgs2 := make([][]byte, len(invoke_parm2))
	for i, arg := range invoke_parm2 {
		queryArgs2[i] = []byte(arg)
	}
	response := stub.InvokeChaincode("cc_arb_dealer", queryArgs2, "vegetablefruitchannel")
	if response.Status != shim.OK {
		errStr := fmt.Sprintf("failed to query chaincode.got error :%s", response.Payload)
		return shim.Error(errStr)
	}
	resstr := string(response.Payload)
	values := strings.Split(resstr, "---")
	if len(values) == 3 {
		return shim.Error("Dealer doesn't work normally!")
	}
	if values[0] != "Fait_x" {
		return shim.Error("The result of getMPCRP is not unexpected!!!")
	}

	var mpcrp mbp.MPCRangeProof
	err = json.Unmarshal([]byte(values[1]), &mpcrp)
	if err != nil {
		return shim.Error(err.Error() + "mpcrp is not unmarshal")
	}
	cx := mpcrp.Cx
	cy := mpcrp.Cy
	cz := mpcrp.Cz

	pro_name := "pro_" + creator
	for i, _ := range info {
		if info[i].Name == pro_name {
			proInfo := info[i]
			prokeyid := proInfo.Keyid
			id := proInfo.Numid
			VecLength := proInfo.VecLength
			m := 2 * proInfo.M
			EC := mbp.NewECPrimeGroupKey(VecLength * m)

			str_prip := args[0]
			var prip []mbp.PrivateParams
			err = json.Unmarshal([]byte(str_prip), &prip)
			if err != nil {
				return shim.Error(err.Error())
			}
			share := make([]mbp.OtherShare, 2)
			share[0] = mbp.ProFinal(prip[0].V, 2*id, m, cy, cz, cx, &prip[0], EC)
			share[1] = mbp.ProFinal(prip[1].V, 2*id+1, m, cy, cz, cx, &prip[1], EC)
			shareJSONasBytes, err := json.Marshal(share)
			if err != nil {
				return shim.Error(err.Error())
			}

			bytes0 := [][]byte{[]byte("otherShare"), shareJSONasBytes}
			err = stub.PutState(prokeyid, bytes.Join(bytes0, []byte("---")))
			if err != nil {
				return shim.Error(err.Error())
			}

			elapsed := time.Since(t0)
			runtime := strconv.FormatFloat(elapsed.Seconds(), 'E', -1, 64)
			return shim.Success([]byte("The operation of prover has finished!---" + "runtime:" + runtime))
		}
	}
	return shim.Error("---The creator of the prover is not exists or correct!!!---Now the unexpected creator is " + pro_name)
}

// get AandS  T1andT2  and othershare
func (t *Arb_Prover) Get_Cur_State(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

func (t *Arb_Prover) Get_State_History(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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
	err1 := shim.Start(new(Arb_Prover))
	if err1 != nil {
		fmt.Printf("error starting simple chaincode:%s \n", err1)
	}
}
