package packngo

import (
	"fmt"
	"strconv"
	"strings"
)

const ipBasePath = "/ips"

// IPService interface defines available IP methods
type IPService interface {
	Assign(deviceID string, assignRequest *AddressField) (*IPAddressAssignment, *Response, error)
	Unassign(assignmentID string) (*Response, error)
	GetReservation(reservationID string) (*IPAddressReservation, *Response, error)
	GetReservationByCIDR(projectID, cidrString string) (*IPAddressReservation, *Response, error)
	GetAssignment(assignmentID string) (*IPAddressAssignment, *Response, error)
	ListReservations(projectID string) ([]IPAddressReservation, *Response, error)
	RequestReservation(projectID string, ipReservationReq *IPReservationRequest) (*AddressField, *Response, error)
	RemoveReservation(ipReservationID string) (*Response, error)
	GetAvailableAddresses(ipReservationID string, r *AvailableRequest) ([]string, *Response, error)
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

// IPAddressReservation is created when user sends IP reservation request for his/her project (considering it's within quota).
type IPAddressReservation struct {
	ipAddressCommon
	Assignments []Href   `json:"assignments"`
	Facility    Facility `json:"facility,omitempty"`
	Available   string   `json:"available"`
	Addon       bool     `json:"addon"`
	Bill        bool     `json:"bill"`
}

// AvailableResponse is a type for listing of avaialable addresses from a reserved block.
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

// AddressField is a type for request/response with dict of type {"address": ... }
type AddressField struct {
	Address string `json:"address"`
}

func (i IPAddressReservation) String() string {
	return Stringify(i)
}

func (i IPAddressAssignment) String() string {
	return Stringify(i)
}

// IPServiceOp implements IPService
type IPServiceOp struct {
	client *Client
}

// GetReservation returns reservation by ID
func (i *IPServiceOp) GetReservation(reservationID string) (*IPAddressReservation, *Response, error) {
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

// GetAssignment returns assignment by ID
func (i *IPServiceOp) GetAssignment(assignmentID string) (*IPAddressAssignment, *Response, error) {
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

type ipReservationRoot struct {
	Reservations []IPAddressReservation `json:"ip_addresses"`
}

// ListReservations provides a list of IP resevations for a single project.
func (i *IPServiceOp) ListReservations(projectID string) ([]IPAddressReservation, *Response, error) {
	path := fmt.Sprintf("%s/%s%s", projectBasePath, projectID, ipBasePath)

	req, err := i.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	reservations := new(ipReservationRoot)
	resp, err := i.client.Do(req, reservations)
	if err != nil {
		return nil, resp, err
	}
	return reservations.Reservations, resp, nil
}

// GetReservationByCIDR returns reservation by CIDR IPv4 net/mask expression, e.g "147.229.20.148/30".
// This is useful upon submitting a reservation request, which returns CIDR of allocated block in exactly this format.
func (i *IPServiceOp) GetReservationByCIDR(projectID, cidrString string) (*IPAddressReservation, *Response, error) {
	cidrSlice := strings.Split(cidrString, "/")
	if len(cidrSlice) != 2 {
		return nil, nil, fmt.Errorf("Invalid CIDR expression: %s", cidrString)
	}
	network := cidrSlice[0]
	cidr, err := strconv.Atoi(cidrSlice[1])
	if err != nil {
		return nil, nil, err
	}
	rs, resp, err := i.ListReservations(projectID)
	if err != nil {
		return nil, resp, err
	}
	for _, r := range rs {
		if r.Network == network && r.CIDR == cidr {
			return &r, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("Couldn't find reservation for CIDR %s", cidrString)

}

func (i *IPServiceOp) deleteFromIP(resourceID string) (*Response, error) {
	path := fmt.Sprintf("%s/%s", ipBasePath, resourceID)

	req, err := i.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := i.client.Do(req, nil)
	return resp, err
}

// Unassign unassigns an IP address from the device to which it is currently assignmed.
// This will remove the relationship between an IP and the device and will make the IP
// address available to be assigned to another device.
func (i *IPServiceOp) Unassign(assignmentID string) (*Response, error) {
	return i.deleteFromIP(assignmentID)
}

// Assign assigns an IP address to a device. The IP address must be in one of the IP ranges assigned to the device’s project.
func (i *IPServiceOp) Assign(deviceID string, assignRequest *AddressField) (*IPAddressAssignment, *Response, error) {
	path := fmt.Sprintf("%s/%s%s", deviceBasePath, deviceID, ipBasePath)

	req, err := i.client.NewRequest("POST", path, assignRequest)

	ipa := new(IPAddressAssignment)
	resp, err := i.client.Do(req, ipa)
	if err != nil {
		return nil, resp, err
	}

	return ipa, resp, err
}

// IPReservationRequest represents the body of a reservation request
type IPReservationRequest struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
	Comments string `json:"comments"`
	Facility string `json:"facility"`
}

// RequestReservation requests more IP space for a project in order to have additional IP addresses to assign to devices
func (i *IPServiceOp) RequestReservation(projectID string, ipReservationReq *IPReservationRequest) (*AddressField, *Response, error) {
	path := fmt.Sprintf("%s/%s%s", projectBasePath, projectID, ipBasePath)

	req, err := i.client.NewRequest("POST", path, &ipReservationReq)
	if err != nil {
		return nil, nil, err
	}

	ip := new(AddressField)
	resp, err := i.client.Do(req, ip)
	if err != nil {
		return nil, resp, err
	}
	return ip, resp, err
}

// RemoveReservation removes an IP reservation from the project.
func (i *IPServiceOp) RemoveReservation(ipReservationID string) (*Response, error) {
	return i.deleteFromIP(ipReservationID)
}

// GetAvailableAddresses lists addresses available from a reserved block
func (i *IPServiceOp) GetAvailableAddresses(ipReservationID string, r *AvailableRequest) ([]string, *Response, error) {
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
