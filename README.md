# Blockchain1.0Demo
	Chaincode文件夹中，包括了chaincode项目代码，执行chaincode的shell脚本和readme，项目的整体功能，以及脚本的使用，请阅读本文档，谢谢。
	项目中包括了main文件:UChain.go以及操作账本的增删改查功能的四个文件夹:insert,delete,update,select
	在main文件中的Init函数，主要初始化了一条数据: key:"adminRole" value:{"username":"admin","password":"123"}
	在Invoke函数中，包括了这些方法:
		1.addUser 新增用户
		2.addProduct 新增产品
		3.updateUser 更新用户
		4.updateProduct 更新产品
		5.deleteUser 删除用户
		6.deleteProduct 删除产品
		7.selectUser 按照key进行用户查询
		8.selectProduct 按照key进行产品查询
		9.selectProductByRange 范围查询产品信息
		10.selectProductByAddress 通过产品地址，运用couchDB的丰富查询功能，进行查询
		11.selectProductByIndex 通过索引来进行查询
	在这里主要讲解新增产品和后三种查询功能，其他的功能，确实简单，请自行学习.
	在新增产品这个函数中，除了通常的将产品信息插入到账本数据库中之外，还新增了将该产品的索引信息插入到账本数据库中，过程如下：首先，你需要命名好索引的名称，在本例中，我命名为myIndex,之后需要确定这个索引可以查出哪些你需要的内容（内容是从该产品中的属性信息进行选取的），本例中，我选取了产品的地址和产品名称.这里需要注意一点，选取的属性信息，第一个是你必须作为参数进行查询的条件，本例中，代码为:addressNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{product.Address, product.Productname}) ,那么你就必须输入产品的地址信息，以进行查询使用

	在范围查询产品信息的这个函数中，运用到了blockchain自带的GetStateByRange方法，我新增了数据如下:apple,apple1,apple2,apple3,grape,orange,peach,商品在数据库里面的顺序为:apple apple1 apple2 apple3 grape orange peach,是按照首字母进行排序的，如果进行范围查询apple1~orange的话，那么返回的结果为:apple1,apple2,apple3,grape 我已经验证通过

	在通过产品地址查询的这个函数中，运用到了couchDB自带的一些查询功能：{"selector":{"address":"tianjin"}} 通过产品地址来查询产品的信息

	在通过索引来查询的这个函数中，通过使用在新增产品中增加的索引信息来进行索引查询。使用blockchain内置的GetStateByPartialCompositeKey方法，通过使用哪个索引作为一个参数，索引的条件信息数组作为第二个参数.
