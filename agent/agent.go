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
	script   *tengo.Script
	compiled *tengo.Compiled
}

// New returns a new agent
func New() *Agent { return &Agent{} }

// NewEXE returns a new executable
func (a *Agent) NewEXE(src string) *Executable {
	return &Executable{
		script: tengo.NewScript([]byte(src)),
	}
}

// SetImports is used to set the imports needed by the tengo script
// these must be stdlib imports or else it will not function
func (e *Executable) SetImports(imports ...string) {
	e.script.SetImports(stdlib.GetModuleMap(imports...))
}

// Run executes the given tengo script
func (e *Executable) Run(ctx context.Context) error {
	compiled, err := e.script.RunContext(ctx)
	if err != nil {
		return err
	}
	e.compiled = compiled
	return nil
}
