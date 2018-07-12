package action

import (
	"github.com/orange-cloudfoundry/bosh-cpi-cloudstack/util"
	"github.com/orange-cloudfoundry/bosh-cpi-cloudstack/config"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"github.com/xanzy/go-cloudstack/cloudstack"
)

func (a CPI) setMetadata(tagType config.Tags, cid string, meta util.MetaMarshal) error {
	params := a.client.Resourcetags.NewCreateTagsParams([]string{cid}, string(tagType), util.ConvertMapToTags(meta))
	_, err := a.client.Resourcetags.CreateTags(params)
	if err != nil {
		return bosherr.WrapErrorf(err, "Setting %s metadata '%s'", tagType, cid)
	}
	return nil
}

func (a CPI) findVmsByName(cid apiv1.VMCID) ([]*cloudstack.VirtualMachine, error) {
	p := a.client.VirtualMachine.NewListVirtualMachinesParams()
	p.SetName(cid.AsString())
	resp, err := a.client.VirtualMachine.ListVirtualMachines(p)
	if err != nil {
		return []*cloudstack.VirtualMachine{}, err
	}
	return resp.VirtualMachines, nil
}

func (a CPI) findVolumesByName(cid apiv1.DiskCID) ([]*cloudstack.Volume, error) {
	p := a.client.Volume.NewListVolumesParams()
	p.SetName(cid.AsString())
	resp, err := a.client.Volume.ListVolumes(p)
	if err != nil {
		return []*cloudstack.Volume{}, err
	}
	return resp.Volumes, nil
}

func (a CPI) findVmId(cid apiv1.VMCID) (string, error) {
	vms, err := a.findVmsByName(cid)
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Can't find vm name '%s'", cid.AsString())
	}
	if len(vms) == 0 {
		return "", bosherr.Errorf("Can't find vm name '%s'", cid.AsString())

	}
	return vms[0].Id, nil
}

func (a CPI) findVolumeId(cid apiv1.DiskCID) (string, error) {
	volumes, err := a.findVolumesByName(cid)
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Can't find disk name '%s'", cid.AsString())
	}
	if len(volumes) == 0 {
		return "", bosherr.Errorf("Can't find disk name '%s'", cid.AsString())
	}
	return volumes[0].Id, nil
}

func (a CPI) findZoneId() (string, error) {
	p := a.client.Zone.NewListZonesParams()
	p.SetName(a.config.CloudStack.DefaultZone)
	resp, err := a.client.Zone.ListZones(p)
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Can't find zone name '%s'", a.config.CloudStack.DefaultZone)
	}
	if len(resp.Zones) == 0 {
		return "", bosherr.Errorf("Can't find zone name '%s'", a.config.CloudStack.DefaultZone)
	}
	return resp.Zones[0].Id, nil
}