package command

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os/exec"
	"time"
)

type Result struct {
	Stdout   string
	Stderr   string
	Pid      int
	ExitCode int
}

func execBufOut(workDir string, timeout int, command string, args ...string) (*Result, error) {
	var result Result
	var err error
	var stdout, stderr bytes.Buffer
	var ctxt, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()
	cmd := exec.CommandContext(ctxt, command, args...)
	if workDir != "" {
		cmd.Dir = workDir
	}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	result.Pid = cmd.Process.Pid
	result.Stdout = stdout.String()
	result.Stderr = stderr.String()
	if err != nil {
		if err.Error() == "signal: killed" {
			result.ExitCode = 137
		} else {
			result.ExitCode = 1
		}
	} else {
		result.ExitCode = 0
	}
	return &result, err
}

func execNoOut(workDir string, timeout int, command string, args ...string) (*Result, error) {
	var result Result
	var err error
	var ctxt, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()
	cmd := exec.CommandContext(ctxt, command, args...)
	if workDir != "" {
		cmd.Dir = workDir
	}
	err = cmd.Run()
	result.Pid = cmd.Process.Pid
	if err != nil {
		if err.Error() == "signal: killed" {
			result.ExitCode = 137
		} else {
			result.ExitCode = 1
		}
	} else {
		result.ExitCode = 0
	}
	return &result, err
}

func execPipeOut(workDir string, timeout int, maxLine int, command string, args ...string) (*Result, error) {
	var err error
	var buf = make([]byte, 128)
	var result Result
	var ctxt, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()

	cmd := exec.CommandContext(ctxt, command, args...)
	if workDir != "" {
		cmd.Dir = workDir
	}
	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		return &result, err
	}
	// 获取标准输出
	stdoutReader := bufio.NewReader(stdout)
	var i = 0
	for {
		n, err := stdoutReader.Read(buf)
		if err != nil || io.EOF == err || i >= maxLine {
			result.Stdout += string(buf[:n])
			break
		}
		i++
		result.Stdout += string(buf[:n])
	}
	_ = stdout.Close()
	// 获取标准错误
	stderrReader := bufio.NewReader(stderr)
	i = 0
	for {
		n, err := stderrReader.Read(buf)
		if err != nil || io.EOF == err || i >= maxLine {
			result.Stderr += string(buf[:n])
			break
		}
		i++
		result.Stderr += string(buf[:n])
	}
	_ = stderr.Close()
	result.Pid = cmd.Process.Pid
	err = cmd.Wait()
	result.Pid = cmd.Process.Pid
	if err != nil {
		if err.Error() == "signal: killed" {
			result.ExitCode = 137
		} else {
			result.ExitCode = 1
		}
	} else {
		result.ExitCode = 0
	}
	return &result, err
}

// 执行系统命令
// :maxLine 最大输出文件行数，-1 全都文件， 0 不输出， >0 指定的行数
// :workDir 命令执行的工作目录， ""表示不指定使用默认
func ExecCommand(workDir string, maxLine int, timeout int, command string, args ...string) (*Result, error) {
	var res = &Result{}
	var err error
	if maxLine < 0 {
		maxLine = -1
	}
	if maxLine > 5000 {
		maxLine = 5000
	}
	switch maxLine {
	case -1:
		res, err = execBufOut(workDir, timeout, command, args...)
	case 0:
		res, err = execNoOut(workDir, timeout, command, args...)
	default:
		res, err = execPipeOut(workDir, timeout, maxLine, command, args...)
	}
	if res == nil {
		res = &Result{"", "", 0, -1}
	}
	return res, err
}
