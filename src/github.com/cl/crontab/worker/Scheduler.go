package worker

import (
	"FoG/src/github.com/cl/crontab/common"
	"fmt"
	"time"
)

//任务调度
type Scheduler struct {
	// jobMgr 将任务事件通过通道进行推送
	jobEventChan chan *common.JobEvent
	// 将任务信息放入table表
	jobPlanTable map[string]*common.JobSchedulePlan
	// 记录执行任务信息
	jobExecutionTable map[string]*common.JobExecuteInfo
	//8 创建获取响应结果的通道
	jobResultChan chan *common.JobExecuteResult
}

var (
	G_scheduler *Scheduler
)

// 1 初始化调度协程
func InitScheduler() (err error) {
	G_scheduler = &Scheduler{
		jobEventChan:      make(chan *common.JobEvent, 1000),
		jobPlanTable:      make(map[string]*common.JobSchedulePlan),
		jobExecutionTable: make(map[string]*common.JobExecuteInfo),
		jobResultChan:     make(chan *common.JobExecuteResult, 1000),
	}

	// 2 启动调度协程
	go G_scheduler.scheduleLoop()

	return
}

// 2 启动调度协程  获取任务，调度任务
func (scheduler *Scheduler) scheduleLoop() {
	var (
		jobEvent      *common.JobEvent
		err           error
		scheduleAfter time.Duration // 间隔时间，休息时间，从当前任务中获取最近调度时间点
		scheduleTimer *time.Timer   // 定时器，到期后执行任务
		jobResult     *common.JobExecuteResult
	)

	// 尝试调度任务--- 会返回可以睡眠的最短时间（下一个任务执行的最短时间）
	scheduleAfter = scheduler.TrySchedule()
	// 调度器的延时定时器
	scheduleTimer = time.NewTimer(scheduleAfter)
	// 从通道中获取到任务
	for {
		select {
		case jobEvent = <-scheduler.jobEventChan:
			// 获取到任务事件， 新建一个处理方法
			scheduler.handleJobEvent(jobEvent, err)
		case <-scheduleTimer.C:
			//最近的任务周期到了，没有到之前会阻塞等待，不会一直空轮询
		case jobResult = <-scheduler.jobResultChan:
			//监听执行结果，并对执行结果进行处理, 将执行完的任务从执行表中移除
			scheduler.HandJobResult(jobResult)
		}
		// 调度任务
		scheduleAfter = scheduler.TrySchedule()
		// 重置调度间隔， 下一次最近任务执行的时间间隔
		scheduleTimer.Reset(scheduleAfter)
	}
}

//2.1 处理 任务事件 ---新增，删除
func (scheduler *Scheduler) handleJobEvent(jobEvent *common.JobEvent, err error) {
	var (
		jobPlan        *common.JobSchedulePlan
		jobExit        bool
		jobExecuteInfo *common.JobExecuteInfo
	)
	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE:
		// 构建执行计划
		if jobPlan, err = common.BuildJobExecutePlan(jobEvent.Job); err != nil {
			return
		}
		// 添加到执行计划表
		scheduler.jobPlanTable[jobEvent.Job.Name] = jobPlan
	case common.JOB_EVENT_DELETE:
		// 普通删除事件，从任务表中删除即可
		if jobPlan, jobExit = scheduler.jobPlanTable[jobEvent.Job.Name]; jobExit {
			delete(scheduler.jobPlanTable, jobEvent.Job.Name)
		}
	case common.JOB_EVENT_KILLER:
		// 强杀事件，取消command的执行-————》 context 的cancel
		// 判断任务是否在执行中
		if jobExecuteInfo, jobExit = scheduler.jobExecutionTable[jobEvent.Job.Name]; jobExit {
			jobExecuteInfo.CancelFunc()
		}
	}
}

// 3 尝试获取间隔休息时间
func (scheduler *Scheduler) TrySchedule() (scheduleAfter time.Duration) {
	var (
		jobPlan  *common.JobSchedulePlan
		now      time.Time
		nearTime *time.Time
	)
	//任务表为空，则默认睡眠1S
	if len(scheduler.jobPlanTable) == 0 {
		scheduleAfter = G_config.ScheduleSleepTime * time.Second
		return
	}
	now = time.Now()
	// 否则的话就遍历当前已有的任务，以最近的时间作为睡眠时间
	for _, jobPlan = range scheduler.jobPlanTable {
		// 下次执行时间如果已经到了，那么就去执行
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			// 到了执行时间，尝试去执行，因为可能任务已经正在执行中
			scheduler.TryStartJob(jobPlan)
			jobPlan.NextTime = jobPlan.Expr.Next(now) // 更新下次执行时间
		}
		/**
		即使当前有任务准备在执行了，那么我们同样会计算此执行中任务的下次执行时间，然后统一进行比较
		最终筛选出最近的时间点
		*/
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
	}
	// 下次调度间隔（最近要执行的任务)
	scheduleAfter = (*nearTime).Sub(now)

	return
}

//6  尝试执行任务   执行时间大于 任务间隔 ,为了防止这样的情况，那么我们就要利用
func (scheduler *Scheduler) TryStartJob(jobPlan *common.JobSchedulePlan) {
	//调度和执行
	var (
		jobExecuteInfo *common.JobExecuteInfo
		jobExecuting   bool
	)
	//如果任务正在执行，那么就跳过本次调度
	if jobExecuteInfo, jobExecuting = scheduler.jobExecutionTable[jobPlan.Job.Name]; jobExecuting {
		fmt.Println(jobPlan.Job.Name, "任务已在执行中...")
		return
	}
	// 构建执行状态
	jobExecuteInfo = common.BuildJobExecuteInfo(jobPlan)

	// 将任务存入到 执行表中
	scheduler.jobExecutionTable[jobPlan.Job.Name] = jobExecuteInfo

	fmt.Println("执行任务:", jobExecuteInfo.Job.Name, jobExecuteInfo.PlanTime, jobExecuteInfo.RealTime)

	// 执行任务--执行模块
	G_executor.ExecuteJob(jobExecuteInfo)
}

// 7 获取任务执行结果
func (scheduler *Scheduler) PushJobResult(jobResult *common.JobExecuteResult) {
	scheduler.jobResultChan <- jobResult
}

//8 对任务执行结果进行处理
func (scheduler *Scheduler) HandJobResult(result *common.JobExecuteResult) {
	//从执行表中删除执行完的任务
	delete(scheduler.jobExecutionTable, result.ExecuteInfo.Job.Name)
}

// 封装推送任务变化事件
func (scheduler *Scheduler) PushJobEvent(jobEvent *common.JobEvent) {
	scheduler.jobEventChan <- jobEvent
}
