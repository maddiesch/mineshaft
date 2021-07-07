package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type Execution struct {
	mu       sync.Mutex
	inPipe   io.WriteCloser
	cmd      exec.Cmd
	stopped  chan struct{}
	stopping bool
}

func StartMinecraftExecution(ctx context.Context, config Config) (*Execution, error) {
	args := []string{
		config.Machine.JavaCommand,
		fmt.Sprintf("-Xmx%s", config.Machine.MaxAllocation),
		fmt.Sprintf("-Xms%s", config.Machine.PreAllocation),
		"-jar",
		config.JarPath(),
		"--nogui",
		"--port",
		config.Server.Port,
	}
	cmd := exec.Cmd{
		Path:   config.Machine.JavaCommand,
		Args:   args,
		Dir:    config.Machine.WorkingDir,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	inPipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	e := &Execution{
		inPipe:  inPipe,
		cmd:     cmd,
		stopped: make(chan struct{}),
	}

	go e._unsafeWait()

	return e, nil
}

func (e *Execution) _unsafeWait() {
	err := e.cmd.Wait()

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.stopping {
		close(e.stopped)
		return
	}

	Logger.Printf("Server stopped unexpectedly! (%s)", err)

	close(e.stopped) // TODO: - Auto Attempt to restart?
}

func (e *Execution) Stop() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.stopping = true

	e.unsafeSend(`stop`)
}

func (e *Execution) unsafeSend(message string) {
	if !strings.HasSuffix(message, "\n") {
		message = message + "\n"
	}

	if _, err := e.inPipe.Write([]byte(message)); err != nil {
		Logger.Printf("Failed to send command -- %v", err)
	}
}

func (e *Execution) Send(message string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.unsafeSend(message)
}

func (e *Execution) Stopped() <-chan struct{} {
	return e.stopped
}
