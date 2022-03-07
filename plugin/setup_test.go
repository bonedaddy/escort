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
		os.Remove("state_exists.json")
	}()
	curDir, err := os.Getwd()
	require.NoError(t, err)
	fooPath := curDir + "/foo"
	barPath := curDir + "/bar"
	bazPath := curDir + "/baz"
	tests := []struct {
		input                  string
		shouldErr              bool
		expectedDataFiles      []string
		expectedDataFileHashes []string
		expectedStateFile      string
		expectedXorKey         string
		wantDataFilesLen       int
	}{
		{`escort`, true, nil, nil, "", "", 0},
		{
			fmt.Sprintf(`escort {
				state_file state_exists.json
				data_files %s %s %s
			}`, fooPath, barPath, bazPath), false,
			// must always be ordered foo, bar ,baz
			[]string{"foo", "bar", "baz"},
			// must always be ordered foo, bar, baz
			[]string{"0|SsvPBwQAAP//\n", "0|SkosAgQAAP//\n", "0|SkqsAgQAAP//\n"},
			"state_exists.json",
			"",
			3,
		},
		{
			fmt.Sprintf(`escort {
				state_file state_exists.json
				data_files %s %s
				xor_key %s
			}`, fooPath, barPath, "z"), false,
			[]string{"foo", "bar"},
			[]string{"0|khEVBQQAAP//\n", "0|kpDmAAQAAP//\n"},
			"state_exists.json",
			"z",
			2,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			os.Remove("state_exists.json")
			defer func() { os.Remove("state_exists.json") }()
			c := caddy.NewTestController("dns", test.input)
			t.Logf("test input: %s\n", test.input)
			escort, err := parse(c)
			if (err != nil) && !test.shouldErr {
				t.Fatalf("expected no error but error found %s", err.Error())
			}
			if err != nil {
				return
			}
			require.NotNil(t, escort)
			checkEntries := func(fromDisk bool) {
				t.Logf("checkEntries from disk %v", fromDisk)
				require.Len(t, escort.state.Entries, test.wantDataFilesLen)
				require.Equal(t, escort.state.xorKey, test.expectedXorKey)
				for filename, entry := range escort.state.Entries {
					require.Len(t, entry.DataSegments, 1)
					if filename == "foo" {
						require.Equal(t, entry.DataSegments[0], test.expectedDataFileHashes[0])
					} else if filename == "bar" {
						require.Equal(t, entry.DataSegments[0], test.expectedDataFileHashes[1])
					} else if filename == "baz" {
						require.Equal(t, entry.DataSegments[0], test.expectedDataFileHashes[2])
					}
				}
			}
			// run checkEntries once against existing state, then again against fresh state
			checkEntries(false)
			escort.state.Save("state_exists.json")
			escort.state, err = LoadOrInitState("state_exists.json", escort.state.xorKey, escort.state.segmentSize, escort.state.segmentIdentifier)
			require.NoError(t, err)
			checkEntries(true)
		})

	}
}
