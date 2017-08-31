package _select

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"bytes"
	"strings"
	"fmt"
	"time"
	"strconv"
)

func SelectUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1{
		return shim.Error("SelectUser parameter count can't less 1")
	}
	var err error
	var username = args[0]
	var key = "User" + username

	userBytes, err := stub.GetState(key)
	if err != nil{
		return shim.Error("Get user data has error"+ err.Error())
	}else if userBytes == nil{
		return shim.Error("User data don't exist"+ err.Error())
	}

	return shim.Success(userBytes)
}

func SelectProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args) < 1{
		return shim.Error("SelectProduct parameter count can't less 1")
	}
	var err error
	var productname = args[0]
	var key = "Product" + productname

	productBytes, err := stub.GetState(key)
	if err != nil{
		return shim.Error("Get product has error"+ err.Error())
	}else  if productBytes == nil{
		return shim.Error("Prodcut data don't exist"+ err.Error())
	}

	return shim.Success(productBytes)
}

//范围查询
func GetProductByRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2{
		return shim.Error("Incorrect number of parameter. Expecting 2")
	}

	startKey := args[0]
	endKey := args[1]

	resultsIterator, err := stub.GetStateByRange(startKey, endKey) // a,a
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	needComma := false //不需要逗号
	for resultsIterator.HasNext(){
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if needComma == true{
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		needComma = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

//通过couchDB丰富的查询功能进行查询
func GetProductByAddress(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args) < 1{
		return shim.Error("GetProductByAddress parameter count can't less 1")
	}
	var err error

	var address = strings.ToLower(args[0])
	queryString := fmt.Sprintf("{\"selector\":{\"address\":\"%s\"}}", address)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil{
		return shim.Error("GetQueryResult has error")
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	needcomma := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error("Next() has error")
		}
		// Add a comma before array members, suppress it for the first array member
		if needcomma == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		needcomma = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

func GetMusicByIndexOfSong(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2{
		return shim.Error("GetMusicByIndexOfSong parameter count can't less 2")
	}
	var err error

	song := args[0]
	words := args[1]

	indexName := "index~music2"
	musicResultsIterator, err := stub.GetStateByPartialCompositeKey(indexName, []string{song, words})
	if err != nil{
		return shim.Error(err.Error())
	}

	defer musicResultsIterator.Close()

	var i int
	var result = ""
	for i = 0; musicResultsIterator.HasNext(); i++{
		responseRange, err := musicResultsIterator.Next()
		if err != nil{
			return shim.Error(err.Error())
		}

		indexName, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil{
			return shim.Error(err.Error())
		}

		returnSong := compositeKeyParts[0]
		returnWords := compositeKeyParts[1]

		result += fmt.Sprintf("found a music from index:%s song:%s words:%s \n", indexName, returnSong, returnWords)
	}

	return shim.Success([]byte(result))
}

func GetMusicByIndexOfStyle(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2{
		return shim.Error("GetMusicByIndexOfStyle parameter count can't less 2")
	}
	var err error

	style := args[0]
	name := args[1]

	indexName := "index~music"
	musicResultsIterator, err := stub.GetStateByPartialCompositeKey(indexName, []string{style, name})
	if err != nil{
		return shim.Error(err.Error())
	}
	defer musicResultsIterator.Close()

	var i int

	var result = ""
	for i = 0; musicResultsIterator.HasNext(); i++ {
		responseRange, err := musicResultsIterator.Next()
		if err != nil{
			return shim.Error(err.Error())
		}

		indexName, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil{
			return shim.Error(err.Error())
		}

		returnStyle := compositeKeyParts[0]
		returnName := compositeKeyParts[1]
		returnSinger := compositeKeyParts[2]

		result += fmt.Sprintf("found a music from index:%s style:%s name:%s singer:%s\n", indexName, returnStyle, returnName, returnSinger)
	}

	return shim.Success([]byte(result))
}

//通过索引进行查询
func GetProductByIndex(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args) < 1{
		return shim.Error("GetProductByIndex parameter count can't less 2")
	}
	var err error
	address := args[0]

	indexName := "myIndex"
	//使用相应产品的索引来进行查询
	productResultsIterator, err := stub.GetStateByPartialCompositeKey(indexName, []string{address})
	if err != nil{
		return shim.Error(err.Error())
	}
	defer productResultsIterator.Close()

	var i int

	var result = ""
	for i = 0; productResultsIterator.HasNext(); i++ {
		responseRange, err := productResultsIterator.Next()
		if err != nil{
			return shim.Error(err.Error())
		}

		indexName, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil{
			return shim.Error(err.Error())
		}

		returnAddress := compositeKeyParts[0]
		returnProductName := compositeKeyParts[1]

		result += fmt.Sprintf("found a product from index:%s productname:%s address:%s \n", indexName, returnProductName, returnAddress)
	}

	return shim.Success([]byte(result))
}

func GetHistoryForProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args) < 1{
		return shim.Error("GetHistoryForProduct parameter count can't less 1")
	}

	productName := args[0]
	key := "Product"+productName
	fmt.Printf("- start getHistoryForProduct: %s\n", productName)

	resultsIterator, err := stub.GetHistoryForKey(key)
	if err != nil{
		return shim.Error("Get history has error"+err.Error())
	}

	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	needComma := false
	for resultsIterator.HasNext(){
		response, err := resultsIterator.Next()
		if err != nil{
			return shim.Error(err.Error())
		}

		if needComma == true{
			buffer.WriteString(",")
		}

		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		needComma = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}