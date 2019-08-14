// Code generated by go-swagger; DO NOT EDIT.

package serviceaccounts

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new serviceaccounts API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for serviceaccounts API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
AddServiceAccountToProject Adds the given service account to the given project
*/
func (a *Client) AddServiceAccountToProject(params *AddServiceAccountToProjectParams, authInfo runtime.ClientAuthInfoWriter) (*AddServiceAccountToProjectCreated, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewAddServiceAccountToProjectParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "addServiceAccountToProject",
		Method:             "POST",
		PathPattern:        "/api/v1/projects/{project_id}/serviceaccounts",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &AddServiceAccountToProjectReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*AddServiceAccountToProjectCreated), nil

}

/*
DeleteServiceAccount Deletes service account for the given project
*/
func (a *Client) DeleteServiceAccount(params *DeleteServiceAccountParams, authInfo runtime.ClientAuthInfoWriter) (*DeleteServiceAccountOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeleteServiceAccountParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "deleteServiceAccount",
		Method:             "DELETE",
		PathPattern:        "/api/v1/projects/{project_id}/serviceaccounts/{serviceaccount_id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &DeleteServiceAccountReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*DeleteServiceAccountOK), nil

}

/*
ListServiceAccounts List Service Accounts for the given project
*/
func (a *Client) ListServiceAccounts(params *ListServiceAccountsParams, authInfo runtime.ClientAuthInfoWriter) (*ListServiceAccountsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListServiceAccountsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "listServiceAccounts",
		Method:             "GET",
		PathPattern:        "/api/v1/projects/{project_id}/serviceaccounts",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListServiceAccountsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*ListServiceAccountsOK), nil

}

/*
UpdateServiceAccount Updates service account for the given project
*/
func (a *Client) UpdateServiceAccount(params *UpdateServiceAccountParams, authInfo runtime.ClientAuthInfoWriter) (*UpdateServiceAccountOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewUpdateServiceAccountParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "updateServiceAccount",
		Method:             "PUT",
		PathPattern:        "/api/v1/projects/{project_id}/serviceaccounts/{serviceaccount_id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &UpdateServiceAccountReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*UpdateServiceAccountOK), nil

}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
