package multiex

import (
	"fmt"
	"os"
	"path/filepath"
)

func init() {
	Register(ExecutorDescribe{Name: "multiex", Function: MultiExUtilsMain, NoInstall: true})
}

func MultiExUtilsMain() {
	if len(os.Args) > 1 && os.Args[1] == "multiex" {
		oldName := os.Args[0]
		os.Args = os.Args[1:]
		os.Args[0] = oldName
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "ls":
			printModules()
		case "help", "usage":
			printUsage()
		case "install":
			createSymlinks()
		}
	} else {
		printUsage()
		printModules()
	}
}

func createSymlinks() {
	binary_path, err := filepath.Abs(os.Args[0])
	if err != nil {
		panic(err)
	}
	binary_path, err = filepath.EvalSymlinks(binary_path)
	if err != nil {
		panic(err)
	}
	basename := filepath.Base(binary_path)
	dirname := filepath.Dir(binary_path)
	for _, module := range executors {
		if module.NoInstall {
			fmt.Printf("Skip create link by func definition (создание ссылки пропущено согласно настроек модуля): '%s'\n", module.Name)
			continue
		}
		link_path := filepath.Join(dirname, module.Name)
		err = os.Symlink(basename, link_path)
		if os.IsExist(err) {
			fmt.Printf("File exists (файл существует): %s\n", link_path)
		} else if err != nil {
			fmt.Printf("Error (ошибка): '%s' while try to create link (при попытке создания ссылки) '%s'\n", err, link_path)
		}
	}
}

func printModules() {
	fmt.Println("List of commands (список команд):")
	for key, _ := range executors {
		fmt.Println(key)
	}
}

func printUsage() {
	fmt.Println(`
Usages:
    multiex [command] args

    Multiex contain multiple independent commands into one executable file - for reduce size of many count small utilities
    with same golang runtime.

    When multiex used - one binary executable file can contain many independent utilities with own rules of work.
    Multiex detect right utility by basename of executable command - as busybox.

    Where multiex may be name of file, symlink or hardlink with filename equals to command name.

    If binary haven't command with same name as call - first argument try usage as command name

    Exec
    ./multiex multiex install
    for create simlinks for all internal commands

    where ./multiex - any name of binary/link doesn't equal of any registered functions or "multiex".
    multiex - name of included worker.

    Other commands for multiex module: ls, help, usage

Использование:
	multies [команда] аргументы

	Multiex позволяет включить множество независимых команд в один исполняемый файл - для сокращения занимаемого места
	большим количеством мелких утилит, разделяющих одно и то же окружение golang.

	Когда используется multiex - один исполняемый файл включает в себя независимые команды с собственными правилами работы.
	Правильная команда для выполнения определяется по имени файла, с которым multiex был вызван - по аналогии с busybox.

	Это имя может быть именем бинарного файла, жесткой или символьной ссылкой, имя которых совпадет с зарегистрированным
	именем функции.

	Если имя выполняемого файла не совпадает ни с одной из зарегистрированных функций - именем функции считается значение
	первого аргумента.
	Если и по этому значению обработчик не найден - выводится справка и список зарегистрированных команд.

	Для создания символьных ссылок в той же папке, где лежит выполняемый файл можно выполнить команду
	./multiex multiex install

	Где ./multiex - любое имя бинарника/ссылки, не совпадающее с зарегистрированными функциями или "multiex".
	multiex - имя встроенного обработчика
	install - команда внутреннему обработчику для создания ссылок

	Другие команды встроенного модуля: ls, help, usage
	`)
}
