package tasks

import (
	"github.com/google/uuid"
	"gotasks/json"
	"time"
)

type ArgsMap map[string]interface{}

type Task struct {
	ID                  string    `json:"task_id"`               //任务id
	CreatedAt           string    `json:"created_at"`            //创建时间
	UpdateAt            string    `json:"update_at"`             //更新时间
	QueueName           string    `json:"queue_name"`            //队列名称
	TaskName            string    `json:"task_name"`             //任务名称(worker端要执行的方法名称)
	ArgsMap             []ArgsMap `json:"args_map"`              //任务完成后返回的参数
	CurrentHandlerIndex int       `json:"current_handler_index"` //任务执行到的handler
	OriginalArgsMap     ArgsMap   `json:"original_args_map"`     //原始参数
	ResultLog           string    `json:"result_log"`            //错误日志
	PanicLog            string    `json:"err_log"`               //异常日志
}

func ToArgsMap(v interface{}) (ArgsMap, error) {
	b, err := json.Json.Marshal(v)
	if err != nil {
		return nil, err
	}

	arg := ArgsMap{}
	err = json.Json.Unmarshal(b, &arg)

	return arg, err
}

func ArgsMapTo(args ArgsMap, v interface{}) error {
	b, err := json.Json.Marshal(args)
	if err != nil {
		return err
	}

	return json.Json.Unmarshal(b, v)
}

func NewTask(queueName string, taskName string, args ArgsMap) *Task {
	now := time.Now().Format("2006-01-02 15:04:05")
	return &Task{
		ID:                  uuid.New().String(),
		CreatedAt:           now,
		UpdateAt:            now,
		QueueName:           queueName,
		TaskName:            taskName,
		ArgsMap:             []ArgsMap{},
		CurrentHandlerIndex: 0,
		OriginalArgsMap:     args,
		ResultLog:           "",
		PanicLog:            "",
	}
}
