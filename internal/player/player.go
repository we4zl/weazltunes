package player

import (
	"errors"
	"io"
	"os/exec"
)

type Player struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	done   chan struct{}
	paused bool
}

func New() *Player {
	return &Player{}
}

func (p *Player) Play(url string) error {
	p.Stop()
	bin, err := exec.LookPath("mpv")
	if err != nil {
		return errors.New("mpv was not found; install mpv to play streams")
	}
	cmd := exec.Command(bin, "--no-video", "--force-window=no", "--input-terminal=no", "--input-file=-", url)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	p.cmd = cmd
	p.stdin = stdin
	p.done = make(chan struct{})
	p.paused = false
	go func() {
		_ = cmd.Wait()
		close(p.done)
	}()
	return nil
}

func (p *Player) TogglePause() (bool, error) {
	if p.cmd == nil || p.cmd.Process == nil {
		return false, errors.New("nothing is playing")
	}
	if p.stdin == nil {
		return p.paused, errors.New("mpv control pipe is not available")
	}
	if _, err := io.WriteString(p.stdin, "cycle pause\n"); err != nil {
		return p.paused, err
	}
	p.paused = !p.paused
	return p.paused, nil
}

func (p *Player) Stop() {
	if p.cmd == nil || p.cmd.Process == nil {
		return
	}
	if p.stdin != nil {
		_, _ = io.WriteString(p.stdin, "quit\n")
		_ = p.stdin.Close()
	}
	_ = p.cmd.Process.Kill()
	if p.done != nil {
		<-p.done
	}
	p.cmd = nil
	p.stdin = nil
	p.done = nil
	p.paused = false
}
