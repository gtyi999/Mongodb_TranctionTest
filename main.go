package main

import (
	"fmt"
	"github.com/globalsign/mgo"
	"time"
	"github.com/globalsign/mgo/bson"
	"github.com/globalsign/mgo/txn"
	"sync"
)

type user_assets struct {
	Id                  bson.ObjectId                         `bson:"_id"`
	Uid                 int64                                 `bson:"uid"`
	Free_service_charge bool                                  `bson:"free_service_charge"` // 是否免手续费
	TokenAuth           string                                `bson:"token_auth"`          // 登录token，用于校验用户身份
	RsaPublicKey        string                                `bson:"rsa_public_key"`      // api交易时候需要用到的公钥，用于校验用户身份
	Asset               map[string]bson.Decimal128            `bson:"asset"`               // 可用资产
	Asset_frozen        map[string]map[string]bson.Decimal128 `bson:"freeze_asset"`        // 冻结资产
	Phone               string                                `bson:"phone"`               // 电话
	TokenDespire        time.Time                             `bson:"token_despire"`       // token过期时间
	NameAuth            int64                                 `bson:"name_auth"`           // 实名认证
	ShamTrade           bool                                  `bson:"sham_trade"`          // 发送伪装交易
	OrderUnlimit        bool                                  `bson:"order_unlimit"`       // 是否有挂单无限流
	DisableTrade        bool                                  `bson:"disable_trade"`       // 是否禁止交易
}

type Account struct {
	ID       bson.ObjectId   `bson:"_id"`
	IDA      int `bson:"ida"`
	Balance int `bson:"balance"`
}

type ID struct {
	FirstName string
	LastName  string
}

func doStuff (wg *sync.WaitGroup, session *mgo.Session,i int){
	defer wg.Done()

	// Copy the session - if needed this will dial a new connection which
	// can later be reused.
	// 复制会话 - 如果需要，将建立新连接，以后可以重复使用。
	// Calling close returns the connection to the pool.
	//调用close将返回与池的连接。
	conn := session.Copy()
	defer conn.Close()

	// Do something(s) with the connection
	_, err:= conn.DB("").C("accountsT").Count()
	if err!=nil{
		fmt.Println("conn.DB.Count err:",err)
		return
	}

	fmt.Println("goroutines:",i)
	c_a:=conn.DB("").C("accountsT")
	result := &Account{}
	err = c_a.Find(bson.M{"ida": 4}).One(&result)
	fmt.Println("result:",result)
	if err!=nil{
		fmt.Println("Mongo find err1:",err)

	}


	//SetDebug启用或禁用调试。
	//txn.SetDebug(false)
    //txn.SetLogger(nil)

	//NewRunner返回一个新的事务运行器，它使用tc来保存其事务。
	//多个事务集合可能存在于单个数据库中，但是给定事务集合中的操作所触及的所有集合必须由它独占。
	//具有相同名称tc但后缀为“.stash”的第二个集合将用于实现insert和remove操作的事务行为。
	runner := txn.NewRunner(c_a)
	ops := []txn.Op{{
		C:      "accountsT",
		Id:     result.ID,
		Assert: bson.M{"balance": bson.M{"$gt": 0}},
		Update: bson.M{"$inc": bson.M{"balance": 100}},
	}}

	err = runner.Run(ops, "", nil)
	if err != nil {
		fmt.Println("00runner.Run err:", err)
	}

}


func doStuff_concurrent (wg *sync.WaitGroup, session *mgo.Session,i int){
	defer wg.Done()

	// Copy the session - if needed this will dial a new connection which
	// can later be reused.
	// 复制会话 - 如果需要，将建立新连接，以后可以重复使用。
	// Calling close returns the connection to the pool.
	//调用close将返回与池的连接。
	conn := session.Copy()
	defer conn.Close()

	// Do something(s) with the connection
	_, err:= conn.DB("").C("accountsT").Count()
	if err!=nil{
		fmt.Println("conn.DB.Count err:",err)
		return
	}

	fmt.Println("goroutines:",i)
	c_a:=conn.DB("").C("accountsT")
	result := &Account{}
	err = c_a.Find(bson.M{"ida": 4}).One(&result)
	fmt.Println("result:",result)
	if err!=nil{
		fmt.Println("Mongo find err1:",err)

	}


	//SetDebug启用或禁用调试。
	//txn.SetDebug(false)
	//txn.SetLogger(nil)

	t1:=time.Now().Unix()

	for j:=0;j<10000;j++ {

		//NewRunner返回一个新的事务运行器，它使用tc来保存其事务。
		//多个事务集合可能存在于单个数据库中，但是给定事务集合中的操作所触及的所有集合必须由它独占。
		//具有相同名称tc但后缀为“.stash”的第二个集合将用于实现insert和remove操作的事务行为。
		runner := txn.NewRunner(c_a)
		ops := []txn.Op{{
			C:      "accountsT",
			Id:     result.ID,
			Assert: bson.M{"balance": bson.M{"$gt": 0}},
			Update: bson.M{"$inc": bson.M{"balance": 100}},
		}}

		err = runner.Run(ops, "", nil)
		if err != nil {
			fmt.Println("00runner.Run err:", err)
		}
	}
	t2:=time.Now().Unix()
	fmt.Println("i:",i,",t2-t1:",t2-t1)

}

