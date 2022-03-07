package plugin

import (
	"fmt"
	"os"
	"testing"

	"github.com/coredns/caddy"
	"github.com/stretchr/testify/require"
)

var (
	fooTestFileData = "foo"
	barTestFileData = "bar"
	bazTestFileData = "baz"
)

func TestPluginParse(t *testing.T) {

}

func TestSetupDnssec(t *testing.T) {
	if err := os.WriteFile("foo", []byte(fooTestFileData), 0644); err != nil {
		t.Fatalf("Failed to write test file: %s", err)
	}
	if err := os.WriteFile("bar", []byte(barTestFileData), 0644); err != nil {
		t.Fatalf("Failed to write test file: %s", err)
	}
	if err := os.WriteFile("baz", []byte(bazTestFileData), 0644); err != nil {
		t.Fatalf("Failed to write test file: %s", err)
	}
	defer func() {
		os.Remove("foo")
		os.Remove("bar")
		os.Remove("baz")
	}()
	curDir, err := os.Getwd()
	require.NoError(t, err)
	fooPath := curDir + "/foo"
	barPath := curDir + "/bar"
	bazPath := curDir + "/baz"

	tests := []struct {
		input             string
		shouldErr         bool
		expectedDataFiles []string
		expectedStateFile string
		expectedXorKey    *byte
	}{
		{`escort`, true, nil, "", nil},
		{
			fmt.Sprintf(`escort {
				state_file state_exists.json
				data_files %s %s %s
			}`, fooPath, barPath, bazPath), false,
			[]string{"foo", "bar", "baz"},
			"state_exists.json",
			nil,
		},
	}

	for i, test := range tests {
		c := caddy.NewTestController("dns", test.input)
		escort, err := parse(c)
		require.NoError(t, err)
		if test.shouldErr && err == nil {
			t.Errorf("Test %d: Expected error but found %s for input %s", i, err, test.input)
		}
		t.Logf("%+v\n", escort)
	}
}
