package main

import (
	"os"
	"testing"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func  TestContract(t *testing.T) {
	os.Setenv("MODE","TEST")
	
	assert := assert.New(t)
	uid := uuid.New().String()

	cc, err := contractapi.NewChaincode(new(Contract))
	assert.Nil(err, "error should be nil")

	stub := shimtest.NewMockStub("TestStub", cc)
	assert.NotNil(stub, "Stub is nil, TestStub creation failed")

	// - - - test Contract:Put function - - - 
	putResp := stub.MockInvoke(uid,[][]byte{
		[]byte("Contract:Put"),
		[]byte("080"),
		[]byte("Fabric"),
		[]byte("fabric@hyperledger"),
		[]byte("INTERNET"),
	})
	assert.EqualValues(OK, putResp.GetStatus(), putResp.GetMessage())
	

	// - - - test Contract:Get function - - - 
	testKey := "080"
	getResp := stub.MockInvoke(uid, [][]byte{
		[]byte("Contract:Get"),
		[]byte(testKey),
	})
	assert.EqualValues(OK, getResp.GetStatus(), getResp.GetMessage())
	assert.NotNil(getResp.Payload, "getResp.Payload should not be nil")
	
	contract := new(ContractObj)
	err = json.Unmarshal(getResp.Payload, contract)
	assert.Nil(err, "json.Unmarshal error should be nil")
	assert.NotNil(contract, "contract should not be nil")

	retrievedKey := contract.Mobile
	assert.EqualValues(testKey, retrievedKey, "testKey and retrievedKey mismatch")
}