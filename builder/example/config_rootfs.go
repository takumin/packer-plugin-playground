package example

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-getter/v2"
	urlhelper "github.com/hashicorp/go-getter/v2/helper/url"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

type RootfsConfig struct {
	RootfsUrls      []string `mapstructure:"rootfs_urls" required:"true"`
	RootfsChecksum  string   `mapstructure:"rootfs_checksum" required:"true"`
	TargetPath      string   `mapstructure:"rootfs_target_path"`
	TargetExtension string   `mapstructure:"rootfs_target_extension"`
}

func (c *RootfsConfig) Prepare(*interpolate.Context) (warnings []string, errs []error) {
	if len(c.RootfsUrls) == 0 {
		errs = append(errs, errors.New("one of rootfs_urls must be specified"))
		return
	}
	if c.TargetExtension == "" {
		extensions := []string{
			"tar",
			"tar.bz2",
			"tar.gz",
			"tar.xz",
			"tar.zst",
			"tbz2",
			"tgz",
			"txz",
			"tzst",
		}
		for _, ext := range extensions {
			if ok, _ := filepath.Match(fmt.Sprintf("*.%s", ext), filepath.Base(c.RootfsUrls[0])); ok {
				c.TargetExtension = ext
				break
			}
		}
		if c.TargetExtension == "" {
			errs = append(errs, errors.New("file extension could not be identified"))
			return
		}
	}
	c.TargetExtension = strings.ToLower(c.TargetExtension)

	if c.RootfsChecksum == "" {
		errs = append(errs, fmt.Errorf("checksum must be specified"))
	} else {
		u, err := urlhelper.Parse(c.RootfsUrls[0])
		if err != nil {
			return warnings, append(errs, fmt.Errorf("url parse: %s", err))
		}

		q := u.Query()
		if c.RootfsChecksum != "" {
			q.Set("checksum", c.RootfsChecksum)
		}
		u.RawQuery = q.Encode()

		wd, err := os.Getwd()
		if err != nil {
			log.Printf("Getwd: %v", err)
		}

		req := &getter.Request{
			Src: u.String(),
			Pwd: wd,
		}
		cksum, err := getter.DefaultClient.GetChecksum(context.TODO(), req)
		if err != nil {
			errs = append(errs, fmt.Errorf("%v in %q", err, req.URL().Query().Get("checksum")))
		} else {
			c.RootfsChecksum = cksum.String()
		}
	}

	return warnings, errs
}
