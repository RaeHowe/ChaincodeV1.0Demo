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
		return shim.Error("Function init input arguments count isn't 2"+args[0]+args[1])
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
		//插入操作
	case "addUser":
		return insert.AddUser(stub, args)
	case "addProduct":
		return insert.AddProduct(stub, args)
	case "addTest":
		return insert.AddTest(stub, args)
		//更新操作
	case "updateUser":
		return update.UpdateUser(stub, args)
	case "updateProduct":
		return update.UpdateProduct(stub, args)
	case "updateTest":
		return update.UpdateTest(stub, args)
		//删除操作
	case "deleteUser":
		return delete.DeleteUser(stub, args)
	case "deleteProduct":
		return delete.DeleteProduct(stub, args)
	case "deleteTest":
		return delete.DeleteTest(stub, args)
		//查询操作
	case "selectUser":
		return _select.SelectUser(stub, args)
	case "selectProduct":
		return _select.SelectProduct(stub, args)
	case "selectTest":
		return _select.SelectTest(stub, args)
	case "selectProductByRange":
		return _select.GetProductByRange(stub, args)
	case "selectProductByAddress":
		return _select.GetProductByAddress(stub, args)
	case "selectHistoryForProduct":
		return _select.GetHistoryForProduct(stub, args)
	case "getTestByIndexOfD":
		return _select.GetTestByIndexOfD(stub, args)
	case "getTestByIndexOfDA":
		return _select.GetTestByIndexOfDA(stub, args)
	case "getTestByIndexOfB":
		return _select.GetTestByIndexOfB(stub, args)
	case "getTestByIndexOfBC":
		return _select.GetTestByIndexOfBC(stub, args)
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