package delete

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
)

//产品表结构
type product struct {
	Productname string `json:"productname"`
	Address 	string `json:"address"`
	Price 		string `json:"price"`
}

func DeleteUser(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args) < 1{
		return shim.Error("DeleteUser argument count can't less 1")
	}
	var err error

	var username = args[0]
	var key = "User" + username
	err = stub.DelState(key)
	if err != nil{
		return shim.Error("Delete user has error"+ err.Error())
	}

	return shim.Success(nil)
}

func DeleteProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args) < 1{
		return shim.Error("DeleteProduct argument count can't less 1")
	}
	var err error
	var productJSON product

	var productname = args[0]
	var key = "Product" + productname

	productBytes, err := stub.GetState(key)
	if err != nil{
		return shim.Error("Get product has error"+err.Error())
	}else if productBytes == nil{
		return shim.Error("The product don't exist"+err.Error())
	}

	err = json.Unmarshal([]byte(productBytes), &productJSON)
	if err != nil{
		return shim.Error("Product obj to JSON has error")
	}

	//delete data
	err = stub.DelState(key)
	if err != nil{
		return shim.Error("Delete product has error"+ err.Error())
	}

	//delete index
	indexname := "address~name"
	addressNameIndexKey, err := stub.CreateCompositeKey(indexname, []string{productJSON.Address})
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.DelState(addressNameIndexKey)
	if err != nil{
		return shim.Error("Delete product index has error"+ err.Error())
	}

	return shim.Success(nil)
}