package packngo

import (
	"fmt"
	"strings"
)

const virtualNetworkBasePath = "/virtual-networks"

// DevicePortService handles operations on a port which belongs to a particular device
type ProjectVirtualNetworkService interface {
	List(*VirtualNetworkListInput) (*VirtualNetworkListOutput, *Response, error)
	Create(*VirtualNetworkCreateInput) (*VirtualNetworkCreateOutput, *Response, error)
	Delete(*VirtualNetworkDeleteInput) (*VirtualNetworkDeleteOutput, *Response, error)
}

type VirtualNetwork struct {
	ID           string `json:"id"`
	Description  string `json:"description,omitempty"`
	Vxlan        int    `json:"vxlan,omitempty"`
	FacilityCode string `json:"facility_code,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	Href         string `json:"href"`
}

type ProjectVirtualNetworkServiceOp struct {
	client *Client
}

type VirtualNetworkListInput struct {
	ProjectId string
	Includes  []string
}

type VirtualNetworkListOutput struct {
	VirtualNetworks []VirtualNetwork `json:"virtual_networks"`
}

func (i *ProjectVirtualNetworkServiceOp) List(input *VirtualNetworkListInput) (*VirtualNetworkListOutput, *Response, error) {
	var path string
	if input.Includes != nil {
		path = fmt.Sprintf("%s/%s%s?include=%s",
			projectBasePath, input.ProjectId, virtualNetworkBasePath, strings.Join(input.Includes, ","))
	} else {
		path = fmt.Sprintf("%s/%s%s",
			projectBasePath, input.ProjectId, virtualNetworkBasePath)
	}
	output := new(VirtualNetworkListOutput)

	resp, err := i.client.DoRequest("GET", path, input, output)
	if err != nil {
		return nil, nil, err
	}

	return output, resp, nil
}

type VirtualNetworkCreateInput struct {
	ProjectId   string `json:"project_id"`
	Description string `json:"description"`
	Facility    string `json:"facility"`
	Vxlan       int    `json:"vxlan"`
	Vlan        int    `json:"vlan"`
}

type VirtualNetworkCreateOutput struct {
	VirtualNetwork VirtualNetwork `json:"virtual_networks"`
}

func (i *ProjectVirtualNetworkServiceOp) Create(input *VirtualNetworkCreateInput) (*VirtualNetworkCreateOutput, *Response, error) {
	// TODO: May need to add timestamp to output from 'post' request
	// for the 'created_at' attribute of VirtualNetwork struct since
	// API response doesn't include it
	path := fmt.Sprintf("%s/%s%s", projectBasePath, input.ProjectId, virtualNetworkBasePath)
	output := new(VirtualNetworkCreateOutput)

	resp, err := i.client.DoRequest("POST", path, input, output)
	if err != nil {
		return nil, nil, err
	}

	return output, resp, nil
}

type VirtualNetworkDeleteInput struct {
	VirtualNetworkId string
}

type VirtualNetworkDeleteOutput struct {
	VirtualNetwork VirtualNetwork `json:"virtual_networks"`
}

func (i *ProjectVirtualNetworkServiceOp) Delete(input *VirtualNetworkDeleteInput) (*VirtualNetworkDeleteOutput, *Response, error) {
	path := fmt.Sprintf("%s/%s", virtualNetworkBasePath, input.VirtualNetworkId)
	output := new(VirtualNetworkDeleteOutput)

	resp, err := i.client.DoRequest("DELETE", path, input, output)
	if err != nil {
		return nil, nil, err
	}

	return output, resp, nil
}
