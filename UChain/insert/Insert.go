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

//产品表结构
type product struct {
	Productname string `json:"productname"`
	Address 	string `json:"address"`
	Price 		string `json:"price"`
}

//["user1",123] ===> 用户名,密码
func AddUser(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args) != 2{
		return shim.Error("Input arguments count isn't 2")
	}
	var err error

	var username = args[0]
	var password = args[1]

	var key = "User"+username

	//验证数据是否已经存在
	resByteArr ,err := stub.GetState(key)
	if err != nil{
		return shim.Error("Failed to get data: " +err.Error())
	}else if resByteArr != nil{
		return shim.Error("This data already exists: " + key)
	}


	user := &user{username, password}
	userJSONBytes, err := json.Marshal(user)
	if err != nil{
		return shim.Error("User obj to byte arr has error" + err.Error())
	}

	err = stub.PutState(key, userJSONBytes)
	if err != nil{
		return shim.Error("Add user has error" + err.Error())
	}

	return shim.Success(nil)
}

//["苹果","天津",3] ===> 名称，产地，斤/元
func AddProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args) != 3{
		return shim.Error("Input arguments count isn't 3")
	}
	var err error

	var productname = args[0]
	var address = args[1]
	var price = args[2]

	var key = "Product"+productname

	product := &product{productname, address, price}
	productJSONBytes, err := json.Marshal(product)
	if err != nil{
		return shim.Error("Product obj to byte arr has error" + err.Error())
	}

	err = stub.PutState(key, productJSONBytes)
	if err != nil{
		return shim.Error("Add product has error"+ err.Error())
	}

	//create index
	indexName := "myIndex"
	addressNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{product.Address, product.Productname})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	stub.PutState(addressNameIndexKey, value)

	return shim.Success(nil)
}