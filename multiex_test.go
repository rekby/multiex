package multiex

import (
	"os"
	"testing"
)

func testInit() {
	executors = make(map[string]ExecutorDescribe)
}

func TestRegister(t *testing.T) {
	testInit()
	var err error

	empty_func := func() {}
	desc := ExecutorDescribe{Name: "t1", Function: empty_func}

	err = Register(desc)
	if err != nil {
		t.Error("Error while register first function")
	}
	if executors["t1"].Name != desc.Name {
		t.Fatal("Error while save ExecutorDescribe")
	}

	if Register(ExecutorDescribe{Name: "t2", Function: empty_func}) != nil {
		t.Error("Must be pass to register same function many times as different names")
	}
	if Register(ExecutorDescribe{Name: "t2", Function: empty_func}) == nil {
		t.Error("Must be error when register executor with same name second time")
	}

	if Register(ExecutorDescribe{Name: "t3", Function: nil}) == nil {
		t.Error("Must be error when register nil function")
	}

	if Register(ExecutorDescribe{Name: "", Function: empty_func}) == nil {
		t.Error("Must be error when register function with empty name")
	}
}

func TestMultiexMain(t *testing.T) {
	testInit()

	c1 := false // Used in many tests
	t1 := func() { c1 = true }
	Register(ExecutorDescribe{Name: "t1", Function: t1})
	os.Args = []string{"t1"}
	Main()
	if !c1 {
		t.Error("Call t1 function")
	}

	c2 := false // Used in many tests
	t2 := func() { c2 = true }
	Register(ExecutorDescribe{Name: "t2", Function: t2})
	os.Args = []string{"t2"}
	Main()
	if !c2 {
		t.Error("Call t2 function")
	}

	c1 = false
	os.Args = []string{"t1"}
	Main()
	if !c1 {
		t.Error("Call t1 second time")
	}

	c1 = false
	os.Args = []string{"/root/t1"}
	Main()
	if !c1 {
		t.Error("Call t1 by full path")
	}

	c2 = false
	os.Args = []string{"/root/t2"}
	Main()
	if !c2 {
		t.Error("Call t2 by full path")
	}

	c1 = false
	os.Args = []string{"./../asd/t1"}
	Main()
	if !c1 {
		t.Error("Call t1 by relative path")
	}

	c2 = false
	os.Args = []string{"./../asd/t2"}
	Main()
	if !c2 {
		t.Error("Call c2 by relative path")
	}

	c1 = false
	os.Args = []string{"any-path", "--multiex-command=t1"}
	Main()
	if !c1 {
		t.Error("Explicit call t1")
	}

	c2 = false
	os.Args = []string{"any-path", "--multiex-command=t2"}
	Main()
	if !c2 {
		t.Error("Explicit call t2")
	}

	c3 := false
	t3 := func() {
		if len(os.Args) != 3 {
			return
		}
		if os.Args[1] != "a1" {
			return
		}
		if os.Args[2] != "a2" {
			return
		}
		c3 = true
	}
	Register(ExecutorDescribe{Name: "t3", Function: t3})
	os.Args = []string{"t3", "a1", "a2"}
	Main()
	if len(os.Args) != 3 || os.Args[0] != "t3" || os.Args[1] != "a1" || os.Args[2] != "a2" {
		t.Error("os.Args preserve when call function t3")
	}
	if !c3 {
		t.Error("Call t3 with args")
	}
}
