package packngo

import (
	"fmt"
)

const portBasePath = "/ports"

// DevicePortService handles operations on a port which belongs to a particular device
type DevicePortService interface {
	Assign(*PortAssignRequest) (*Port, *Response, error)
	Unassign(*PortAssignRequest) (*Port, *Response, error)
	Bond(*BondRequest) (*Port, *Response, error)
	Disbond(*DisbondRequest) (*Port, *Response, error)
	ConvertToLayerTwo(string) (*Port, *Response, error)
	GetBondedPort(string) (*Port, error)
	GetPortByName(string, string) (*Port, error)
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
	PortID           string `json:"id"`
	VirtualNetworkID string `json:"vnid"`
}

type BondRequest struct {
	PortID     string `json:"id"`
	BulkEnable bool   `json:"bulk_enable"`
}

type DisbondRequest struct {
	PortID      string `json:"id"`
	BulkDisable bool   `json:"bulk_disable"`
}

func (i *DevicePortServiceOp) GetBondedPort(deviceID string) (*Port, error) {
	device, _, err := i.client.Devices.Get(deviceID)
	if err != nil {
		return nil, err
	}
	for _, port := range device.NetworkPorts {
		if port.Type == "NetworkBondPort" {
			return &port, nil
		}
	}

	return nil, fmt.Errorf("No bonded port found in device %s", deviceID)
}

func (i *DevicePortServiceOp) GetPortByName(deviceID, name string) (*Port, error) {
	device, _, err := i.client.Devices.Get(deviceID)
	if err != nil {
		return nil, err
	}
	for _, port := range device.NetworkPorts {
		if port.Name == name {
			return &port, nil
		}
	}

	return nil, fmt.Errorf("Port %s not found in device %s", name, deviceID)
}

func (i *DevicePortServiceOp) Assign(par *PortAssignRequest) (*Port, *Response, error) {
	path := fmt.Sprintf("%s/%s/assign", portBasePath, par.PortID)
	return i.portAction(path, par)
}

func (i *DevicePortServiceOp) Unassign(par *PortAssignRequest) (*Port, *Response, error) {
	path := fmt.Sprintf("%s/%s/unassign", portBasePath, par.PortID)
	return i.portAction(path, par)
}

func (i *DevicePortServiceOp) Bond(br *BondRequest) (*Port, *Response, error) {
	path := fmt.Sprintf("%s/%s/bond", portBasePath, br.PortID)
	return i.portAction(path, br)
}

func (i *DevicePortServiceOp) Disbond(dr *DisbondRequest) (*Port, *Response, error) {
	path := fmt.Sprintf("%s/%s/disbond", portBasePath, dr.PortID)
	return i.portAction(path, dr)
}

func (i *DevicePortServiceOp) portAction(path string, req interface{}) (*Port, *Response, error) {
	port := new(Port)

	resp, err := i.client.DoRequest("POST", path, req, port)
	if err != nil {
		return nil, resp, err
	}

	return port, resp, err
}

func (i *DevicePortServiceOp) ConvertToLayerTwo(portID string) (*Port, *Response, error) {
	path := fmt.Sprintf("%s/%s/convert/layer-2", portBasePath, portID)
	port := new(Port)

	resp, err := i.client.DoRequest("POST", path, nil, port)
	if err != nil {
		return nil, resp, err
	}

	return port, resp, err
}
