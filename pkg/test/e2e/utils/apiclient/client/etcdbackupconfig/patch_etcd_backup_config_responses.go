// Code generated by go-swagger; DO NOT EDIT.

package etcdbackupconfig

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"k8c.io/kubermatic/v2/pkg/test/e2e/utils/apiclient/models"
)

// PatchEtcdBackupConfigReader is a Reader for the PatchEtcdBackupConfig structure.
type PatchEtcdBackupConfigReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PatchEtcdBackupConfigReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewPatchEtcdBackupConfigOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewPatchEtcdBackupConfigUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewPatchEtcdBackupConfigForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewPatchEtcdBackupConfigDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewPatchEtcdBackupConfigOK creates a PatchEtcdBackupConfigOK with default headers values
func NewPatchEtcdBackupConfigOK() *PatchEtcdBackupConfigOK {
	return &PatchEtcdBackupConfigOK{}
}

/* PatchEtcdBackupConfigOK describes a response with status code 200, with default header values.

EtcdBackupConfig
*/
type PatchEtcdBackupConfigOK struct {
	Payload *models.EtcdBackupConfig
}

func (o *PatchEtcdBackupConfigOK) Error() string {
	return fmt.Sprintf("[PATCH /api/v2/projects/{project_id}/clusters/{cluster_id}/etcdbackupconfigs/{ebc_id}][%d] patchEtcdBackupConfigOK  %+v", 200, o.Payload)
}
func (o *PatchEtcdBackupConfigOK) GetPayload() *models.EtcdBackupConfig {
	return o.Payload
}

func (o *PatchEtcdBackupConfigOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.EtcdBackupConfig)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPatchEtcdBackupConfigUnauthorized creates a PatchEtcdBackupConfigUnauthorized with default headers values
func NewPatchEtcdBackupConfigUnauthorized() *PatchEtcdBackupConfigUnauthorized {
	return &PatchEtcdBackupConfigUnauthorized{}
}

/* PatchEtcdBackupConfigUnauthorized describes a response with status code 401, with default header values.

EmptyResponse is a empty response
*/
type PatchEtcdBackupConfigUnauthorized struct {
}

func (o *PatchEtcdBackupConfigUnauthorized) Error() string {
	return fmt.Sprintf("[PATCH /api/v2/projects/{project_id}/clusters/{cluster_id}/etcdbackupconfigs/{ebc_id}][%d] patchEtcdBackupConfigUnauthorized ", 401)
}

func (o *PatchEtcdBackupConfigUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewPatchEtcdBackupConfigForbidden creates a PatchEtcdBackupConfigForbidden with default headers values
func NewPatchEtcdBackupConfigForbidden() *PatchEtcdBackupConfigForbidden {
	return &PatchEtcdBackupConfigForbidden{}
}

/* PatchEtcdBackupConfigForbidden describes a response with status code 403, with default header values.

EmptyResponse is a empty response
*/
type PatchEtcdBackupConfigForbidden struct {
}

func (o *PatchEtcdBackupConfigForbidden) Error() string {
	return fmt.Sprintf("[PATCH /api/v2/projects/{project_id}/clusters/{cluster_id}/etcdbackupconfigs/{ebc_id}][%d] patchEtcdBackupConfigForbidden ", 403)
}

func (o *PatchEtcdBackupConfigForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewPatchEtcdBackupConfigDefault creates a PatchEtcdBackupConfigDefault with default headers values
func NewPatchEtcdBackupConfigDefault(code int) *PatchEtcdBackupConfigDefault {
	return &PatchEtcdBackupConfigDefault{
		_statusCode: code,
	}
}

/* PatchEtcdBackupConfigDefault describes a response with status code -1, with default header values.

errorResponse
*/
type PatchEtcdBackupConfigDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the patch etcd backup config default response
func (o *PatchEtcdBackupConfigDefault) Code() int {
	return o._statusCode
}

func (o *PatchEtcdBackupConfigDefault) Error() string {
	return fmt.Sprintf("[PATCH /api/v2/projects/{project_id}/clusters/{cluster_id}/etcdbackupconfigs/{ebc_id}][%d] patchEtcdBackupConfig default  %+v", o._statusCode, o.Payload)
}
func (o *PatchEtcdBackupConfigDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *PatchEtcdBackupConfigDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
