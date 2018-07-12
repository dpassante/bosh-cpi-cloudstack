package action

import (
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

func (a CPI) RebootVM(cid apiv1.VMCID) error {
	a.client.AsyncTimeout(a.config.CloudStack.Timeout.Reboot)

	id, err := a.findVmId(cid)
	if err != nil {
		return bosherr.WrapErrorf(err, "Rebooting vm %s", cid.AsString())
	}

	p := a.client.VirtualMachine.NewRebootVirtualMachineParams(id)
	_, err = a.client.VirtualMachine.RebootVirtualMachine(p)
	if err != nil {
		return bosherr.WrapErrorf(err, "Rebooting vm %s", cid.AsString())
	}
	return nil
}