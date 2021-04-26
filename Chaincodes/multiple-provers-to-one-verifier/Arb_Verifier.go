package main

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/rand"
	"mbp"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Arb_Verifier struct{}

func isPowerOfTwo(n int) bool {
	return n > 0 && n&(n-1) == 0
}

// Verifier is controller to control and check the identity of prover and dealer
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

func (t *Arb_Verifier) Init(stub shim.ChaincodeStubInterface) pb.Response {
	err := stub.PutState("SessionNum", []byte("0"))
	if err != nil {
		return shim.Error("The operation of putstate is error")
	}
	return shim.Success([]byte("Success invoke and not opter!!"))
}

func (t *Arb_Verifier) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fn, args := stub.GetFunctionAndParameters()
	if fn == "Init_Setup" {
		return t.Init_Setup(stub, args)
	} else if fn == "Get_Setup_Info" {
		return t.Get_Setup_Info(stub, args)
	} else if fn == "Get_KeyID_Sess_Num" {
		return t.Get_KeyID_Sess_Num(stub, args)
	} else if fn == "Ver_Prf" {
		return t.Ver_Prf(stub, args)
	}
	return shim.Error("Recevied unkown function invocation")
}

// Only the admin of verifier can write the controltable
func (t *Arb_Verifier) Init_Setup(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	t0 := time.Now()
	creator := GetCreator(stub)
	if creator == "" {
		return shim.Error("The operator of Getcreator is failed!")
	}

	// local msp or fabric-ca start
	if (creator != "Admin@supervision.qklfood.com") && (creator != "supervision_admin") {
		return shim.Error("The verfier is not a admin, so it can not modify the information!---creator:" + creator)
	}

	result, err := stub.GetState("SessionNum")
	if err != nil {
		return shim.Error(err.Error())
	}
	index, err := strconv.Atoi(string(result))
	index += 1
	verifier_keyid := strconv.Itoa(index)

	// args: Veclength xxx arbrange xxx prover xxx  xxx  dealer xxx
	length := len(args)
	if args[0] != "VecLength" || args[2] != "ArbRange" {
		return shim.Error("The setting of args is error!!!")
	}

	VecLength, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}

	// process arbitrary range
	arbrange := strings.Split(args[3], "-")
	range0, err := strconv.ParseUint(arbrange[0], 10, 64)
	if err != nil {
		return shim.Error(err.Error())
	}
	range1, err := strconv.ParseUint(arbrange[1], 10, 64)
	if err != nil {
		return shim.Error(err.Error())
	}
	var MinRange, MaxRange uint64
	if range1 > range0 {
		MinRange = range0
		MaxRange = range1
	} else {
		MinRange = range1
		MaxRange = range0
	}

	// the number of prover
	m := length - 7
	if !isPowerOfTwo(m) || !isPowerOfTwo(VecLength) {
		return shim.Error("The length of vector or the number of party should be the power of 2")
	}
	rand.Seed(time.Now().UnixNano())
	info := make([]IdentityStore, m+1)
	// prover information
	for i := 0; i < m; i++ {
		name := "pro_" + args[i+5]
		num := rand.Intn(1000000)
		info[i].Name = name
		info[i].Keyid = name + strconv.Itoa(num)
		info[i].Numid = i
		info[i].VecLength = VecLength
		info[i].M = m
		info[i].MinRange = MinRange
		info[i].MaxRange = MaxRange
	}

	//dealer information
	name := "deal_" + args[m+6]
	num := rand.Intn(1000)
	info[m].Name = name
	info[m].Keyid = name + strconv.Itoa(num)
	info[m].VecLength = VecLength
	info[m].M = m
	info[m].MinRange = MinRange
	info[m].MaxRange = MaxRange

	infoJSONasBytes, err := json.Marshal(info)
	if err != nil {
		return shim.Error("The marshal of info is error")
	}
	err = stub.PutState(verifier_keyid, infoJSONasBytes)
	if err != nil {
		return shim.Error("The operation of putstate is error")
	}
	err = stub.PutState("SessionNum", []byte(verifier_keyid))
	if err != nil {
		return shim.Error("The operation of putstate is error")
	}

	elapsed := time.Since(t0)
	runtime := strconv.FormatFloat(elapsed.Seconds(), 'E', -1, 64)
	return shim.Success([]byte("---Init_Setup is successful---runtime:" + runtime))
}

// Everyone can read the controltable
func (t *Arb_Verifier) Get_Setup_Info(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	info, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(info)
}

func (t *Arb_Verifier) Get_KeyID_Sess_Num(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Expecting 0")
	}
	info, err := stub.GetState("SessionNum")
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(info)
}

func (t *Arb_Verifier) Ver_Prf(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	t0 := time.Now()
	length := len(args)
	if (length != 0) && (length != 1) {
		return shim.Error("Incorrect number of arguments. Expecting 0 or 1")
	}

	creator := GetCreator(stub)
	if creator == "" {
		return shim.Error("The operator of Getcreator is failed!")
	}
	// Admin@supervision.qklfood.com   User1@supervision.qklfood.com
	// supervision_admin               supervision_user1
	// Register user need to be admitted by the fabric-ca of supervision.
	if !strings.Contains(creator, "supervision") {
		return shim.Error("The verfier is not a admin, so it can not modify the information!---creator:" + creator)
	}

	var verify_keyid string
	if length == 0 {
		result, err := stub.GetState("SessionNum")
		if err != nil {
			return shim.Error(err.Error())
		}
		verify_keyid = string(result)
	} else if length == 1 {
		verify_keyid = args[0]
	}

	var info []IdentityStore
	// get dealer info： the keyid of dealer n  m
	result, err := stub.GetState(verify_keyid)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = json.Unmarshal(result, &info)
	if err != nil {
		return shim.Error(err.Error())
	}

	dealInfo := info[len(info)-1]
	//dealname := dealInfo.Name
	dealkeyid := dealInfo.Keyid
	VecLength := dealInfo.VecLength
	m := 2 * dealInfo.M
	EC := mbp.NewECPrimeGroupKey(VecLength * m)

	// using dealname to get proof
	invoke_parm := []string{"Get_MPC_Range_Prf", dealkeyid}
	queryArgs := make([][]byte, len(invoke_parm))
	for i, arg := range invoke_parm {
		queryArgs[i] = []byte(arg)
	}
	response := stub.InvokeChaincode("cc_arb_dealer", queryArgs, "vegetablefruitchannel")
	if response.Status != shim.OK {
		errStr := fmt.Sprintf("failed to query chaincode.got error :%s", response.Payload)
		return shim.Error(errStr)
	}

	resstr := string(response.Payload)
	values := strings.Split(resstr, "---")
	if values[1] != verify_keyid {
		return shim.Error("Dealer doesn't work normally!")
	}

	var mpcrp mbp.MPCRangeProof
	err = json.Unmarshal([]byte(values[2]), &mpcrp)
	if err != nil {
		return shim.Error(err.Error())
	}

	ok := mbp.MPCVerify(&mpcrp, EC)
	elapsed := time.Since(t0)
	runtime := strconv.FormatFloat(elapsed.Seconds(), 'E', -1, 64)
	if ok {
		return shim.Success([]byte("---The verification of mpc range proof is true!!!---" + "runtime:" + runtime))
	} else {
		return shim.Error("--- MPC Range Proof FAILURE ---" + "runtime:" + runtime)
	}
}

func main() {
	err1 := shim.Start(new(Arb_Verifier))
	if err1 != nil {
		fmt.Printf("error starting simple chaincode:%s \n", err1)
	}
}
