package example

import (
	"context"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
)

type StepExtract struct {
}

func (s *StepExtract) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	return multistep.ActionContinue
}

func (s *StepExtract) Cleanup(state multistep.StateBag) {
}
