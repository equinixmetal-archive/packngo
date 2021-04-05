package packngo

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
)

const (
	apiURLEnvVar      = "PACKET_API_URL"
	packngoAccTestVar = "PACKNGO_TEST_ACTUAL_API"
	testProjectPrefix = "PACKNGO_TEST_DELME_2d768716_"
	testFacilityVar   = "PACKNGO_TEST_FACILITY"
	testMetroVar      = "PACKNGO_TEST_METRO"
	testPlanVar       = "PACKNGO_TEST_PLAN"
	testRecorderEnv   = "PACKNGO_TEST_RECORDER"

	testRecorderRecord   = "record"
	testRecorderPlay     = "play"
	testRecorderDisabled = "disabled"

	recorderDefaultMode = recorder.ModeDisabled

	// defaults should be available to most users
	testFacilityDefault   = "ny5"
	testFacilityAlternate = "dc13"
	testMetroDefault      = "sv"
	testPlanDefault       = "c3.small.x86"
	testOS                = "ubuntu_18_04"
)

func testPlan() string {
	envPlan := os.Getenv(testPlanVar)
	if envPlan != "" {
		return envPlan
	}
	return testPlanDefault
}

func testMetro() string {
	envMet := os.Getenv(testMetroVar)
	if envMet != "" {
		return envMet
	}
	return testMetroDefault
}

func testFacility() string {
	envFac := os.Getenv(testFacilityVar)
	if envFac != "" {
		return envFac
	}
	return testFacilityDefault
}

func randString8() string {
	// test recorder needs replayable names, not randoms
	mode, _ := testRecordMode()
	if mode != recorder.ModeDisabled {
		return "testrand"
	}

	n := 8
	rand.Seed(time.Now().UnixNano())
	letterRunes := []rune("acdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// MockClient makes it simpler to test the Client
type MockClient struct {
	fnNewRequest          func(method, path string, body interface{}) (*http.Request, error)
	fnDo                  func(req *http.Request, v interface{}) (*Response, error)
	fnDoRequest           func(method, path string, body, v interface{}) (*Response, error)
	fnDoRequestWithHeader func(method string, headers map[string]string, path string, body, v interface{}) (*Response, error)
}

var _ requestDoer = &MockClient{}

// NewRequest uses the mock NewRequest function
func (mc *MockClient) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	return mc.fnNewRequest(method, path, body)
}

// Do uses the mock Do function
func (mc *MockClient) Do(req *http.Request, v interface{}) (*Response, error) {
	return mc.fnDo(req, v)
}

// DoRequest uses the mock DoRequest function
func (mc *MockClient) DoRequest(method, path string, body, v interface{}) (*Response, error) {
	return mc.fnDoRequest(method, path, body, v)
}

// DoRequestWithHeader uses the mock DoRequestWithHeader function
func (mc *MockClient) DoRequestWithHeader(method string, headers map[string]string, path string, body, v interface{}) (*Response, error) {
	return mc.fnDoRequestWithHeader(method, headers, path, body, v)
}

// setupWithProject returns a client, project id, and teardown function
// configured for a new project with a test recorder for the named test
func setupWithProject(t *testing.T) (*Client, string, func()) {
	c, stopRecord := setup(t)
	p, _, err := c.Projects.Create(&ProjectCreateRequest{Name: testProjectPrefix + randString8()})
	if err != nil {
		t.Fatal(err)
	}

	return c, p.ID, func() {
		_, err := c.Projects.Delete(p.ID)
		if err != nil {
			panic(fmt.Errorf("while deleting %s: %s", p, err))
		}
		stopRecord()
	}

}

func skipUnlessAcceptanceTestsAllowed(t *testing.T) {
	if os.Getenv(packngoAccTestVar) == "" {
		t.Skipf("%s is not set", packngoAccTestVar)
	}
}

// testRecorder creates the named recorder
func testRecorder(t *testing.T, name string, mode recorder.Mode) (*recorder.Recorder, func()) {
	r, err := recorder.NewAsMode(path.Join("fixtures", name), mode, nil)
	if err != nil {
		t.Fatal(err)
	}

	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "X-Auth-Token")
		return nil
	})

	return r, func() {
		if err := r.Stop(); err != nil {
			t.Fatal(err)
		}
	}
}

func testRecordMode() (recorder.Mode, error) {
	modeRaw := os.Getenv(testRecorderEnv)
	mode := recorderDefaultMode

	switch strings.ToLower(modeRaw) {
	case testRecorderRecord:
		mode = recorder.ModeRecording
	case testRecorderPlay:
		mode = recorder.ModeReplaying
	case "":
		// no-op
	case testRecorderDisabled:
		// no-op
	default:
		return mode, fmt.Errorf("invalid %s mode: %s", testRecorderEnv, modeRaw)
	}
	return mode, nil
}

func setup(t *testing.T) (*Client, func()) {
	name := t.Name()
	apiToken := os.Getenv(authTokenEnvVar)
	if apiToken == "" {
		t.Fatalf("If you want to run packngo test, you must export %s.", authTokenEnvVar)
	}

	mode, err := testRecordMode()
	if err != nil {
		t.Fatal(err)
	}
	apiURL := os.Getenv(apiURLEnvVar)
	if apiURL == "" {
		apiURL = baseURL
	}
	r, stopRecord := testRecorder(t, name, mode)
	httpClient := http.DefaultClient
	httpClient.Transport = r
	c, err := NewClientWithBaseURL("packngo test", apiToken, httpClient, apiURL)
	if err != nil {
		t.Fatal(err)
	}

	return c, stopRecord
}

