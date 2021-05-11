package packngo

import "runtime/debug"

// Version of the packngo package
var Version = "(devel)"

const packagePath = "github.com/packethost/packngo"

// init finds packngo in the dependency so the package Version can be properly
// reflected in API UserAgent headers and client introspection
func init() {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	for _, d := range bi.Deps {
		if d.Path == packagePath {
			Version = d.Version
			if d.Replace != nil {
				Version = d.Replace.Version
			}
			break
		}
	}
}
