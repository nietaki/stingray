package state

import (
	"encoding/json"
	"fmt"
	"hash/maphash"
	"os"
	"reflect"
	"strings"
)

type StateManager[T comparable] struct {
	Data T

	name     string
	hashSeed maphash.Seed
	lastHash uint64
}

// MAIN INTERFACE

func NewStateManager[T comparable]() *StateManager[T] {
	ret := StateManager[T]{}
	ret.CheckChanged() // initialize lastHash
	return &ret
}

func (this *StateManager[T]) CheckChanged() bool {
	curHash := this.computeHash()
	if curHash != this.lastHash {
		this.lastHash = curHash
		return true
	}
	return false
}

func (this *StateManager[T]) SetName(name string) {
	this.name = name
}

// CORE FUNCTIONALITY

func (this *StateManager[T]) computeHash() uint64 {
	return maphash.Comparable(this.hashSeed, this.Data)
}

// SAVING AND LOADING

func (this *StateManager[T]) Save() error {
	path := this.dataPath()
	data, err := json.MarshalIndent(this.Data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (this *StateManager[T]) Load() error {
	bytes, err := os.ReadFile(this.dataPath())
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, &this.Data)
}

func (this *StateManager[T]) LoadOrDefault(df T) {
	err := this.Load()
	if err != nil {
		this.Data = df
	}
}

func (this *StateManager[T]) GetName() string {
	if this.name != "" {
		return this.name
	}
	rt := reflect.TypeFor[T]()
	packagePath := rt.PkgPath()
	lastPackageBit := packagePath[strings.LastIndex(packagePath, "/")+1:]
	structName := rt.Name()
	return fmt.Sprintf("%s-%s", lastPackageBit, structName)
}

func (this *StateManager[T]) dataPath() string {
	return fmt.Sprintf("state/%s.json", this.GetName())
}
