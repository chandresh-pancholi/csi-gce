package driver

import (
	"context"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/golang/protobuf/ptypes/wrappers"
	"log"
)

func (d *Driver) GetPluginInfo(ctx context.Context, req *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	log.Printf("GetPluginInfo: called with args %+v", *req)
	return &csi.GetPluginInfoResponse{
		Name:         DriverName,
		VendorVersion: "0.1.0",
	}, nil
}


func (d *Driver) GetPluginCapabilities(ctx context.Context, req *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	log.Printf("GetPluginCapabilities: called with args %+v", *req)
	resp := &csi.GetPluginCapabilitiesResponse{
		Capabilities: []*csi.PluginCapability{
			{
				Type: &csi.PluginCapability_Service_{
					Service: &csi.PluginCapability_Service{
						Type: csi.PluginCapability_Service_CONTROLLER_SERVICE,
					},
				},
			},
			{
				Type: &csi.PluginCapability_Service_{
					Service: &csi.PluginCapability_Service{
						//Type: csi.PluginCapability_Service_VOLUME_ACCESSIBILITY_CONSTRAINTS,
						Type: csi.PluginCapability_Service_ACCESSIBILITY_CONSTRAINTS,
					},
				},
			},
		},
	}

	return resp, nil
}

func (d *Driver) Probe(ctx context.Context, req *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	log.Printf("Probe: called with args %+v", *req)
	return &csi.ProbeResponse{
		Ready: &wrappers.BoolValue{
			Value: true,
		},
	}, nil
}