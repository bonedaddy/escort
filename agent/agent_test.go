package agent

import (
	"testing"

	"github.com/Shopify/go-lua"
	"github.com/stretchr/testify/require"
)

var (
	testLuaProgram = `
		print("Hello World")
	`
)

func TestAgent(t *testing.T) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	require.NoError(t, lua.DoString(l, testLuaProgram))
}

/*
https://github.com/d5/tengo/blob/master/docs/interoperability.md
https://github.com/Shopify/go-lua
*/
