package player

import (
	"errors"
	"os/exec"
)

type Player struct {
	cmd     *exec.Cmd
	control mpvControl
	done    chan struct{}
	paused  bool
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
	control := newMPVControl()
	args := []string{"--no-video", "--force-window=no", "--input-terminal=no"}
	args = append(args, control.args()...)
	args = append(args, url)
	cmd := exec.Command(bin, args...)
	if err := cmd.Start(); err != nil {
		control.close()
		return err
	}
	p.cmd = cmd
	p.control = control
	p.done = make(chan struct{})
	p.paused = false
	go func() {
		_ = cmd.Wait()
		control.close()
		close(p.done)
	}()
	return nil
}

func (p *Player) TogglePause() (bool, error) {
	if p.cmd == nil || p.cmd.Process == nil {
		return false, errors.New("nothing is playing")
	}
	if p.control == nil {
		return p.paused, errors.New("mpv control pipe is not available")
	}
	if err := p.control.command(`{"command":["cycle","pause"]}` + "\n"); err != nil {
		return p.paused, err
	}
	p.paused = !p.paused
	return p.paused, nil
}

func (p *Player) Stop() {
	if p.cmd == nil || p.cmd.Process == nil {
		return
	}
	if p.control != nil {
		_ = p.control.command(`{"command":["quit"]}` + "\n")
	}
	_ = p.cmd.Process.Kill()
	if p.done != nil {
		<-p.done
	}
	p.cmd = nil
	p.control = nil
	p.done = nil
	p.paused = false
}
