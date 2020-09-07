package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"gopkg.in/alecthomas/kingpin.v2"

	// "github.com/ethereum/go-ethereum/accounts/abi"
	// "github.com/hayeah/sol/jsonabi"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hayeah/sol/etherscan"

	"github.com/pkg/errors"
	"github.com/qtumproject/solar/abi"
	"github.com/qtumproject/solar/jsonabi"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
)

var toAddress = kingpin.Arg("toAddress", "contract address to call").Required().String()
var methodName = kingpin.Arg("methodName", "contract method to call").Required().String()
var methodArgs = kingpin.Arg("args", "method args as JSON array").Default("").String()

var fromAddress = kingpin.Flag("fromAddress", "method args as JSON array").Short('f').Default("0x0000000000000000000000000000000000000000").String()

var etherScanAPIKey = kingpin.Flag("etherscanKey", "etherscan API key").Envar("ETHERSCAN_KEY").Default("6YSRMSWDJX6QTH9W76ASK2I2N17XHPG648").String()

// I have a feeling that consensys is going to come knocking
var ethRPCURL = kingpin.Flag("ethRPCURL", "ETH JSON-RPC url").Envar("ETHRPC_URL").Default("https://mainnet.infura.io/v3/0834aab6765f4fcfbb7254e18e1cc212").String()

func run() error {
	ethRPCURL := *ethRPCURL
	// tether: https://etherscan.io/token/0xdac17f958d2ee523a2206206994597c13d831ec7

	// https://api.etherscan.io/api?module=contract&action=getabi&address=0xBB9bc244D798123fDe783fCc1C72d3Bb8C189413&apikey=6YSRMSWDJX6QTH9W76ASK2I2N17XHPG648

	escan := etherscan.APIClient{Key: *etherScanAPIKey}
	abistr, err := escan.GetABI(*toAddress)
	if err != nil {
		return err
	}

	// abidata, err := ioutil.ReadFile("tether.abi.json")
	// if err != nil {
	// 	return err
	// }
	// abistr := string(abidata)

	eabi, err := abi.JSON(bytes.NewReader([]byte(abistr)))
	if err != nil {
		return err
	}

	oeabi, err := ethabi.JSON(bytes.NewReader([]byte(abistr)))
	if err != nil {
		return err
	}

	// calldata, err := eabi.Pack("name")
	// if err != nil {
	// 	return err
	// }

	fromAddr := *fromAddress
	toAddr := *toAddress

	methodName := *methodName
	methodArgs := *methodArgs

	method, found := eabi.Methods[methodName]
	if !found {
		return errors.Errorf("method not found: %s", method)
	}

	var argsData []byte

	if methodArgs != "" {
		// jsonInputs := []byte(methodArgs)
		argsData, err = jsonabi.EncodeJSONValues(method.Inputs, []byte(methodArgs))
		if err != nil {
			return errors.Wrap(err, "encode JSON using ABI")
		}
	}

	calldata := append(method.Id(), argsData...)

	// fmt.Printf("calldata %x\n", calldata)

	ethc, err := ethclient.Dial(ethRPCURL)
	if err != nil {
		return err
	}

	to := common.HexToAddress(toAddr)

	returndata, err := ethc.CallContract(context.Background(), ethereum.CallMsg{
		From: common.HexToAddress(fromAddr),
		To:   &to,
		Data: calldata,
	}, nil) // default to latest block

	if err != nil {
		return err
	}

	// for method.Outputs
	var outputs interface{}

	// fmt.Printf("returndata %x\n", returndata)

	oeabi.Unpack(&outputs, methodName, returndata)

	// switch v := outputs.(type) {
	// case *big.Int:
	// 	// big.NewRat(a int64, b int64)
	// 	bf := new(big.Float)
	// 	bf.SetInt(v)
	// 	bf.Quo(bf, big.NewFloat(1e6))
	// 	bf.SetPrec(10)

	// 	outputs = bf
	// }

	// fmt.Println(reflect.TypeOf(outputs).String(), outputs)
	fmt.Println(outputs)

	return nil
}

func main() {
	kingpin.Parse()

	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}
