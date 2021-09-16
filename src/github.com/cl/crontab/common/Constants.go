package common

const (
	//保存任务的目录
	JOB_SAVE_DIR = "/cron/jobs/"

	// 保存任务事件
	JOB_EVENT_SAVE = 1

	// 删除任务事件
	JOB_EVENT_DELETE = 2

	// 任务强杀事件
	JOB_EVENT_KILLER = 3

	// 强杀事件
	JOB_KILLER_DIR = "/cron/killer/"

	// 锁
	JOB_LOCK_DIR = "/cron/lock/"
)
