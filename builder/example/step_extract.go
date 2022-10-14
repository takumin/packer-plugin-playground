package example

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type StepExtract struct {
	RootfsPathKey       string
	WorkingDirectoryKey string
}

func (s *StepExtract) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	rootfs_path := state.Get(s.RootfsPathKey).(string)

	if rootfs_path == "" {
		ui.Error(fmt.Sprintf("'%s' must be set.", s.RootfsPathKey))
		state.Put("error", fmt.Errorf("'%s' not set", s.RootfsPathKey))
		return multistep.ActionHalt
	}

	wd, err := os.MkdirTemp("", "packer_rootfs")
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}
	state.Put(s.WorkingDirectoryKey, wd)
	ui.Say(fmt.Sprintf("Working Directory: %s", wd))

	// #nosec G204
	out, err := exec.CommandContext(ctx, "rootlesskit", "tar", "-xvf", rootfs_path, "-C", wd).CombinedOutput()
	if err != nil {
		state.Put("error", err)
		s.Cleanup(state)
		return multistep.ActionHalt
	}
	state.Put("rootfs_extract", strings.Split(string(out), "\n"))
	ui.Say("Extract Rootfs Archive")

	return multistep.ActionContinue
}

func (s *StepExtract) Cleanup(state multistep.StateBag) {
	if err := os.RemoveAll(state.Get(s.WorkingDirectoryKey).(string)); err != nil {
		state.Put("error", err)
	}
}
