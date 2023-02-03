## Usage

### Quick Start
```shell
# step 1, start http server
# start http server for metamask-tool-ui
go run main.go

# step 2, sign and generate tx
# open chrome, visit http://localhost:8080
# import wallet, private key = 0xdfcb02b38ac1bc221b51cb4bec373236ae673f5524d030cef4551dbd58bb0d25
# connect wallet 
# set bnb amount, click  `Send BNB on GNFD` Button
# copy the generated `txRawBytesHex` value

# step 3, send tx to local gnfd chain
# prepare for js tool to verify
npm install ts-node -g
cd gnfd-js-tool && npm install
ts-node src/gnfd-send-tx.ts <txRawBytesHex copied from step2> 


```
