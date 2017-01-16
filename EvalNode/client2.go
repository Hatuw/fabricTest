package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/container"
	//"path/filepath"
	//	"github.com/hyperledger/fabric/core"
	//	"github.com/hyperledger/fabric/core/peer"
	//"github.com/spf13/viper"
	//"io/ioutil"
	"math/rand"
	"os"
	//"strings"
	"time"
	//"errors"
	"golang.org/x/net/context"
	//"github.com/golang/protobuf/jsonpb"
	"github.com/hyperledger/fabric/core/crypto"
	"github.com/hyperledger/fabric/core/crypto/primitives"
	pb "github.com/hyperledger/fabric/protos"
)

const (
	localStore string = "/var/hyperledger/production/client/"
)

type Transaction struct {
	Jsonrpc string `json:"jsonrpc,omitempty"`
	Method  string `json:"method,omitempty"`
	Params  params `json:"params,omitempty"`
	ID      int    `json:"id,omitempty"`
}

type params struct {
	Type          int               `json:"type,omitempty"`
	ChaincodeID   map[string]string `json:"chaincodeID,omitempty"`
	CtorMsg       ctorMsg           `json:"ctorMsg"`
	SecureContext string            `josn:"secureContext,omitempty"`
}

type ctorMsg struct {
	Function string   `json:"function,omitempty"`
	Args     []string `json:"args,omitempty"`
}

func RandomId() int {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	return r.Intn(1000000)
}

// this is for normal resp trasnaction upon http
func MakeATransaction() (*bytes.Buffer, error) {
	t := &Transaction{
		"2.0",
		"deploy",
		params{
			1,
			map[string]string{"path": "github.com/hyperledger/fabric/examples/chaincode/go/HelloWorld"},
			ctorMsg{"init", []string{"Hello,World"}},
			"diego"},
		RandomId(),
	}
	b, err := json.Marshal(t)
	if err != nil {
		fmt.Printf("error raised: %v\n", err)
		return nil, err
	}
	return bytes.NewBuffer(b), nil
}

// this is only for pb.chaincodSpec
func MakeAChaincodeSpec() (*pb.ChaincodeSpec, error) {
	var spec pb.ChaincodeSpec
	//var spec2 pb.ChaincodeSpec
	/*
		function := "init"
		args := []string{"Hello,Pig", "stupid"}
		buff1, buff2 := &bytes.Buffer{}, &bytes.Buffer{}
		gob.NewEncoder(buff1).Encode(function)
		gob.NewEncoder(buff2).Encode(args)
		byte1, byte2 := buff1.Bytes(), buff2.Bytes()
		fmt.Println(byte1, byte2)
		c1 := &pb.ChaincodeInput{Args: [][]byte{byte1, byte2}}
		fmt.Println(c1)

		spec.Type = 1 //type
		spec.ChaincodeID = &pb.ChaincodeID{Path: "github.com/hyperledger/fabric/examples/chaincode/go/HelloWorld"}
		spec.CtorMsg = c1
	*/

	t := &params{
		1,
		map[string]string{"path": "github.com/hyperledger/fabric/examples/chaincode/go/HelloWorld"},
		ctorMsg{"init", []string{"Hello, Pig"}},
		"diego"}

	b, err := json.Marshal(t)
	if err != nil {
		fmt.Printf("Error raised: %v", err)
		return nil, err
	}

	tmp, err := ioutil.ReadAll(bytes.NewBuffer(b))
	if err != nil {
		fmt.Println("Read error: %v", err)
		return nil, err
	}
	//fmt.Println(b, bytes.NewBuffer(b))
	err = json.Unmarshal(tmp, &spec2)
	if err != nil {
		fmt.Printf("pb unmarshal error: %v", err)
		os.Exit(0)
	}

	//fmt.Println(spec2)
	return &spec, nil
}

func getChaincodeBytes(context context.Context, spec *pb.ChaincodeSpec) (*pb.ChaincodeDeploymentSpec, error) {
	var codePackageBytes []byte
	var err error
	codePackageBytes, err = container.GetChaincodePackageBytes(spec)
	if err != nil {
		return nil, err
	}
	chaincodeDeploymentSpec := &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec, CodePackage: codePackageBytes}
	return chaincodeDeploymentSpec, nil
}

func getMetadata(chaincodeSpec *pb.ChaincodeSpec) ([]byte, error) {
	return chaincodeSpec.Metadata, nil
}

