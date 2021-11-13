package models

import (
	"fmt"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	"os/exec"
	"runtime"
)

// TestReadable 测试操作系统是否为可读操作系统
func TestReadable() {
	system := runtime.GOOS
	switch system {
	case "linux":
		{
			//检查设备是否可写入，若不可写入则重启设备
			touch := exec.Command("touch", "-a", "/opt/readonly_test")
			reboot := exec.Command("reboot")
			err := touch.Run()
			if err != nil {
				mid.Log.Info(fmt.Sprintf("可读性异常，准备重启:%v", err))
				err = reboot.Run()
				if err != nil {
					mid.Log.Info(fmt.Sprintf("err:%v", err))
				}
			} else {
				rm := exec.Command("rm", "-rf", "/opt/readonly_test")
				err = rm.Run()
				mid.Log.Info("系统无异常......")
				if err != nil {
					mid.Log.Info(fmt.Sprintf("err:%v", err))
				}
			}
		}
	case "windows":
		{
			mid.Log.Info("WIN不需要开机可读性检测.....")
		}

	}
}
