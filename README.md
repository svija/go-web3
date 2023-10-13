<h1 >go-web3 for Falcon</h1>


<h1 >Falcon Func:</h1>

```
func F_bool_to_str(in bool) string{
    if in {
        return "true"
    }
    return "false"
}

func F_int(str string)(rst int){
    fl,_:=strconv.Atoi(str)
    return fl
}

func F_int_to_str(in int)(rst string){
    return strconv.Itoa(in)
}

func F_float(str string)(rst float64){
    fl,_:=strconv.ParseFloat(str, 64)
    return fl
}

func F_float_to_str(str float64)(rst string){
    return strconv.FormatFloat(str, 'f', 12, 64)
}


func F_float_to_str8(str float64)(rst string){
    return strconv.FormatFloat(str, 'f', 8, 64)
}

func F_uint64_to_str(in uint64)(rst string){
    rst=strconv.FormatUint(in, 10)
    return rst
}

func F_int64_to_str(in int64)(rst string){
    rst=strconv.FormatInt(in, 10)
    return rst
}
```


default fork from:@chenzhijie
## API
- [NewWeb3()](#NewWeb3)
- [SetChainId(chainId int64)](#setchainidchainid-int64)
- [SetAccount(privateKey string) error](#setaccountprivatekey-string-error)
- [GetBlockNumber()](#GetBlockNumber)
- [GetNonce(addr common.Address, blockNumber *big.Int) (uint64, error)](#getnonceaddr-commonaddress-blocknumber-bigint-uint64-error)
- [NewContract(abiString string, contractAddr ...string) (*Contract, error)](#newcontractabistring-string-contractaddr-string-contract-error)
- [Call(methodName string, args ...interface{}) (interface{}, error)](#callmethodname-string-args-interface-interface-error)
- [EncodeABI(methodName string, args ...interface{}) ([]byte, error)](#encodeabimethodname-string-args-interface-byte-error)
- [SendRawTransaction(to common.Address,amount *big.Int,gasLimit uint64,gasPrice *big.Int,data []byte) (common.Hash, error) ](#sendrawtransactionto-commonaddressamount-bigintgaslimit-uint64gasprice-bigintdata-byte-commonhash-error)

### NewWeb3()

Creates a new web3 instance with http provider.

```golang
// change to your rpc provider
var rpcProviderURL = "https://rpc.flashbots.net"
web3, err := web3.NewWeb3(rpcProviderURL)
if err != nil {
    panic(err)
}
```


### GetBlockNumber()

Get current block number.

```golang
blockNumber, err := web3.Eth.GetBlockNumber()
if err != nil {
    panic(err)
}
fmt.Println("Current block number: ", blockNumber)
// => Current block number:  11997285
```


### SetChainId(chainId int64)

Setup chainId for different network.

```golang
web3.Eth.SetChainId(1)
```


### SetAccount(privateKey string) error

Setup default account with privateKey (hex format)

```golang
pv, err := crypto.GenerateKey()
if err != nil {
    panic(err)
}
privateKey := hex.EncodeToString(crypto.FromECDSA(pv))
err := web3.Eth.SetAccount(privateKey)
if err != nil {
    panic(err)
}
```


### GetNonce(addr common.Address, blockNumber *big.Int) (uint64, error)

Get transaction nonce for address

```golang
nonce, err := web3.Eth.GetNonce(web3.Eth.Address(), nil)
if err != nil {
    panic(err)
}
fmt.Println("Latest nonce: ", nonce)
// => Latest nonce: 1 
```

### NewContract(abiString string, contractAddr ...string) (*Contract, error)

Init contract api

```golang
abiString := `[
	{
		"constant": true,
		"inputs": [],
		"name": "totalSupply",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	}
]`
contractAddr := "0x6B175474E89094C44Da98b954EedeAC495271d0F" // contract address
contract, err := web3.Eth.NewContract(abiString, contractAddr)
if err != nil {
    panic(err)
}
```

### Call(methodName string, args ...interface{}) (interface{}, error)

Contract call method

```golang

totalSupply, err := contract.Call("totalSupply")
if err != nil {
    panic(err)
}
fmt.Printf("Total supply %v\n", totalSupply)

// => Total supply  10000000000
```

### EncodeABI(methodName string, args ...interface{}) ([]byte, error)

EncodeABI data

```golang

data, err := contract.EncodeABI("balanceOf", web3.Eth.Address())
if err != nil {
    panic(err)
}
fmt.Printf("Data %x\n", data)

// => Data 70a08231000000000000000000000000c13a163dd812ed7eb8bb9152651054eae5ee0999 
```

### SendRawTransaction(to common.Address,amount *big.Int,gasLimit uint64,gasPrice *big.Int,data []byte) (common.Hash, error) 

Send transaction

```golang

txHash, err := web3.Eth.SendRawTransaction(
    common.HexToAddress(tokenAddr),
    big.NewInt(0),
    gasLimit,
    web3.Utils.ToGWei(1),
    approveInputData,
)
if err != nil {
    panic(err)
}
fmt.Printf("Send approve tx hash %v\n", txHash)

// => Send approve tx hash  0x837136c8b6f34b519c049d1cf703d3bba47d32f6801c25d83d0113bdc0e6936a 
```

## Examples

- **[Chain API](./examples/chain/chain.go)**
- **[Contract API](./examples/contract/erc20.go)**
- **[EIP1559 API](./examples/eip1559/main.go)**
