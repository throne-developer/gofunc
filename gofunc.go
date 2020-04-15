package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

/* 单个函数测试工具，节省编写 xxx_test.go 文件的时间
在被执行函数所在目录执行： gofunc yourTestFunc
若有输入参数，执行： gofunc 'yourTestFunc(a, b, ...)'
若需要打印返回值，添加 -r 参数
*/
var (
	help        = flag.Bool("h", false, "print usage")
	printReturn = flag.Bool("r", false, "print return values")
)

func main() {
	flag.Parse()
	if *help || len(flag.Args()) != 1 {
		usage()
		return
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("get current directory err: ", err.Error())
		return
	}

	funcname := flag.Args()[0]
	if !strings.Contains(funcname, "(") {
		funcname += "()"
	}

	pkgName := getPackageName(dir)
	if pkgName == "" {
		fmt.Println("cannot found golang package name")
		return
	}

	testFuncName := "TestGoFunc"
	content := genTestFileContent(pkgName, testFuncName, funcname, *printReturn)
	tempFile, writeErr := writeTempTestFile(content, dir)
	defer deleteFile(tempFile)

	if writeErr != nil {
		fmt.Println("writeTempTestFile err: ", err.Error())
		return
	}

	executeGoTest(testFuncName)
}

func usage() {
	fmt.Println("usage:  gofunc yourTestFunc")
	fmt.Println("        gofunc 'yourTestFunc(arg1, arg2)'")
	fmt.Println("        gofunc -r yourTestFunc")
	fmt.Println("        gofunc -r 'yourTestFunc(arg1, arg2)'")
	fmt.Println()
	flag.PrintDefaults()
}

func getPackageName(dir string) string { /*读取任意一个go文件的包名*/
	pkgName := ""
	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".go") {
			if name, found := getPackageNameFromGoFile(path); found {
				pkgName = name
				return errors.New("ok")
			}
		}
		return nil
	})
	return pkgName
}

func getPackageNameFromGoFile(filePath string) (name string, found bool) { /*读取go文件的包名*/
	lines, err := loadFile(filePath)
	if err != nil {
		return
	}

	for _, line := range lines {
		if strings.HasPrefix(line, "package") {
			name = line
			found = true
			break
		}
	}
	return
}

func genTestFileContent(packageName, testFuncName, funcname string, printResult bool) string {
	templateContent := `%s
import (
	"fmt"
	"testing"
	"math"
	"strconv"
	"strings"
	"time"
)
func %s(t *testing.T) {
	if false {
		time.Now()
		math.Floor(0.0)
		strconv.Itoa(1)
		strings.TrimSpace("")
	}
	
	fmt.Print(">> ")
	%s
}`

	if printResult {
		funcname = fmt.Sprintf(`fmt.Println(%s)`, funcname)
	}
	return fmt.Sprintf(templateContent, packageName, testFuncName, funcname)
}

func writeTempTestFile(content string, dir string) (filePath string, writeErr error) {
	filePath = fmt.Sprintf(`%s\temp_gofunc_test.go`, dir)
	writeErr = writeFileWithLines(filePath, []string{content})
	return
}

func executeGoTest(testFuncName string) {
	cmd := exec.Command("go", "test", "-test.run", testFuncName)
	var output, outErr bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &outErr
	e := cmd.Run()
	if e != nil {
		fmt.Println("run error :" + e.Error())
		fmt.Println(output.String())
		fmt.Println(outErr.String())
	} else {
		fmt.Println(output.String())
	}
}
