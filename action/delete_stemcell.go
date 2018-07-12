package action

import (
	"fmt"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

// DeleteStemcell - Delete CS template matching given stemcell name
func (a CPI) DeleteStemcell(cid apiv1.StemcellCID) error {
	// TODO [xmt]  "featured", "self", "selfexecutable","sharedexecutable","executable", and "community" ?
	listTplP := a.client.Template.NewListTemplatesParams("self")
	listTplP.SetName(cid.AsString())

	templatesRes, err := a.client.Template.ListTemplates(listTplP)
	if err != nil {
		return bosherr.WrapErrorf(err, "[delete_stemcell] error while listing templates matching stemcell name '%s'", cid.AsString())
	}

	if templatesRes.Count == 0 {
		return fmt.Errorf("[delete_stemcell] could not find any template matching stemcell name '%s'", cid.AsString())
	} else if templatesRes.Count != 1 {
		return fmt.Errorf("[delete_stemcell] found multiple templates matching stemcell name '%s'", cid.AsString())
	}

	zoneid, err := a.findZoneId()
	if err != nil {
		return err
	}

	template := templatesRes.Templates[0]
	deleteP := a.client.Template.NewDeleteTemplateParams(template.Id)
	deleteP.SetZoneid(zoneid)
	_, err = a.client.Template.DeleteTemplate(deleteP)
	if err != nil {
		return bosherr.WrapErrorf(err, "[delete_stemcell] error while deleteing stemcell '%s'", cid.AsString())
	}
	return nil
}