package agent

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	tengoHelloWorld = `
	fmt := import("fmt")
	fmt.println("hello world")
	`
)

func TestAgent(t *testing.T) {
	agent := New()
	exe := agent.NewEXE(tengoHelloWorld)
	require.NoError(t, exe.Run(context.Background()))
}

/*
https://github.com/d5/tengo/blob/master/docs/interoperability.md
https://github.com/Shopify/go-lua
*/
