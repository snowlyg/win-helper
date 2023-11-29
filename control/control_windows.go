//go:build windows
// +build windows

package control

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/v3/process"
	"github.com/snowlyg/helper/dir"
	helper "github.com/snowlyg/win-helper"
)

func Install(srvName, execPath, displayName, systemName, pwd string) error {
	status, err := helper.ServiceStatus(srvName)
	if err != nil {
		return fmt.Errorf("get error msg %w", err)
	}

	if status == helper.StatusRunning {
		return fmt.Errorf("%s is running", srvName)
	}

	if status == helper.StatusUninstall {
		return helper.ServiceInstall(srvName, execPath, displayName, systemName, pwd)
	}

	return nil
}

func Status(srvName string) (helper.Status, error) {
	status, err := helper.ServiceStatus(srvName)
	if err != nil {
		return helper.StatusUnknown, fmt.Errorf("get service status  %w", err)
	}
	return status, nil
}

func ProcessId(srvName string) (uint32, error) {
	processId, err := helper.ServiceProcessId(srvName)
	if err != nil {
		return processId, err
	}
	return processId, nil
}

func Stop(srvName string) error {
	status, err := helper.ServiceStatus(srvName)
	if err != nil {
		return fmt.Errorf("get error msg %w", err)
	}

	if status != helper.StatusRunning {
		return nil
	}

	restop := 3
	for restop > 0 {
		go func() {
			err := helper.ServiceStop(srvName)
			if err != nil {
				log.Println(err)
			}
		}()
		log.Println(restop)
		time.Sleep(1 * time.Second)
		restop--
	}

	status, err = helper.ServiceStatus(srvName)
	if err != nil {
		log.Println(err)
		return err
	}

	if status != helper.StatusStopped {
		return errors.New("服务未停止")
	}

	pid, _ := dir.ReadString("./.pid")

	ppid, _ := strconv.Atoi(pid)
	if ppid == 0 {
		processId, err := helper.ServiceProcessId(srvName)
		if err != nil {
			log.Println(err)
			return err
		}
		if err == nil {
			dir.WriteString("./.pid", strconv.FormatUint(uint64(processId), 10))
		}
		ppid = int(processId)
	}

	if b, _ := process.PidExists(int32(ppid)); b {
		ps, _ := process.Processes()
		if len(ps) > 0 {
			for _, p := range ps {
				var parentPid int32
				if parent, err := p.Parent(); err == nil {
					parentPid = parent.Pid
				}
				if p.Pid != int32(ppid) && parentPid != int32(ppid) {
					continue
				}
				err = p.Kill()
				if err != nil {
					log.Println(err)
				}
			}
		}
	}

	return nil
}

func Start(srvName string) error {
	status, err := helper.ServiceStatus(srvName)
	if err != nil {
		return fmt.Errorf("get service status  %w", err)
	}

	if status == helper.StatusRunning {
		return nil
	}

	if status == helper.StatusUninstall {
		return fmt.Errorf("service uninstall")
	}

	restart := 3
	for restart > 0 {
		err := helper.ServiceStart(srvName)
		if err != nil {
			log.Println(err)
		}
		status, _ = helper.ServiceStatus(srvName)
		if status == helper.StatusRunning {
			processId, err := helper.ServiceProcessId(srvName)
			if err == nil {
				dir.WriteString("./.pid", strconv.FormatUint(uint64(processId), 10))
			}
			return nil
		}
		restart--
		log.Println("启动失败1次")
		continue
	}

	return errors.New("启动失败")
}

func Uninstall(srvName string) error {
	status, err := helper.ServiceStatus(srvName)
	if err != nil {
		return fmt.Errorf("service status get error %w", err)
	}

	if status == helper.StatusUninstall {
		return nil
	}

	return helper.ServiceUninstall(srvName)
}
