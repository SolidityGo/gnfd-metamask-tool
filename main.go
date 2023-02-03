package main

import (
	"encoding/json"
	"fmt"
	client "github.com/bnb-chain/gnfd-go-sdk/client/grpc"
	"github.com/bnb-chain/gnfd-go-sdk/keys"
	"github.com/bnb-chain/gnfd-go-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"log"
	"net/http"
)

type TxInfo struct {
	Amount    int64  `json:"amount"`
	Signature string `json:"signature"`
}

type Account struct {
	Address string `json:"address"`
}

type AccountInfo struct {
	AccountNumber uint64 `json:"accountNumber"`
	Sequence      uint64 `json:"sequence"`
}

var gnfdCli client.GreenfieldClient

var (
	GrpcConn = "localhost:9090"
	ChainId  = "greenfield_9000-121"
)

func main() {
	km, _ := keys.NewPrivateKeyManager("dfcb02b38ac1bc221b51cb4bec373236ae673f5524d030cef4551dbd58bb0d25")
	gnfdCli = client.NewGreenfieldClientWithKeyManager(GrpcConn, ChainId, km)

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./gnfd-metamask-ui"))
	mux.Handle("/", http.StripPrefix("/", fs))
	mux.HandleFunc("/tx", txHandler)
	mux.HandleFunc("/account", accountHandler)
	http.ListenAndServe(":8080", mux)
}

func txHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("--------txHandler")

	decoder := json.NewDecoder(req.Body)
	var t TxInfo
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	log.Println("sig from js", t.Signature)

	sendTokenReq := types.SendTokenRequest{
		Token:     "bnb",
		Amount:    t.Amount,
		ToAddress: "0x0000000000000000000000000000000000000001",
	}

	txBytes, err := gnfdCli.GenerateTx(t.Signature, sendTokenReq)
	if err != nil {
		panic(err)
	}

	txHex := hexutil.Encode(txBytes)

	log.Println("sent tx raw bytes to frontend", txHex)

	fmt.Fprintf(w, txHex)
}

func accountHandler(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var acc Account
	err := decoder.Decode(&acc)
	if err != nil {
		panic(err)
	}
	acct, err := gnfdCli.Account(acc.Address)
	if err != nil {
		panic(err)
	}
	accountNum := acct.GetAccountNumber()
	accountSeq := acct.GetSequence()

	accInfo := AccountInfo{
		AccountNumber: accountNum,
		Sequence:      accountSeq,
	}

	result, err := json.Marshal(accInfo)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, string(result))
}
