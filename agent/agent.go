package agent

import (
	"context"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
)

// Agent bundles together all agent functionality
type Agent struct {
}

// Executable is a single instance of an executable to run
type Executable struct {
	src      string
	script   *tengo.Script
	compiled *tengo.Compiled
}

// New returns a new agent
func New() *Agent { return &Agent{} }

// NewEXE returns a new executable
// we declare the function like this to enable obfuscating its name
func NewEXE(agent *Agent, src string) *Executable {
	script := tengo.NewScript([]byte(src))
	// register all stdlib imports
	script.SetImports(stdlib.GetModuleMap(stdlib.AllModuleNames()...))
	return &Executable{
		src:    src,
		script: script,
	}
}

// RunEXE executes the given tengo script
func RunEXE(ctx context.Context, exe *Executable) error {
	compiled, err := exe.script.RunContext(ctx)
	if err != nil {
		return err
	}
	exe.compiled = compiled
	return nil
}
