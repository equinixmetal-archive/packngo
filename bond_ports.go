package packngo

import "fmt"

const NetworkTypeHybridBonded = "hybrid-bonded"

func (d *Device) GetBondNetworkType(portName string) string {
	for _, p := range d.NetworkPorts {
		if p.Name == portName {
			return p.NetworkType
		}
	}
	return ""
}

func (d *Device) GetEthPortsInBond(name string) []*Port {
	ports := []*Port{}
	for _, port := range d.NetworkPorts {
		if port.Bond != nil && port.Bond.Name == name {
			ports = append(ports, &port)
		}
	}
	return ports
}

func BondStateTransitionNecessary(type1, type2 string) bool {
	if type1 == type2 {
		return false
	}
	if type1 == NetworkTypeHybridBonded && type2 == NetworkTypeL3 {
		return false
	}
	if type2 == NetworkTypeHybridBonded && type1 == NetworkTypeL3 {
		return false
	}
	return true
}

func (i *DevicePortServiceOp) BondToNetworkType(deviceID, bondPortName, targetType string) (*Device, error) {
	d, _, err := i.client.Devices.Get(deviceID, nil)
	if err != nil {
		return nil, err
	}

	curType := d.GetBondNetworkType(bondPortName)

	if !BondStateTransitionNecessary(curType, targetType) {
		return nil, fmt.Errorf("Bond doesn't need to be converted from %s to %s", curType, targetType)
	}

	err = i.ConvertDeviceBond(d, bondPortName, targetType)
	if err != nil {
		return nil, err
	}

	d, _, err = i.client.Devices.Get(deviceID, nil)

	if err != nil {
		return nil, err
	}

	finalType := d.GetNetworkType()

	if BondStateTransitionNecessary(finalType, targetType) {
		return nil, fmt.Errorf(
			"Failed to convert %s on device %s from %s to %s. New type was %s",
			bondPortName, deviceID, curType, targetType, finalType)

	}
	return d, err
}

func (i *DevicePortServiceOp) ConvertDeviceBond(d *Device, bondPortName, targetType string) error {
	bondPort, err := d.GetPortByName(bondPortName)
	if err != nil {
		return err
	}

	if targetType == NetworkTypeL3 || targetType == NetworkTypeHybridBonded {
		_, _, err := i.Bond(bondPort, false)
		if err != nil {
			return err
		}
		_, _, err = i.PortToLayerThree(d.ID, bondPortName)
		if err != nil {
			return err
		}
		// device needs to be refreshed, the bond and convert calls might bond eths
		d, _, err := i.client.Devices.Get(d.ID, nil)
		if err != nil {
			return err
		}
		for _, p := range d.GetEthPortsInBond(bondPortName) {
			_, _, err := i.Bond(p, false)
			if err != nil {
				return err
			}
		}
	}
	if targetType == NetworkTypeHybrid {
		_, _, err := i.Bond(bondPort, false)
		if err != nil {
			return err
		}
		_, _, err = i.PortToLayerThree(d.ID, bondPortName)
		if err != nil {
			return err
		}

		// device needs to be refreshed, the bond and convert calls might bond eths
		d, _, err := i.client.Devices.Get(d.ID, nil)
		if err != nil {
			return err
		}
		ethLatter := d.GetEthPortsInBond(bondPortName)[1]

		_, _, err = i.Disbond(ethLatter, false)
		if err != nil {
			return err
		}
	}
	if targetType == NetworkTypeL2Individual {
		_, _, err := i.PortToLayerTwo(d.ID, bondPortName)
		if err != nil {
			return err
		}
		// device needs to be refreshed, the convert call might break the bond
		d, _, err := i.client.Devices.Get(d.ID, nil)
		if err != nil {
			return err
		}
		bondPort, err := d.GetPortByName(bondPortName)
		if err != nil {
			return err
		}
		_, _, err = i.Disbond(bondPort, true)
		if err != nil {
			return err
		}
	}
	if targetType == NetworkTypeL2Bonded {
		_, _, err := i.PortToLayerTwo(d.ID, bondPortName)
		if err != nil {
			return err
		}
		// device needs to be refreshed, the convert call might break the bond
		d, _, err := i.client.Devices.Get(d.ID, nil)
		if err != nil {
			return err
		}
		for _, p := range d.GetEthPortsInBond(bondPortName) {
			_, _, err := i.Bond(p, false)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
