package packngo

import (
	"fmt"
	"strings"
)

const portBasePath = "/ports"

// DevicePortService handles operations on a port which belongs to a particular device
type DevicePortService interface {
	Assign(*PortAssignRequest) (*Port, *Response, error)
	Unassign(*PortAssignRequest) (*Port, *Response, error)
	AssignNative(*PortAssignRequest) (*Port, *Response, error)
	UnassignNative(string) (*Port, *Response, error)
	Bond(*Port, bool) (*Port, *Response, error)
	Disbond(*Port, bool) (*Port, *Response, error)
	DeviceToNetworkType(string, string) (*Device, error)
	DeviceNetworkType(string) (string, error)
	PortToLayerTwo(string, string) (*Port, *Response, error)
	PortToLayerThree(string, string) (*Port, *Response, error)
	GetPortByName(string, string) (*Port, error)
	GetOddEthPorts(*Device) (map[string]*Port, error)
	GetAllEthPorts(*Device) (map[string]*Port, error)
	ConvertDevice(*Device, string) error
}

type PortData struct {
	MAC    string `json:"mac"`
	Bonded bool   `json:"bonded"`
}

type BondData struct {
	ID string `json:"id"`

	// Name of the port interface for the bond ("bond0")
	Name string `json:"name"`
}

// Port is a hardware port associated with a reserved or instanciated hardware
// device.
type Port struct {
	ID string `json:"id"`

	// Type is either "NetworkBondPort" for bond ports or "NetworkPort" for
	// bondable ethernet ports
	Type string `json:"type"`

	// Name of the interface for this port (such as "bond0" or "eth0")
	Name string `json:"name"`

	Data PortData `json:"data"`

	// Indicates whether or not the bond can be broken on the port (when applicable).
	DisbondOperationSupported bool `json:"disbond_operation_supported,omitempty"`

	// NetworkType is either of layer2-bonded, layer2-individual, layer3,
	// hybrid, hybrid-bonded
	NetworkType string `json:"network_type,omitempty"`

	// The Native VLAN attached to the port
	// <https://metal.equinix.com/developers/docs/layer2-networking/native-vlan>
	NativeVirtualNetwork *VirtualNetwork `json:"native_virtual_network"`

	// VLANs attached to the port
	AttachedVirtualNetworks []VirtualNetwork `json:"virtual_networks"`

	Bond *BondData `json:"bond"`
}

type AddressRequest struct {
	AddressFamily int  `json:"address_family"`
	Public        bool `json:"public"`
}

