package example

import (
	"context"
	"errors"

	"github.com/hashicorp/go-getter/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/multistep/commonsteps"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

const BuilderId = "example.builder"

type Builder struct {
	config Config
	runner multistep.Runner
}

func (b *Builder) ConfigSpec() hcldec.ObjectSpec {
	return b.config.FlatMapstructure().HCL2Spec()
}

func (b *Builder) Prepare(raws ...interface{}) (generatedVars []string, warnings []string, err error) {
	warnings, err = b.config.Prepare(raws...)
	if err != nil {
		return nil, warnings, err
	}
	return nil, warnings, nil
}

func (b *Builder) Run(ctx context.Context, ui packer.Ui, hook packer.Hook) (packer.Artifact, error) {
	// https://github.com/solo-io/packer-plugin-arm-image/blob/227d49f89dd8c20a86e19787c4c67ca6b5f13201/pkg/builder/builder.go#L239
	getter.Decompressors = make(map[string]getter.Decompressor)

	steps := []multistep.Step{
		&commonsteps.StepDownload{
			Description: "Rootfs",
			ResultKey:   "rootfs_path",
			Url:         b.config.RootfsUrls,
			Checksum:    b.config.RootfsChecksum,
			TargetPath:  b.config.TargetPath,
			Extension:   b.config.TargetExtension,
		},
		&StepExtract{},
		new(commonsteps.StepProvision),
	}

	state := new(multistep.BasicStateBag)
	state.Put("hook", hook)
	state.Put("ui", ui)

	b.runner = commonsteps.NewRunnerWithPauseFn(steps, b.config.PackerConfig, ui, state)
	b.runner.Run(ctx, state)

	if err, ok := state.GetOk("error"); ok {
		return nil, err.(error)
	}

	if _, ok := state.GetOk(multistep.StateCancelled); ok {
		return nil, errors.New("build was cancelled")
	}

	if _, ok := state.GetOk(multistep.StateHalted); ok {
		return nil, errors.New("build was halted")
	}

	return &Artifact{}, nil
}
