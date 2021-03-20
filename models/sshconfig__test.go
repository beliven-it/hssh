package models

import (
	"os"
	"testing"
)

func createSSHFakeFile(path string) {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	file.WriteString(`
# HSSH start managed
#Include config.hssh.d/*
 Include config.personal.d/*
# Include kotor/*
# HSSH end managed
Include nova/*
Include proxima/*
	`)
}

func TestSSHConfigTakeIncludes(t *testing.T) {
	fakePath := "/tmp/config"
	createSSHFakeFile(fakePath)

	i := NewSSHConfig(fakePath)

	includes := i.GetIncludes()

	if len(includes) != 3 {
		t.Errorf("Invalid include count %d", len(includes))
	}
}
