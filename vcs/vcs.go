package vcs

import "runtime/debug"

// Version extract the revision from the git when it is build.
func Version() string {
	var revision string
	var modified bool

	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				revision = s.Value
			case "vcs.modified":
				modified = s.Value == "true"
			}
		}
	}

	if modified {
		return revision + "-dirty"
	}

	return revision
}
