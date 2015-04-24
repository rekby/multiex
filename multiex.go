package multiex

import (
	"errors"
	"os"
	"path"
	"reflect"
)

// Struct for describe included commands
// ATTENTION: don't use short initialization form of structure (without field names) - it can have additional fields in
// the feauture.

// Структура для описания подключаемой команды
// ВНИМАНИЕ: не используйте короткую форму инициализации структуры (без указания имен полей) - структура может
// дополняться в будущем и использование короткой формы сломает код при обновлении библиотеки
type ExecutorDescribe struct {
	Function  func() // Pointer for function. Указатель на функцию.
	Name      string // Name os command what was called. Имя команды, по которому функция будет вызываться.
	NoInstall bool   // Don't create symlink when multiex install called. Не создавать ссылку при вызове команды установки.
}

// Task for function registration
// Описание задачи для регистрации функции
type registerQueueTask struct {
	module ExecutorDescribe
	err    chan error
}

// Golang hasn't explicit gurantee about sequense exec of init function. Create chan for gurantee of registration.
// Golang не гарантирует строго последовательного вызова функции init. Создаем канал для гарантии поочередной регистрации функций.
var registerQueue = make(chan registerQueueTask)

// Map of funcctions: name - describe
// Карта функций: имя - описание
var executors = make(map[string]ExecutorDescribe)

// Call executor-function by description
// Вызвать функцию-исполнитель по описанию
func callWorker(module ExecutorDescribe) {
	f := reflect.ValueOf(module.Function)
	f.Call([]reflect.Value{})
}

// Start goroutines for register workers.
// Запускает горутины для регистрации обработчиков
func startRegisterWorker() {
	var registerTask registerQueueTask

	// infinite loop for add modules from queue
	// бесконечный цикл для добавления функций из очереди
	go func() {
		// This defer function ensure infinite goroutine for sequence registration of workers
		// Эта отложенная функция гарантирует непрерывное выполнение горутины для последовательной регистрации обработчиков

		defer func() {
			if err := recover(); err != nil && registerTask.err != nil {
				registerTask.err <- err.(error)
			}
			startRegisterWorker() // restore infinite loop. Восстановить цикл обработки
		}()

		for {
			registerTask = <-registerQueue
			if _, has := executors[registerTask.module.Name]; has {
				registerTask.err <- errors.New("Module has registered already (модуль уже зарегистрирован): '" + registerTask.module.Name + "'")
			} else {
				executors[registerTask.module.Name] = registerTask.module
				registerTask.err <- nil
			}
		}
	}()
}

// Register function as module worker
// Зарегистрировать функцию как обработчик
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
// Функция для вызова из выполняемого файла для распознавания и вызова итогового обработчика
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
	startRegisterWorker()
}
