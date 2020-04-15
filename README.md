# gofunc
简单好用的go test工具

#### 使用场景
用于快速测试某个函数或表达式，比如有个函数用于删除slice里的一个元素：
```golang
func DeleteSliceInt(data []int, i int) []int {
    return append(data[:i], data[i+1:]...)
}
```
正常的单元测试需要新建 xxx_test.go 文件，编写 Testxxx 函数，这样操作有时候过于繁琐。

为简化步骤，开发了工具命令 gofunc ，在命令行模式下调用 DeleteSliceInt 实现单元测试，并输出返回值：
```bash
$ gofunc -r 'DeleteSliceInt([]int{0,1,2,3,4,5}, 3)'
>> [0 1 2 4 5]
PASS
ok      0.089s
```

#### 使用方法
- git clone到本地，执行go build，得到exe
- 复制 gofunc 程序到系统bin目录
- ***切换到测试函数所在目录***，执行gofunc命令

```bash
- 若测试函数
- 无参数，无返回值： gofunc yourFuncName
- 有参数，无返回值： gofunc 'yourFuncName(arg1, arg2)'
- 无参数，有返回值： gofunc -r yourFuncName
- 有参数，有返回值： gofunc -r 'yourFuncName(arg1, arg2)'
```

#### 表达式测试
```bash
$ gofunc -r 'time.Now().Format("2006-01-02 15:04:05")'
>> 2020-04-15 14:01:48
$ gofunc -r 'int(math.Floor(float64(59)/60 + 0.5))'
>> 1
$ gofunc -r "strconv.FormatFloat(3.1415926, 'f', 2, 64)"
>> 3.14
$ gofunc -r 'strings.ReplaceAll("who make test easier", "who", "gofunc")'
>> gofunc make test easier

目前仅支持fmt、time、math、strconv、strings包
```
