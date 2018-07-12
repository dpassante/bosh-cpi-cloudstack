package action

import (
	"bytes"
	"crypto/tls"
	"fmt"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"github.com/orange-cloudfoundry/bosh-cpi-cloudstack/config"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

// CreateStemcell - Create CS template from given stemcell
//
// 1. read cloud-properties
// 2. generate template name from random id
// 3. request CS an upload token
// 4. push image to CS recieved endpoint
func (a CPI) CreateStemcell(imagePath string, cp apiv1.StemcellCloudProps) (apiv1.StemcellCID, error) {
	csProp := CloudStackCloudProperties{}
	err := cp.As(&csProp)
	if err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapErrorf(err, "[create_stemcell] error while reading cloud_properties")
	}
	if err = csProp.Validate(); err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapErrorf(err, "[create_stemcell] unable to validate cloud_properties")
	}

	// TODO [xmt]: handle light stemcell properly
	if len(csProp.LightTemplate) != 0 {
		return apiv1.StemcellCID{}, bosherr.WrapErrorf(err, "[create_stemcell] not handling light stemcell yet")
	}

	id := uuid.NewV4()
	name := fmt.Sprintf(config.TemplateNameFormat, id)
	parts := strings.Split(name, "-")
	name = strings.Join(parts[0:4], "-")

	zoneid, err := a.findZoneId()
	if err != nil {
		return apiv1.StemcellCID{}, err
	}

	ostypeid, err := a.findOsTypeId(a.config.CloudStack.Stemcell.OsType)
	if err != nil {
		return apiv1.StemcellCID{}, err
	}

	// TODO [xmt]: check disk format
	params := a.client.Template.NewGetUploadParamsForTemplateParams(
		name,
		config.TemplateFormat,
		config.Hypervisor,
		name,
		ostypeid,
		zoneid)
	params.SetIsextractable(true)
	params.SetRequireshvm(*a.config.CloudStack.Stemcell.RequiresHvm)
	params.SetBits(64)

	a.logger.Debug("create_stemcell", "requesting upload parameters : %#v", params)
	res, err := a.client.Template.GetUploadParamsForTemplate(params)
	if err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapErrorf(err, "[create_stemcell] could not get upload parameters")
	}

	request, err := NewFileUploadRequest(res.PostURL, "file", imagePath)
	if err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapErrorf(err, "[create_stemcell] could not prepare upload request for '%s'", imagePath)
	}

	request.Header.Set("X-signature", res.Signature)
	request.Header.Set("X-metadata", res.Metadata)
	request.Header.Set("X-expires", res.Expires)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	// client := &http.Client{}

	a.logger.Debug("create_stemcell", "uploading template to : %s", res.PostURL)
	uploadRes, err := client.Do(request)
	if err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapErrorf(err, "[create_stemcell] error while uploading file '%s'", imagePath)
	}

	if uploadRes.StatusCode != 200 {
		bodyBytes, _ := ioutil.ReadAll(uploadRes.Body)
		return apiv1.StemcellCID{}, fmt.Errorf("[create_stemcell] error while uploading file '%s' : %s", imagePath, string(bodyBytes))
	}

	a.logger.Debug("create_stemcell", "generated template id '%s' for stemcell '%s'", res.Id)
	return apiv1.NewStemcellCID(res.Id), nil
}

// NewFileUploadRequest -
func NewFileUploadRequest(uri string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileContents)
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return http.NewRequest("POST", uri, body)
}
