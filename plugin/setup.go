package plugin

import (
	"errors"
	"strings"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/v2fly/v2ray-core/v4/common/net"
)

// config syntax

/*

	`state_file` is a JSON a file containing tricked data_files to avoid
	having to repeteadly trick data at startup. the name of the file is used
	as the hostname of the dns record. for example `data_files /foo/bar` would
	use `bar.example.com` as the record to query. `data_files /foo/bar.baz` would use
	`foo.baz.example.com` as the record to query. optionally `xor_key ..` can be used
	to specify a single letter to use as an encryption key to perform XOR encryption
	of the data before tricking it. this shouldn't be considered safe/secure for transmissiong
	of secrets/sensitive information. it's merely used as a simple method of temporarily increasing
	the difficulty of signature analysis/

	no encryption:

	escort {
		state_file /path/to/state.json
		data_files /foo/bar /foo/bar.baz
	}

	using `a` as the XOR encryption key:

	escort {
		state_file /path/to/state.json
		data_files /foo/bar /foo/bar.baz
		xor_key a
	}
*/

func init() { plugin.Register("escort", setup) }

func setup(c *caddy.Controller) error {
	escort, err := parse(c)
	if err != nil {
		return plugin.Error("escort", err)
	}
	// finally add the plugin to coredns
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		escort.Next = next
		return escort
	})
	return nil
}

func parse(c *caddy.Controller) (*Escort, error) {
	var (
		stateFile string
		dataFiles []string
		xorKey    *byte = nil
	)
	for c.Next() {
		if c.NextBlock() {
			for {
				switch c.Val() {
				case "state_file":
					stateFile = c.RemainingArgs()[0]
				case "data_files":
					dataFiles = c.RemainingArgs()
				case "xor_key":
					args := c.RemainingArgs()
					xorKeyStr := args[0]
					if len(xorKeyStr) > 0 {
						return nil, plugin.Error("escort", errors.New("xor_key length greater than 1"))
					}
					a := xorKeyStr[0]
					xorKey = &a
				}
			}
		}
	}
	state, err := NewStateSynchronized(stateFile, xorKey, dataFiles)
	if err != nil {
		// dont wrap error with plugin.Error as the function does that
		return nil, err
	}
	return &Escort{state: state}, nil
}

// ensures a valid host is being used.
func isValidHost(host string) error {
	if strings.HasPrefix(host, "http://") {
		host = strings.TrimPrefix(host, "http://")
	} else {
		host = strings.TrimPrefix(host, "https://")
	}
	_, err := net.ParseDestination(host)
	return err
}
