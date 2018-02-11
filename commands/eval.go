package commands

import (
	"fmt"
	"strings"

	"yuzu"

	"github.com/robertkrimen/otto"
	// only importing it in a library-type-level for otto to autoload underscore, also god damn you linters
	_ "github.com/robertkrimen/otto/underscore"
)

var vm = otto.New()

// Eval f
type Eval struct{}

// IsOwnerOnly f
func (Eval) IsOwnerOnly() bool {
	return true
}

// Help f
func (Eval) Help() [2]string {
	return [2]string{"Executes javascript code", "<code>"}
}

// Process f
func (Eval) Process(context yuzu.Context) {
	code := strings.Join(context.Args, " ")
	vm.Set("context", context)
	value, err := vm.Run(code)
	if err != nil {
		embed := yuzu.NewEmbed("")
		embed.Field("Input", "```js\n"+code+"```", false)
		embed.Field("Output", "```"+fmt.Sprintf("%s", value)+"```", false)
		context.SayEmbed(embed)
		return
	}
	embed := yuzu.NewEmbed("")
	embed.Field("Input", "```js\n"+code+"```", false)
	embed.Field("Output", "```js\n"+fmt.Sprintf("%s", value)+"```", false)
	context.SayEmbed(embed)
}
