package multiex

import (
	"errors"
	"os"
	"path"
	"reflect"
)

// Struct for describe included commands
// ATTENTION: don't use short initialization form of structure - it can have additional fields in the feauture.
type ExecutorDescribe struct {
	Function  func() // Pointer for function
	Name      string // Name os command what was called
	NoInstall bool   // Don't create symlink when multiex install called
}

type registerQueueTask struct {
	module ExecutorDescribe
	err    chan error
}

var registerQueue = make(chan registerQueueTask)
var executors = make(map[string]ExecutorDescribe)

func callWorker(module ExecutorDescribe) {
	f := reflect.ValueOf(module.Function)
	f.Call([]reflect.Value{})
}

func registerWorker() {
	var registerTask registerQueueTask

	// This is lifetime goroutine
	defer func() {
		if registerTask.err != nil {
			err := recover()
			registerTask.err <- err.(error)
		}
		registerWorker() // restore infinite loop
	}()

	// infinite loop for add modules from queue
	for {
		registerTask = <-registerQueue
		if _, has := executors[registerTask.module.Name]; has {
			registerTask.err <- errors.New("Module has registered already: '" + registerTask.module.Name + "'")
		} else {
			executors[registerTask.module.Name] = registerTask.module
			registerTask.err <- nil
		}
	}
}

// Register function as module worker
func Register(module ExecutorDescribe) error {
	if module.Function == nil {
		return errors.New("Function can't be nil")
	}
	if module.Name == "" {
		return errors.New("Name can't be empty")
	}
	errChan := make(chan error)
	registerQueue <- registerQueueTask{module: module, err: errChan}
	return <-errChan
}

// Function to call from executable module for start process dispatch call
func Main() {
	var commandName string

	commandName = os.Args[0]
	commandName = path.Base(commandName)
	module, has := executors[commandName]
	if !has {
		// check explicit command name in first argument
		if len(os.Args) > 1 {
			// restore original args while exit
			oldArgs := make([]string, len(os.Args))
			copy(oldArgs, os.Args)
			defer func() { os.Args = oldArgs }()

			commandName = os.Args[1]
			os.Args = os.Args[1:]
			os.Args[0] = oldArgs[0] // preserve path
		} else {
		}

		module, has = executors[commandName]
	}

	if has {
		callWorker(module)
	} else {
		printUsage()
		printModules()
		return
	}
}

func init() {
	go registerWorker()
}
