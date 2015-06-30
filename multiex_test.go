package multiex

import (
	"os"
	"path/filepath"
	"testing"
)

func testInit() {
	executors = make(map[string]ExecutorDescribe)
	SetInstallPrefix("")
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

func TestInstall(t *testing.T) {
	testInit()
	dir_path, _ := filepath.Abs(os.Args[0])
	dir_path = filepath.Dir(dir_path)
	os.Args = []string{os.Args[0], "install"}
	defer func() { os.Args = []string{os.Args[0]} }()

	SetInstallPrefix("testprefix_")

	f := func() {}
	Register(ExecutorDescribe{Name: "test-1", Function: f})
	defer func() { os.Remove(filepath.Join(dir_path, "testprefix_test-1")) }()
	MultiExUtilsMain()
	if _, err := os.Stat(filepath.Join(dir_path, "testprefix_test-1")); os.IsNotExist(err) {
		t.Error("Create simple link")
	}

	Register(ExecutorDescribe{Name: "test-2", Function: f})
	Register(ExecutorDescribe{Name: "test-3", Function: f})
	defer func() { os.Remove(filepath.Join(dir_path, "testprefix_test-21")); os.Remove(filepath.Join(dir_path, "testprefix_test-3")) }()
	MultiExUtilsMain()
	if _, err := os.Stat(filepath.Join(dir_path, "testprefix_test-2")); os.IsNotExist(err) {
		t.Error("Create second link after first")
	}
	if _, err := os.Stat(filepath.Join(dir_path, "testprefix_test-3")); os.IsNotExist(err) {
		t.Error("Create third link with second in same time")
	}

	Register(ExecutorDescribe{Name: "test-noninstall", Function: f, NoInstall: true})
	defer func() { os.Remove(filepath.Join(dir_path, "testprefix_test-noninstall")) }()
	MultiExUtilsMain()
	if _, err := os.Stat(filepath.Join(dir_path, "testprefix_test-noninstall")); os.IsExist(err) {
		t.Error("Don't create executor with NoInstall flag")
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
	os.Args = []string{"any-path", "t1"}
	Main()
	if !c1 {
		t.Error("Explicit call t1")
	}

	c2 = false
	os.Args = []string{"any-path", "t2"}
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

	c4 := false
	t4 := func() {
		if len(os.Args) != 3 {
			return
		}
		if os.Args[0] != "any-path" {
			return
		}
		if os.Args[1] != "a1" {
			return
		}
		if os.Args[2] != "a2" {
			return
		}
		c4 = true
	}
	Register(ExecutorDescribe{Name: "t4", Function: t4})
	os.Args = []string{"any-path", "t4", "a1", "a2"}
	Main()
	if len(os.Args) != 4 || os.Args[0] != "any-path" || os.Args[1] != "t4" || os.Args[2] != "a1" || os.Args[3] != "a2" {
		t.Error("os.Args preserve when call function t4 by explicit name")
	}
	if !c4 {
		t.Error("Call t4 with args by explicit name")
	}

}
