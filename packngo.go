package packngo

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

const (
	packetTokenEnvVar = "PACKET_AUTH_TOKEN"
	libraryVersion    = "0.3.0"
	baseURL           = "https://api.packet.net/"
	userAgent         = "packngo/" + libraryVersion
	mediaType         = "application/json"
	debugEnvVar       = "PACKNGO_DEBUG"

	headerRateLimit     = "X-RateLimit-Limit"
	headerRateRemaining = "X-RateLimit-Remaining"
	headerRateReset     = "X-RateLimit-Reset"
)

var redirectsErrorRe = regexp.MustCompile(`stopped after \d+ redirects\z`)

// GetOptions are options common to Packet API GET requests
type GetOptions struct {
	// Includes are a list of fields to expand in the request results.
	//
	// For resources that contain collections of other resources, the Packet API
	// will only return the `Href` value of these resources by default. In
	// nested API Go types, this will result in objects that have zero values in
	// all fiends except their `Href` field. When an object's associated field
	// name is "included", the returned fields will be Uumarshalled into the
	// nested object. Field specifiers can use a dotted notation up to three
	// references deep. (For example, "memberships.projects" can be used in
	// ListUsers.)
	Includes []string `url:"includes,omitempty"`

	// Excludes reduce the size of the API response by removing nested objects
	// that may be returned.
	//
	// The default behavior of the Packet API is to "exclude" fields, but some
	// API endpoints have an "include" behavior on certain fields. Nested Go
	// types unmarshalled into an "excluded" field will only have a values in
	// their `Href` field.
	Excludes []string `url:"excludes,omitempty"`
}

// GetOptions returns GetOptions from GetOptions (and is nil-receiver safe)
func (g *GetOptions) GetOptions() *GetOptions {
	getOpts := GetOptions{}
	if g != nil {
		getOpts.Includes = g.Includes
		getOpts.Excludes = g.Excludes
	}
	return &getOpts
}

// ListOptions are options common to Packet API paginated GET requests
type ListOptions struct {
	// avoid embedding GetOptions (packngo-breaking-change) for now

	// Includes are a list of fields to expand in the request results.
	Includes []string `url:"includes,omitempty"`

	// Excludes reduce the size of the API response by removing nested objects
	// that may be returned.
	Excludes []string `url:"excludes,omitempty"`

	// Page is the page of results to retrieve for paginated result sets
	Page int `url:"page,omitempty"`

	// PerPage is the number of results to return per page for paginated result
	// sets,
	PerPage int `url:"per_page,omitempty"`

	// The device plan
	Plan string `url:"plan,omitempty"`

	// The device state
	State string `url:"state,omitempty"`

	// The device facility code
	FacilityCode string `url:"facility,omitempty"`

	// Device reservation status
	Reserved *bool `url:"reserved,omitempty"`
}

// GetOptions returns GetOptions from ListOptions (and is nil-receiver safe)
func (l *ListOptions) GetOptions() *GetOptions {
	getOpts := GetOptions{}
	if l != nil {
		getOpts.Includes = l.Includes
		getOpts.Excludes = l.Excludes
	}
	return &getOpts
}

// SearchOptions are options common to API GET requests that include a
// multi-field search filter.
type SearchOptions struct {
	// avoid embedding GetOptions (for similar behavior to ListOptions)

	// Includes are a list of fields to expand in the request results.
	Includes []string `url:"includes,omitempty"`

	// Excludes reduce the size of the API response by removing nested objects
	// that may be returned.
	Excludes []string `url:"excludes,omitempty"`

	// Search is a special API query parameter that, for resources that support
	// it, will filter results to those with any one of various fields matching
	// the supplied keyword.  For example, a resource may have a defined search
	// behavior matches either a name or a fingerprint field, while another
	// resource may match entirely different fields.  Search is currently
	// implemented for SSHKeys and uses an exact match.
	Search string `url:"search,omitempty"`
}

// GetOptions returns GetOptions from ListOptions (and is nil-receiver safe)
func (s *SearchOptions) GetOptions() *GetOptions {
	getOpts := GetOptions{}
	if s != nil {
		getOpts.Includes = s.Includes
		getOpts.Excludes = s.Excludes
	}
	return &getOpts
}

// OptionsGetter provides GetOptions
type OptionsGetter interface {
	GetOptions() *GetOptions
}

