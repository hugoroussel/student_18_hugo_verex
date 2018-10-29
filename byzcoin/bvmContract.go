package byzcoin

import (
	"errors"

	"github.com/dedis/cothority/byzcoin"
	"github.com/dedis/darc"
	"github.com/dedis/protobuf"
)

/*
spawn: create a new EVM and initialize the database structure
invoke:createAccount
invoke:sendCommand
invoke:mintCoins - send that many coins directly to the account, out of nowhere
*/

var contractBvmID = "bvm"

func contractBvm(cdb byzcoin.CollectionView, inst byzcoin.Instruction, cIn []byzcoin.Coin) (scs []byzcoin.StateChange, cOut []byzcoin.Coin, err error) {
	cOut = cIn

	err = inst.VerifyDarcSignature(cdb)
	if err != nil {
		return
	}
	var value []byte
	var darcID darc.ID
	value, _, darcID, err = cdb.GetValues(inst.InstanceID.Slice())
	if err != nil {
		return
	}

	switch inst.GetType() {
	case byzcoin.SpawnType:
		bvm := spawnEvm()
		cs := NewContractStruct(inst.Spawn.Args)
		var csBuf []byte
		csBuf, err = protobuf.Encode(&cs)
		if err != nil {
			return
		}
		scs = []byzcoin.StateChange{
			byzcoin.NewStateChange(byzcoin.Create, instID, contractBvmID, csBuf, darcID),
		}
		return

	case byzcoin.InvokeType:

		/*
			case "createAccount":
				fmt.Println("Creating account")
			case "sendCommand":
				fmt.Println("Sending command")
			case "mintCoin":
				fmt.Println("Sending command")
		*/
	}

	err = errors.New("didn't find any instructions")
	return

}

func NewContractStruct(args byzcoin.Arguments) KeyValueData {
	cs := KeyValueData{}
	for _, kv := range args {
		cs.Storage = append(cs.Storage, KeyValue{kv.Name, kv.Value})
	}
	return cs

}