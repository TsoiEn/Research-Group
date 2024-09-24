package credentials

import (
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Student struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Age        int    `json:"age"`
	Birthday   string `json:"birthday"`
	Credential string `json:"credential"`
}

type chaincode interface {
	Init(stub ChaincodeStubInterface) pb.Response
	Invoke(stub ChaincodeStubInterface) pb.Response
}

type ChaincodeStubInterface interface {
	InvokeChaincode(chaincodeName string, args [][]byte, channel string) pb.Response
	GetState(key string) ([]byte, error)
	PutState(key string, value []byte) error
	DelState(key string) error
	GetStateByRange(startKey, endKey string) (StateQueryIteratorInterface, error)
	GetTxTimestamp() (*timestamp.Timestamp, error)
	GetTxID() string
	GetChannelID() string
}
