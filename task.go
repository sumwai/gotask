package gotask

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
)

type (
	Task struct {
		Tasks  []Item
		debug  bool
		logger *log.Logger
	}
	Item struct {
		Name string
		Func Func
	}
	Func   func(...any) any
	Params []any
	Exit   struct {
		Code    int
		Message string
		Data    any
	}
)

func (task_params Params) Parse(need_params ...any) error {
	if len(need_params) > len(task_params) {
		return fmt.Errorf("params length mismatch, want: %d, but: %d", len(need_params), len(task_params))
	}
	for i := range need_params {
		pv := reflect.ValueOf(task_params[i])
		paramsv := reflect.ValueOf(need_params[i])
		if paramsv.Type().Kind() != reflect.Pointer {
			return fmt.Errorf("params #" + strconv.Itoa(i) + " must be pointer")
		}
		if paramsv.Elem().Type() != pv.Type() {
			return fmt.Errorf(
				"params #" + strconv.Itoa(i) + " type mismatch" + ":" +
					paramsv.Elem().Type().String() + " != " + pv.Type().String(),
			)
		}
		if !paramsv.Elem().CanSet() {
			return fmt.Errorf("params #" + strconv.Itoa(i) + " cannot set")
		}
		paramsv.Elem().Set(pv)
	}
	return nil
}

func New() *Task {
	logger := log.New(os.Stdout, "[Task] ", log.Ldate|log.Ltime|log.Lmsgprefix)
	return &Task{
		logger: logger,
	}
}

func (t *Task) DebugMode(debug bool) {
	t.debug = debug
}

func (t *Task) AddTask(name string, fn Func) {
	t.Tasks = append(t.Tasks, Item{
		Name: name,
		Func: fn,
	})
}

func (t *Task) Run(params ...any) any {
	var result any
	for _, task := range t.Tasks {
		t.Log("[Start]", task.Name)
		t.Log("[Params]", params)
		result = task.Func(params...)
		if e, ok := result.(Exit); ok {
			t.Log("[Exit]", task.Name, e.Message)
			result = e.Data
			break
		}
		t.Log("[Done]", task.Name)
		t.Log("[Result]", result)
		t.Log()
		if result == nil {
			params = []any{}
			continue
		}
		if p, ok := result.(Params); ok {
			params = p
			continue
		}
		params = []any{result}
	}
	return result
}

func (t *Task) Log(msg ...any) {
	if !t.debug {
		return
	}
	t.logger.Println(msg...)
}
