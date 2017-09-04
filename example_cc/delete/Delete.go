package delete

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"github.com/example_cc/insert"
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

func DeleteTest(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1{
		return shim.Error("DeleteTest argument count isn't 1")
	}
	var err error
	var testJSON insert.Test

	var A = args[0]
	var key = "Test" + A

	testBytes, err := stub.GetState(key)
	if err != nil{
		return shim.Error(err.Error()) //Get test has error
	}else if testBytes == nil{
		return shim.Error(err.Error()) //The test don't exist
	}

	err = json.Unmarshal(testBytes, &testJSON)
	if err != nil{
		return shim.Error(err.Error()) //Test obj to JSON has error
	}

	//删除指定数据信息
	err = stub.DelState(key)
	if err != nil{
		return shim.Error(err.Error()) //Delete test has error
	}

	//删除指定数据的相关索引信息===>myIndexOfBC
	indexName := "myIndexOfBC"
	testIndexKey, err := stub.CreateCompositeKey(indexName, []string{testJSON.B, testJSON.C})
	if err != nil{
		return shim.Error(err.Error())
	}

	err = stub.DelState(testIndexKey)
	if err != nil{
		return shim.Error(err.Error()) //Delete test index has error
	}

	//删除指定数据的相关索引信息===>myIndexOfDA
	indexName2 := "myIndexOfDA"
	testIndexKey2, err := stub.CreateCompositeKey(indexName2, []string{testJSON.D, testJSON.A})
	if err != nil{
		return shim.Error(err.Error())
	}

	err = stub.DelState(testIndexKey2)
	if err != nil{
		return shim.Error(err.Error()) //Delete test index has error
	}


	return shim.Success(nil)
}

func DeleteProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args) != 1{
		return shim.Error("DeleteProduct argument count isn't 1")
	}
	var err error

	var productname = args[0]
	var key = "Product" + productname

	productBytes, err := stub.GetState(key)
	if err != nil{
		return shim.Error("Get product has error"+err.Error())
	}else if productBytes == nil{
		return shim.Error("The product don't exist"+err.Error())
	}

	//delete data
	err = stub.DelState(key)
	if err != nil{
		return shim.Error("Delete product has error"+ err.Error())
	}

	return shim.Success(nil)
}