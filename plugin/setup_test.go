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
	z := byte("z"[0])
	tests := []struct {
		input             string
		shouldErr         bool
		expectedDataFiles []string
		expectedStateFile string
		expectedXorKey    *byte
		wantDataFilesLen  int
	}{
		{`escort`, true, nil, "", nil, 0},
		{
			fmt.Sprintf(`escort {
				state_file state_exists.json
				data_files %s %s %s
			}`, fooPath, barPath, bazPath), false,
			[]string{"foo", "bar", "baz"},
			"state_exists.json",
			nil,
			3,
		},
		{
			fmt.Sprintf(`escort {
				state_file state_exists.json
				data_files %s %s
				xor_key %s
			}`, fooPath, barPath, "z"), false,
			[]string{"foo", "bar"},
			"state_exists.json",
			&z,
			2,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			c := caddy.NewTestController("dns", test.input)
			t.Logf("test input: %s\n", test.input)
			escort, err := parse(c)
			if test.shouldErr && err == nil {
				t.Errorf("Test %d: Expected error but found %s for input %s", i, err, test.input)
			} else if test.shouldErr {
				return
			}
			require.NotNil(t, escort)
			t.Logf("%+v\n", escort)
			require.Len(t, escort.state.Entries, test.wantDataFilesLen)
			if test.expectedXorKey == nil {
				require.Nil(t, escort.state.xorKey)
			} else {
				require.Equal(t, *escort.state.xorKey, *test.expectedXorKey)
			}
		})

	}
}
