package plugin

import (
	"bytes"
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
	// maps filename -> entry
	Entries           map[string]StateEntry
	xorKey            string
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
	xorKey string,
	segmentSize int,
	segmentIdentifier string,
) (*State, error) {
	fh, err := os.Open(stateFilePath)
	if err != nil && os.IsNotExist(err) {
		state := new(State)
		state.xorKey = xorKey
		state.segmentIdentifier = segmentIdentifier
		state.segmentSize = segmentSize
		state.Entries = make(map[string]StateEntry)
		return state, nil
	}
	data, err := ioutil.ReadAll(fh)
	if err != nil {
		return nil, err
	}
	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	if state.Entries == nil {
		state.Entries = make(map[string]StateEntry)
	}
	state.xorKey = xorKey
	state.segmentSize = segmentSize
	state.segmentIdentifier = segmentIdentifier
	return &state, nil
}

func NewStateSynchronized(
	stateFilePath string,
	xorKey string,
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

	if err := state.Synchronise(dataFilesMap, stateFilePath); err != nil {
		return nil, plugin.Error("escort", err)
	}
	return state, state.Synchronise(dataFilesMap, stateFilePath)
}

func (s *State) Synchronise(
	// maps filename -> filepath
	dataFiles map[string]string,
	stateFilePath string,
) error {
	var xorKey *byte = nil
	if s.xorKey != "" {
		if len(s.xorKey) != 1 {
			return plugin.Error("escort", fmt.Errorf("xor_key %s length %v not equal to 1", s.xorKey, len(s.xorKey)))
		}
		_xor := byte(s.xorKey[0])
		xorKey = &_xor
	}
	core := pkg.NewCore(xorKey, s.segmentIdentifier, s.segmentSize)
	s.mx.Lock()
	defer s.mx.Unlock()
	var newEntries []StateEntry
	for filename, filepath := range dataFiles {
		_, ok := s.Entries[filename]
		if ok {
			continue
		}

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
	for _, entry := range newEntries {
		s.Entries[entry.FileName] = entry
	}
	if len(newEntries) > 0 {
		// save the config
		return s.Save(stateFilePath)
	}
	return nil
}

func (s *State) Save(filepath string) error {
	data, err := json.Marshal(&s)
	if err != nil {
		return err
	}
	var pretty bytes.Buffer
	if err := json.Indent(&pretty, data, "", "\t"); err != nil {
		return err
	}
	return ioutil.WriteFile(filepath, pretty.Bytes(), os.ModePerm)
}
