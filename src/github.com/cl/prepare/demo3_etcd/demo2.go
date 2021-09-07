package main

//
//import (
//	"fmt"
//	clientv3 "go.etcd.io/etcd/client/v3"
//	"time"
//)
//
//func main(){
//	var(
//		config clientv3.Config
//		client *clientv3.Client
//		err error
//		putR *clientv3.PutResponse
//		getR *clientv3.GetResponse
//		lease clientv3.Lease
//		leaseR *clientv3.LeaseGrantResponse
//	)
//
//	config = clientv3.Config{
//		Endpoints: []string{"localhost:2379"},
//		DialTimeout: 5*time.Second,
//	}
//
//	// 建立连接
//	if client,err =  clientv3.New(config); err != nil{
//		fmt.Println("连接失败:",err)
//	}
//
//	//建立租约
//
//}
