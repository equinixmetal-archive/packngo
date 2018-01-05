package packngo

import (
	"fmt"
)

const portBasePath = "/ports"

type Port struct {
	ID                      string           `json:"id"`
	Type                    string           `json:"type"`
	Name                    string           `json:"name"`
	AttachedVirtualNetworks []VirtualNetwork `json:"virtual_networks"`
}

type VirtualNetwork struct {
	ID           string `json:"id,omitempty"`
	Description  string `json:"description,omitempty"`
	Vxlan        int    `json:"vxlan,omitempty"`
	FacilityCode string `json:"facility_code,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	Href         string `json:"href"`
}

type PortAssignInput struct {
	DeviceId         string
	PortId           string
	VirtualNetworkId int `json:"vnid"`
}

type PortAssignOutput struct {
	PortId          string           `json:"id"`
	VirtualNetworks []VirtualNetwork `json:"virtual_networks"`
}

type PortUnassignInput struct {
	DeviceId         string
	PortId           string
	VirtualNetworkId int `json:"vnid"`
}

type PortUnassignOutput struct {
	PortId          string           `json:"id"`
	VirtualNetworks []VirtualNetwork `json:"virtual_networks"`
}

type DevicePortServiceOp struct {
	client *Client
}

// DevicePortService handles operations on a port which belongs to a particular device
type DevicePortService interface {
	Assign(*PortAssignInput) (*PortAssignOutput, *Response, bool, error)
	Unassign(*PortUnassignInput) (*PortUnassignOutput, *Response, bool, error)
	GetBondedPort(string) (*Port, bool, error)
}

// Assign associates virtual networks to a port
func (i *DevicePortServiceOp) Assign(input *PortAssignInput) (*PortAssignOutput, *Response, bool, error) {
	// First get the device information in order to determine if this is the first vlan assigned to this port.
	// Requires a conversion to layer-2
	device, _, err := i.client.Devices.GetWith(input.DeviceId, []string{"virtual_networks"})

	// No network ports for this device so no-op
	if len(device.NetworkPorts) == 0 {
		return nil, nil, false, err
	}
	for index := range device.NetworkPorts {
		if port := device.NetworkPorts[index]; port.ID != input.PortId || port.hasVirtualNetwork(input.VirtualNetworkId) {
			continue
		} else {
			if len(port.AttachedVirtualNetworks) == 0 {
				// convert to layer-3 (and attach vlan)
				return i.convertLayerTwo(input)
			} else {
				// not the first vlan, so attach without converting
				return i.assignVirtualNetwork(input)
			}
		}
	}
	return nil, nil, false, err
}

func (i *DevicePortServiceOp) assignVirtualNetwork(input *PortAssignInput) (*PortAssignOutput, *Response, bool, error) {
	path := fmt.Sprintf("%s/%s/assign", portBasePath, input.PortId)
	assignOutput := new(PortAssignOutput)

	resp, err := i.client.DoRequest("POST", path, input, assignOutput)
	if err != nil {
		return nil, resp, false, err
	}

	return assignOutput, resp, true, err
}

func (i *DevicePortServiceOp) convertLayerTwo(input *PortAssignInput) (*PortAssignOutput, *Response, bool, error) {
	path := fmt.Sprintf("%s/%s/convert/layer-2", portBasePath, input.PortId)
	assignOutput := new(PortAssignOutput)

	resp, err := i.client.DoRequest("POST", path, input, assignOutput)
	if err != nil {
		return nil, resp, false, err
	}

	return assignOutput, resp, true, err
}

func (p *Port) hasVirtualNetwork(vnid int) bool {
	for i := range p.AttachedVirtualNetworks {
		if p.AttachedVirtualNetworks[i].Vxlan == vnid {
			return true
		}
	}
	return false
}

func (i *DevicePortServiceOp) Unassign(input *PortUnassignInput) (*PortUnassignOutput, *Response, bool, error) {
	path := fmt.Sprintf("%s/%s/unassign", portBasePath, input.PortId)
	unassignOutput := new(PortUnassignOutput)

	resp, err := i.client.DoRequest("POST", path, input, unassignOutput)
	if err != nil {
		return nil, resp, false, err
	}

	return unassignOutput, resp, true, err
}

func (i *DevicePortServiceOp) GetBondedPort(deviceId string) (*Port, bool, error) {
	device, _, err := i.client.Devices.Get(deviceId)
	if len(device.NetworkPorts) == 0 {
		return nil, false, err
	}
	for index := range device.NetworkPorts {

		if port := device.NetworkPorts[index]; port.Type == "NetworkBondPort" {
			return &port, true, nil
		}
	}
	return nil, false, err
}
