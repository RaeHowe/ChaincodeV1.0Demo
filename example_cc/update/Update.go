package update

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"github.com/example_cc/insert"
)

//用户表结构
type user struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
}

//产品表结构
type product struct {
	Productname string `json:"productname"`
	Address 	string `json:"address"`
	Price 		string `json:"price"`
}


//["raehowe", 321] ===> 原始用户名,新用户名,密码
func UpdateUser(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args) != 2{
		return shim.Error("UpdateUser arguments count isn't 2")
	}
	var err error
	var username = args[0]
	var key = "User" + username
	var newPassword = args[1]

	userJSONBytes, err := stub.GetState(key)
	if err != nil{
		return shim.Error("Get user data has error" + err.Error())
	}else if userJSONBytes == nil{
		return shim.Error("The user data don't exist" + err.Error())
	}

	userObj := user{}
	err = json.Unmarshal(userJSONBytes, &userObj)
	if err != nil{
		return shim.Error("JSON to obj has error" + err.Error())
	}

	userObj.Password = newPassword
	newJSONBytes, err := json.Marshal(userObj)
	if err != nil{
		return shim.Error("Obj to JSON has error" + err.Error())
	}
	err = stub.PutState(key, newJSONBytes)
	if err != nil{
		return shim.Error("Save data has error" + err.Error())
	}

	return shim.Success(nil)
}

//还没有加索引的更新
func UpdateTest(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4{
		return shim.Error("UpdateTest arguments count isn't 4")
	}
	var err error
	var A = args[0]
	var B = args[1]
	var C = args[2]
	var D = args[3]

	var key = "Test" + A
	testJSONBytes, err := stub.GetState(key)
	if err != nil{
		return shim.Error(err.Error()) //Get test data has error
	}else if testJSONBytes == nil{
		return shim.Error(err.Error()) //The test data don't exist
	}

	testObj := insert.Test{}
	err = json.Unmarshal(testJSONBytes, &testObj)
	if err != nil{
		return shim.Error(err.Error()) //JSON to obj has error
	}

	testObj.B = B
	testObj.C = C
	testObj.D = D

	newJSONBytes, err := json.Marshal(testObj)
	if err != nil{
		return shim.Error(err.Error()) //Obj to JSON has error
	}
	err = stub.PutState(key, newJSONBytes)
	if err != nil{
		return shim.Error(err.Error()) //Save data has error
	}

	return shim.Success(nil)
}

func UpdateProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args) < 3{
		return shim.Error("UpdateProduct arguments count isn't 3")
	}
	var err error
	var productName = args[0]
	var key = "Product" + productName
	var address = args[1]
	var price = args[2]

	productJSONBytes, err := stub.GetState(key)
	if err != nil{
		return shim.Error("Get product data has error" + err.Error())
	}else if productJSONBytes == nil{
		return shim.Error("The product data don't exist" + err.Error())
	}

	productObj := product{}
	err = json.Unmarshal(productJSONBytes, &productObj)
	if err != nil{
		return shim.Error("JSON to obj has error" + err.Error())
	}

	productObj.Address = address
	productObj.Price = price
	newJSONBytes, err := json.Marshal(productObj)
	if err != nil{
		return shim.Error("Obj to JSON has error" + err.Error())
	}
	err = stub.PutState(key, newJSONBytes)
	if err != nil{
		return shim.Error("Save data has error" + err.Error())
	}

	return shim.Success(nil)
}
