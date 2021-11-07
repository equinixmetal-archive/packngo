package metadata

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
)

const BaseURL = "https://metadata.platformequinix.com"

func GetMetadata() (*CurrentDevice, error) {
	return GetMetadataFromURL(BaseURL)
}

func GetMetadataFromURL(baseURL string) (*CurrentDevice, error) {
	res, err := http.Get(baseURL + "/metadata")
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	var result struct {
		Error string `json:"error"`
		*CurrentDevice
	}
	if err := json.Unmarshal(b, &result); err != nil {
		if res.StatusCode >= 400 {
			return nil, errors.New(res.Status)
		}
		return nil, err
	}
	if result.Error != "" {
		return nil, errors.New(result.Error)
	}
	return result.CurrentDevice, nil
}

func GetUserData() ([]byte, error) {
	return GetUserDataFromURL(BaseURL)
}

func GetUserDataFromURL(baseURL string) ([]byte, error) {
	res, err := http.Get(baseURL + "/userdata")
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	return b, err
}

type AddressFamily int

const (
	IPv4 = AddressFamily(4)
	IPv6 = AddressFamily(6)
)

type AddressInfo struct {
	ID          string        `json:"id"`
	Family      AddressFamily `json:"address_family"`
	Enabled     bool          `json:"enabled"`
	Public      bool          `json:"public"`
	Management  bool          `json:"management"`
	Address     net.IP        `json:"address"`
	NetworkMask net.IP        `json:"netmask"`
	Gateway     net.IP        `json:"gateway"`
	NetworkBits int           `json:"cidr"`
	Network     net.IP        `json:"network"`
	ParentBlock struct {
		Network net.IP `json:"network"`
		NetMask net.IP `json:"netmask"`
		CIDR    int    `json:"cidr"`
		Href    string `json:"href"`
	} `json:"parent_block"`
	CreatedAt Timestamp `json:"created_at,omitempty"`
}

type BondingMode int

const (
	BondingBalanceRR    = BondingMode(0)
	BondingActiveBackup = BondingMode(1)
	BondingBalanceXOR   = BondingMode(2)
	BondingBroadcast    = BondingMode(3)
	BondingLACP         = BondingMode(4)
	BondingBalanceTLB   = BondingMode(5)
	BondingBalanceALB   = BondingMode(6)
)

var bondingModeStrings = map[BondingMode]string{
	BondingBalanceRR:    "balance-rr",
	BondingActiveBackup: "active-backup",
	BondingBalanceXOR:   "balance-xor",
	BondingBroadcast:    "broadcast",
	BondingLACP:         "802.3ad",
	BondingBalanceTLB:   "balance-tlb",
	BondingBalanceALB:   "balance-alb",
}

func (m BondingMode) String() string {
	if str, ok := bondingModeStrings[m]; ok {
		return str
	}
	return fmt.Sprintf("%d", m)
}

type PrivateSubnet struct {
	// TODO(displague): Discover the fields
}

// Storage is identical to packngo.CPR
type Storage struct {
	Disks []struct {
		Device     string `json:"device"`
		WipeTable  bool   `json:"wipeTable"`
		Partitions []struct {
			Label  string `json:"label"`
			Number int    `json:"number"`
			Size   string `json:"size"`
		} `json:"partitions"`
	} `json:"disks"`
	Raid []struct {
		Devices []string `json:"devices"`
		Level   string   `json:"level"`
		Name    string   `json:"name"`
	} `json:"raid,omitempty"`
	Filesystems []struct {
		Mount struct {
			Device string `json:"device"`
			Format string `json:"format"`
			Point  string `json:"point"`
			Create struct {
				Options []string `json:"options"`
			} `json:"create"`
		} `json:"mount"`
	} `json:"filesystems"`
}

