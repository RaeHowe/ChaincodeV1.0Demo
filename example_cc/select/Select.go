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
		return shim.Error(err.Error()) //Get user data has error
	}else if userBytes == nil{
		return shim.Error(err.Error()) //User data don't exist
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
		return shim.Error(err.Error()) //Get product has error
	}else  if productBytes == nil{
		return shim.Error(err.Error()) //Product data don't exist
	}

	return shim.Success(productBytes)
}

func SelectTest(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1{
		return shim.Error("SelectTest parameter count isn't 1")
	}
	var err error
	var A = args[0]
	var key = "Test" + A

	testBytes, err := stub.GetState(key)
	if err != nil{
		return shim.Error(err.Error()) //Get test has error
	}else if testBytes == nil{
		return shim.Error(err.Error()) //Test data don't exist
	}

	return shim.Success(testBytes)
}

//范围查询
func GetProductByRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2{
		return shim.Error("GetProductByRange parameter count isn't 2")
	}

	startKey := args[0]
	endKey := args[1]

	resultsIterator, err := stub.GetStateByRange(startKey, endKey) // Productapple Productpear
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
		return shim.Error(err.Error()) //GetQueryResult has error
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

//通过索引来查询test对象
func GetTestByIndexOfD(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1{
		return shim.Error("GetTestByIndexOfD parameter count isn't 1")
	}
	var err error
	var D = args[0]

	indexName := "myIndexOfDA"
	testResultsIterator, err := stub.GetStateByPartialCompositeKey(indexName, []string{D})
	if err != nil{
		return shim.Error(err.Error())
	}
	defer testResultsIterator.Close()

	var i int
	var result = "["
	var needComma = false
	for i = 0; testResultsIterator.HasNext(); i++{
		result += "{"
		responseRange, err := testResultsIterator.Next()
		if err != nil{
			return shim.Error(err.Error())
		}

		indexName, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil{
			return shim.Error(err.Error())
		}

		returnD := compositeKeyParts[0]
		returnA := compositeKeyParts[1]

		result += fmt.Sprintf("data%d: Get Test obj from index:%s D:%s A:%s \n", i, indexName, returnD, returnA)

		result += "}"
		if needComma == true {
			result += ","
		}
		needComma = true
	}

	result += "]"

	return shim.Success([]byte(result))
}

func GetTestByIndexOfDA(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2{
		return shim.Error("GetTestByIndex2Of1 parameter count isn't 2")
	}
	var err error
	var D = args[0]
	var A = args[1]

	indexName := "myIndexOfDA"
	testResultsIterator, err := stub.GetStateByPartialCompositeKey(indexName, []string{D, A})
	if err != nil{
		return shim.Error(err.Error())
	}
	defer testResultsIterator.Close()

	var i int
	var result = "["
	var needComma = false
	for i = 0; testResultsIterator.HasNext(); i++{
		result += "{"
		responseRange, err := testResultsIterator.Next()
		if err != nil{
			return shim.Error(err.Error())
		}

		indexName, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil{
			return shim.Error(err.Error())
		}

		returnD := compositeKeyParts[0]
		returnA := compositeKeyParts[1]

		result += fmt.Sprintf("data%d: Get Test obj from index:%s D:%s A:%s \n", i, indexName, returnD, returnA)

		result += "}"
		if needComma == true {
			result += ","
		}
		needComma = true
	}

	result += "]"

	return shim.Success([]byte(result))
}

func GetTestByIndexOfB(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1{
		return shim.Error("GetTestByIndex2Of2 parameter count isn't 2")
	}
	var err error
	var B = args[0]

	indexName := "myIndexOfBC"
	testResultsIterator, err := stub.GetStateByPartialCompositeKey(indexName, []string{B})
	if err != nil{
		return shim.Error(err.Error())
	}
	defer testResultsIterator.Close()

	var i int
	var result = "["
	var needComma = false
	for i = 0; testResultsIterator.HasNext(); i++{
		result += "{"
		responseRange, err := testResultsIterator.Next()
		if err != nil{
			return shim.Error(err.Error())
		}

		indexName, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil{
			return shim.Error(err.Error())
		}

		returnB := compositeKeyParts[0]
		returnC := compositeKeyParts[1]

		result += fmt.Sprintf("data%d: Get Test obj from index:%s B:%s C:%s \n", i, indexName, returnB, returnC)

		result += "}"
		if needComma == true {
			result += ","
		}
		needComma = true
	}

	result += "]"

	return shim.Success([]byte(result))
}

func GetTestByIndexOfBC(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2{
		return shim.Error("GetTestByIndex parameter count isn't 2")
	}
	var err error
	var B = args[0]
	var C = args[1]

	indexName := "myIndexOfBC"
	testResultsIterator, err := stub.GetStateByPartialCompositeKey(indexName, []string{B, C})
	if err != nil{
		return shim.Error(err.Error())
	}
	defer testResultsIterator.Close()

	var i int
	var result = "["
	var needComma = false
	for i = 0; testResultsIterator.HasNext(); i++{
		result += "{"
		responseRange, err := testResultsIterator.Next()
		if err != nil{
			return shim.Error(err.Error())
		}

		indexName, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil{
			return shim.Error(err.Error())
		}

		returnB := compositeKeyParts[0]
		returnC := compositeKeyParts[1]

		result += fmt.Sprintf("data%d: Get Test obj from index:%s B:%s C:%s \n", i, indexName, returnB, returnC)

		result += "}"
		if needComma == true {
			result += ","
		}
		needComma = true
	}

	result += "]"

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