func makeSureGetOptionsInclude(g *GetOptions, s string) *GetOptions {
	if g == nil {
		return &GetOptions{Includes: []string{s}}
	}
	if !contains(g.Includes, s) {
		g.Includes = append(g.Includes, s)
	}
	return g
}

func makeSureListOptionsInclude(l *ListOptions, s string) *ListOptions {
	if l == nil {
		return &ListOptions{Includes: []string{s}}
	}
	if !contains(l.Includes, s) {
		l.Includes = append(l.Includes, s)
	}
	return l
}

type paramsReady interface {
	Params() url.Values
}

// compile-time assertions that paramsReady is implemented
var (
	_ paramsReady = (*GetOptions)(nil)
	_ paramsReady = (*ListOptions)(nil)
	_ paramsReady = (*SearchOptions)(nil)
)

// urlQuery generates a URL query string ("?foo=bar") from any object that
// implements the paramsReady interface
func urlQuery(p paramsReady) string {
	return p.Params().Encode()
}

// Params generates URL values from GetOptions fields
func (g *GetOptions) Params() url.Values {
	params := url.Values{}
	if g == nil {
		return params
	}
	if len(g.Includes) != 0 {
		params.Set("include", strings.Join(g.Includes, ","))
	}
	if len(g.Excludes) != 0 {
		params.Set("exclude", strings.Join(g.Excludes, ","))
	}

	return params
}

// Params generates URL values from ListOptions fields
func (l *ListOptions) Params() url.Values {
	if l == nil {
		return url.Values{}
	}
	params := l.GetOptions().Params()

	if l.Page != 0 {
		params.Set("page", fmt.Sprintf("%d", l.Page))
	}
	if l.PerPage != 0 {
		params.Set("per_page", fmt.Sprintf("%d", l.PerPage))
	}

	if l.FacilityCode != "" {
		params.Set("facility", l.FacilityCode)
	}

	if l.Plan != "" {
		params.Set("plan", l.Plan)
	}

	if l.State != "" {
		params.Set("state", l.State)
	}

	if l.Reserved != nil {
		params.Set("reserved", strconv.FormatBool(*l.Reserved))
	}

	return params
}

// Params generates a URL values from SearchOptions fields
func (s *SearchOptions) Params() url.Values {
	if s == nil {
		return url.Values{}
	}

	params := s.GetOptions().Params()
	params.Set("search", s.Search)
	return params
}

// meta contains pagination information
type meta struct {
	Self           *Href `json:"self"`
	First          *Href `json:"first"`
	Last           *Href `json:"last"`
	Previous       *Href `json:"previous,omitempty"`
	Next           *Href `json:"next,omitempty"`
	Total          int   `json:"total"`
	CurrentPageNum int   `json:"current_page"`
	LastPageNum    int   `json:"last_page"`
}

// Response is the http response from api calls
type Response struct {
	*http.Response
	Rate
}

// Href is an API link
type Href struct {
	Href string `json:"href"`
}

func (r *Response) populateRate() {
	// parse the rate limit headers and populate Response.Rate
	if limit := r.Header.Get(headerRateLimit); limit != "" {
		r.Rate.RequestLimit, _ = strconv.Atoi(limit)
	}
	if remaining := r.Header.Get(headerRateRemaining); remaining != "" {
		r.Rate.RequestsRemaining, _ = strconv.Atoi(remaining)
	}
	if reset := r.Header.Get(headerRateReset); reset != "" {
		if v, _ := strconv.ParseInt(reset, 10, 64); v != 0 {
			r.Rate.Reset = Timestamp{time.Unix(v, 0)}
		}
	}
}

