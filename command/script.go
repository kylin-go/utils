package command

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func setInterpreterFromScript(line *string) string {
	re := regexp.MustCompile(`^#!\s*(?P<name>[^#]*)`)
	match := re.FindStringSubmatch(*line)
	if len(match) > 0 {
		return strings.Trim(match[1], "\n")
	}
	return ""
}

func setInterpreterFromPath(scriptPath *string) string {
	re := regexp.MustCompile(".*\\.(sh|py)$")
	match := re.FindStringSubmatch(*scriptPath)
	if len(match) > 0 {
		if match[1] == "sh" {
			return "sh"
		}
		if match[1] == "py" {
			return "python"
		}
	}
	return ""
}

func setExecUser(line *string) string {
	re := regexp.MustCompile(`^#\s*user\s+(?P<name>[^#]*)(.*)`)
	match := re.FindStringSubmatch(*line)
	if len(match) > 0 {
		return match[1]
	}
	return ""
}

// 封装exec.command封装功能
// 1. 读取脚本头 #! <interpreter>，或者脚本扩展名(sh, py)获取脚本解释器,如果脚本设置，不解析脚本名称；如果脚本内容没有指定，读取脚本扩展名
// 2. 读取本头 #user <username>，获取执行脚本的用户，如果execUser参数不为空，则使用指定的用户
func ExecScript(scriptPath, execUser string, args ...string) (*exec.Cmd, error) {
	var (
		user        string
		interpreter string
		params      []string
		command     string
	)
	f, err := os.Open(scriptPath)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReader(f)
	defer func() { _ = f.Close() }()
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF || line[:1] != "#" {
			break
		}
		if execUser == "" {
			if user == "" {
				user = setExecUser(&line)
			}
		}
		if interpreter == "" {
			interpreter = setInterpreterFromScript(&line)
		}
	}
	if execUser != "" {
		user = execUser
	}
	if interpreter == "" {
		interpreter = setInterpreterFromPath(&scriptPath)
	}
	if interpreter == "" {
		return nil, errors.New("没有设置脚本解释器，设置方法：\n    1.脚本名称为*.sh或*.py;\n    2.脚本头信息增加#!/usr/bin/python")
	}
	if user == "" {
		command = interpreter
		params = append([]string{scriptPath}, args...)
	} else {
		command = "su"
		params = append([]string{"-", user, interpreter, scriptPath}, args...)
	}
	return exec.Command(command, params...), nil
}

func StdExecScript(scriptPath, execUser string, args ...string) (stdout, stderr bytes.Buffer, err error) {
	cmd, err := ExecScript(scriptPath, execUser, args...)
	if err != nil {
		return
	}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Start()
	if err != nil {
		return
	}
	err = cmd.Wait()
	return
}
