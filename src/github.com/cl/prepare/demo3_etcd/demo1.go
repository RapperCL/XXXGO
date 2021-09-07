package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		err    error
		putR   *clientv3.PutResponse
		getR   *clientv3.GetResponse
		delR   *clientv3.DeleteResponse
	)

	fmt.Printf("开始连接")
	config = clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	kv = clientv3.NewKV(client)

	if putR, err = kv.Put(context.TODO(), "te", "v1", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("K,V的id：" + string(putR.Header.Revision))
		if putR.PrevKv != nil {
			fmt.Println("K,V值：" + string(putR.PrevKv.Value))
		}

	}

	// 获取
	if getR, err = kv.Get(context.TODO(), "te"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("te的值:", string(getR.Kvs[0].Value))
	}

	//删除
	if delR, err = kv.Delete(context.TODO(), "te"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("删除：", delR.Header.Revision)
	}

}
