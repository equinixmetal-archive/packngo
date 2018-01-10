package packngo

import (
	"fmt"
)

const portBasePath = "/ports"

// DevicePortService handles operations on a port which belongs to a particular device
type DevicePortService interface {
	Assign(*PortAssignRequest) (*PortAssignResponse, *Response, bool, error)
	Unassign(*PortUnassignRequest) (*PortUnassignResponse, *Response, bool, error)
	Bond(*PortBondRequest) (*Port, *Response, error)
	Disbond(*PortDisbondRequest) (*Port, *Response, error)
	GetBondedPort(string) (*Port, bool, error)
}

type Port struct {
	ID                      string           `json:"id"`
	Type                    string           `json:"type"`
	Name                    string           `json:"name"`
	AttachedVirtualNetworks []VirtualNetwork `json:"virtual_networks"`
}

type DevicePortServiceOp struct {
	client *Client
}

type PortAssignRequest struct {
	DeviceID         string
	PortID           string
	VirtualNetworkID int `json:"vnid"`
}

type PortAssignResponse struct {
	PortID          string           `json:"id"`
	VirtualNetworks []VirtualNetwork `json:"virtual_networks"`
}

// Assign associates virtual networks to a port
func (i *DevicePortServiceOp) Assign(input *PortAssignRequest) (*PortAssignResponse, *Response, bool, error) {
	// First get the device information in order to determine if this is the first VLAN assigned to this port.
	// Requires a conversion to layer-2
	device, _, err := i.client.Devices.GetWith(input.DeviceID, []string{"virtual_networks"})

	for _, port := range device.NetworkPorts {
		if port.ID != input.PortID || port.hasVirtualNetwork(input.VirtualNetworkID) {
			continue
		}
		if len(port.AttachedVirtualNetworks) == 0 {
			// convert to layer-2 (and attach vlan)
			return i.convertToLayerTwo(input)
		} else {
			// not the first VLAN, so attach without converting
			return i.assignVirtualNetwork(input)
		}
	}

	return nil, nil, false, err
}

type PortUnassignRequest struct {
	DeviceID         string
	PortID           string
	VirtualNetworkID int `json:"vnid"`
}

type PortUnassignResponse struct {
	PortID          string           `json:"id"`
	VirtualNetworks []VirtualNetwork `json:"virtual_networks"`
}

func (i *DevicePortServiceOp) Unassign(input *PortUnassignRequest) (*PortUnassignResponse, *Response, bool, error) {
	path := fmt.Sprintf("%s/%s/unassign", portBasePath, input.PortID)
	unassignResponse := new(PortUnassignResponse)

	resp, err := i.client.DoRequest("POST", path, input, unassignResponse)
	if err != nil {
		return nil, resp, false, err
	}

	return unassignResponse, resp, true, err
}

type PortBondRequest struct {
	PortID     string
	BulkEnable bool
}

func (i *DevicePortServiceOp) Bond(input *PortBondRequest) (*Port, *Response, error) {
	path := fmt.Sprintf("%s/%s/bond", portBasePath, input.PortID)
	if input.BulkEnable {
		path += "?bulk_enable=true"
	}
	output := new(Port)

	resp, err := i.client.DoRequest("POST", path, input, output)
	if err != nil {
		return nil, nil, err
	}

	return output, resp, nil
}

type PortDisbondRequest struct {
	PortID      string
	BulkDisable bool
}

func (i *DevicePortServiceOp) Disbond(input *PortDisbondRequest) (*Port, *Response, error) {
	path := fmt.Sprintf("%s/%s/bond", portBasePath, input.PortID)
	if input.BulkDisable {
		path += "?bulk_disable=true"
	}
	output := new(Port)

	resp, err := i.client.DoRequest("POST", path, input, output)
	if err != nil {
		return nil, nil, err
	}

	return output, resp, nil
}

func (i *DevicePortServiceOp) GetBondedPort(deviceID string) (*Port, bool, error) {
	device, _, err := i.client.Devices.Get(deviceID)
	for _, port := range device.NetworkPorts {
		if port.Type == "NetworkBondPort" {
			return &port, true, nil
		}
	}

	return nil, false, err
}

// Private helper methods
func (i *DevicePortServiceOp) assignVirtualNetwork(input *PortAssignRequest) (*PortAssignResponse, *Response, bool, error) {
	path := fmt.Sprintf("%s/%s/assign", portBasePath, input.PortID)
	assignResponse := new(PortAssignResponse)

	resp, err := i.client.DoRequest("POST", path, input, assignResponse)
	if err != nil {
		return nil, resp, false, err
	}

	return assignResponse, resp, true, err
}

func (i *DevicePortServiceOp) convertToLayerTwo(input *PortAssignRequest) (*PortAssignResponse, *Response, bool, error) {
	path := fmt.Sprintf("%s/%s/convert/layer-2", portBasePath, input.PortID)
	assignResponse := new(PortAssignResponse)

	resp, err := i.client.DoRequest("POST", path, input, assignResponse)
	if err != nil {
		return nil, resp, false, err
	}

	return assignResponse, resp, true, err
}

func (p *Port) hasVirtualNetwork(vnid int) bool {
	for i := range p.AttachedVirtualNetworks {
		if p.AttachedVirtualNetworks[i].Vxlan == vnid {
			return true
		}
	}

	return false
}
