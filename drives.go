// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.
package firecracker

import (
	"strconv"

	models "github.com/firecracker-microvm/firecracker-go-sdk/client/models"
)

const rootDriveName = "root_drive"

// DrivesBuilder is a builder that will build an array of drives used to set up
// the firecracker microVM. The DriveID will be an incrementing number starting
// at one
type DrivesBuilder struct {
	rootDrive models.Drive
	drives    []models.Drive
}

// NewDrivesBuilder will return a new DrivesBuilder with a given rootfs.
func NewDrivesBuilder(rootDrivePath string) DrivesBuilder {
	return DrivesBuilder{}.WithRootDrive(rootDrivePath)
}

// DriveOpt represents an optional function used to allow for specific
// customization of the models.Drive structure.
type DriveOpt func(*models.Drive)

// WithRootDrive will set the given builder with the a new root path. The root
// drive will be set to read and write by default.
func (b DrivesBuilder) WithRootDrive(rootDrivePath string, opts ...DriveOpt) DrivesBuilder {
	b.rootDrive = models.Drive{
		DriveID:      String(rootDriveName),
		PathOnHost:   &rootDrivePath,
		IsRootDevice: Bool(true),
		IsReadOnly:   Bool(false),
	}

	for _, opt := range opts {
		opt(&b.rootDrive)
	}

	return b
}

// AddDrive will add a new drive to the given builder.
func (b DrivesBuilder) AddDrive(path string, readOnly bool, opts ...DriveOpt) DrivesBuilder {
	drive := models.Drive{
		DriveID:      String(strconv.Itoa(len(b.drives))),
		PathOnHost:   &path,
		IsRootDevice: Bool(false),
		IsReadOnly:   &readOnly,
	}

	for _, opt := range opts {
		opt(&drive)
	}

	b.drives = append(b.drives, drive)
	return b
}

// VhostUserDriveOpt represents an optional function used to customize a
// vhost-user-block drive. It is a distinct type from DriveOpt so that options
// only valid for path-backed drives (e.g. WithCacheType, WithIoEngine) cannot
// be applied to a vhost-user drive at compile time. Socket is mutually
// exclusive with PathOnHost/CacheType/IoEngine in the Firecracker API.
type VhostUserDriveOpt func(*models.Drive)

// AddVhostUserDrive will add a new vhost-user-block drive to the given builder.
// A vhost-user drive is backed by a vhost-user-block backend listening on the
// given unix socket instead of a host file path. Only the WithVhost* options
// are accepted, so PathOnHost/CacheType/IoEngine can never be set on such a
// drive.
func (b DrivesBuilder) AddVhostUserDrive(socket string, readOnly bool, opts ...VhostUserDriveOpt) DrivesBuilder {
	drive := models.Drive{
		DriveID:      String(strconv.Itoa(len(b.drives))),
		Socket:       &socket,
		IsRootDevice: Bool(false),
		IsReadOnly:   &readOnly,
	}

	for _, opt := range opts {
		opt(&drive)
	}

	b.drives = append(b.drives, drive)
	return b
}

// Build will construct an array of drives with the root drive at the very end.
func (b DrivesBuilder) Build() []models.Drive {
	return append(b.drives, b.rootDrive)
}

// WithDriveID sets the ID of the drive
func WithDriveID(id string) DriveOpt {
	return func(d *models.Drive) {
		d.DriveID = String(id)
	}
}

// WithReadOnly sets the drive read-only
func WithReadOnly(flag bool) DriveOpt {
	return func(d *models.Drive) {
		d.IsReadOnly = Bool(flag)
	}
}

// WithPartuuid sets the unique ID of the boot partition
func WithPartuuid(uuid string) DriveOpt {
	return func(d *models.Drive) {
		d.Partuuid = uuid
	}
}

// WithRateLimiter sets the rate limitter of the drive
func WithRateLimiter(limiter models.RateLimiter) DriveOpt {
	return func(d *models.Drive) {
		d.RateLimiter = &limiter
	}
}

// WithCacheType sets the cache strategy for the block device
func WithCacheType(cacheType string) DriveOpt {
	return func(d *models.Drive) {
		d.CacheType = String(cacheType)
	}
}

// WithIoEngine sets the io engine of the drive
// Defaults to Sync, Async is in developer preview at the moment
// https://github.com/firecracker-microvm/firecracker/blob/v1.1.0/docs/api_requests/block-io-engine.md
func WithIoEngine(ioEngine string) DriveOpt {
	return func(d *models.Drive) {
		d.IoEngine = String(ioEngine)
	}
}

// WithVhostDriveID sets the ID of a vhost-user-block drive.
func WithVhostDriveID(id string) VhostUserDriveOpt {
	return func(d *models.Drive) {
		d.DriveID = String(id)
	}
}

// WithVhostReadOnly sets a vhost-user-block drive read-only.
func WithVhostReadOnly(flag bool) VhostUserDriveOpt {
	return func(d *models.Drive) {
		d.IsReadOnly = Bool(flag)
	}
}

// WithVhostPartuuid sets the unique ID of the boot partition on a
// vhost-user-block drive.
func WithVhostPartuuid(uuid string) VhostUserDriveOpt {
	return func(d *models.Drive) {
		d.Partuuid = uuid
	}
}

// WithVhostRateLimiter sets the rate limiter of a vhost-user-block drive.
func WithVhostRateLimiter(limiter models.RateLimiter) VhostUserDriveOpt {
	return func(d *models.Drive) {
		d.RateLimiter = &limiter
	}
}
