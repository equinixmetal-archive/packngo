package packngo

import "fmt"

const (
	volumeBasePath      = "/storage"
	attachmentsBasePath = "/attachments"
)

// VolumeService interface defines available Volume methods
type VolumeService interface {
	Get(string) (*Volume, *Response, error)
	Update(*VolumeUpdateRequest) (*Volume, *Response, error)
	Delete(string) (*Response, error)
	Create(*VolumeCreateRequest) (*Volume, *Response, error)
	Attach(string, string) (*VolumeAttachment, *Response, error)
	Detach(string, string) (*Response, error)
	GetAttachment(string) (*VolumeAttachment, *Response, error)
	DeleteAttachment(string) (*Response, error)
}

// Volume represents a volume
type Volume struct {
	ID               string            `json:"id"`
	Name             string            `json:"name,omitempty"`
	Description      string            `json:"description,omitempty"`
	Size             int               `json:"size,omitempty"`
	State            string            `json:"state,omitempty"`
	Locked           bool              `json:"locked,omitempty"`
	BillingCycle     string            `json:"billing_cycle,omitempty"`
	Created          string            `json:"created_at,omitempty"`
	Updated          string            `json:"updated_at,omitempty"`
	Href             string            `json:"href,omitempty"`
	SnapshotPolicies []*SnapshotPolicy `json:"snapshot_policies,omitempty"`
	Attachments      []*Attachment     `json:"attachments,omitempty"`
	Plan             *Plan             `json:"plan,omitempty"`
	Facility         *Facility         `json:"facility,omitempty"`
	Project          *Project          `json:"project,omitempty"`
}

// SnapshotPolicy used to execute actions on volume
type SnapshotPolicy struct {
	ID                string `json:"id"`
	Href              string `json:"href"`
	SnapshotFrequency string `json:"snapshot_frequency,omitempty"`
	SnapshotCount     int    `json:"snapshot_count,omitempty"`
}

func (v Volume) String() string {
	return Stringify(v)
}

// VolumeCreateRequest type used to create a Packet volume
type VolumeCreateRequest struct {
	Size             int               `json:"size"`
	BillingCycle     string            `json:"billing_cycle"`
	ProjectID        string            `json:"project_id"`
	PlanID           string            `json:"plan_id"`
	FacilityID       string            `json:"facility_id"`
	Description      string            `json:"description,omitempty"`
	SnapshotPolicies []*SnapshotPolicy `json:"snapshot_policies,omitempty"`
}

func (v VolumeCreateRequest) String() string {
	return Stringify(v)
}

// VolumeUpdateRequest type used to update a Packet volume
type VolumeUpdateRequest struct {
	ID          string `json:"id"`
	Description string `json:"description,omitempty"`
	Plan        string `json:"plan,omitempty"`
}

// VolumeAttachRequest type used to attach a Packet volume to a device
type VolumeAttachRequest struct {
	DeviceID string `json:"device_id"`
}

// Link is API resource link
type Link struct {
	Href string `json:"href"`
}

// Attachment is a helper type to parse
type Attachment struct {
	Device         `json:"device"`
	AttachmentID   string `json:"id"`
	AttachmentHref string `json:"href"`
}

// VolumeAttachment is a type from Packet API
type VolumeAttachment struct {
	Href       string `json:"href"`
	ID         string `json:"id"`
	VolumeHref Link   `json:"volume"`
	DeviceHref Link   `json:"device"`
}

func (v VolumeUpdateRequest) String() string {
	return Stringify(v)
}

// VolumeServiceOp implements VolumeService
type VolumeServiceOp struct {
	client *Client
}

// Get returns a volume by id
func (v *VolumeServiceOp) Get(volumeID string) (*Volume, *Response, error) {
	path := fmt.Sprintf("%s/%s?include=facility,snapshot_policies,attachments.device", volumeBasePath, volumeID)
	req, err := v.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	volume := new(Volume)
	resp, err := v.client.Do(req, volume)
	if err != nil {
		return nil, resp, err
	}

	return volume, resp, err
}

// Update updates a volume
func (v *VolumeServiceOp) Update(updateRequest *VolumeUpdateRequest) (*Volume, *Response, error) {
	path := fmt.Sprintf("%s/%s", volumeBasePath, updateRequest.ID)
	req, err := v.client.NewRequest("PATCH", path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	volume := new(Volume)
	resp, err := v.client.Do(req, volume)
	if err != nil {
		return nil, resp, err
	}

	return volume, resp, err
}

// Delete deletes a volume
func (v *VolumeServiceOp) Delete(volumeID string) (*Response, error) {
	path := fmt.Sprintf("%s/%s", volumeBasePath, volumeID)

	req, err := v.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := v.client.Do(req, nil)

	return resp, err
}

// Create creates a new volume for a project
func (v *VolumeServiceOp) Create(createRequest *VolumeCreateRequest) (*Volume, *Response, error) {
	url := fmt.Sprintf("%s/%s%s", projectBasePath, createRequest.ProjectID, volumeBasePath)
	req, err := v.client.NewRequest("POST", url, createRequest)
	if err != nil {
		return nil, nil, err
	}

	volume := new(Volume)
	resp, err := v.client.Do(req, volume)
	if err != nil {
		return nil, resp, err
	}

	return volume, resp, err
}

// Attach volume to a device
func (v *VolumeServiceOp) Attach(volumeID, deviceID string) (*VolumeAttachment, *Response, error) {
	url := fmt.Sprintf("%s/%s%s", volumeBasePath, volumeID, attachmentsBasePath)
	volAttachRequest := VolumeAttachRequest{DeviceID: deviceID}
	req, err := v.client.NewRequest("POST", url, volAttachRequest)
	if err != nil {
		return nil, nil, err
	}
	volumeAttachment := new(VolumeAttachment)
	resp, err := v.client.Do(req, volumeAttachment)
	if err != nil {
		return nil, resp, err
	}
	return volumeAttachment, resp, nil
}

// GetAttachment gets attachment by id
func (v *VolumeServiceOp) GetAttachment(attachmentID string) (*VolumeAttachment, *Response, error) {
	path := fmt.Sprintf("%s%s/%s", volumeBasePath, attachmentsBasePath, attachmentID)
	req, err := v.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	volumeAttachment := new(VolumeAttachment)
	resp, err := v.client.Do(req, volumeAttachment)
	if err != nil {
		return nil, resp, err
	}

	return volumeAttachment, resp, nil
}

// DeleteAttachment deletes attachment by id
func (v *VolumeServiceOp) DeleteAttachment(attachmentID string) (*Response, error) {
	path := fmt.Sprintf("%s%s/%s", volumeBasePath, attachmentsBasePath, attachmentID)
	req, err := v.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := v.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Detach detaches volume from device by Device and Volume IDs
func (v *VolumeServiceOp) Detach(volumeID, deviceID string) (*Response, error) {
	vol, _, err := v.Get(volumeID)
	if err != nil {
		return nil, err
	}
	attachmentToRemoveURL := ""
	for _, a := range vol.Attachments {
		if a.ID == deviceID {
			attachmentToRemoveURL = a.AttachmentHref
			break
		}
	}
	if attachmentToRemoveURL == "" {
		return nil, fmt.Errorf("can not detach, volume %s is not attached to device %s",
			volumeID, deviceID)
	}
	req, err := v.client.NewRequest("DELETE", attachmentToRemoveURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := v.client.Do(req, nil)

	return resp, err
}
