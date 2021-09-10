package master

import (
	"FoG/src/github.com/cl/crontab/common"
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

/**
etcd操作类
*/
type JobMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	//单例
	G_jobMgr *JobMgr
)

//初始化管理器
func InitJobMgr() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		lease  clientv3.Lease
	)
	config = clientv3.Config{
		Endpoints:   G_config.EtcdEndpoints,
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Second,
	}

	//建立客户端
	if client, err = clientv3.New(config); err != nil {
		fmt.Println("建立连接失败", err)
		return
	}

	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	//赋值单例
	G_jobMgr = &JobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	return
}

// 存储任务的方法
func (jobMgr *JobMgr) SaveJob(job *common.Job) (oldJob *common.Job, err error) {
	var (
		jobKey    string
		jobValue  []byte
		putResp   *clientv3.PutResponse
		oldJobObj common.Job
	)
	jobKey = "/cron/jobs/" + job.Name
	// 存入etcd时，需要将job序列化
	if jobValue, err = json.Marshal(job); err != nil {
		return
	}
	//保存到etcd
	if putResp, err = jobMgr.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		return
	}

	if putResp.PrevKv != nil {
		// 对旧值做反序列化
		if err = json.Unmarshal(putResp.PrevKv.Value, &oldJobObj); err != nil {
			err = nil
			return
		}
		oldJob = &oldJobObj
	}
	return
}

// 查询任务的方法
func (jobMgr *JobMgr) ListJobs() (jobList []*common.Job, err error) {
	var (
		dirKey  string
		getResp *clientv3.GetResponse
		kvPair  *mvccpb.KeyValue
		job     *common.Job
	)

	// 获取任务目录,  将任务目录枚举
	dirKey = common.JOB_SAVE_DIR

	// 获取目录下的所有任务信息
	if getResp, err = jobMgr.kv.Get(context.TODO(), dirKey, clientv3.WithPrefix()); err != nil {
		goto ERR
	}

	//初始化数组空间, 指针类型
	jobList = make([]*common.Job, 0)

	//遍历所有的查询结果，并进行反序列化
	for _, kvPair = range getResp.Kvs {
		job = &common.Job{}
		if err = json.Unmarshal(kvPair.Value, job); err != nil {
			err = nil
			continue
		}
		// 指针重新赋值
		jobList = append(jobList, job)
	}

	return
ERR:
	fmt.Println("获取任务错误;", err)
	return
}

// 删除服务, 删除对应的服务，并返回
func (jobMgr *JobMgr) DelJob(jobName string) (jobList []*common.Job, err error) {
	// 获取要删除的key
	var (
		delKey  string
		delResp *clientv3.DeleteResponse
		kvPair  *mvccpb.KeyValue
		job     *common.Job
	)
	delKey = common.JOB_SAVE_DIR + jobName
	if delResp, err = jobMgr.kv.Delete(context.TODO(), delKey, clientv3.WithPrevKV()); err != nil {
		return
	}
	//获取响应
	for _, kvPair = range delResp.PrevKvs {
		job = &common.Job{}
		// 反序列化
		if err = json.Unmarshal(kvPair.Value, &job); err != nil {
			fmt.Println("返回删除的数据时，序列化失败", kvPair.Value)
			err = nil
			continue
		}
		jobList = append(jobList, job)
	}
	return
}

// 强杀服务,给对应的key设置1s的租约，让其过期
// 强杀任务，结束掉当前运行的任务
func (jobMgr *JobMgr) KillJob(jobName string) (err error) {
	var (
		leaseResp *clientv3.LeaseGrantResponse
		leaseId   clientv3.LeaseID
		killKey   string
	)

	killKey = common.JOB_KILLER_DIR + jobName
	// 设置1s的租约
	if leaseResp, err = jobMgr.lease.Grant(context.TODO(), 1); err != nil {
		return
	}
	leaseId = leaseResp.ID
	// 给对应的key绑定 此租约
	if _, err = jobMgr.kv.Put(context.TODO(), killKey, "", clientv3.WithLease(leaseId)); err != nil {
		return
	}
	return
}
