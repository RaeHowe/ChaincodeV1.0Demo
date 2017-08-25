jq --version > /dev/null 2>&1
if [ $? -ne 0 ]; then
	echo "Please Install 'jq' https://stedolan.github.io/jq/ to execute this script"
	echo
	exit 1
fi
starttime=$(date +%s)

echo "POST request Enroll on Org1  ..."
echo
ORG1_TOKEN=$(curl -s -X POST \
  http://localhost:4000/users \
  -H "content-type: application/x-www-form-urlencoded" \
  -d 'username=Jim&orgName=org1')
echo $ORG1_TOKEN
ORG1_TOKEN=$(echo $ORG1_TOKEN | jq ".token" | sed "s/\"//g")
echo
echo "ORG1 token is $ORG1_TOKEN"
echo
echo "POST request Enroll on Org2 ..."
echo
ORG2_TOKEN=$(curl -s -X POST \
  http://localhost:4000/users \
  -H "content-type: application/x-www-form-urlencoded" \
  -d 'username=Barry&orgName=org2')
echo $ORG2_TOKEN
ORG2_TOKEN=$(echo $ORG2_TOKEN | jq ".token" | sed "s/\"//g")
echo
echo "ORG2 token is $ORG2_TOKEN"
echo


echo "POST request Create channel  ..."
echo
curl -s -X POST \
  http://localhost:4000/channels \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"channelName":"mychannel",
	"channelConfigPath":"../artifacts/channel/mychannel.tx"
}'
echo
echo
sleep 5
echo "POST request Join channel on Org1"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/peers \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["localhost:7051","localhost:7056"]
}'
echo
echo

echo "POST request Join channel on Org2"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/peers \
  -H "authorization: Bearer $ORG2_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["localhost:8051","localhost:8056"]
}'
echo
echo


echo "POST Install chaincode on Org1"
echo
curl -s -X POST \
  http://localhost:4000/chaincodes \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
  "peers": ["localhost:7051","localhost:7056"],
  "chaincodeName":"mycc",
  "chaincodePath":"github.com/example_cc",
  "chaincodeVersion":"v0"
}'
echo
echo


echo "POST Install chaincode on Org2"
echo
curl -s -X POST \
  http://localhost:4000/chaincodes \
  -H "authorization: Bearer $ORG2_TOKEN" \
  -H "content-type: application/json" \
  -d '{
  "peers": ["localhost:8051","localhost:8056"],
  "chaincodeName":"mycc",
  "chaincodePath":"github.com/example_cc",
  "chaincodeVersion":"v0"
}'
echo
echo

echo "POST instantiate chaincode on peer1 of Org1"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
  "chaincodeName":"mycc",
  "chaincodeVersion":"v0",
  "functionName":"init",
  "args":["admin","1qasw23ed"]
}'
echo
echo

echo "插入两个用户: raehowe, peter"

echo "POST invoke chaincode on peers of Org1 and Org2"
echo
TRX_ID=$(curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
  "peers": ["localhost:7051", "localhost:8051"],
  "fcn":"addUser",
  "args":["raehowe", "123"]
}')
echo "Transacton ID is $TRX_ID"
echo

echo "POST invoke chaincode on peers of Org1 and Org2"
echo
TRX_ID=$(curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
  "peers": ["localhost:7051", "localhost:8051"],
  "fcn":"addUser",
  "args":["peter", "321"]
}')
echo "Transacton ID is $TRX_ID"
echo

echo "插入五种产品: apple, pear, orange, peach, grape"

echo "POST invoke chaincode on peers of Org1 and Org2"
echo
TRX_ID=$(curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
  "peers": ["localhost:7051", "localhost:8051"],
  "fcn":"addProduct",
  "args":["apple", "tianjin", "3"]
}')
echo "Transacton ID is $TRX_ID"
echo

echo "POST invoke chaincode on peers of Org1 and Org2"
echo
TRX_ID=$(curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
  "peers": ["localhost:7051", "localhost:8051"],
  "fcn":"addProduct",
  "args":["pear", "beijing", "2"]
}')
echo "Transacton ID is $TRX_ID"
echo

echo "POST invoke chaincode on peers of Org1 and Org2"
echo
TRX_ID=$(curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
  "peers": ["localhost:7051", "localhost:8051"],
  "fcn":"addProduct",
  "args":["orange", "shandong", "4"]
}')
echo "Transacton ID is $TRX_ID"
echo

