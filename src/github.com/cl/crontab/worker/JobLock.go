package worker

import (
	"FoG/src/github.com/cl/crontab/common"
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type JobLock struct {
	//基于etcd实现分布式锁
	kv clientv3.KV
	lease clientv3.Lease
	jobName string

	leaseId clientv3.LeaseID  // 租约ID,记录之后，后续取消
	cancelFunc context.CancelFunc // 取消函数 ，后续基于该函数来取消续约
	isLocked bool // 是否上锁成功

}

//1 初始化
func InitJobLock(jobName string,kv clientv3.KV,lease clientv3.Lease)(jobLock *JobLock){
	jobLock = &JobLock{
		kv: kv,
		lease:lease,
		jobName: jobName,
	}
	return
}

//2 尝试上锁 基于 lease + txn实现
func(jobLock *JobLock) TryLock()(err error){
	var(
		leaseGrantResp *clientv3.LeaseGrantResponse
		ctx context.Context
		cancelFunc context.CancelFunc
		leaseId clientv3.LeaseID
		keepRespChan <- chan *clientv3.LeaseKeepAliveResponse
		txn clientv3.Txn
		lockKey string
		txnResp *clientv3.TxnResponse
	)
	ctx,cancelFunc = context.WithCancel(context.TODO())

	// 1 创建租约
	if leaseGrantResp,err = jobLock.lease.Grant(context.TODO(),5);err != nil{
		return
	}
	//2 获取id
	leaseId = leaseGrantResp.ID
	//3 自动续租
	if keepRespChan, err = jobLock.lease.KeepAlive(ctx,leaseId); err != nil{
		return
	}

	//4 处理租约应答
	go func() {
		var(
			keepResp *clientv3.LeaseKeepAliveResponse
		)
		for{
			select{
			case keepResp = <- keepRespChan:
				if keepResp == nil{
					goto END
				}
			}
		}
	END:
		fmt.Println("-------租约到期")
	}()

	// 5 创建事务
	txn = jobLock.kv.Txn(context.TODO())

	// 锁
	lockKey =  common.JOB_LOCK_DIR + jobLock.jobName

	// 尝试获取锁
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey),"=",0)).
		Then(clientv3.OpPut(lockKey,"",clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(lockKey))

	// 提交事务
	if txnResp,err = txn.Commit(); err != nil{
		goto FAIL
	}

	// 6 成功返回，失败就释放租约
	if !txnResp.Succeeded{
		err = common.ERR_LOCK_ALREADY_REQUIRED
		goto FAIL
	}
	// 7 抢锁成功
	jobLock.leaseId = leaseId

    return
 FAIL:
		cancelFunc() // 取消自动续约
		jobLock.lease.Revoke(context.TODO(),leaseId) // 释放租约
		return
}

// 3 释放锁
func(jobLock *JobLock) Unlock(){
	if jobLock.isLocked{
		jobLock.cancelFunc()
		jobLock.lease.Revoke(context.TODO(),jobLock.leaseId)
	}
}