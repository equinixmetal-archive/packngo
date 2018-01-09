package packngo

import (
	"fmt"
	"strings"
)

const portBasePath = "/ports"

// DevicePortService handles operations on a port which belongs to a particular device
type DevicePortService interface {
	Assign(*PortAssignInput) (*PortAssignOutput, *Response, bool, error)
	Unassign(*PortUnassignInput) (*PortUnassignOutput, *Response, bool, error)
	Bond(*PortBondInput) (*PortBondOutput, *Response, error)
	Disbond(*PortDisbondInput) (*PortDisbondOutput, *Response, error)
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

type PortAssignInput struct {
	DeviceId         string
	PortId           string
	VirtualNetworkId int `json:"vnid"`
}

type PortAssignOutput struct {
	PortId          string           `json:"id"`
	VirtualNetworks []VirtualNetwork `json:"virtual_networks"`
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

type PortUnassignInput struct {
	DeviceId         string
	PortId           string
	VirtualNetworkId int `json:"vnid"`
}

type PortUnassignOutput struct {
	PortId          string           `json:"id"`
	VirtualNetworks []VirtualNetwork `json:"virtual_networks"`
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

type PortBondInput struct {
	PortId     string
	BulkEnable bool
}

type PortBondOutput struct {
	Port Port
}

func (i *DevicePortServiceOp) Bond(input *PortBondInput) (*PortBondOutput, *Response, error) {
	var path string
	if input.BulkEnable {
		path = fmt.Sprintf("%s/%s/bond?bulk_enable=true", portBasePath, input.PortId)
	} else {
		path = fmt.Sprintf("%s/%s/bond", portBasePath, input.PortId)
	}
	output := new(PortBondOutput)

	resp, err := i.client.DoRequest("POST", path, input, output)
	if err != nil {
		return nil, nil, err
	}

	return output, resp, nil
}

type PortDisbondInput struct {
	PortId      string
	BulkDisable bool
}

type PortDisbondOutput struct {
	Port Port
}

func (i *DevicePortServiceOp) Disbond(input *PortDisbondInput) (*PortDisbondOutput, *Response, error) {
	var path string
	if input.BulkDisable {
		path = fmt.Sprintf("%s/%s/bond?bulk_disable=true", portBasePath, input.PortId)
	} else {
		path = fmt.Sprintf("%s/%s/bond", portBasePath, input.PortId)
	}
	output := new(PortDisbondOutput)

	resp, err := i.client.DoRequest("POST", path, input, output)
	if err != nil {
		return nil, nil, err
	}

	return output, resp, nil
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

// Private helper methods
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
