package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

/**
作用： 实现指定方法在实例的生命周期内只会被执行一次，常用在单例中。
那么其实我们有另外的方法来实现：通过单例模式（加锁实现）
*/

func first() {
	fmt.Println("第一个方法实现")
}
func second() {
	fmt.Println("第二个方法实现")
}

/**
  通过下面的方法，我们可以发现Do方法可以被多次调用，但只会执行一次对应的方法；

  那么它们是如何实现的呢？
具体的实现也很简单，按照我们前面的思路，通过加锁实现，但是每次加锁，会比较影响性能，
于是为了防止每次加，Once内部维护了一个数据结构
type Once struct {
	done uint32
	m    Mutex
}
 定义一个变量来实现： 执行时，通过判断done的值是否等于0，如果等于0，那么就代表当前方法没有被执行，于是就会执行f
  不等于0，就执行返回，不执行函数f。
通过这种方式就可以避免多次加锁了

*/
func main() {
	var once sync.Once
	once.Do(first)
	once.Do(second)
	fmt.Println("方法执行完成")

	// 利用锁，实现方法的单次执行
	conn := getConn()
	if conn == nil {
		panic("conn is nil")
	}
	fmt.Println("connection:", conn)
}

/**
例子，利用互斥锁来保证线程安全
*/
var connMu sync.Mutex
var conn net.Conn

func getConn() net.Conn {
	connMu.Lock()
	defer connMu.Unlock()

	// 返回已创建好的连接
	if conn != nil {
		return conn
	}

	// 创建连接
	conn, _ = net.DialTimeout("tcp", "baidu.com:80", 10*time.Second)

	return conn
}