type BackToL3Request struct {
	RequestIPs []AddressRequest `json:"request_ips"`
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

func (i *DevicePortServiceOp) GetPortByName(deviceID, name string) (*Port, error) {
	device, _, err := i.client.Devices.Get(deviceID, nil)
	if err != nil {
		return nil, err
	}
	return device.GetPortByName(name)
}

func (i *DevicePortServiceOp) Assign(par *PortAssignRequest) (*Port, *Response, error) {
	return i.client.Ports.Assign(par)
}

func (i *DevicePortServiceOp) AssignNative(par *PortAssignRequest) (*Port, *Response, error) {
	return i.client.Ports.AssignNative(par)
}

func (i *DevicePortServiceOp) UnassignNative(portID string) (*Port, *Response, error) {
	return i.client.Ports.UnassignNative(portID)
}

func (i *DevicePortServiceOp) Unassign(par *PortAssignRequest) (*Port, *Response, error) {
	return i.client.Ports.Unassign(par)
}

func (i *DevicePortServiceOp) Bond(p *Port, be bool) (*Port, *Response, error) {
	if p.Data.Bonded {
		return p, nil, nil
	}

	return i.client.Ports.Bond(p.ID, be)
}

func (i *DevicePortServiceOp) Disbond(p *Port, bd bool) (*Port, *Response, error) {
	if !p.Data.Bonded {
		return p, nil, nil
	}
	return i.client.Ports.Disbond(p.ID, bd)
}

func (i *DevicePortServiceOp) PortToLayerTwo(deviceID, portName string) (*Port, *Response, error) {
	p, err := i.GetPortByName(deviceID, portName)
	if err != nil {
		return nil, nil, err
	}
	if strings.HasPrefix(p.NetworkType, "layer2") {
		return p, nil, nil
	}

	return i.client.Ports.ConvertToLayerTwo(p.ID)
}

func (i *DevicePortServiceOp) PortToLayerThree(deviceID, portName string) (*Port, *Response, error) {
	p, err := i.GetPortByName(deviceID, portName)
	if err != nil {
		return nil, nil, err
	}
	if (p.NetworkType == NetworkTypeL3) || (p.NetworkType == NetworkTypeHybrid) {
		return p, nil, nil
	}

	ips := []AddressRequest{
		{AddressFamily: 4, Public: true},
		{AddressFamily: 4, Public: false},
		{AddressFamily: 6, Public: true},
	}

	return i.client.Ports.ConvertToLayerThree(p.ID, ips)
}

func (i *DevicePortServiceOp) DeviceNetworkType(deviceID string) (string, error) {
	d, _, err := i.client.Devices.Get(deviceID, nil)
	if err != nil {
		return "", err
	}
	return d.GetNetworkType(), nil
}

func (i *DevicePortServiceOp) GetAllEthPorts(d *Device) (map[string]*Port, error) {
	d, _, err := i.client.Devices.Get(d.ID, nil)
	if err != nil {
		return nil, err
	}
	return d.GetPhysicalPorts(), nil
}

func (i *DevicePortServiceOp) GetOddEthPorts(d *Device) (map[string]*Port, error) {
	d, _, err := i.client.Devices.Get(d.ID, nil)
	if err != nil {
		return nil, err
	}
	ret := map[string]*Port{}
	eth1, err := d.GetPortByName("eth1")
	if err != nil {
		return nil, err
	}
	ret["eth1"] = eth1

	eth3, err := d.GetPortByName("eth3")
	if err != nil {
		return ret, nil
	}
	ret["eth3"] = eth3
	return ret, nil

}

func (i *DevicePortServiceOp) ConvertDevice(d *Device, targetType string) error {
	bondPorts := d.GetBondPorts()

	if targetType == NetworkTypeL3 {
		// TODO: remove vlans from all the ports
		for _, p := range bondPorts {
			_, _, err := i.Bond(p, false)
			if err != nil {
				return err
			}
		}
		_, _, err := i.PortToLayerThree(d.ID, "bond0")
		if err != nil {
			return err
		}
		allEthPorts, err := i.GetAllEthPorts(d)
		if err != nil {
			return err
		}
		for _, p := range allEthPorts {
			_, _, err := i.Bond(p, false)
			if err != nil {
				return err
			}
		}
	}
	if targetType == NetworkTypeHybrid {
		for _, p := range bondPorts {
			_, _, err := i.Bond(p, false)
			if err != nil {
				return err
			}
		}

		_, _, err := i.PortToLayerThree(d.ID, "bond0")
		if err != nil {
			return err
		}

		// ports need to be refreshed before bonding/disbonding
		oddEthPorts, err := i.GetOddEthPorts(d)
		if err != nil {
			return err
		}

		for _, p := range oddEthPorts {
			_, _, err := i.Disbond(p, false)
			if err != nil {
				return err
			}
		}
	}
	if targetType == NetworkTypeL2Individual {
		_, _, err := i.PortToLayerTwo(d.ID, "bond0")
		if err != nil {
			return err
		}
		for _, p := range bondPorts {
			_, _, err = i.Disbond(p, true)
			if err != nil {
				return err
			}
		}
	}
	if targetType == NetworkTypeL2Bonded {

		for _, p := range bondPorts {
			_, _, err := i.PortToLayerTwo(d.ID, p.Name)
			if err != nil {
				return err
			}
		}
		allEthPorts, err := i.GetAllEthPorts(d)
		if err != nil {
			return err
		}
		for _, p := range allEthPorts {
			_, _, err := i.Bond(p, false)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (i *DevicePortServiceOp) DeviceToNetworkType(deviceID string, targetType string) (*Device, error) {
	d, _, err := i.client.Devices.Get(deviceID, nil)
	if err != nil {
		return nil, err
	}

	curType := d.GetNetworkType()

	if curType == targetType {
		return nil, fmt.Errorf("Device already is in state %s", targetType)
	}
	err = i.ConvertDevice(d, targetType)
	if err != nil {
		return nil, err
	}

	d, _, err = i.client.Devices.Get(deviceID, nil)

	if err != nil {
		return nil, err
	}

	finalType := d.GetNetworkType()

	if finalType != targetType {
		return nil, fmt.Errorf(
			"Failed to convert device %s from %s to %s. New type was %s",
			deviceID, curType, targetType, finalType)

	}
	return d, err
}
