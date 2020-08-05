package go_utils

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

func ExecCommand(shell string, writers ...io.Writer) error {

	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/sh", "-c", shell)

	//StdoutPipe方法返回一个在命令Start后与命令标准输出关联的管道。Wait方法获知命令结束后会关闭这个管道，一般不需要显式的关闭该管道。
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return fmt.Errorf("执行sh出错 %s", err)
	}

	cmd.Start()
	//创建一个流来读取管道内内容，这里逻辑是通过一行一行的读取的
	reader := bufio.NewReader(stdout)

	//实时循环读取输出流中的一行内容
	for {
		line, _, err2 := reader.ReadLine()
		if err2 != nil || io.EOF == err2 {
			break
		}
		out := string(line)
		for _, w := range writers {
			fmt.Fprint(w, out)
		}
	}

	//阻塞直到该命令执行完成，该命令必须是被Start方法开始执行的
	err = cmd.Wait()
	return err
}
