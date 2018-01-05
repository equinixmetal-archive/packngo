package packngo

import (
	"fmt"
	"testing"
)

func TestVirtualNetworkAttachToPort(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	fmt.Println("Not yet implemented")

	// TODO: Figure out how to mock test this.
	// attempt to assign virtual network to a non-bonded port                     (assert failure)
	// assign virtual network to bonded port                                      (assert success)
	// attempt to assign same virtual network to previous port                    (assert failure)
	// unassign virtual network from bonded port                                  (assert success)
	// attempt to unassign the same virtual network from the previous bonded port (assert failure)
}