func CreateDeployTx(chaincodeDeploymentSpec *pb.ChaincodeDeploymentSpec, uuid string, nonce []byte, attrs ...string) (*pb.Transaction, error) {
	tx, err := pb.NewChaincodeDeployTransaction(chaincodeDeploymentSpec, uuid)
	if err != nil {
		return nil, err
	}
	tx.Metadata, err = getMetadata(chaincodeDeploymentSpec.GetChaincodeSpec())
	if err != nil {
		return nil, err
	}

	if nonce == nil {
		tx.Nonce, err = primitives.GetRandomNonce()
		if err != nil {
			return nil, err
		}
	} else {
		tx.Nonce = nonce
	}

	//handle confidentiality
	fmt.Println(chaincodeDeploymentSpec.ChaincodeSpec.ConfidentialityLevel)
	return tx, nil

}

func Deploy(ctx context.Context, spec *pb.ChaincodeSpec) (*string, error) {
	chaincodeDeploymentSpec, err := getChaincodeBytes(ctx, spec)
	if err != nil {
		return nil, err
	}

	transID := chaincodeDeploymentSpec.ChaincodeSpec.ChaincodeID.Name

	var tx *pb.Transaction
	var sec crypto.Client

	sec, err = crypto.InitClient(spec.SecureContext, nil)
	defer crypto.CloseClient(sec)
	spec.SecureContext = ""

	if err != nil {
		return nil, err
	}

	tx, err = sec.NewChaincodeDeployTransaction(chaincodeDeploymentSpec, transID, spec.Attributes...)
	if err != nil {
		return nil, err
	}

	fmt.Println(tx)
	return &transID, err

}

func main() {
	//configuration
	//for viper testing
	/*
		viper.SetEnvPrefix("core")
		viper.AutomaticEnv()
		replacer := strings.NewReplacer(".", "_")
		viper.SetEnvKeyReplacer(replacer)
		viper.SetConfigName("core")
		viper.AddConfigPath("/opt/gopath/src/github.com/hyperledger/fabric/peer/")
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("error raise: %v", err))
		}
		//viper.Set("peer.fileSystemPath", filepath.Join("/", "var", "hyperledger", "production"))
		err = core.CacheConfiguration()
		if err != nil {
			panic(fmt.Errorf("error raise: %v", err))
		}
		//viper.AddConfigPath("/hyperledger/peer/")
		//viper.AddConfigPath("/opt/gopath/src/github.com/hyperledger/fabric/peer/")
		//fmt.Println(viper.GetString("peer.fileSystemPath"))
		//fmt.Println(viper.GetString("peer.gomaxprocs"))
		fmt.Println(viper.GetBool("security.enabled"))
		fmt.Println(peer.SecurityEnabled())
		//	fmt.Println(string(viper.GetString("peer.validator.consensus.plugin")))

		//fmt.Println(viper.GetString("chaincode.mode") == chaincode.DevModeUserRunsChaincode)
		//fmt.Println(viper.GetBool("security.privacy"))
		//fmt.Println(viper.GetBool("security.enabled"))

		//define the devop server
		//var serverDevops pb.DevopsServer
		//serverDevops = //use underlying *core.Devops
		//var spec pb.ChaincodeSpec
		//	t, err := MakeATransaction()
		//transId := new(string)
	*/
	spec, err := MakeAChaincodeSpec()
	if err != nil {
		os.Exit(0)
	}
	fmt.Println(spec)

	/*
		chaincodeDeploymentSpec, err := getChaincodeBytes(context.Background(), spec)
		if err != nil {
			os.Exit(0)
		}

		tx, err := CreateDeployTx(chaincodeDeploymentSpec, chaincodeDeploymentSpec.ChaincodeSpec.ChaincodeID.Name, []byte{}, spec.Attributes...)
		if err != nil {
			os.Exit(0)
		}
		fmt.Println(tx.Timestamp)
	*/
	/*
		transId, err = Deploy(context.Background(), spec)

		if err != nil {
			fmt.Printf("Error raised: %v", err)
			os.Exit(0)
		}

		fmt.Println(transId)
	*/
	/*
		chaincodeDeploymentSpec, err := getChaincodeBytes(context.Background(), spec)
		if err != nil {
			fmt.Printf("error raised: %v", err)
		}
	*/
	//fmt.Println(chaincodeDeploymentSpec.ChaincodeSpec.ChaincodeID.Name)
}
