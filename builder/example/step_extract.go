package example

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type StepExtract struct {
}

func (s *StepExtract) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	rootfs_path := state.Get("rootfs_path").(string)

	if rootfs_path == "" {
		ui.Error("'rootfs_path' must be set.")
		state.Put("error", fmt.Errorf("'rootfs_path' not set"))
		return multistep.ActionHalt
	}

	ui.Say(fmt.Sprintf("Rootfs_path to %q", rootfs_path))

	return multistep.ActionContinue
}

func (s *StepExtract) Cleanup(state multistep.StateBag) {
}
