package multiex

import (
	"fmt"
	"os"
    "path/filepath")

func init() {
	Register(ExecutorDescribe{Name: "multiex", Function: MultiExUtilsMain})
}

func MultiExUtilsMain() {
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
        link_path := filepath.Join(dirname, module.Name)
        err = os.Symlink(basename, link_path)
        if os.IsExist(err) {
            fmt.Printf("File exists: %s\n", link_path)
        } else if err != nil {
            fmt.Printf("Error: '%s' while try to create link '%s'\n", err, link_path)
        }
    }
}

func printModules() {
	fmt.Println("List of commands:")
	for key, _ := range executors {
		fmt.Println(key)
	}
}

func printUsage() {
	fmt.Println(`
Usages:
    multiex [--multiex-command=...] args

    Multiex contain multiple independent commands into one executable file - for reduce size of many count small utilities
    with same golang runtime.

    When multiex used - one binary executable file can contain many independent utilities with own rules of work.
    Multiex detect right utility by basename of executable command - as busybox.

    Where multiex may be name of file, symlink or hardlink with filename equals to command name. It is usual case and
    you can omit multiex-command parameter.

    If param --multiex-command exist - name of executed file is ignored.

    Exec
    multiex --multiex-command=multiex install
    for create simlinks for all internal commands
	`)
}