// ErrorResponse is the http response used on errors
type ErrorResponse struct {
	Response    *http.Response
	Errors      []string `json:"errors"`
	SingleError string   `json:"error"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, strings.Join(r.Errors, ", "), r.SingleError)
}

// Client is the base API Client
type Client struct {
	client *retryablehttp.Client
	debug  bool

	BaseURL *url.URL

	UserAgent     string
	ConsumerToken string
	APIKey        string

	RateLimit Rate

	// Packet Api Objects
	APIKeys                APIKeyService
	BGPConfig              BGPConfigService
	BGPSessions            BGPSessionService
	Batches                BatchService
	CapacityService        CapacityService
	DeviceIPs              DeviceIPService
	DevicePorts            DevicePortService
	Devices                DeviceService
	Emails                 EmailService
	Events                 EventService
	Facilities             FacilityService
	Hardware               HardwareService
	HardwareReservations   HardwareReservationService
	Notifications          NotificationService
	OperatingSystems       OSService
	Organizations          OrganizationService
	Plans                  PlanService
	ProjectIPs             ProjectIPService
	ProjectVirtualNetworks ProjectVirtualNetworkService
	Projects               ProjectService
	SSHKeys                SSHKeyService
	SpotMarket             SpotMarketService
	SpotMarketRequests     SpotMarketRequestService
	TwoFactorAuth          TwoFactorAuthService
	Users                  UserService
	VPN                    VPNService
	VolumeAttachments      VolumeAttachmentService
	Volumes                VolumeService
}

// requestDoer provides methods for making HTTP requests and receiving the
// response, errors, and a structured result
//
// This interface is used in *ServiceOp as a mockable alternative to a full
// Client object.
type requestDoer interface {
	NewRequest(method, path string, body interface{}) (*retryablehttp.Request, error)
	Do(req *retryablehttp.Request, v interface{}) (*Response, error)
	DoRequest(method, path string, body, v interface{}) (*Response, error)
	DoRequestWithHeader(method string, headers map[string]string, path string, body, v interface{}) (*Response, error)
}

// NewRequest inits a new http request with the proper headers
func (c *Client) NewRequest(method, path string, body interface{}) (*retryablehttp.Request, error) {
	// relative path to append to the endpoint url, no leading slash please
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	// json encode the request body, if any
	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := retryablehttp.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Close = true

	req.Header.Add("X-Auth-Token", c.APIKey)
	req.Header.Add("X-Consumer-Token", c.ConsumerToken)

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

// Do executes the http request
func (c *Client) Do(req *retryablehttp.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	response := Response{Response: resp}
	response.populateRate()
	if c.debug {
		dumpResponse(response.Response)
	}
	c.RateLimit = response.Rate

	err = checkResponse(resp)
	// if the response is an error, return the ErrorResponse
	if err != nil {
		return &response, err
	}

	if v != nil {
		// if v implements the io.Writer interface, return the raw response
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return &response, err
			}
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return &response, err
			}
		}
	}

	return &response, err
}

func dumpResponse(resp *http.Response) {
	o, _ := httputil.DumpResponse(resp, true)
	strResp := string(o)
	reg, _ := regexp.Compile(`"token":(.+?),`)
	reMatches := reg.FindStringSubmatch(strResp)
	if len(reMatches) == 2 {
		strResp = strings.Replace(strResp, reMatches[1], strings.Repeat("-", len(reMatches[1])), 1)
	}
	log.Printf("\n=======[RESPONSE]============\n%s\n\n", strResp)
}

func dumpRequest(req *retryablehttp.Request) {
	o, _ := httputil.DumpRequestOut(req.Request, false)
	strReq := string(o)
	reg, _ := regexp.Compile(`X-Auth-Token: (\w*)`)
	reMatches := reg.FindStringSubmatch(strReq)
	if len(reMatches) == 2 {
		strReq = strings.Replace(strReq, reMatches[1], strings.Repeat("-", len(reMatches[1])), 1)
	}
	bbs, _ := req.BodyBytes()
	log.Printf("\n=======[REQUEST]=============\n%s%s\n", strReq, string(bbs))
}

// DoRequest is a convenience method, it calls NewRequest followed by Do
// v is the interface to unmarshal the response JSON into
func (c *Client) DoRequest(method, path string, body, v interface{}) (*Response, error) {
	req, err := c.NewRequest(method, path, body)
	if c.debug {
		dumpRequest(req)
	}
	if err != nil {
		return nil, err
	}
	return c.Do(req, v)
}

// DoRequestWithHeader same as DoRequest
func (c *Client) DoRequestWithHeader(method string, headers map[string]string, path string, body, v interface{}) (*Response, error) {
	req, err := c.NewRequest(method, path, body)
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	if c.debug {
		dumpRequest(req)
	}
	if err != nil {
		return nil, err
	}
	return c.Do(req, v)
}

// NewClient initializes and returns a Client
func NewClient() (*Client, error) {
	apiToken := os.Getenv(packetTokenEnvVar)
	if apiToken == "" {
		return nil, fmt.Errorf("you must export %s", packetTokenEnvVar)
	}
	c := NewClientWithAuth("packngo lib", apiToken, nil)
	return c, nil

}

// NewClientWithAuth initializes and returns a Client, use this to get an API Client to operate on
// N.B.: Packet's API certificate requires Go 1.5+ to successfully parse. If you are using
// an older version of Go, pass in a custom http.Client with a custom TLS configuration
// that sets "InsecureSkipVerify" to "true"
func NewClientWithAuth(consumerToken string, apiKey string, httpClient *retryablehttp.Client) *Client {
	client, _ := NewClientWithBaseURL(consumerToken, apiKey, httpClient, baseURL)
	return client
}

func PacketRetryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	// do not retry on context.Canceled or context.DeadlineExceeded
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	if err != nil {
		if v, ok := err.(*url.Error); ok {
			// Don't retry if the error was due to too many redirects.
			if redirectsErrorRe.MatchString(v.Error()) {
				return false, nil
			}

			// Don't retry if the error was due to TLS cert verification failure.
			if _, ok := v.Err.(x509.UnknownAuthorityError); ok {
				return false, nil
			}
		}

		// The error is likely recoverable so retry.
		return true, nil
	}

	// Check the response code. We retry on 500-range responses to allow
	// the server time to recover, as 500's are typically not permanent
	// errors and may relate to outages on the server side. This will catch
	// invalid response codes as well, like 0 and 999.
	//if resp.StatusCode == 0 || (resp.StatusCode >= 500 && resp.StatusCode != 501) {
	//	return true, nil
	//}

	return false, nil
}

// NewClientWithBaseURL returns a Client pointing to nonstandard API URL, e.g.
// for mocking the remote API
func NewClientWithBaseURL(consumerToken string, apiKey string, httpClient *retryablehttp.Client, apiBaseURL string) (*Client, error) {
	if httpClient == nil {
		// Don't fall back on http.DefaultClient as it's not nice to adjust state
		// implicitly. If the client wants to use http.DefaultClient, they can
		// pass it in explicitly.
		httpClient = retryablehttp.NewClient()
		httpClient.RetryWaitMin = time.Second
		httpClient.RetryWaitMax = 30 * time.Second
		httpClient.RetryMax = 10
		httpClient.CheckRetry = PacketRetryPolicy
	}

	u, err := url.Parse(apiBaseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{client: httpClient, BaseURL: u, UserAgent: userAgent, ConsumerToken: consumerToken, APIKey: apiKey}
	c.APIKeys = &APIKeyServiceOp{client: c}
	c.BGPConfig = &BGPConfigServiceOp{client: c}
	c.BGPSessions = &BGPSessionServiceOp{client: c}
	c.Batches = &BatchServiceOp{client: c}
	c.CapacityService = &CapacityServiceOp{client: c}
	c.DeviceIPs = &DeviceIPServiceOp{client: c}
	c.DevicePorts = &DevicePortServiceOp{client: c}
	c.Devices = &DeviceServiceOp{client: c}
	c.Emails = &EmailServiceOp{client: c}
	c.Events = &EventServiceOp{client: c}
	c.Facilities = &FacilityServiceOp{client: c}
	c.Hardware = &HardwareServiceOp{client: c}
	c.HardwareReservations = &HardwareReservationServiceOp{client: c}
	c.Notifications = &NotificationServiceOp{client: c}
	c.OperatingSystems = &OSServiceOp{client: c}
	c.Organizations = &OrganizationServiceOp{client: c}
	c.Plans = &PlanServiceOp{client: c}
	c.ProjectIPs = &ProjectIPServiceOp{client: c}
	c.ProjectVirtualNetworks = &ProjectVirtualNetworkServiceOp{client: c}
	c.Projects = &ProjectServiceOp{client: c}
	c.SSHKeys = &SSHKeyServiceOp{client: c}
	c.SpotMarket = &SpotMarketServiceOp{client: c}
	c.SpotMarketRequests = &SpotMarketRequestServiceOp{client: c}
	c.TwoFactorAuth = &TwoFactorAuthServiceOp{client: c}
	c.Users = &UserServiceOp{client: c}
	c.VPN = &VPNServiceOp{client: c}
	c.VolumeAttachments = &VolumeAttachmentServiceOp{client: c}
	c.Volumes = &VolumeServiceOp{client: c}
	c.debug = os.Getenv(debugEnvVar) != ""

	return c, nil
}

func checkResponse(r *http.Response) error {
	// return if http status code is within 200 range
	if c := r.StatusCode; c >= 200 && c <= 299 {
		// response is good, return
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	// if the response has a body, populate the message in errorResponse
	if err != nil {
		return err
	}

	if len(data) > 0 {
		err = json.Unmarshal(data, errorResponse)
		if err != nil {
			return err
		}
	}

	return errorResponse
}
