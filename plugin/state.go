package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"

	"github.com/bonedaddy/escort/pkg"
	"github.com/coredns/coredns/plugin"
)

type State struct {
	Entries           []StateEntry
	xorKey            *byte
	segmentSize       int
	segmentIdentifier string
	mx                sync.RWMutex
}

type StateEntry struct {
	// the absolute file path
	// /foo/bar
	FilePath string
	// the file name
	// bar
	FileName     string
	DataSegments []string
}

func LoadOrInitState(
	stateFilePath string,
	xorKey *byte,
	segmentSize int,
	segmentIdentifier string,
) (*State, error) {
	fh, err := os.Open(stateFilePath)
	if err != nil && os.IsNotExist(err) {
		state := new(State)
		state.xorKey = xorKey
		state.segmentIdentifier = segmentIdentifier
		state.segmentSize = segmentSize
		return state, nil
	}
	data, err := ioutil.ReadAll(fh)
	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

func NewStateSynchronized(
	stateFilePath string,
	xorKey *byte,
	dataFiles []string,
) (*State, error) {
	state, err := LoadOrInitState(stateFilePath, xorKey, 250, "|")
	if err != nil {
		return nil, plugin.Error("escort", err)
	}
	var dataFilesMap = make(map[string]string)
	for _, dataFile := range dataFiles {
		_, filename := path.Split(dataFile)
		dataFilesMap[filename] = dataFile
	}
	return state, state.Synchronise(dataFilesMap)
}

func (s *State) Synchronise(
	// maps filename -> filepath
	dataFiles map[string]string,
) error {
	core := pkg.NewCore(s.xorKey, s.segmentIdentifier, s.segmentSize)
	s.mx.Lock()
	defer s.mx.Unlock()
	var newEntries []StateEntry
	for filename, filepath := range dataFiles {
		found := false
		for _, entry := range s.Entries {
			if entry.FileName == filename {
				found = true
				break
			}
		}
		if !found {
			fileData, err := ioutil.ReadFile(filepath)
			if err != nil {
				return plugin.Error("escort", errors.New("failed to load data to trick"))
			}
			dataSegments, err := core.Trick(fileData)
			if err != nil {
				return plugin.Error("escort", fmt.Errorf("failed to trick data %s", err.Error()))
			}
			newEntries = append(newEntries, StateEntry{
				FilePath:     filepath,
				FileName:     filename,
				DataSegments: dataSegments,
			})
		}
	}
	s.Entries = append(s.Entries, newEntries...)
	return nil
}
