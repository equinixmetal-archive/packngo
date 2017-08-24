package packngo

import (
	"fmt"
	"strconv"
	"strings"
)

const ipBasePath = "/ips"

// DeviceIPService handles assignment of addresses from reserved blocks to instances in a project.
type DeviceIPService interface {
	Assign(deviceID string, assignRequest *AddressStruct) (*IPAddressAssignment, *Response, error)
	Unassign(assignmentID string) (*Response, error)
	Get(assignmentID string) (*IPAddressAssignment, *Response, error)
}

// ProjectIPService handles reservation of IP address blocks for a project.
type ProjectIPService interface {
	Get(reservationID string) (*IPAddressReservation, *Response, error)
	GetByCIDR(projectID, cidr string) (*IPAddressReservation, *Response, error)
	List(projectID string) ([]IPAddressReservation, *Response, error)
	Request(projectID string, ipReservationReq *IPReservationRequest) (*AddressStruct, *Response, error)
	Remove(ipReservationID string) (*Response, error)
	AvailableAddresses(ipReservationID string, r *AvailableRequest) ([]string, *Response, error)
}

type ipAddressCommon struct {
	ID            string `json:"id"`
	Address       string `json:"address"`
	Gateway       string `json:"gateway"`
	Network       string `json:"network"`
	AddressFamily int    `json:"address_family"`
	Netmask       string `json:"netmask"`
	Public        bool   `json:"public"`
	CIDR          int    `json:"cidr"`
	Created       string `json:"created_at,omitempty"`
	Updated       string `json:"updated_at,omitempty"`
	Href          string `json:"href"`
}

// IPAddressReservation is created when user sends IP reservation request for a project (considering it's within quota).
type IPAddressReservation struct {
	ipAddressCommon
	Assignments []Href   `json:"assignments"`
	Facility    Facility `json:"facility,omitempty"`
	Available   string   `json:"available"`
	Addon       bool     `json:"addon"`
	Bill        bool     `json:"bill"`
}

// AvailableResponse is a type for listing of available addresses from a reserved block.
type AvailableResponse struct {
	Available []string `json:"available"`
}

// AvailableRequest is a type for listing available addresses from a reserved block.
type AvailableRequest struct {
	CIDR int `json:"cidr"`
}

// IPAddressAssignment is created when an IP address from reservation block is assigned to a device.
type IPAddressAssignment struct {
	ipAddressCommon
	AssignedTo Href `json:"assignments"`
}

// IPReservationRequest represents the body of a reservation request.
type IPReservationRequest struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
	Comments string `json:"comments"`
	Facility string `json:"facility"`
}

// AddressStruct is a helper type for request/response with dict like {"address": ... }
type AddressStruct struct {
	Address string `json:"address"`
}

func deleteFromIP(client *Client, resourceID string) (*Response, error) {
	path := fmt.Sprintf("%s/%s", ipBasePath, resourceID)

	req, err := client.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req, nil)
	return resp, err
}

func (i IPAddressReservation) String() string {
	return Stringify(i)
}

func (i IPAddressAssignment) String() string {
	return Stringify(i)
}

// DeviceIPServiceOp is interface for IP-address assignment methods.
type DeviceIPServiceOp struct {
	client *Client
}

// Unassign unassigns an IP address from the device to which it is currently assigned.
// This will remove the relationship between an IP and the device and will make the IP
// address available to be assigned to another device.
func (i *DeviceIPServiceOp) Unassign(assignmentID string) (*Response, error) {
	return deleteFromIP(i.client, assignmentID)
}

// Assign assigns an IP address to a device.
// The IP address must be in one of the IP ranges assigned to the device’s project.
func (i *DeviceIPServiceOp) Assign(deviceID string, assignRequest *AddressStruct) (*IPAddressAssignment, *Response, error) {
	path := fmt.Sprintf("%s/%s%s", deviceBasePath, deviceID, ipBasePath)

	req, err := i.client.NewRequest("POST", path, assignRequest)

	ipa := new(IPAddressAssignment)
	resp, err := i.client.Do(req, ipa)
	if err != nil {
		return nil, resp, err
	}

	return ipa, resp, err
}

