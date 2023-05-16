package utils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/google/shlex"
)

// CommandResult to return the information
type CommandResult struct {
	Stdout                    string `json:"stdout"`
	RC                        int    `json:"rc"`
	ExecutionUnixTimestampSec int64  `json:"execution_unix_timestamp_sec"`
}

// Unified exit codes
const (
	Ok            = 0
	Warning       = 1
	Critical      = 2
	Unknown       = 3
	Timeout       = 124
	NotExecutable = 126
	NotFound      = 127
)

type CommandArgs struct {
	Command string
	Timeout time.Duration
	Shell   string
	Stdin   string
}

var (
	commandSysproc *syscall.SysProcAttr = &syscall.SysProcAttr{
		// Run all processes in an own process group to be able to kill the process group and all child processes
		Setpgid: true,
	}
)

func handleCommandError(arg0 string, err error) int {
	if os.IsNotExist(err) { // does not work with windows
		return NotFound
	}

	if os.IsPermission(err) { // does not work with windows
		return NotExecutable
	}

	return Unknown
}

func killProcessGroup(p *os.Process) error {
	if p.Pid == -1 {
		return errors.New("os: process already released")
	}
	if p.Pid == 0 {
		return errors.New("os: process not initialized")
	}
	sig := os.Kill
	s, ok := sig.(syscall.Signal)
	if !ok {
		return errors.New("os: unsupported signal type")
	}

	// Kill to negative pid number kills the process group (only if Setpgid=true)
	if e := syscall.Kill(-p.Pid, s); e != nil {
		if e == syscall.ESRCH {
			return errors.New("os: process already finished")
		}
		return e
	}
	return nil
}

func parseCommand(command, shell string) ([]string, string, error) {
	if shell == "" {
		args, err := shlex.Split(command)
		return args, "", err
	}

	return []string{shell}, command, nil
}

// RunCommand in shell style with timeout on every platform
func RunCommand(ctx context.Context, commandArgs CommandArgs) (*CommandResult, error) {
	result := &CommandResult{
		ExecutionUnixTimestampSec: time.Now().Unix(),
	}

	env := map[string]string{}

	processEnv := make([]string, 0)

	for k, v := range env {
		processEnv = append(processEnv, fmt.Sprintf("%s=%s", k, v))
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	ctxTimeout, cancel := context.WithTimeout(ctx, commandArgs.Timeout)
	defer cancel()

	args, stdin, err := parseCommand(commandArgs.Command, commandArgs.Shell)
	if err != nil {
		result.RC = Unknown
		result.Stdout = err.Error()

		return result, err
	}

	if commandArgs.Stdin != "" {
		// User passed data to put on stdin so put this data on stdin !
		stdin = commandArgs.Stdin
	}

	outputBuf := &bytes.Buffer{}
	stdinBuf := bytes.NewBufferString(stdin)

	c := exec.CommandContext(ctxTimeout, args[0], args[1:]...)
	c.Env = processEnv
	c.Stdout = outputBuf
	c.Stderr = outputBuf
	c.Stdin = stdinBuf

	c.SysProcAttr = commandSysproc

	// Do not hang forever
	// https://github.com/golang/go/issues/18874
	// https://github.com/golang/go/issues/22610
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctxTimeout.Done()
		switch ctxTimeout.Err() {
		case context.DeadlineExceeded:
			if c.Process != nil {
				//Kill process because of timeout
				// nolint:errcheck
				killProcessGroup(c.Process)
			}
		case context.Canceled:
			//Process exited gracefully
			if c.Process != nil {
				// nolint:errcheck
				killProcessGroup(c.Process)
			}
		}
	}()
	err = c.Run()

	if ctxTimeout.Err() == context.DeadlineExceeded {
		result.Stdout = fmt.Sprintf("Command %s timed out after %s seconds", strings.Join(args, " "), commandArgs.Timeout.String())
		result.RC = Timeout
		return result, err
	}

	if err != nil && c.ProcessState == nil {
		rc := handleCommandError(args[0], err)
		switch rc {
		case NotFound:
			result.Stdout = fmt.Sprintf("No such file or directory: '%s'", strings.Join(args, " "))
		case NotExecutable:
			result.Stdout = fmt.Sprintf("File not executable: '%s'", strings.Join(args, " "))
		default:
			result.Stdout = fmt.Sprintf("Unknown error: %s Command: '%s'", err.Error(), strings.Join(args, " "))
		}
		result.RC = rc
		return result, err
	}

	//No errors on command execution
	result.Stdout = outputBuf.String()
	result.RC = Unknown

	state := c.ProcessState
	if status, ok := state.Sys().(syscall.WaitStatus); ok {
		result.RC = status.ExitStatus()
	}

	return result, nil
}
