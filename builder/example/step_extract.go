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
	RootfsKey string
	ResultKey string
}

func (s *StepExtract) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	rootfs_path := state.Get(s.RootfsKey).(string)

	if rootfs_path == "" {
		ui.Error(fmt.Sprintf("'%s' must be set.", s.RootfsKey))
		state.Put("error", fmt.Errorf("'%s' not set", s.RootfsKey))
		return multistep.ActionHalt
	}

	wd, err := os.MkdirTemp("", "packer_rootfs")
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}
	state.Put(s.ResultKey, wd)
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
	if err := os.RemoveAll(state.Get(s.ResultKey).(string)); err != nil {
		state.Put("error", err)
	}
}
