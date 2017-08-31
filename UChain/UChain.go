package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/example_cc/select"
	"github.com/example_cc/insert"
	"github.com/example_cc/update"
	"github.com/example_cc/delete"
	"encoding/json"
)

var logger = shim.NewLogger("example_cc")

type SimpleChaincode struct {

}

//在此函数中，进行数据的初始化操作，录入一些必要的初始数据 eg: key:adminRole  value {"username":"admin","password":"123"}
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response{
	_, args := stub.GetFunctionAndParameters() //获取到调用Init函数的方法名称和参数
	if len(args) != 2{
		return shim.Error("Input arguments count isn't 2"+args[0]+args[1])
	}
	var err error
	var key = "adminRole"
	var obj = make(map[string]interface{})
	obj["username"] = args[0]
	obj["password"] = args[1]

	tmpJsonArr, err := json.Marshal(obj)
	if err != nil{
		return shim.Error(err.Error())
	}

	err = stub.PutState(key, tmpJsonArr)
	if err != nil{
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	funcName, args := stub.GetFunctionAndParameters()
	switch funcName {
	case "addUser":
		return insert.AddUser(stub, args)
	case "addProduct":
		return insert.AddProduct(stub, args)
	case "addMusic":
		return insert.AddMusic(stub, args)
	case "updateUser":
		return update.UpdateUser(stub, args)
	case "updateProduct":
		return update.UpdateProduct(stub, args)
	case "deleteUser":
		return delete.DeleteUser(stub, args)
	case "deleteProduct":
		return delete.DeleteProduct(stub, args)
	case "selectUser":
		return _select.SelectUser(stub, args)
	case "selectProduct":
		return _select.SelectProduct(stub, args)
	case "selectProductByRange":
		return _select.GetProductByRange(stub, args)
	case "selectProductByAddress":
		return _select.GetProductByAddress(stub, args)
	case "selectProductByIndex":
		return _select.GetProductByIndex(stub, args)
	case "selectHistoryForProduct":
		return _select.GetHistoryForProduct(stub, args)
	case "selectMusicByIndexOfStyle":
		return _select.GetMusicByIndexOfStyle(stub, args)
	case "selectMusicByIndexOfSong":
		return _select.GetMusicByIndexOfSong(stub, args)
	default:
		return shim.Error("Can't find the function name")
	}
	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}