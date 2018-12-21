package driver

import (
	"context"
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"os/exec"
)

const (
	// default file system type to be used when it is not provided
	defaultFsType = "ext4"
)

func (d *Driver) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	volumeID := req.GetVolumeId()
	if len(volumeID) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID not provided")
	}
	// Check that target_path is created by CO and is a directory
	target := req.GetStagingTargetPath()
	if len(target) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Staging target not provided")
	}

	source, ok := req.PublishContext["devicePath"]
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "Device path not provided")
	}

	volCap := req.GetVolumeCapability()
	if volCap == nil {
		return nil, status.Error(codes.InvalidArgument, "Volume capability not provided")
	}



	// TODO: consider replacing IsLikelyNotMountPoint by IsNotMountPoint
	notMnt, err := d.mounter.Interface.IsLikelyNotMountPoint(target)
	if err != nil {
		if os.IsNotExist(err) {
			if errMkDir := d.mounter.Interface.MakeDir(target); errMkDir != nil {
				msg := fmt.Sprintf("could not create target dir %q: %v", target, errMkDir)
				return nil, status.Error(codes.Internal, msg)
			}
			notMnt = true
		} else {
			msg := fmt.Sprintf("could not determine if %q is valid mount point: %v", target, err)
			return nil, status.Error(codes.Internal, msg)
		}
	}

	if !notMnt {
		msg := fmt.Sprintf("target %q is not a valid mount point", target)
		return nil, status.Error(codes.InvalidArgument, msg)
	}
	// Get fs type that the volume will be formatted with
	attributes := req.GetVolumeContext()
	fsType, exists := attributes["fsType"]
	if !exists || fsType == "" {
		fsType = defaultFsType
	}

	// FormatAndMount will format only if needed
	log.Printf("NodeStageVolume: formatting %s and mounting at %s", source, target)
	_, err = exec.Command("gcsfuse", volumeID, target).Output()
	if err != nil {
		log.Fatal(err)
	}

	err = d.mounter.FormatAndMount(source, target, fsType, nil)
	if err != nil {
		msg := fmt.Sprintf("could not format %q and mount it at %q", source, target)
		return nil, status.Error(codes.Internal, msg)
	}

	return &csi.NodeStageVolumeResponse{}, nil

	return nil, nil

}

func (d *Driver) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	log.Printf("NodeUnstageVolume: called with args %+v", *req)
	volumeID := req.GetVolumeId()
	if len(volumeID) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID not provided")
	}

	target := req.GetStagingTargetPath()
	if len(target) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Staging target not provided")
	}

	log.Printf("NodeUnstageVolume: unmounting %s", target)
	err := d.mounter.Interface.Unmount(target)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not unmount target %q: %v", target, err)
	}

	return &csi.NodeUnstageVolumeResponse{}, nil
}

func (d *Driver) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	log.Printf("NodePublishVolume: called with args %+v", *req)
	volumeID := req.GetVolumeId()
	if len(volumeID) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID not provided")
	}

	source := req.GetStagingTargetPath()
	if len(source) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Staging target not provided")
	}

	target := req.GetTargetPath()
	if len(target) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Target path not provided")
	}

	volCap := req.GetVolumeCapability()
	if volCap == nil {
		return nil, status.Error(codes.InvalidArgument, "Volume capability not provided")
	}

	//if !d.isValidVolumeCapabilities([]*cmd.VolumeCapability{volCap}) {
	//	return nil, status.Error(codes.InvalidArgument, "Volume capability not supported")
	//}

	options := []string{"bind"}
	if req.GetReadonly() {
		options = append(options, "ro")
	}

	log.Printf("NodePublishVolume: creating dir %s", target)
	if err := d.mounter.Interface.MakeDir(target); err != nil {
		return nil, status.Errorf(codes.Internal, "Could not create dir %q: %v", target, err)
	}

	log.Printf("NodePublishVolume: mounting %s at %s", source, target)
	if err := d.mounter.Interface.Mount(source, target, "ext4", options); err != nil {
		os.Remove(target)
		return nil, status.Errorf(codes.Internal, "Could not mount %q at %q: %v", source, target, err)
	}

	return &csi.NodePublishVolumeResponse{}, nil
}

func (d *Driver) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	log.Printf("NodeUnpublishVolume: called with args %+v", *req)
	volumeID := req.GetVolumeId()
	if len(volumeID) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID not provided")
	}

	target := req.GetTargetPath()
	if len(target) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Target path not provided")
	}

	log.Printf("NodeUnpublishVolume: unmounting %s", target)
	err := d.mounter.Interface.Unmount(target)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not unmount %q: %v", target, err)
	}

	return &csi.NodeUnpublishVolumeResponse{}, nil
}

func (d *Driver) NodeGetVolumeStats(ctx context.Context, req *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "NodeGetVolumeStats is not implemented yet")

}
func (d *Driver) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	log.Printf("NodeGetCapabilities: called with args %+v", *req)
	var caps []*csi.NodeServiceCapability
	for _, cap := range d.nodeCaps {
		c := &csi.NodeServiceCapability{
			Type: &csi.NodeServiceCapability_Rpc{
				Rpc: &csi.NodeServiceCapability_RPC{
					Type: cap,
				},
			},
		}
		caps = append(caps, c)
	}
	return &csi.NodeGetCapabilitiesResponse{Capabilities: caps}, nil
}

func (d *Driver) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	//log.Infof("NodeGetInfo: called with args %+v", *req)
	//m := d.cloud.GetMetadata()
	//
	//topology := &cmd.Topology{
	//	Segments: map[string]string{topologyKey: m.GetAvailabilityZone()},
	//}
	//
	//return &cmd.NodeGetInfoResponse{
	//	NodeId:             m.GetInstanceID(),
	//	AccessibleTopology: topology,
	//}, nil

	return  &csi.NodeGetInfoResponse{}, nil
}

func verifyTargetDir(target string) error {
	if target == "" {
		return status.Error(codes.InvalidArgument,
			"target path required")
	}

	tgtStat, err := os.Stat(target)
	if err != nil {
		if os.IsNotExist(err) {
			return status.Errorf(codes.FailedPrecondition,
				"target: %s not pre-created", target)
		}
		return status.Errorf(codes.Internal,
			"failed to stat target, err: %s", err.Error())
	}

	// This check is mandated by the spec, but this would/should fail if the
	// volume has a block accessType. Maybe staging isn't intended to be used
	// with block? That would make sense you cannot share the volume for block.
	if !tgtStat.IsDir() {
		return status.Errorf(codes.FailedPrecondition,
			"existing path: %s is not a directory", target)
	}

	return nil
}
