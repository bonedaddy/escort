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
	exe := NewEXE(agent, tengoHelloWorld)
	require.NoError(t, RunEXE(context.Background(), exe))
}

/*
https://github.com/d5/tengo/blob/master/docs/interoperability.md
https://github.com/Shopify/go-lua
*/
