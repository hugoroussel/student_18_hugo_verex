package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/vm/runtime"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

func main() {

	minimum_token := "608060405234801561001057600080fd5b506103e9806100206000396000f300608060405260043610610041576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806323b872dd14610046575b600080fd5b34801561005257600080fd5b506100b1600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506100b3565b005b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614151515610158576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260058152602001807f6572726f7200000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16141515156101fc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260058152602001807f6572726f7200000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b6000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205481111515156102b2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260058152602001807f6572726f7200000000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b806000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054036000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550806000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054016000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055505050505600a165627a7a72305820ec7a6a0fdf572ea6069e09fd8164494d64f83012242ee89746fc8c032f91ef0a0029"

	simple := "608060405234801561001057600080fd5b5060cc8061001f6000396000f300608060405260043610603f576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063943640c3146044575b600080fd5b348015604f57600080fd5b5060566098565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6000339050905600a165627a7a723058201c5103516279d76f524d63922319382752c879d703433cb1d6396d1570fe55ec0029"

	//crowdsale_ethereum := "60806040526000600760006101000a81548160ff0219169083151502179055506000600760016101000a81548160ff02191690831515021790555034801561004657600080fd5b5060405160a080610a728339810180604052810190808051906020019092919080519060200190929190805190602001909291908051906020019092919080519060200190929190505050846000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550670de0b6b3a76400008402600181905550603c83024201600381905550670de0b6b3a7640000820260048190555080600560006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505050505061091e806101546000396000f300608060405260043610610099576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806301cb3b201461027c57806329dcb0cf1461029357806338af3eed146102be5780636e66f6e91461031557806370a082311461036c5780637a3a0e84146103c35780637b3e5e7b146103ee578063a035b1fe14610419578063fd6b7ef814610444575b6000600760019054906101000a900460ff161515156100b757600080fd5b34905080600660003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254019250508190555080600260008282540192505081905550600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb336004548481151561016357fe5b046040518363ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050600060405180830381600087803b1580156101e957600080fd5b505af11580156101fd573d6000803e3d6000fd5b505050507fe842aea7a5f1b01049d752008c53c52890b1a6daf660cf39e8eec506112bbdf633826001604051808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200182151515158152602001935050505060405180910390a150005b34801561028857600080fd5b5061029161045b565b005b34801561029f57600080fd5b506102a861053b565b6040518082815260200191505060405180910390f35b3480156102ca57600080fd5b506102d3610541565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561032157600080fd5b5061032a610566565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561037857600080fd5b506103ad600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061058c565b6040518082815260200191505060405180910390f35b3480156103cf57600080fd5b506103d86105a4565b6040518082815260200191505060405180910390f35b3480156103fa57600080fd5b506104036105aa565b6040518082815260200191505060405180910390f35b34801561042557600080fd5b5061042e6105b0565b6040518082815260200191505060405180910390f35b34801561045057600080fd5b506104596105b6565b005b600354421015156105395760015460025410151561051d576001600760006101000a81548160ff0219169083151502179055507fec3f991caf7857d61663fd1bba1739e04abd4781238508cde554bb849d790c856000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600254604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390a15b6001600760016101000a81548160ff0219169083151502179055505b565b60035481565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60066020528060005260406000206000915090505481565b60015481565b60025481565b60045481565b6000600354421015156108ef57600760009054906101000a900460ff16151561076757600660003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205490506000600660003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506000811115610766573373ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f1935050505015610720577fe842aea7a5f1b01049d752008c53c52890b1a6daf660cf39e8eec506112bbdf633826000604051808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200182151515158152602001935050505060405180910390a1610765565b80600660003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055505b5b5b600760009054906101000a900460ff1680156107cf57503373ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16145b156108ee576000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc6002549081150290604051600060405180830381858888f19350505050156108d1577fe842aea7a5f1b01049d752008c53c52890b1a6daf660cf39e8eec506112bbdf66000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff166002546000604051808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200182151515158152602001935050505060405180910390a16108ed565b6000600760006101000a81548160ff0219169083151502179055505b5b5b505600a165627a7a72305820a9815687e2c10588ae8bc058d862ddec1ec17fe63396b6ccc4d265b79f16a8640029"

	public_key, _ := getKeys()

	//To call newEVM function the following parameters are needed :
	//ctx Context, statedb StateDB, chainconfig *params.ChainConfig, vmconfig Config

	//func NewDatabase(db ethdb.Database) Database
	/*
		NewDatabase creates a backing store for state. The returned database is safe for concurrent use and retains cached trie nodes in memory.
		The pool is an optional intermediate trie-node memory pool between the low level storage layer and the high level trie abstraction.
	*/

	fmt.Println("Setting up database")
	sdb, _ := getDB()

	fmt.Println("Setting up balances")
	sdb.SetBalance(common.HexToAddress(public_key), big.NewInt(1000000000000))
	fmt.Println("Sending", sdb.GetBalance(common.HexToAddress(public_key)), "eth to", common.HexToAddress(public_key).Hex())

	fmt.Println("Setting up context")
	test := func(vm.StateDB, common.Address, *big.Int) bool {
		return true
	}

	transfer := func(vm.StateDB, common.Address, common.Address, *big.Int) {
		log.Println("tried to transfer")
	}

	gethash := func(uint64) common.Hash {
		return common.HexToHash("0x0000000000000000000000000000000000000000")
	}

	//Context creation
	ctx := vm.Context{CanTransfer: test, Transfer: transfer, GetHash: gethash, Origin: common.HexToAddress(public_key), GasPrice: big.NewInt(1), Coinbase: common.HexToAddress(public_key), GasLimit: 10000000000, BlockNumber: big.NewInt(0), Time: big.NewInt(1), Difficulty: big.NewInt(1)}

	fmt.Println("===== Through runtime =====")

	fmt.Println("Creation of contract")
	create_ret, contract_addr, _, err := runtime.Create(common.Hex2Bytes(minimum_token), getConfig())
	if err != nil {
		fmt.Println("Contract deployment unsuccessful")
		fmt.Println(err)
	} else {
		fmt.Println("Successful contract deployment")
	}

	fmt.Println("Return of contract", create_ret)
	fmt.Println("Address of contract", contract_addr.Hex())

	fmt.Println("===== End runtime =====")

	fmt.Println("Setting up VMs")
	bvm := vm.NewEVM(ctx, sdb, getChainConfig(), getVMConfig())
	//func NewEnv(cfg *Config) *vm.EVM
	bvm1 := runtime.NewEnv(getConfig())

	fmt.Println("===== Through vm =====")
	//(ret []byte, contractAddr common.Address, leftOverGas uint64, err error)
	accountRef := &vm.AccountRef{} // s{common.HexToAddress(public_key)}
	ret, addrContract, leftOverGas, err := bvm1.Create(accountRef, common.Hex2Bytes(simple), 100000000, big.NewInt(0))
	if err != nil {
		fmt.Println("Contract deployment unsuccessful")
		fmt.Println(err)
	} else {
		fmt.Println("Successful contract deployment")
	}
	fmt.Println("Return of contract", common.Bytes2Hex(ret))
	fmt.Println("Left over gas : ", leftOverGas)
	fmt.Println("Contract address", addrContract.Hex())

	fmt.Println("===== End vm =====")

	bvmInterpreter := vm.NewEVMInterpreter(bvm, getVMConfig())
	fmt.Println(bvmInterpreter.CanRun(common.Hex2Bytes(minimum_token)))

	bvm1Interpreter := vm.NewEVMInterpreter(bvm1, getVMConfig())
	fmt.Println(bvm1Interpreter.CanRun(common.Hex2Bytes(minimum_token)))

	//func (in *EVMInterpreter) Run(contract *Contract, input []byte, readOnly bool) (ret []byte, err error)

	/*contract := &vm.Contract{
		CallerAddress: common.HexToAddress("0xE420b7546D387039dDaD2741a688CbEBD2578363"),
		Code:          common.Hex2Bytes(minimum_token),
		CodeHash:      common.HexToHash("0x0000000000000000000000000000000000000000"),
		//CodeAddr:      &(common.HexToAddress("0xE420b7546D387039dDaD2741a688CbEBD2578363")),
		Input: nil,
		Gas:   1000000,
	}*/

	//fmt.Println(ret)
	//fmt.Println(addrContract)

	//new_contract := vm.NewContract(accountRef, contract, big.NewInt(10000), 1)

}
