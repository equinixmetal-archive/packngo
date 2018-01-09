package packngo

import (
	"fmt"
	"testing"
)

func TestVirtualNetworks(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	fmt.Println("Not yet implemented")
	// TODO: Test several bad inputs to ensure rejection without adverse affects
	// List virtual networks using bad project ID
	// Ensure there are zero for the fake project
	// Create virtual network with bad POST body parameters
	// Ensure create failed
	// Ensure zero virtual networks still
	// Create virtual network
	// List virtual networks and ensure length is one
	// Delete virtual network
	// Ensure zero virtual networks attached to project
}
