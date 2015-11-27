package plugin

import (
	"os/exec"

	"github.com/mgalvey/otto"
)

type process struct {
	cmd *exec.Cmd
}

func (p *process) Kill() {
	if proc := p.cmd.Process; proc != nil {
		proc.Kill()
	}
}

func (p *Plugin) apiProcessNew(call otto.FunctionCall) otto.Value {
	callback := call.Argument(0)
	command := call.Argument(1).String()

	args := make([]string, len(call.ArgumentList)-2)
	for i, arg := range call.ArgumentList[2:] {
		args[i] = arg.String()
	}

	proc := &process{
		cmd: exec.Command(command, args...),
	}

	go func() {
		var str string
		bytes, err := proc.cmd.Output()
		if err == nil {
			if bytes != nil {
				str = string(bytes)
			}
			p.callValue(callback, proc.cmd.ProcessState.Success(), str)
		} else {
			p.callValue(callback, false, "")
		}
	}()

	ret, _ := p.state.ToValue(proc)
	return ret
}
