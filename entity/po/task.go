package po

import (
	"encoding/json"
	"github.com/samber/lo"
	"sort"
	"time"
)

type TaskStatus string

const (
	TaskActionMusic  = "MUSIC"
	TaskActionLyrics = "LYRICS"
)

const (
	TaskStatusNotStart   TaskStatus = "NOT_START"
	TaskStatusSubmitted             = "SUBMITTED"
	TaskStatusQueued                = "QUEUED"
	TaskStatusInProgress            = "IN_PROGRESS"
	TaskStatusFailure               = "FAILURE"
	TaskStatusSuccess               = "SUCCESS"
	TaskStatusUnknown               = "UNKNOWN"
)

func (status TaskStatus) IsDone() bool {
	return status == TaskStatusSuccess ||
		status == TaskStatusFailure ||
		status == TaskStatusUnknown
}

// Task 任务, 用于记录任务的执行情况，用于轮询任务的执行情况
type Task struct {
	BaseModel
	TaskID     string     `json:"task_id" gorm:"type:varchar(50);index"`  // 第三方id，不一定有
	Platform   string     `json:"platform" gorm:"type:varchar(30);index"` // 平台
	Action     string     `json:"action" gorm:"type:varchar(40);index"`   // 任务类型, song, lyrics, description-mode
	Status     TaskStatus `json:"status" gorm:"type:varchar(20);index"`   // 任务状态, submitted, queueing, processing, success, failed
	FailReason string     `json:"fail_reason"`
	SubmitTime int64      `json:"submit_time" gorm:"index"`
	StartTime  int64      `json:"start_time" gorm:"index"`
	FinishTime int64      `json:"finish_time" gorm:"index"`
	//Progress   string     `json:"progress" gorm:"type:varchar(20);index"`
	SearchItem string `json:"search_item" gorm:"type:varchar(100);index"` // 搜索项
	Data       any    `json:"data" gorm:"type:json"`
}

type TaskWithData[T TaskData] struct {
	BaseModel
	TaskID     string     `json:"task_id" gorm:"type:varchar(50);index"`
	Platform   string     `json:"platform" gorm:"type:varchar(30);index"` // 平台
	Action     string     `json:"action" gorm:"type:varchar(40);index"`   // 任务类型, song, lyrics, description-mode
	Status     TaskStatus `json:"status" gorm:"type:varchar(20);index"`   // 任务状态, submitted, queueing, processing, success, failed
	FailReason string     `json:"fail_reason"`
	SubmitTime int64      `json:"submit_time" gorm:"index"`
	StartTime  int64      `json:"start_time" gorm:"index"`
	FinishTime int64      `json:"finish_time" gorm:"index"`
	//Progress   string     `json:"progress" gorm:"type:varchar(20);index"`
	SearchItem string `json:"search_item" gorm:"type:varchar(100);index"` // 搜索项
	Data       T      `json:"data" gorm:"type:json"`
}

func GetTaskByTaskID[T TaskData](uuid string) (*TaskWithData[T], bool, error) {
	var task TaskWithData[T]
	err := DB.Model(&Task{}).
		Where("task_id = ?", uuid).
		First(&task).Error
	exist, err := RecordExist(err)
	return &task, exist, err
}

func GetTaskByTaskIDs[T TaskData](uuid []string) ([]TaskWithData[T], error) {
	var tasks []TaskWithData[T]
	err := DB.Model(&Task{}).
		Where("task_id in (?)", uuid).
		Find(&tasks).Error
	return tasks, err
}

func UpdateTask[T TaskData](id int64, task *TaskWithData[T]) error {
	return DB.Model(&Task{}).
		Where("id = ?", id).
		Updates(task).Error
}

func AddTasks[T TaskData](tasks []TaskWithData[T]) error {
	for i := range tasks {
		if tasks[i].Status == "" {
			tasks[i].Status = TaskStatusNotStart
		}
		if tasks[i].SubmitTime == 0 {
			tasks[i].SubmitTime = time.Now().Unix()
		}
	}
	for _, chunk := range lo.Chunk(tasks, 50) {
		err := DB.Model(&Task{}).Create(&chunk).Error
		if err != nil {
			return err
		}
	}
	return nil
}

type TaskQuery struct {
	BaseQuery
	TaskID    string
	TaskIDs   []string
	Action    string
	Actions   []string
	Status    TaskStatus
	NotStatus TaskStatus
	UnFinish  bool
}

func GetTaskByQuery(query TaskQuery) ([]Task, error) {
	var tasks []Task
	tx := DB.Model(&Task{})

	if query.ID != 0 {
		tx = tx.Where("id = ?", query.ID)
	}
	if len(query.IDs) != 0 {
		tx = tx.Where("id in (?)", query.IDs)
	}
	if query.TaskID != "" {
		tx = tx.Where("task_id = ?", query.TaskID)
	}
	if len(query.TaskIDs) != 0 {
		tx = tx.Where("task_id in (?)", query.TaskIDs)
	}
	if query.Action != "" {
		tx = tx.Where("action = ?", query.Action)
	}
	if len(query.Actions) != 0 {
		tx = tx.Where("action in (?)", query.Actions)
	}
	if query.Status != "" {
		tx = tx.Where("status = ?", query.Status)
	}
	if query.NotStatus != "" {
		tx = tx.Where("status != ?", query.NotStatus)
	}
	if query.UnFinish {
		tx = tx.Where("finish_time <= ?", 0)
	}
	err := tx.Find(&tasks).Error
	return tasks, err
}

func CheckTaskNeedUpdate[T TaskData](oldTask, newTask *TaskWithData[T]) bool {
	if oldTask.Status != newTask.Status {
		return true
	}
	if oldTask.SubmitTime != newTask.SubmitTime {
		return true
	}
	if oldTask.FailReason != newTask.FailReason {
		return true
	}
	if oldTask.StartTime != newTask.StartTime {
		return true
	}
	if oldTask.FinishTime != newTask.FinishTime {
		return true
	}
	/*if oldTask.Progress != newTask.Progress {
		return true
	}*/

	oldData, _ := json.Marshal(oldTask.Data)
	newData, _ := json.Marshal(newTask.Data)

	sort.Slice(oldData, func(i, j int) bool {
		return oldData[i] < oldData[j]
	})
	sort.Slice(newData, func(i, j int) bool {
		return newData[i] < newData[j]
	})

	if string(oldData) != string(newData) {
		return true
	}

	return false
}

func (t *TaskWithData[T]) TaskFailed(reason string) error {
	t.FinishTime = time.Now().Unix()
	t.Status = TaskStatusFailure
	t.FailReason = reason

	return DB.Model(&Task{}).
		Where("id = ?", t.ID).
		Updates(t).Error
}