func doStuffTest (wg *sync.WaitGroup, session *mgo.Session){
	defer wg.Done()

	// Copy the session - if needed this will dial a new connection which
	// can later be reused.
	// 复制会话 - 如果需要，将建立新连接，以后可以重复使用。
	// Calling close returns the connection to the pool.
	//调用close将返回与池的连接。
	conn := session.Copy()
	defer conn.Close()

	// Do something(s) with the connection
	_, err:= conn.DB("").C("accountsT").Count()
	if err!=nil{
		fmt.Println("conn.DB.Count err:",err)
		return
	}

	c_a:=conn.DB("").C("accountsT")
	result := &Account{}
	err = c_a.Find(bson.M{"ida": 4}).One(&result)
	fmt.Println("result:",result)
	if err!=nil{
		fmt.Println("Mongo find err1:",err)

	}



	//timer := time.NewTimer(time.Second * 1)

	t6:=time.Now().Unix()
	for j:=0;j<100000;j++ {

		//NewRunner返回一个新的事务运行器，它使用tc来保存其事务。
		//多个事务集合可能存在于单个数据库中，但是给定事务集合中的操作所触及的所有集合必须由它独占。
		//具有相同名称tc但后缀为“.stash”的第二个集合将用于实现insert和remove操作的事务行为。
		runner := txn.NewRunner(c_a)
		ops := []txn.Op{{
			C:      "accountsT",
			Id:     result.ID,
			Assert: bson.M{"balance": bson.M{"$gt": 0}},
			Update: bson.M{"$inc": bson.M{"balance": 100}},
		}}

		err = runner.Run(ops, "", nil)
		if err != nil {
			fmt.Println("00runner.Run err:", err)
		}

		//fmt.Println("i:",i)
		//select {
		//case <-timer.C:
		//	// call timed out
		//	fmt.Println("tick time out")
		//	return
		//}
	}

	t7:=time.Now().Unix()
	fmt.Println("t7-t6=",t7-t6)

}



const  (
	pool_num=4
)

