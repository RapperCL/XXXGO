package worker

import (
	"FoG/src/github.com/cl/crontab/common"
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type JobMgr struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

var (
	G_jobMgr *JobMgr
)

// 初始化任务管理器
func InitJobMgr() (err error) {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		lease   clientv3.Lease
		watcher clientv3.Watcher
	)
	config = clientv3.Config{
		Endpoints:   G_config.EtcdEndpoints,
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		return
	}

	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	watcher = clientv3.NewWatcher(client)

	// 赋值单例
	G_jobMgr = &JobMgr{
		client:  client,
		kv:      kv,
		lease:   lease,
		watcher: watcher,
	}

	// 启动任务监听
	G_jobMgr.watchJobs()

	return
}

// 监听任务变化
func (jobMgr *JobMgr) watchJobs() (err error) {

	var (
		getResp            *clientv3.GetResponse
		job                *common.Job
		kvPair             *mvccpb.KeyValue
		jobEvent           *common.JobEvent
		watchStartRevision int64
		watchChan          clientv3.WatchChan
		watchResp          clientv3.WatchResponse
		watchEvent         *clientv3.Event
		jobName            string
	)
	// 1 获取当前任务目录下的所有任务信息
	if getResp, err = jobMgr.kv.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix()); err != nil {
		return
	}
	// 2 任务信息
	for _, kvPair = range getResp.Kvs {
		if job, err = common.UnpackJob(kvPair.Value); err != nil {
			//转换任务结构体，然后推送给任务调度器
			jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
			// 同步给scheduler job
		}
	}

	// 2 从当前revision向后监听变化事件
	// 监听协程
	go func() {
		// 从当前的下一个版本开始监听变化
		watchStartRevision = getResp.Header.Revision + 1
		// 监听 前缀的后续变化
		watchChan = jobMgr.watcher.Watch(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT: // 任务保存事件
					if job, err = common.UnpackJob(watchEvent.Kv.Value); err != nil {
						continue
					}
					// 构建一个更新Event
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)

				case mvccpb.DELETE:
					jobName = common.ExtractJobName(string(watchEvent.Kv.Key))

					job = &common.Job{Name: jobName}

					// 构建一个删除 事件
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_DELETE, job)
				}
				// 将变化的jobEvent 推送给调度处理器
			}
		}
	}()
	return
}
