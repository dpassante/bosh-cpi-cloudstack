package fakes

import (
	apiv1 "github.com/cppforlife/bosh-cpi-go/apiv1"

	bwcstem "github.com/orange-cloudfoundry/bosh-cpi-cloudstack/stemcell"
	bwcvm "github.com/orange-cloudfoundry/bosh-cpi-cloudstack/vm"
)

type FakeCreator struct {
	CreateAgentID     apiv1.AgentID
	CreateStemcell    bwcstem.Stemcell
	CreateProps       bwcvm.VMProps
	CreateNetworks    apiv1.Networks
	CreateEnvironment apiv1.VMEnv
	CreateVM          bwcvm.VM
	CreateErr         error
}

func (c *FakeCreator) Create(
	agentID apiv1.AgentID, stemcell bwcstem.Stemcell, props bwcvm.VMProps,
	networks apiv1.Networks, env apiv1.VMEnv) (bwcvm.VM, error) {

	c.CreateAgentID = agentID
	c.CreateProps = props
	c.CreateStemcell = stemcell
	c.CreateNetworks = networks
	c.CreateEnvironment = env

	return c.CreateVM, c.CreateErr
}