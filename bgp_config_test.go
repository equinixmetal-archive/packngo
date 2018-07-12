package packngo

import (
	"testing"
)

var configID string
var configRequest CreateBGPConfigRequest

func TestAccCreateBGPConfig(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, projectID, _ := setupWithProject(t)

	configRequest = CreateBGPConfigRequest{
		DeploymentType: "local",
		Asn:            65000,
		Md5:            "c3RhY2twb2ludDIwMTgK",
	}

	_, err := c.BGPConfig.Create(projectID, configRequest)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccGetBgpConfig(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c := setup(t)

	config, _, err := c.BGPConfig.Get(projectID)
	if err != nil {
		t.Fatal(err)
	}

	if config == nil {
		t.Fatal("BGP config not retrieved")
	}

	if config.Md5 != configRequest.Md5 {
		t.Fatal("BGP config is not set up properly")
	}

	if config.Asn != configRequest.Asn {
		t.Fatal("BGP config is not set up properly")
	}

	_, err = c.BGPConfig.Delete(config.ID)
	if err != nil {
		t.Fatal(err)
	}
	projectTeardown(c)
}
