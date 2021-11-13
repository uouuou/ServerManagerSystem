package models

import (
	"bufio"
	"fmt"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	"io"
	"net/http"
	"os/exec"
	"strings"
)

// CheckCommandExists 检查命令是否存在
func CheckCommandExists(command string) bool {
	if _, err := exec.LookPath(command); err != nil {
		return false
	}
	return true
}

// RunWebShell 运行网上的脚本
func RunWebShell(webShellPath string) {
	if !strings.HasPrefix(webShellPath, "http") && !strings.HasPrefix(webShellPath, "https") {
		fmt.Printf("shell path must start with http or https!")
		return
	}
	resp, err := http.Get(webShellPath)
	if err != nil {
		mid.Log.Error(fmt.Sprintf("%v err: %v", mid.RunFuncName(), err))
	}
	defer resp.Body.Close()
	installShell, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	_ = ExecCommand(string(installShell))
}

// ExecCommand 运行命令并实时查看运行结果
func ExecCommand(command string) error {
	cmd := exec.Command("bash", "-c", command)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		mid.Log.Error(fmt.Sprintf("%v Error:The command is err: %v", mid.RunFuncName(), err))
		return err
	}
	ch := make(chan string, 100)
	stdoutScan := bufio.NewScanner(stdout)
	stderrScan := bufio.NewScanner(stderr)
	go func() {
		for stdoutScan.Scan() {
			line := stdoutScan.Text()
			ch <- line
		}
	}()
	go func() {
		for stderrScan.Scan() {
			line := stderrScan.Text()
			ch <- line
		}
	}()
	var err error
	go func() {
		err = cmd.Wait()
		if err != nil && !strings.Contains(err.Error(), "exit status") {
			mid.Log.Error(fmt.Sprintf("%v wait:%v", mid.RunFuncName(), err))
		}
		close(ch)
	}()
	for line := range ch {
		fmt.Println(line)
	}
	return err
}

// ExecCommandWithResult 运行命令并获取结果(传入ERR)
func ExecCommandWithResult(command string) (res string, err error) {
	out, err := exec.Command("bash", "-c", command).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		mid.Log.Error(fmt.Sprintf("%v Command:%v err:%v", mid.RunFuncName(), command, err))
		return string(out), err
	}
	return string(out), err
}

// ExecCommandNoErr 运行命令并获取结果
func ExecCommandNoErr(command string) bool {
	err := exec.Command("bash", "-c", command).Run()
	if err != nil {
		return false
	}
	return true
}
