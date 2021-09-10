package worker

import (
	"FoG/src/github.com/cl/crontab/common"
	"time"
)

//任务调度
type Scheduler struct {
	// jobMgr 将任务事件通过通道进行推送
	jobEventChan chan *common.JobEvent
	// 将任务信息放入table表
	jobPlanTable map[string]*common.JobSchedulePlan
}

var (
	G_scheduler *Scheduler
)

// 1 初始化调度协程
func InitScheduler() (err error) {
	G_scheduler = &Scheduler{
		jobEventChan: make(chan *common.JobEvent, 1000),
		jobPlanTable: make(map[string]*common.JobSchedulePlan),
	}

	// 2 启动调度协程
	return
}

// 2 启动调度协程  获取任务，调度任务
func (scheduler *Scheduler) scheduleLoop() {
	var (
		jobEvent      *common.JobEvent
		err           error
		schcduleAfter time.Duration // 间隔时间，休息时间，从当前任务中获取最近调度时间点
		scheduleTimer *time.Timer   // 定时器，到期后执行任务
	)

	//调度器的延时定时器
	scheduleTimer = time.NewTimer()
	// 从通道中获取到任务
	for {
		select {
		case jobEvent = <-scheduler.jobEventChan:
			// 获取到任务事件， 新建一个处理方法
			scheduler.handleJobEvent(jobEvent, err)

		}
	}
}

//2.1 处理 任务事件 ---新增，删除
func (scheduler *Scheduler) handleJobEvent(jobEvent *common.JobEvent, err error) {
	var (
		jobPlan *common.JobSchedulePlan
		jobExit bool
	)
	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE:
		// 构建执行任务
		if jobPlan, err = common.BuildJobExecuteInfo(jobEvent.Job); err != nil {
			return
		}
		// 添加到执行计划表
		scheduler.jobPlanTable[jobEvent.Job.Name] = jobPlan
	case common.JOB_EVENT_DELETE:
		// 普通删除事件，从任务表中删除即可
		if jobPlan, jobExit = scheduler.jobPlanTable[jobEvent.Job.Name]; jobExit {
			delete(scheduler.jobPlanTable, jobEvent.Job.Name)
		}
	}
}

// 3尝试获取间隔休息时间
func (scheduler *Scheduler) TrySchedule() (scheduleAfter time.Duration) {

	//任务表为空，则默认睡眠1S
	if len(scheduler.jobPlanTable) == 0 {
		scheduleAfter = G_config.ScheduleSleepTime * time.Second
	}

	return
}
