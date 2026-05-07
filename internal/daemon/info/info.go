package info

import (
	"time"
)

type Info interface {
	// Version returns the version of the daemon.
	Version() string

	// IsVersionChanged returns true if the version of the daemon has changed.
	IsVersionChanged() bool

	// IsSettingsChanged returns true if the settings of the daemon has changed.
	IsSettingsChanged() bool

	// BuiltWith returns the build information of the daemon.
	BuiltWith() string
	// StartTime returns the start time of the daemon.
	StartTime() time.Time
	// Uptime returns the uptime of the daemon.
	Uptime() time.Duration

	// Metadata returns the metadata of the daemon.
	Metadata() MetadataContainer
}

type info struct {
	version          string
	isVersionChanged bool

	isSettingsChanged bool

	builtWith string
	startTime time.Time

	metadata MetadataContainer
}

func New(
	version string,
	isVersionChanged bool,

	builtWith string,
	startTime time.Time,

	isSettingsChanged bool,

	metadata MetadataContainer,
) Info {
	return &info{
		version:          version,
		isVersionChanged: isVersionChanged,

		isSettingsChanged: isSettingsChanged,

		builtWith: builtWith,
		startTime: startTime,
		metadata:  metadata,
	}
}

func (i *info) Version() string {
	return i.version
}

func (i *info) IsVersionChanged() bool {
	return i.isVersionChanged
}

func (i *info) IsSettingsChanged() bool {
	return i.isSettingsChanged
}

func (i *info) BuiltWith() string {
	return i.builtWith
}

func (i *info) StartTime() time.Time {
	return i.startTime
}

func (i *info) Uptime() time.Duration {
	return time.Since(i.StartTime())
}

func (i *info) Metadata() MetadataContainer {
	return i.metadata
}
