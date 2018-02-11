package commands

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"yuzu"
	"yuzu/logger"
)

// Exec execute things
type Exec struct{}

// IsOwnerOnly f
func (Exec) IsOwnerOnly() bool {
	return true
}

// Help f
func (Exec) Help() [2]string {
	return [2]string{"Executes a command in the terminal", ""}
}

// Process f
func (Exec) Process(ctx yuzu.Context) {
	if len(ctx.Args) == 0 {
		_, e := ctx.Say("Please provide some arguments to execute.")
		if e != nil {
			logger.ERROR.L(fmt.Sprintf("error sending message, %s", e))
			return
		}
		return
	}

	suffix := strings.Join(ctx.Args, " ")
	out, err := exec.Command(suffix).Output()
	if err != nil {
		ctx.Say("```md\n" + err.Error() + "```")
		return
	}
	n := bytes.IndexByte(out, 0)
	st := string(out[:n])

	ctx.Say("```md\n" + st + "```")
}