func projectTeardown(c *Client) {
	ps, _, err := c.Projects.List(nil)
	if err != nil {
		panic(fmt.Errorf("while teardown: %s", err))
	}
	for _, p := range ps {
		if strings.HasPrefix(p.Name, testProjectPrefix) {
			_, err := c.Projects.Delete(p.ID)
			if err != nil {
				panic(fmt.Errorf("while deleting %s: %s", p, err))
			}
		}
	}
}

func organizationTeardown(c *Client) {
	ps, _, err := c.Organizations.List(nil)
	if err != nil {
		panic(fmt.Errorf("while teardown: %s", err))
	}
	for _, p := range ps {
		if strings.HasPrefix(p.Name, testProjectPrefix) {
			_, err := c.Organizations.Delete(p.ID)
			if err != nil {
				panic(fmt.Errorf("while deleting %s: %s", p, err))
			}
		}
	}
}

func TestAccInvalidCredentials(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c := NewClientWithAuth("packngo test", "wrongApiToken", nil)
	_, r, expectedErr := c.Projects.List(nil)
	matched, err := regexp.MatchString(".*Invalid.*", expectedErr.Error())
	if err != nil {
		t.Fatalf("Err while matching err string from response err %s: %s", expectedErr, err)
	}
	if r.StatusCode != 401 {
		t.Fatalf("Expected 401 as response code, got: %d", r.StatusCode)
	}

	if !matched {
		t.Fatalf("Unexpected error string: %s", expectedErr)
	}

}

func Test_dumpDeprecation(t *testing.T) {
	type args struct {
		resp *http.Response
	}
	tests := []struct {
		name   string
		args   args
		logged string
	}{
		{
			name: "Deprecation",
			args: args{
				resp: &http.Response{
					Header: http.Header{
						"Deprecation": {
							"Sat, 1 Aug 2020 23:59:59 GMT",
						},
						"Link": {
							"<https://api.example.com/deprecation>; rel=\"deprecation\"; type=\"text/html\"",
						},
					},
					Request: &http.Request{
						Method: "POST",
						URL:    &url.URL{Path: "/deprecated"},
					},
				},
			},
			logged: "WARNING: \"POST /deprecated\" reported deprecation on Sat, 1 Aug 2020 23:59:59 GMT\nWARNING: See <https://api.example.com/deprecation> for deprecation details",
		},
		{
			name: "Sunset",
			args: args{
				resp: &http.Response{
					Header: http.Header{
						"Sunset": {
							"Sat, 1 Aug 2020 23:59:59 GMT",
						},
						"Link": {
							"<https://api.example.com/sunset>; rel=\"sunset\"; type=\"text/html\"",
						},
					},
					Request: &http.Request{
						Method: "GET",
						URL:    &url.URL{Path: "/sunset"},
					},
				},
			},
			logged: "WARNING: \"GET /sunset\" reported sunsetting on Sat, 1 Aug 2020 23:59:59 GMT\nWARNING: See <https://api.example.com/sunset> for sunset details",
		},
		{
			name: "DeprecateAndSunset",
			args: args{
				resp: &http.Response{
					Header: http.Header{
						"Sunset": {
							"Sat, 1 Aug 2020 23:59:59 GMT",
						},
						"Deprecation": {
							"true",
						},
						// comma separated header value and repeated header
						"Link": {
							"<https://api.example.com/deprecation/field-a>; rel=\"deprecation\"; type=\"text/html\"",
							"<https://api.example.com/sunset/value-a>; rel=\"sunset\"; type=\"text/html\"",
							"<https://api.example.com/sunset>; rel=\"sunset\"; type=\"text/html\",<https://api.example.com/deprecation>; rel=\"deprecation\"; type=\"text/html\"",
						},
					},
					Request: func() *http.Request {
						body := bytes.NewReader([]byte("{\"sunset\":true,\"deprecated\":true}"))
						r, _ := http.NewRequest(
							http.MethodPost,
							"/deprecate-and-sunset", body)
						return r
					}(),
				},
			},
			// only the comma separate header is returned by Header.Get()
			logged: "WARNING: \"POST /deprecate-and-sunset\" reported deprecation\nWARNING: \"POST /deprecate-and-sunset\" reported sunsetting on Sat, 1 Aug 2020 23:59:59 GMT\nWARNING: See <https://api.example.com/deprecation/field-a> for deprecation details\nWARNING: See <https://api.example.com/sunset/value-a> for sunset details\nWARNING: See <https://api.example.com/sunset> for sunset details\nWARNING: See <https://api.example.com/deprecation> for deprecation details",
		},
		{
			name: "None",
			args: args{
				resp: &http.Response{
					Header: http.Header{},
				},
			},
			logged: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logged := &bytes.Buffer{}
			log.SetOutput(logged)
			f := log.Flags()
			log.SetFlags(0)
			defer func() {
				log.SetOutput(os.Stderr)
				log.SetFlags(f)
			}()
			dumpDeprecation(tt.args.resp)
			got := strings.TrimSpace(logged.String())
			if got != tt.logged {
				t.Logf("%s failed; got %q, want %q", t.Name(), got, tt.logged)
				t.Fail()
			}
		})
	}
}
