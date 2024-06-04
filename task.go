package main

import (
	"suno-api/common"
	"suno-api/entity/po"
	"time"
)

var queue = make(chan Task, 100)

type Task struct {
	Action string
	ID     string
}

func AddTask(task Task) {
	queue <- task
}

func startTaskWorker() {
	common.SafeGoroutine(func() {
		for task := range queue {
			switch task.Action {
			case po.TaskActionMusic:
				t := task
				safeRun(func() {
					common.LoopTask(func() (bool, *common.RelayError) {
						done, relayErr := loopFetchTask(t.ID)
						if relayErr != nil {
							common.Logger.Errorw("loopFetchTask error", "err", relayErr)
							return false, relayErr
						}
						if done {
							common.Logger.Infow("任务完成", "task", t.ID)
							return true, nil
						}
						return false, nil
					}, 5*time.Second, 10*time.Minute, true)
				})
			}
		}
	})
}

func safeRun(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			common.Logger.Errorw("safeRun panic", "err", err)
		}
	}()

	fn()
}

func recoverTasks() {
	tasks, err := po.GetTaskByQuery(po.TaskQuery{
		UnFinish: true,
	})
	if err != nil {
		common.Logger.Fatalw("recoverTasks error", "err", err)
	}
	for _, t := range tasks {
		AddTask(Task{
			Action: t.Action,
			ID:     t.TaskID,
		})
	}

}
