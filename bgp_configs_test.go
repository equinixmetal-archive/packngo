package packngo

import (
	"testing"
)

func TestAccBGPConfig(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, projectID, projectDestroy := setupWithProject(t)
	defer projectDestroy()

	configRequest := CreateBGPConfigRequest{
		DeploymentType: "local",
		Asn:            65000,
		Md5:            "c3RhY2twb2ludDIwMTgK",
	}

	_, err := c.BGPConfig.Create(projectID, configRequest)
	if err != nil {
		t.Fatal(err)
	}

	config, _, err := c.BGPConfig.Get(projectID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if config.Md5 != configRequest.Md5 {
		t.Fatal("BGP config is not set up properly")
	}

	if config.Asn != configRequest.Asn {
		t.Fatal("BGP config is not set up properly")
	}

	// _, err = c.BGPConfig.Delete(config.ID)
	// if err != nil {
	// 	t.Fatal(err)
	// }
}
