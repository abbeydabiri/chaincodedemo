package main

import (
	"encoding/json"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//ContractObj ...
type ContractObj struct {
	Mobile    string    `json:"mobile"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Network   string    `json:"network"`
	Created   time.Time `json:"created"`
	Createdby string    `json:"createdby"`
}

//ContractHistory ...
type ContractHistory struct {
	TxID      string       `json:"txID"`
	Timestamp time.Time    `json:"timestamp"`
	Contract  *ContractObj `json:"contract"`
}

//Contract for handling writing and reading from the world state
type Contract struct {
	contractapi.Contract
}

//Put adds a new key with value to the world state
func (contract *Contract) Put(ctx contractapi.TransactionContextInterface, mobile string, name string, email string, network string) {

	if mobile == "" {
		// err = errors.New("Mobile can not be empty")
		return
	}

	if name == "" {
		// err = errors.New("Name can not be empty")
		return
	}

	obj := new(ContractObj)
	obj.Mobile = mobile
	obj.Name = name
	obj.Email = email
	obj.Network = network

	// if obj.Created, err = GetTimestamp(ctx); err != nil {
	// 	return
	// }

	// if obj.Createdby, err = GetCallerID(ctx); err != nil {
	// 	return
	// }

	key := mobile
	objBytes, _ := json.Marshal(obj)
	ctx.GetStub().PutState(key, []byte(objBytes))
	return
}

//Get retrieves the value linked to a key from the world state
func (contract *Contract) Get(ctx contractapi.TransactionContextInterface, key string) *ContractObj {

	existingObj, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil
		// return nil, err
	}

	if existingObj == nil {
		return nil
		// return nil, fmt.Errorf("Cannot read world state pair with key %s. Does not exist", key)
	}

	ContractObj := new(ContractObj)
	if err := json.Unmarshal(existingObj, ContractObj); err != nil {
		return nil
		// return nil, fmt.Errorf("Data retrieved from world state for key %s was not of type ContractObj", key)
	}
	return ContractObj
}

//History retrieves the history linked to a key from the world state
func (contract *Contract) History(ctx contractapi.TransactionContextInterface, key string) []ContractHistory {

	iter, err := ctx.GetStub().GetHistoryForKey(key)
	if err != nil {
		return nil
		// return nil, err
	}
	defer func() { _ = iter.Close() }()

	var results []ContractHistory
	for iter.HasNext() {
		state, err := iter.Next()
		if err != nil {
			return nil
			// return nil, err
		}

		entryObj := new(ContractObj)
		if errNew := json.Unmarshal(state.Value, entryObj); errNew != nil {
			return nil
			// return nil, errNew
		}

		entry := ContractHistory{
			TxID:      state.GetTxId(),
			Timestamp: time.Unix(state.GetTimestamp().GetSeconds(), 0),
			Contract:  entryObj,
		}

		results = append(results, entry)
	}
	return results
}
