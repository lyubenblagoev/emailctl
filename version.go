package emailctl

import "fmt"
import "bytes"

// Version holds application version information.
type Version struct {
	// Major is the major version number
	Major int
	// Minor is the minor version number
	Minor int
	// Patch is the patch version number
	Patch int
}

// String returns a version string in format MAJOR.MINOR.PATCH
func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// FullVersion returns the complete version string for emailctl.
func (v Version) FullVersion() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("emailctl version %s\n", v.String()))
	buffer.WriteString("Postfix REST Server API V1 Client")
	return buffer.String()
}
