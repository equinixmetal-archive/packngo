package packngo

import (
	"fmt"
	"strings"
)

const virtualNetworkBasePath = "/virtual-networks"

// DevicePortService handles operations on a port which belongs to a particular device
type ProjectVirtualNetworkService interface {
	List(*VirtualNetworkListRequest) (*VirtualNetworkListResponse, *Response, error)
	Create(*VirtualNetworkCreateRequest) (*VirtualNetworkCreateResponse, *Response, error)
	Delete(*VirtualNetworkDeleteRequest) (*VirtualNetworkDeleteResponse, *Response, error)
}

type VirtualNetwork struct {
	ID           string `json:"id"`
	Description  string `json:"description,omitempty"`
	VXLAN        int    `json:"vxlan,omitempty"`
	FacilityCode string `json:"facility_code,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	Href         string `json:"href"`
}

type ProjectVirtualNetworkServiceOp struct {
	client *Client
}

type VirtualNetworkListRequest struct {
	ProjectID string
	Includes  []string
}

type VirtualNetworkListResponse struct {
	VirtualNetworks []VirtualNetwork `json:"virtual_networks"`
}

func (i *ProjectVirtualNetworkServiceOp) List(input *VirtualNetworkListRequest) (*VirtualNetworkListResponse, *Response, error) {
	path := fmt.Sprintf("%s/%s%s", projectBasePath, input.ProjectID, virtualNetworkBasePath)
	if input.Includes != nil {
		path += fmt.Sprintf("?include=%s", strings.Join(input.Includes, ","))
	}
	output := new(VirtualNetworkListResponse)

	resp, err := i.client.DoRequest("GET", path, input, output)
	if err != nil {
		return nil, nil, err
	}

	return output, resp, nil
}

type VirtualNetworkCreateRequest struct {
	ProjectID   string `json:"project_id"`
	Description string `json:"description"`
	Facility    string `json:"facility"`
	VXLAN       int    `json:"vxlan"`
	VLAN        int    `json:"vlan"`
}

type VirtualNetworkCreateResponse struct {
	VirtualNetwork VirtualNetwork `json:"virtual_networks"`
}

func (i *ProjectVirtualNetworkServiceOp) Create(input *VirtualNetworkCreateRequest) (*VirtualNetworkCreateResponse, *Response, error) {
	// TODO: May need to add timestamp to output from 'post' request
	// for the 'created_at' attribute of VirtualNetwork struct since
	// API response doesn't include it
	path := fmt.Sprintf("%s/%s%s", projectBasePath, input.ProjectID, virtualNetworkBasePath)
	output := new(VirtualNetworkCreateResponse)

	resp, err := i.client.DoRequest("POST", path, input, output)
	if err != nil {
		return nil, nil, err
	}

	return output, resp, nil
}

type VirtualNetworkDeleteRequest struct {
	VirtualNetworkID string
}

type VirtualNetworkDeleteResponse struct {
	VirtualNetwork VirtualNetwork `json:"virtual_networks"`
}

func (i *ProjectVirtualNetworkServiceOp) Delete(input *VirtualNetworkDeleteRequest) (*VirtualNetworkDeleteResponse, *Response, error) {
	path := fmt.Sprintf("%s/%s", virtualNetworkBasePath, input.VirtualNetworkID)
	output := new(VirtualNetworkDeleteResponse)

	resp, err := i.client.DoRequest("DELETE", path, input, output)
	if err != nil {
		return nil, nil, err
	}

	return output, resp, nil
}