func main(){
	//0-----协程并发的例子
	// This example shows the best practise for concurrent use of a mgo session.
	//此示例显示了并发使用mgo会话的最佳实践。
	// Internally mgo maintains a connection pool, dialling new connections as
	// required.
	//内部mgo维护连接池，根据需要建立新连接。
	// Some general suggestions:
	//一些一般性建议：
	// 		- Define a struct holding the original session, database name and
	// 			collection name instead of passing them explicitly.
	//定义包含原始会话，数据库名称和集合名称的结构，而不是显式传递它们。
	// 		- Define an interface abstracting your data access instead of exposing
	// 			mgo to your application code directly.
	//定义一个抽象数据访问的接口，而不是直接将mgo暴露给您的应用程序代码。
	// 		- Limit concurrency at the application level, not with SetPoolLimit().
    //   限制应用程序级别的并发性，而不是使用SetPoolLimit（）。
	// This will be our concurrent worker
	//这将是我们的并发工作线程

	//0-----协程并发的例子
	// Dial a connection to Mongo - this creates the connection pool 与Mongo的连接 - 这将创建连接池
	//session, err := mgo.Dial("xietianran:xtr123456@127.0.0.1:27017/trade")
	//if err != nil {
	//	panic(err)
	//}

	// Concurrently do things, passing the session to the worker
	//wg := &sync.WaitGroup{}
	//for i := 0; i < 10; i++ {
	//	wg.Add(1)
	//	go doStuff(wg, session,i)
	//}
	//wg.Wait()
	//session.Close()
	//


	//01--测试一秒能执行多少次事务？  一个连接  一秒做20次事务（upadte 和insert) 左右     10万次用时 489秒
	// Dial a connection to Mongo - this creates the connection pool 与Mongo的连接 - 这将创建连接池
	//session, err := mgo.Dial("xietianran:xtr123456@127.0.0.1:27017/trade")
	//if err != nil {
	//	panic(err)
	//}
	//wg := &sync.WaitGroup{}
	////WaitGroup等待完成goroutine的集合。主goroutine调用Add来设置要等待的goroutines的数量。
	////然后每个goroutine运行并在完成后调用Done。同时，Wait可以用来阻塞 直到所有goroutine完成。
	////首次使用后，不得复制WaitGroup。
	//wg.Add(1)
	//go doStuffTest(wg, session)
	//wg.Wait()
	//session.Close()
	//fmt.Println("end main")


	//02----多个连接的情况的事务的效率 ？--10协程 10个连接  每个协程处理10000个事务 平均用时970秒  平均每秒处理10个事务
	// Dial a connection to Mongo - this creates the connection pool 与Mongo的连接 - 这将创建连接池
	var 	sessions         []*mgo.Session
	i:=0
	for i=0;i<10;i++ {
		session, err := mgo.Dial("xietianran:xtr123456@192.168.163.213:27017/trade")
		if err != nil {
			panic(err)
		}
		sessions=append(sessions,session)
	}
	wg := &sync.WaitGroup{}
	for i = 0; i < 10; i++ {
		wg.Add(1)
		go doStuff_concurrent(wg, sessions[i],i)
	}
	wg.Wait()
	for i=0;i<10;i++ {
		sessions[i].Close()
	}
	fmt.Println("end main")



	//fmt.Println("hello main")
	//url:="xietianran:xtr123456@192.168.163.208:27017/trade"
	//session, err := mgo.Dial(url)
	//if err!=nil{
	//	fmt.Println("mgo.Dial err:",err,",session:",session)
	//}
	//
	//c := session.DB("trade").C("user_assets")
	//defer session.Clone()
	////0---mongo 查找
	//result := &user_assets{}
	//err = c.Find(bson.M{"uid": 1014}).One(&result)
	//if err!=nil{
	//	fmt.Println("Mongo find err1:",err)
	//
	//}
	////obId1:=result.Id
	////fmt.Printf("result1:%v\n",result)
	//
	//err = c.Find(bson.M{"uid": 1015}).One(&result)
	//if err!=nil{
	//	fmt.Println("Mongo find err2:",err)
	//
	//}
	////obId2:=result.Id
	////fmt.Printf("result2:%v\n",result)
	////fmt.Println("obId2:",obId2,",obId1:",obId1)

	//1--mongo Update 事务
	//runner := txn.NewRunner(c)
	//
	//ops := []txn.Op{{
	//	C:      "user_assets",
	//	Id:     obId1,
	//	Assert: bson.M{"asset.btc": bson.M{"$gte": 100}},
	//	Update: bson.M{"$inc": bson.M{"asset.btc": -100}},
	//}, {
	//	C:      "user_assets",
	//	Id:     obId2,
	//	Assert: bson.M{"asset.btc": bson.M{"$gte": 1000001}},
	//	Update: bson.M{"$inc": bson.M{"asset.btc": 3}},
	//}}
	//id := bson.NewObjectId() // Optional
	//err = runner.Run(ops, id, nil)
	//if err!=nil{
	//	fmt.Println("runner.Run err:",err)
	//}

	//2--mongo Insert 事务
	//result.Uid=88888888
	//runner := txn.NewRunner(c)
	//ops := []txn.Op{{
	//	C:      "user_assets",
	//	Id:    bson.NewObjectId(),
	//	Assert: txn.DocMissing,
	//	Insert: result,
	//}}
	//id := bson.NewObjectId() // Optional
	//err = runner.Run(ops, id, nil)
	//if err!=nil{
	//	fmt.Println("runner.Run err:",err)
	//}

	//3--mongo Remove 事务
	//runner := txn.NewRunner(c)
	//ops := []txn.Op{{
	//	C:      "user_assets",
	//	Id:    result.Id,
	//	Assert: txn.DocExists,
	//	Remove: true,
	//}}
	//id := bson.NewObjectId() // Optional
	//err = runner.Run(ops, id, nil)
	//if err!=nil{
	//	fmt.Println("runner.Run err:",err)
	//}


	//4--测试是否存在文档
	//c1 := session.DB("trade").C("accountsT")
	//err = c1.Insert(bson.M{"_id": 0, "balance": 300})
	//if err!=nil{
	//	fmt.Println("c1.Insert err:",err)
	//}
	//
	//runner := txn.NewRunner(c1)
	//exists := []txn.Op{{
	//	C:      "accountsT",
	//	Id:     0,
	//	Assert: txn.DocExists,
	//}}
	//missing := []txn.Op{{
	//	C:      "accountsT",
	//	Id:     0,
	//	Assert: txn.DocMissing,
	//}}
	//
	//err = runner.Run(exists, "", nil)
	//if err!=nil{
	//	fmt.Println("1exists runner.Run err:",err)
	//}
	//err = runner.Run(missing, "", nil)
	//if err!=nil{
	//	fmt.Println("2missing runner.Run err:",err)
	//}
	//
	//err = c1.RemoveId(0)
	//if err!=nil{
	//	fmt.Println("c1.RemoveId err:",err)
	//}
	//
	//err = runner.Run(exists, "", nil)
	//if err!=nil{
	//	fmt.Println("3runner.Run err:",err)
	//}
	//err = runner.Run(missing, "", nil)
	//if err!=nil{
	//	fmt.Println("4 runner.Run err:",err)
	//}

	//5 TestInsert
	//c1 := session.DB("trade").C("accountsT")
	//err = c1.Insert(bson.M{"_id": 0, "balance": 300})
	//if err!=nil{
	//	fmt.Println("c1.Insert err:",err)
	//}
	//
	//runner := txn.NewRunner(c1)
	//ops := []txn.Op{{
	//	C:      "accountsT",
	//	Id:     0,
	//	Insert: bson.M{"balance": 200},
	//}}
	//
	//err = runner.Run(ops, "", nil)
	//if err!=nil{
	//	fmt.Println("00runner.Run err:",err)
	//}
	//
	//var account Account
	//err = c1.FindId(0).One(&account)
	//if err!=nil{
	//	fmt.Println("c1.FindId err:",err)
	//}
	//
	//ops[0].Id = 1
	//err = runner.Run(ops, "", nil)
	//if err!=nil{
	//	fmt.Println("11account:",account)
	//	fmt.Println("11runner.Run err:",err)
	//}
	//
	//err =c1.FindId(1).One(&account)
	//if err!=nil{
	//	fmt.Println("22account:",account)
	//	fmt.Println("22runner.Run err:",err)
	//}

	//6 TestInsertStructId
	//c1 := session.DB("trade").C("accountsT")
	//ops := []txn.Op{{
	//	C:      "accountsT",
	//	Id:     ID{FirstName: "John", LastName: "Jones"},
	//	Assert: txn.DocMissing,
	//	Insert: bson.M{"balance": 200},
	//}, {
	//	C:      "accountsT",
	//	Id:     ID{FirstName: "Sally", LastName: "Smith"},
	//	Assert: txn.DocMissing,
	//	Insert: bson.M{"balance": 800},
	//}}
	//runner := txn.NewRunner(c1)
	//
	//err = runner.Run(ops, "", nil)
	//if err!=nil{
	//	fmt.Println("1TestInsertStructId runner.Run err:",err)
	//}
	//n, _ := c1.Find(nil).Count()
	//fmt.Println("2TestInsertStructId c1.Find Count:",n)


	//7 TestRemove
	//c1 := session.DB("trade").C("accountsT")
	//runner := txn.NewRunner(c1)
	//err = c1.Insert(bson.M{"_id": 0, "balance": 300})
	//if err!=nil{
	//	fmt.Println("1TestRemove c1.Insert err:",err)
	//}
	//
	//ops := []txn.Op{{
	//	C:      "accountsT",
	//	Id:     0,
	//	Remove: true,
	//}}
	//
	//err = runner.Run(ops, "", nil)
	//if err!=nil{
	//	fmt.Println("2TestRemove runner.Run:",err)
	//}
	//
	//err = c1.FindId(0).One(nil)
	//if err!=nil{
	//	fmt.Println("3TestRemove c1.FindId:",err)
	//}
	//
	//err = runner.Run(ops, "", nil)
	//if err!=nil{
	//	fmt.Println("4TestRemove runner.Run:",err)
	//}

	//8 TestUpdate
	//c1 := session.DB("trade").C("accountsT")
	//runner := txn.NewRunner(c1)
	//err = c1.Insert(bson.M{"_id": 0, "balance": 200})
	//if err!=nil{
	//	fmt.Println("1TestUpdate c1.Insert err:",err)
	//}
	//err = c1.Insert(bson.M{"_id": 1, "balance": 200})
	//if err!=nil{
	//	fmt.Println("2TestUpdate c1.Insert err:",err)
	//}
	//
	//ops := []txn.Op{{
	//	C:      "accountsT",
	//	Id:     0,
	//	Update: bson.M{"$inc": bson.M{"balance": 100}},
	//}}
	//
	//err = runner.Run(ops, "", nil)
	//if err!=nil{
	//	fmt.Println("3TestUpdate runner.Run err:",err)
	//}
	//
	//var account Account
	//err = c1.FindId(0).One(&account)
	//if account.Balance==300{
	//	fmt.Println("4TestUpdate account.Balance==300,update success")
	//}
	//if err!=nil{
	//	fmt.Println("4TestUpdate c1.FindId err:",err)
	//}
	//
	//ops[0].Id = 1
	//err = runner.Run(ops, "", nil)
	//if err!=nil{
	//	fmt.Println("5TestUpdate runner.Run err:",err)
	//}
	//err = c1.FindId(1).One(&account)
	//if account.Balance==300{
	//	fmt.Println("6TestUpdate account.Balance==300,update success")
	//}


	//9 TestInsertUpdate
	//c1 := session.DB("trade").C("accountsT")
	//runner := txn.NewRunner(c1)
	//ops := []txn.Op{{
	//	C:      "accountsT",
	//	Id:     0,
	//	Insert: bson.M{"_id": 0, "balance": 200},
	//	Assert:txn.DocMissing, //这个很重要 如不加这个字段每次都会执行事务成功 即下面的update会成功 加这个字段就只会执行成功一次
	//}, {
	//	C:      "accountsT",
	//	Id:     0,
	//	Update: bson.M{"$inc": bson.M{"balance": 100}},
	//}}
	//
	//err = runner.Run(ops, "", nil)
	//if err!=nil{
	//	fmt.Println("1TestInsertUpdate runner.Run err:",err)
	//}
	//
	//var account Account
	//err = c1.FindId(0).One(&account)
	//if account.Balance==300{
	//	fmt.Println("2TestInsertUpdate account.Balance==300,success")
	//}
	//
	//err = runner.Run(ops, "", nil)
	//if err!=nil{
	//	fmt.Println("3TestInsertUpdate runner.Run err:",err)
	//}
	//
	//err = c1.FindId(0).One(&account)
	//if err!=nil{
	//	fmt.Println("4TestInsertUpdate c1.FindId err:",err)
	//}
	//fmt.Println("5TestInsertUpdate account:",account)

	//10 TestUpdateInsert
	//c1 := session.DB("trade").C("accountsT")
	//err = c1.Insert(bson.M{"_id": 0, "balance": 200})
	//if err!=nil{
	//	fmt.Println("TestUpdateInsert c1.Insert err:",err)
	//}
	//runner := txn.NewRunner(c1)
	//ops := []txn.Op{{
	//	C:      "accountsT",
	//	Id:     0,
	//	Update: bson.M{"$inc": bson.M{"balance": 100}},
	//	Assert: bson.M{"_id": bson.M{"$gte": 0}},
	//}, {
	//	C:      "accountsT",
	//	Id:     1,
	//	Insert: bson.M{"_id": 1, "balance": 200},
	//	Assert: txn.DocMissing,
	//}}
	//
	//err = runner.Run(ops, "", nil)
	//if err!=nil{
	//	fmt.Println("1TestUpdateInsert runner.Run err: ",err)
	//}

	//11 TestInsertRemoveInsert
	//c1 := session.DB("trade").C("accountsT")
	//runner := txn.NewRunner(c1)
	//ops := []txn.Op{{
	//	C:      "accountsT",
	//	Id:     0,
	//	Insert: bson.M{"_id": 0, "balance": 200},
	//	Assert:txn.DocMissing,
	//}, {
	//	C:      "accountsT",
	//	Id:     0,
	//	Remove: true,
	//}, {
	//	C:      "accountsT",
	//	Id:     0,
	//	Insert: bson.M{"_id": 0, "balance": 300},
	//	Assert:txn.DocMissing,
	//}}
	//
	//err = runner.Run(ops, "", nil)
	//if err!=nil{
	//	fmt.Println("1TestInsertRemoveInsert runner.Run,err:",err)
	//}




	//-------------------------------------------------------------------------------------
	//txn.Op 参数说明如下
	//Assert
	//Assert optionally holds a query document that is used to test the operation document at the time the transaction is
	//going to be applied.
	//	断言可选地保存一个查询文档，该文档用于在将要应用事务时测试操作文档。
	//The assertions for all operations in a transaction are tested before any changes take place,and the transaction is entirely aborted if any of them fails.
	//	在发生任何更改之前，将对事务中所有操作的断言进行测试，如果任何更改失败，则事务将完全中止。
	//This is also the only way to prevent a transaction from being being applied (the transaction continues despite the outcome of Insert, Update, and Remove).
	//这也是阻止应用事务的唯一方法（尽管插入，更新和删除结果，事务仍在继续）。


	//
	//0
	//The Insert, Update and Remove fields describe the mutation intended by the operation.
	//	At most one of them may be set per operation.
	//	If none are set, Assert must be set and the operation becomes a read-only test.
	//“插入”，“更新”和“删除”字段描述了操作所需的突变。
	//每次操作最多可以设置其中一个。
	//如果未设置，则必须设置Assert并且操作变为只读测试。
	//1
	//Insert holds the document to be inserted at the time the transaction is applied.
	//	The Id field will be inserted into the document automatically as its _id field.
	//	The transaction will continue even if the document already  exists.
	//	Use Assert with txn.DocMissing if the insertion is required.
	//	Insert在应用事务时保存要插入的文档。
	//Id字段将作为其_id字段自动插入到文档中。
	//即使文档已存在，交易仍将继续。
	//如果需要插入，请使用带有txn.DocMissing的Assert。
	//2
	//Update holds the update document to be applied at the time the transaction is applied.
	//	The transaction will continue even if a document with Id is missing.
	//	Use Assert to test for the document presence or its contents.
	//	Update保存要在应用事务时应用的更新文档。
	//即使缺少具有Id的文档，事务仍将继续。
	//使用Assert测试文档是否存在或其内容。
	//3
	//Remove indicates whether to remove the document with Id.
	//	The transaction continues even if the document doesn't yet exist at the time the transaction is applied.
	//	Use Assert with txn.DocExists to make sure it will be removed.
	//	删除表示是否删除带有Id的文档。
	//即使在应用事务时文档尚不存在，事务仍将继续。
	//与txn.DocExists一起使用Assert以确保它将被删除。


	//4 txn-revno
	//The document revision is the value of the txn-revno field after the change has been applied.
	//	Negative values indicate the document was not present in the collection.
	//	Revisions will not change when updates or removes are applied to missing documents or inserts are attempted when the document isn't present.
	// 文档修订版是应用更改后txn-revno字段的值。
	//负值表示文档未出现在集合中。
	//当对缺失文档应用更新或删除时，或者在文档不存在时尝试插入时，修订不会更改。


	//PurgeMissing
	//PurgeMissing removes from collections any state that refers to transaction documents that
	//for whatever reason have been lost from the system (removed by accident or lost in a hard crash, for example).
	//PurgeMissing从集合中删除任何引用事务文档的状态，该事务文档无论出于何种原因从系统中丢失（例如，意外删除或在硬崩溃中丢失）。
	//This method should very rarely be needed, if at all, and should never be used during the normal operation of an application. Its purpose is to put a system that has seen unavoidable corruption back in a working state.
	//如果有的话，应该很少需要这种方法，并且在应用程序的正常操作期间不应该使用该方法。 其目的是将一个已经看到不可避免的腐败的系统重新置于工作状态。
	//


	//type Chaos struct
	//	// Chaos holds parameters for the failure injection mechanism.
	//	混沌保存故障注入机制的参数。
	//	KillChance is the 0.0 to 1.0 chance that a given checkpoint within the algorithm will raise an interruption that will stop the procedure.
	//	KillChance是算法中给定检查点引发中断的0.0到1.0几率，该中断将停止该过程。
	//
	//	SlowdownChance is the 0.0 to 1.0 chance that a given checkpoint within the algorithm will be delayed by Slowdown before continuing.
	//	SlowdownChance是算法中给定检查点在继续之前被Slowdown延迟的0.0到1.0的几率。
	//
	//	If Breakpoint is set, the above settings will only affect thenamed breakpoint.
	//	如果设置了断点，则上述设置仅影响已设置的断点。


	// SetChaos将故障注入参数设置为c。
	//func SetChaos（c Chaos）

}



