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

//用于测试index的表结构
type music struct {
	Name string `json:"name"`
	Words string `json:"words"`
	Song string `json:"song"`
	Singer string `json:"singer"`
	Time string `json:"time"`
	Style string `json:"style"`
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

func AddMusic(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 6{
		return shim.Error("Input arguments count isn't 3")
	}
	var err error

	var name = args[0]
	var words = args[1]
	var song = args[2]
	var singer = args[3]
	var time = args[4]
	var style = args[5]

	var key = "Music"+name

	music := &music{name,words,song,singer,time,style}
	musicJSONBytes, err := json.Marshal(music)
	if err != nil{
		return shim.Error("Music obj to byte arr has error"+ err.Error())
	}

	err = stub.PutState(key,musicJSONBytes)
	if err != nil{
		return shim.Error("Add product has error"+ err.Error())
	}

	//create index
	indexName := "index~music"
	musicIndexKey, err := stub.CreateCompositeKey(indexName, []string{music.Style, music.Name, music.Singer})
	if err != nil{
		return shim.Error(err.Error())
	}

	indexName2 := "index~music2"
	musicIndexKey2, err := stub.CreateCompositeKey(indexName2, []string{music.Song, music.Words})
	if err != nil{
		return shim.Error(err.Error())
	}

	value := []byte{0x00}
	stub.PutState(musicIndexKey, value)

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
	//使用产品地址和产品名称组成组合键以用于查询操作
	addressNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{product.Address, product.Productname})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	stub.PutState(addressNameIndexKey, value)

	return shim.Success(nil)
}