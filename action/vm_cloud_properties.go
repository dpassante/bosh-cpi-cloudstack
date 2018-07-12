package action

import (
	"fmt"
)

type CloudStackCloudProperties struct {
	Name            string `json:"name"`
	Version         string `json:"version"`
	Infrastructure  string `json:"infrastructure"`
	Hypervisor      string `json:"hypervisor"`
	Disk            int    `json:"disk"`
	DiskFormat      string `json:"disk_format"`
	ContainerFormat string `json:"container_format"`
	OsType          string `json:"os_type"`
	OsDistro        string `json:"os_distro"`
	Architecture    string `json:"architecture"`
	AutoDiskConfig  bool   `json:"auto_disk_config"`
	LightTemplate   string `json:"light_template"`
}

func (cc CloudStackCloudProperties) Validate() error {
	if cc.Infrastructure != "cloudstack" {
		return fmt.Errorf("infrastructure '%s' is not supported (must be cloudstack)", cc.Infrastructure)
	}
	if cc.Architecture != "x86_64" {
		return fmt.Errorf("architecture '%s' is not supported (must be x86_64)", cc.Architecture)
	}
	if cc.Hypervisor != "xen" {
		return fmt.Errorf("hypervisor '%s' is not supported (must be xen)", cc.Architecture)
	}
	if cc.OsType != "linux" {
		return fmt.Errorf("os_type '%s' is not supported (must be linux)", cc.OsType)
	}
	return nil
}