// Get returns assignment by ID.
func (i *DeviceIPServiceOp) Get(assignmentID string) (*IPAddressAssignment, *Response, error) {
	path := fmt.Sprintf("%s/%s", ipBasePath, assignmentID)

	req, err := i.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	ipa := new(IPAddressAssignment)
	resp, err := i.client.Do(req, ipa)
	if err != nil {
		return nil, resp, err
	}

	return ipa, resp, err
}

// ProjectIPServiceOp is interface for IP assignment methods.
type ProjectIPServiceOp struct {
	client *Client
}

// Get returns reservation by ID.
func (i *ProjectIPServiceOp) Get(reservationID string) (*IPAddressReservation, *Response, error) {
	path := fmt.Sprintf("%s/%s", ipBasePath, reservationID)

	req, err := i.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	ipr := new(IPAddressReservation)
	resp, err := i.client.Do(req, ipr)
	if err != nil {
		return nil, resp, err
	}

	return ipr, resp, err
}

// List provides a list of IP resevations for a single project.
func (i *ProjectIPServiceOp) List(projectID string) ([]IPAddressReservation, *Response, error) {
	path := fmt.Sprintf("%s/%s%s", projectBasePath, projectID, ipBasePath)

	req, err := i.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}
	type ipReservationRoot struct {
		Reservations []IPAddressReservation `json:"ip_addresses"`
	}

	reservations := new(ipReservationRoot)
	resp, err := i.client.Do(req, reservations)
	if err != nil {
		return nil, resp, err
	}
	return reservations.Reservations, resp, nil
}

// GetByCIDR returns reservation by CIDR IPv4 net/mask expression, e.g "147.229.20.148/30".
// This is useful upon submitting a reservation request, which returns CIDR of allocated block in exactly this format.
func (i *ProjectIPServiceOp) GetByCIDR(projectID, cidr string) (*IPAddressReservation, *Response, error) {
	cidrSlice := strings.Split(cidr, "/")
	if len(cidrSlice) != 2 {
		return nil, nil, fmt.Errorf("invalid CIDR expression: %s", cidr)
	}
	network := cidrSlice[0]
	subnet, err := strconv.Atoi(cidrSlice[1])
	if err != nil {
		return nil, nil, err
	}
	rs, resp, err := i.List(projectID)
	if err != nil {
		return nil, resp, err
	}
	for i, r := range rs {
		if r.Network == network && r.CIDR == subnet {
			return &rs[i], resp, nil
		}
	}
	return nil, resp, fmt.Errorf("couldn't find reservation for CIDR %s", cidr)

}

// Request requests more IP space for a project in order to have additional IP addresses to assign to devices.
func (i *ProjectIPServiceOp) Request(projectID string, ipReservationReq *IPReservationRequest) (*AddressStruct, *Response, error) {
	path := fmt.Sprintf("%s/%s%s", projectBasePath, projectID, ipBasePath)

	req, err := i.client.NewRequest("POST", path, ipReservationReq)
	if err != nil {
		return nil, nil, err
	}

	ip := new(AddressStruct)
	resp, err := i.client.Do(req, ip)
	if err != nil {
		return nil, resp, err
	}
	return ip, resp, err
}

// Remove removes an IP reservation from the project.
func (i *ProjectIPServiceOp) Remove(ipReservationID string) (*Response, error) {
	return deleteFromIP(i.client, ipReservationID)
}

// AvailableAddresses lists addresses available from a reserved block
func (i *ProjectIPServiceOp) AvailableAddresses(ipReservationID string, r *AvailableRequest) ([]string, *Response, error) {
	path := fmt.Sprintf("%s/%s/available", ipBasePath, ipReservationID)

	req, err := i.client.NewRequest("GET", path, r)
	if err != nil {
		return nil, nil, err
	}

	ar := new(AvailableResponse)
	resp, err := i.client.Do(req, ar)
	if err != nil {
		return nil, resp, err
	}
	return ar.Available, resp, nil

}
