package main

import (
	"bytes"
	"log"
	"sync"
)

var buffers = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func GetBuffer() *bytes.Buffer {
	// 下面.(* 这是什么写法
	return buffers.Get().(*bytes.Buffer)
}

func PutBuffer(buf *bytes.Buffer) {
	// 转换缓存池，变为可写？
	buf.Reset()
	buffers.Put(buf)
}
func get(gp *sync.WaitGroup) {
	defer gp.Done()
	buffers.Get()
}

func main() {
	//var gp  sync.WaitGroup
	//gp.Add(10)
	//for i :=10;i>0;i--{
	//	go get(&gp)
	//}
	//gp.Wait()
	//by := buffers.Get();
	//fmt.Println("从缓存中取出数据:",by);
	//b := new(bytes.Buffer)
	//buffers.Put(b)
	te()

}

func te() {
	var pipe = &sync.Pool{New: func() interface{} { return "i am best" }}

	val := "i am best"

	pipe.Put(val)

	log.Println("取值", pipe.Get())

	log.Println("再次取值", pipe.Get())
}
