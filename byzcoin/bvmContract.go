package byzcoin

import (
	"errors"
	"github.com/dedis/cothority/byzcoin"
	"github.com/dedis/cothority/darc"
	"github.com/dedis/onet/log"
	"github.com/dedis/protobuf"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)


var ContractBvmID = "bvm"
var nilAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")

type contractBvm struct {
	byzcoin.BasicContract
	ES
}

func contractBvmFromBytes(in []byte) (byzcoin.Contract, error) {
	cv := &contractBvm{}
	err := protobuf.Decode(in, &cv.ES)
	if err != nil {
		return nil, err
	}
	return cv, nil
}

//Spawn deploys an EVM
func (c *contractBvm) Spawn(rst byzcoin.ReadOnlyStateTrie, inst byzcoin.Instruction, coins []byzcoin.Coin) (sc []byzcoin.StateChange, cout []byzcoin.Coin, err error) {
	cout = coins
	es := c.ES
	memdb, db, _, err := spawnEvm()
	if err != nil{
		return nil, nil, err
	}
	es.RootHash, err = db.Commit(true)
	if err != nil {
		return nil, nil, err
	}
	err = db.Database().TrieDB().Commit(es.RootHash, true)
	if err != nil {
		return nil, nil, err
	}
	es.DbBuf, err = memdb.Dump()
	esBuf, err := protobuf.Encode(&es)
	// Then create a StateChange request with the data of the instance. The
	// InstanceID is given by the DeriveID method of the instruction that allows
	// to create multiple instanceIDs out of a given instruction in a pseudo-
	// random way that will be the same for all nodes.
	sc = []byzcoin.StateChange{
		byzcoin.NewStateChange(byzcoin.Create, inst.DeriveID(""), ContractBvmID, esBuf, darc.ID(inst.InstanceID.Slice())),
	}
	return
}

//Invoke provides three instructions : display, credit and transaction
func (c *contractBvm) Invoke(rst byzcoin.ReadOnlyStateTrie, inst byzcoin.Instruction, coins []byzcoin.Coin) (sc []byzcoin.StateChange, cout []byzcoin.Coin, err error) {
	cout = coins
	var darcID darc.ID
	_, _, _, darcID, err = rst.GetValues(inst.InstanceID.Slice())
	if err != nil {
		return
	}
	es := c.ES
	switch inst.Invoke.Command {

	case "display":
		addressBuf := inst.Invoke.Args.Search("address")
		if addressBuf == nil {
			return nil, nil, errors.New("no address provided")
		}
		address := common.HexToAddress(string(addressBuf))
		_, db, err := getDB(es)
		if err !=nil {
			return nil, nil, err
		}
		ret := db.GetBalance(address)
		if ret == big.NewInt(0) {
			log.LLvl1(address.Hex(), "balance empty")
		}
		log.LLvl1( address.Hex(), "balance", ret ,"wei")
		return nil, nil, nil

	case "credit":
		addressBuf := inst.Invoke.Args.Search("address")
		if addressBuf == nil {
			return nil, nil, errors.New("no address provided")
		}
		address := common.HexToAddress(string(addressBuf))
		memdb, db, err := getDB(es)
		if err != nil {
			return nil, nil, err
		}
		//By default credit, credits 5*1e18 wei. To change this, add a new parameter to the byzcoin transaction with the desired value
		db.SetBalance(address, big.NewInt(1e18*5))
		log.LLvl1(address.Hex(), "credited 5 eth")

		//Commits the general stateDb
		es.RootHash, err = db.Commit(true)
		if err != nil {
			return nil, nil ,err
		}

		//Commits the low level trieDB
		err = db.Database().TrieDB().Commit(es.RootHash, true)
		if err != nil {
			return nil, nil, err
		}

		//Saves the general Ethereum State
		es.DbBuf, err = memdb.Dump()
		if err != nil {
			return nil, nil, err
		}

		//Save the Ethereum structure
		esBuf, err := protobuf.Encode(&es)
		if err != nil {
			return nil, nil , err
		}
		sc = []byzcoin.StateChange{
			byzcoin.NewStateChange(byzcoin.Update, inst.InstanceID,
				ContractBvmID, esBuf, darcID),
		}

	case "transaction":
		memdb, db, err := getDB(es)
		if err != nil{
			return nil, nil, err
		}
		txBuffer := inst.Invoke.Args.Search("tx")
		if txBuffer == nil {
			log.LLvl1("no transaction provided in byzcoin transaction")
			return nil, nil, err
		}
		var ethTx types.Transaction
		err = ethTx.UnmarshalJSON(txBuffer)
		if err != nil {
			return nil, nil, err
		}
		transactionReceipt, err := sendTx(&ethTx, db)
		if err != nil {
			log.ErrFatal(err)
			return nil, nil, err
		}

		if transactionReceipt.ContractAddress.Hex() != nilAddress.Hex() {
			log.LLvl1("contract deployed at:", transactionReceipt.ContractAddress.Hex(), "tx status:", transactionReceipt.Status, "(0/1 fail/success)", "gas used:", transactionReceipt.GasUsed, "tx receipt:", transactionReceipt.TxHash.Hex())
		} else {
			log.LLvl1("tx status:", transactionReceipt.Status, "(0/1 fail/success)", "gas used:", transactionReceipt.GasUsed, "tx receipt:", transactionReceipt.TxHash.Hex())
		}


		//Commits the general stateDb
		es.RootHash, err = db.Commit(true)
		if err != nil {
			return nil, nil, err
		}

		//Commits the low level trieDB
		err = db.Database().TrieDB().Commit(es.RootHash, true)
		if err != nil {
			return nil, nil, err
		}

		//Saves the general Ethereum State
		es.DbBuf, err = memdb.Dump()
		if err != nil {
			return nil, nil, err
		}

		//Save the Ethereum structure
		esBuf, err := protobuf.Encode(&es)
		if err != nil {
			return nil, nil , err
		}
		sc = []byzcoin.StateChange{
			byzcoin.NewStateChange(byzcoin.Update, inst.InstanceID,
				ContractBvmID, esBuf, darcID),
		}
	default :
		err = errors.New("Contract can only display, credit and receive transactions")
		return

	}
	return
}

//sendTx is a helper function that applies the signed transaction to the EVM
func sendTx(tx *types.Transaction, db *state.StateDB) (*types.Receipt, error){

	//get parameters defined in params
	chainconfig := getChainConfig()
	config := getVMConfig()

	// GasPool tracks the amount of gas available during execution of the transactions in a block.
	gp := new(core.GasPool).AddGas(uint64(1e18))
	usedGas := uint64(0)
	ug := &usedGas

	// ChainContext supports retrieving headers and consensus parameters from the
	// current blockchain to be used during transaction processing.
	var bc core.ChainContext
	// Header represents a block header in the Ethereum blockchain.
	var header  *types.Header
	header = &types.Header{
		Number: big.NewInt(0),
		Difficulty: big.NewInt(0),
		ParentHash: common.Hash{0},
		Time: big.NewInt(0),
	}

	receipt, usedGas, err := core.ApplyTransaction(chainconfig, bc, &nilAddress, gp, db, header, tx, ug, config)
	if err !=nil {
		log.Error()
		return nil, err
	}
	return receipt, nil
}



type ES struct {
	DbBuf []byte
	RootHash common.Hash
}