type Specs struct {
	CPUs []struct {
		Count int    `json:"count"`
		Type  string `json:"type"`
	} `json:"cpus"`
	Memory struct {
		Total string `json:"total"`
	} `json:"memory"`
	Drives []struct {
		Count int `json:"count"`

		// Size is the disk size (example: "480GB")
		Size string `json:"size"`

		// Type is the disk type (example: "SSD")
		Type string `json:"type"`
	} `json:"drives"`

	NICs []struct {
		Count int `json:"count"`

		// Type of NIC (example: "10Gbps")
		Type string `json:"type"`
	} `json:"nics"`

	GPU []struct {
		Count int `json:"count"`

		// Type is the type of GPU (example: "Intel HD Graphics P630")
		Type string `json:"type"`
	} `json:"gpu"`

	Features struct {
		RAID bool `json:"raid"`
		TXT  bool `json:"txt"`
		UEFI bool `json:"uefi"`
	} `json:"features"`
}

type IPNet net.IPNet
type CurrentDevice struct {
	ID             string          `json:"id"`
	Hostname       string          `json:"hostname"`
	IQN            string          `json:"iqn"`
	Plan           string          `json:"plan"`
	Class          string          `json:"class"`
	Facility       string          `json:"facility"`
	PrivateSubnets []IPNet         `json:"private_subnets"`
	Tags           []string        `json:"tags"`
	SSHKeys        []string        `json:"ssh_keys"`
	OS             OperatingSystem `json:"operating_system"`
	Network        NetworkInfo     `json:"network"`
	Storage        Storage         `json:"storage"`
	Volumes        []VolumeInfo    `json:"volumes"`
	Specs
	SwitchShortID string      `json:"switch_short_id"`
	APIURL        string      `json:"api_url"`
	PhoneHomeURL  string      `json:"phone_home_url"`
	UserStateURL  string      `json:"user_state_url"`
	CustomData    interface{} `json:"customdata"`

	// This is available, but is actually inaccurate, currently:
	//   APIBaseURL string          `json:"api_url"`
}

type InterfaceInfo struct {
	Name string `json:"name"`
	MAC  string `json:"mac"`
	Bond string `json:"bond"`
}

func (i *InterfaceInfo) ParseMAC() (net.HardwareAddr, error) {
	return net.ParseMAC(i.MAC)
}

type NetworkInfo struct {
	Interfaces []InterfaceInfo `json:"interfaces"`
	Addresses  []AddressInfo   `json:"addresses"`

	Bonding struct {
		Mode            BondingMode `json:"mode"`
		LinkAggregation string      `json:"link_aggregation"`
		MAC             string      `json:"mac"`
	} `json:"bonding"`
}

func (n *NetworkInfo) BondingMode() BondingMode {
	return n.Bonding.Mode
}

type LicenseActivation struct {
	State string `json:"state"`
}

type OperatingSystem struct {
	Slug              string            `json:"slug"`
	Distro            string            `json:"distro"`
	Version           string            `json:"version"`
	LicenseActivation LicenseActivation `json:"license_activation"`
	ImageTag          string            `json:"image_tag"`
}

type VolumeInfo struct {
	Name string   `json:"name"`
	IQN  string   `json:"iqn"`
	IPs  []net.IP `json:"ips"`

	Capacity struct {
		Size int    `json:"size,string"`
		Unit string `json:"unit"`
	} `json:"capacity"`
}

// Timestamp represents a time that can be unmarshalled from a JSON string
// formatted as either an RFC3339 or Unix timestamp. All
// exported methods of time.Time can be called on Timestamp.
type Timestamp struct {
	time.Time
}

func (t Timestamp) String() string {
	return t.Time.String()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// IPNet is expected in RFC4632 or  RFC4291 format.
func (i *IPNet) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	_, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		return err
	}
	if ipnet == nil {
		return nil
	}
	i.IP = ipnet.IP
	i.Mask = ipnet.Mask
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// Time is expected in RFC3339 or Unix format.
func (t *Timestamp) UnmarshalJSON(data []byte) (err error) {
	str := string(data)
	i, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		t.Time = time.Unix(i, 0)
	} else {
		t.Time, err = time.Parse(`"`+time.RFC3339+`"`, str)
	}
	return
}

// Equal reports whether t and u are equal based on time.Equal
func (t Timestamp) Equal(u Timestamp) bool {
	return t.Time.Equal(u.Time)
}
