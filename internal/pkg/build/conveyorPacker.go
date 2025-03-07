// Copyright (c) 2018-2025, Sylabs Inc. All rights reserved.
// Copyright (c) Contributors to the Apptainer project, established as
//   Apptainer a Series of LF Projects LLC.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package build

import (
	"context"
	"fmt"

	"github.com/sylabs/singularity/v4/internal/pkg/build/sources"
	"github.com/sylabs/singularity/v4/internal/pkg/ociimage"
	"github.com/sylabs/singularity/v4/pkg/build/types"
)

// Conveyor is responsible for downloading from remote sources (library, shub, docker...).
type Conveyor interface {
	Get(context.Context, *types.Bundle) error
}

// Packer is the type which is responsible for installing the chroot directory,
// metadata directory, and potentially other files/directories within the Bundle.
type Packer interface {
	Pack(context.Context) (*types.Bundle, error)
}

// ConveyorPacker describes an interface that a ConveyorPacker type must implement.
type ConveyorPacker interface {
	Conveyor
	Packer
}

// NewConveyorPacker returns a valid ConveyorPacker for the given image definition.
func NewConveyorPacker(def types.Definition) (ConveyorPacker, error) {
	bs, ok := def.Header["bootstrap"]
	if !ok {
		return nil, fmt.Errorf("no bootstrap specification found")
	}

	switch bs {
	case "library":
		return &sources.LibraryConveyorPacker{}, nil
	case "oras":
		return &sources.OrasConveyorPacker{}, nil
	case "shub":
		return &sources.ShubConveyorPacker{}, nil
	case ociimage.SupportedTransport(bs):
		return &sources.OCIConveyorPacker{}, nil
	case "busybox":
		return &sources.BusyBoxConveyorPacker{}, nil
	case "debootstrap":
		return &sources.DebootstrapConveyorPacker{}, nil
	case "arch":
		return &sources.ArchConveyorPacker{}, nil
	case "localimage":
		return &sources.LocalConveyorPacker{}, nil
	case "yum", "dnf":
		return &sources.YumConveyorPacker{}, nil
	case "zypper":
		return &sources.ZypperConveyorPacker{}, nil
	case "scratch":
		return &sources.ScratchConveyorPacker{}, nil
	case "":
		return nil, fmt.Errorf("no bootstrap specification found")
	default:
		return nil, fmt.Errorf("invalid build source %q", def.Header["bootstrap"])
	}
}
