package insert

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
)

//用户表结构
type user struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
}

//产品表结构，主要用于测试范围查询和历史数据查询
type product struct {
	ProductName string `json:"productName"`
	Address 	string `json:"address"`
	Price 		string `json:"price"`
}

//测试表结构，主要用于查询索引查询
type Test struct {
	A string `json:"a"`
	B string `json:"b"`
	C string `json:"c"`
	D string `json:"d"`
} 


//["user1",123] ===> 用户名,密码
func AddUser(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args) != 2{
		return shim.Error("Function \"AddUser\" input arguments count isn't 2")
	}
	var err error

	var username = args[0]
	var password = args[1]

	var key = "User"+username

	//验证数据是否已经存在
	resByteArr ,err := stub.GetState(key)
	if err != nil{
		return shim.Error(err.Error()) //Failed to get data
	}else if resByteArr != nil{
		return shim.Error("This data already exists:" + key)
	}

	user := &user{username, password}
	userJSONBytes, err := json.Marshal(user)
	if err != nil{
		return shim.Error(err.Error()) //User obj to byte arr has error
	}

	err = stub.PutState(key, userJSONBytes)
	if err != nil{
		return shim.Error( err.Error()) //Add user has error
	}

	return shim.Success(nil)
}

// ["1","2","1","1"] ===> A,B,C,D
func AddTest(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4{
		return shim.Error("Input argument count isn't 4")
	}
	var err error

	var A = args[0]
	var B = args[1]
	var C = args[2]
	var D = args[3]

	var key = "Test" + A

	//验证数据唯一性
	resByteArr, err := stub.GetState(key)
	if err != nil{
		return shim.Error(err.Error())
	}else if resByteArr != nil{
		return shim.Error("This data already exists:" + key)
	}

	test := &Test{A,B,C, D}
	testJSONBytes, err := json.Marshal(test)
	if err != nil{
		return shim.Error(err.Error()) //Test obj to byte arr has error
	}

	err = stub.PutState(key, testJSONBytes)
	if err != nil{
		return shim.Error(err.Error()) //Add test has error
	}

	//#############以对象中的B、C作为compositeKey进行查询#################
	indexName := "myIndexOfBC"
	testIndexKey, err := stub.CreateCompositeKey(indexName, []string{test.B, test.C})
	if err != nil{
		return shim.Error(err.Error())
	}

	//验证compositeKey的唯一性
	var message = ""
	resKeyArr, err := stub.GetState(testIndexKey)
	if err != nil{
		return shim.Error(err.Error())
	}

	if resKeyArr != nil{ //证明该索引已经存在了
			message = "Add object success, but add this index has error, because this index already exist"
	}else { //保存该索引信息
		value := []byte{0x00}
		err = stub.PutState(testIndexKey, value)
		if err != nil{
			return shim.Error(err.Error())
		}
	}
	//###############以对象中的D、A作为compositeKey进行查询###############
	indexName2 := "myIndexOfDA"
	testIndexKey2, err := stub.CreateCompositeKey(indexName2, []string{test.D, test.A})
	if err != nil{
		return shim.Error(err.Error())
	}

	//验证compositeKey的唯一性
	resKeyArr2, err := stub.GetState(testIndexKey2)
	if err != nil{
		return shim.Error(err.Error())
	}

	if resKeyArr2 != nil{
		message = "Add object success, but add this index has error, because this index already exist"
	}else {
		value2 := []byte{0x00}
		err = stub.PutState(testIndexKey2, value2)
		if err != nil{
			return shim.Error(err.Error())
		}
	}

	return shim.Success([]byte(message))
}

//["苹果","天津",3] ===> 名称，产地，斤/元
func AddProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args) != 3{
		return shim.Error("Function \"AddProduct\" input arguments count isn't 3")
	}
	var err error

	var productName = args[0]
	var address = args[1]
	var price = args[2]

	var key = "Product"+productName

	//验证数据的唯一性
	resByteArr, err := stub.GetState(key)
	if err != nil{
		return shim.Error(err.Error())
	}else if resByteArr != nil{
		return shim.Error("This data already exists:" + key)
	}

	product := &product{productName, address, price}
	productJSONBytes, err := json.Marshal(product)
	if err != nil{
		return shim.Error(err.Error()) //Product obj to byte arr has error
	}

	err = stub.PutState(key, productJSONBytes)
	if err != nil{
		return shim.Error(err.Error()) //Add product has error
	}

	return shim.Success(nil)
}