echo "POST invoke chaincode on peers of Org1 and Org2"
echo
TRX_ID=$(curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
  "peers": ["localhost:7051", "localhost:8051"],
  "fcn":"addProduct",
  "args":["peach", "tianjin", "1"]
}')
echo "Transacton ID is $TRX_ID"
echo

echo "POST invoke chaincode on peers of Org1 and Org2"
echo
TRX_ID=$(curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
  "peers": ["localhost:7051", "localhost:8051"],
  "fcn":"addProduct",
  "args":["grape", "xinjiang", "3"]
}')
echo "Transacton ID is $TRX_ID"
echo

echo "新增3个产品用于测试: apple1, apple2, apple3"

echo "POST invoke chaincode on peers of Org1 and Org2"
echo
TRX_ID=$(curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
  "peers": ["localhost:7051", "localhost:8051"],
  "fcn":"addProduct",
  "args":["apple1", "jiangsu", "10"]
}')
echo "Transacton ID is $TRX_ID"
echo

echo "POST invoke chaincode on peers of Org1 and Org2"
echo
TRX_ID=$(curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
  "peers": ["localhost:7051", "localhost:8051"],
  "fcn":"addProduct",
  "args":["apple2", "gansu","11"]
}')
echo "Transacton ID is $TRX_ID"
echo

echo "POST invoke chaincode on peers of Org1 and Org2"
echo
TRX_ID=$(curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
  "peers": ["localhost:7051", "localhost:8051"],
  "fcn":"addProduct",
  "args":["apple3", "tianjin", "12"]
}')
echo "Transacton ID is $TRX_ID"
echo

echo "删除一个用户: peter 目前剩下: raehowe"

echo "POST invoke chaincode on peers of Org1 and Org2"
echo
TRX_ID=$(curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
  "peers": ["localhost:7051", "localhost:8051"],
  "fcn":"deleteUser",
  "args":["peter"]
}')
echo "Transacton ID is $TRX_ID"
echo

echo "删除一个产品: peach 目前剩下: apple, apple1, apple2, apple3, pear, orange, grape"

echo "POST invoke chaincode on peers of Org1 and Org2"
echo
TRX_ID=$(curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
  "peers": ["localhost:7051", "localhost:8051"],
  "fcn":"deleteProduct",
  "args":["peach"]
}')
echo "Transacton ID is $TRX_ID"
echo

echo "查询用户: raehowe"

echo "GET query chaincode on peer1 of Org1"
echo
curl -s -X GET \
  "http://localhost:4000/channels/mychannel/chaincodes/mycc?peer=peer1&fcn=selectUser&args=%5B%22raehowe%22%5D" \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json"
echo
echo

echo "查询产品: apple"

echo "GET query chaincode on peer1 of Org1"
echo
curl -s -X GET \
  "http://localhost:4000/channels/mychannel/chaincodes/mycc?peer=peer1&fcn=selectProduct&args=%5B%22apple%22%5D" \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json"
echo
echo

echo "范围查询: startKey:apple1 endKey: orange"
echo "商品顺序我猜测为:apple apple1 apple2 apple3 grape orange peach 如果进行范围查询apple1~orange的话，那么返回的结果确实为:apple1,apple2,apple3,grape 我已经验证通过"

echo "GET query chaincode on peer1 of Org1"
echo
curl -s -X GET \
  "http://localhost:4000/channels/mychannel/chaincodes/mycc?peer=peer1&fcn=selectProductByRange&args=%5B%22Productapple1%22%2C%22Productorange%22%5D" \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json"
echo
echo

echo "指定查询: 查询地址为tianjin的产品信息"

echo "GET query chaincode on peer1 of Org1"
echo
curl -s -X GET \
  "http://localhost:4000/channels/mychannel/chaincodes/mycc?peer=peer1&fcn=selectProductByAddress&args=%5B%22tianjin%22%5D" \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json"
echo
echo

echo "通过Index查询"

echo "GET query chaincode on peer1 of Org1"
echo
curl -s -X GET \
  "http://localhost:4000/channels/mychannel/chaincodes/mycc?peer=peer1&fcn=selectProductByIndex&args=%5B%22tianjin%22%5D" \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json"
echo
echo

