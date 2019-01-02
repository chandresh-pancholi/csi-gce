package driver

import (
	"context"
	"github.com/chandresh-pancholi/csi-gce/pkg/util"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	"k8s.io/kubernetes/pkg/util/mount"
	"net"
)

type Driver struct {
	endpoint string
	nodeID   string
	version  string

	srv   *grpc.Server

	mounter *mount.SafeFormatAndMount

	volumeCaps     []csi.VolumeCapability_AccessMode
	controllerCaps []csi.ControllerServiceCapability_RPC_Type
	nodeCaps       []csi.NodeServiceCapability_RPC_Type
}

func NewDriver(mounter *mount.SafeFormatAndMount, endPoint, nodeId string) *Driver  {
	if mounter == nil {
		 mounter = newSafeMounter()
	}


	return &Driver{
		endpoint: endPoint,
		//cloud:    cloud,
		//mounter:  mounter,
		nodeID: nodeId,
		volumeCaps: []csi.VolumeCapability_AccessMode{
			{
				Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
			},
		},
		controllerCaps: []csi.ControllerServiceCapability_RPC_Type{
			csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
			csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME,
		},
		nodeCaps: []csi.NodeServiceCapability_RPC_Type{
			csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
		},
	}
}

func newSafeMounter() *mount.SafeFormatAndMount  {
	return &mount.SafeFormatAndMount{
		Interface: mount.New(""),
		Exec: mount.NewOsExec(),
	}
}

func (d *Driver) Run() error {
	scheme, addr, err := util.ParseEndpoint(d.endpoint)
	if err != nil {
		return err
	}

	listener, err := net.Listen(scheme, addr)
	if err != nil {
		return err
	}

	logErr := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			glog.Errorf("GRPC error: %v", err)
		}
		return resp, err
	}
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(logErr),
	}
	d.srv = grpc.NewServer(opts...)

	csi.RegisterIdentityServer(d.srv, d)
	csi.RegisterControllerServer(d.srv, d)
	csi.RegisterNodeServer(d.srv, d)

	glog.Infof("Listening for connections on address: %#v", listener.Addr())
	return d.srv.Serve(listener)
}