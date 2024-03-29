//go:generate packer-sdc mapstructure-to-hcl2 -type Config

package example

import (
	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	Comm                communicator.Config `mapstructure:",squash"`
	RootfsConfig        `mapstructure:",squash"`

	ctx interpolate.Context
}

func (c *Config) Prepare(raws ...interface{}) ([]string, error) {
	err := config.Decode(c, &config.DecodeOpts{
		PluginType:         BuilderId,
		Interpolate:        true,
		InterpolateContext: &c.ctx,
	}, raws...)
	if err != nil {
		return nil, err
	}

	if c.Comm.Type == "" {
		c.Comm.Type = "rootlesskit"
	}

	var errs *packer.MultiError
	warnings := make([]string, 0)

	rootfsWarnings, rootfsErrs := c.RootfsConfig.Prepare(&c.ctx)
	warnings = append(warnings, rootfsWarnings...)
	errs = packer.MultiErrorAppend(errs, rootfsErrs...)

	if errs != nil && len(errs.Errors) > 0 {
		return warnings, errs
	}

	return warnings, nil
}
