package packngo

const osBasePath = "/operating-systems"

type OS struct {
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	Distro  string `json:"distro"`
	Version int    `json:"version"`
}
func (o OS) String() string {
	return Stringify(o)
}
