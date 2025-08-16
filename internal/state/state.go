package state

import (
	"encoding/json"
	"os"
	"time"
)

// State holds metadata for the cache, including
// the last update timestamp and a map of file hashes.
type State struct {
	LastUpdated int64
	Data        map[string]string
}

var defaultPath = ".nox-state.json"

// SetPath updates the default file path used for saving and loading state
func SetPath(path string) {
	defaultPath = path
}

// Touch updates the LastUpdated timestamp to the current time.
func (s *State) Touch() {
	s.LastUpdated = time.Now().Unix()
}

// Load reads the state from the state file
func Load() (*State, error) {
	return loadFromFile(defaultPath)
}

// Save writes the given State to the state file.
// Overwrites any existing state file.
func Save(state *State) error {
	return saveToFile(defaultPath, state)
}

// loadFromFile reads the state JSON from the specified file path.
func loadFromFile(path string) (*State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return &State{Data: make(map[string]string)}, nil
	}
	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

// saveToFile writes the State as JSON to the specified file path.
func saveToFile(path string, state *State) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
