package grpcurl

import (
	"errors"
	"github.com/golang/glog"
	"golang.org/x/net/context"
	"url-shortener/protos"
)

type ShortenerServer struct {
}

func (*ShortenerServer) ShortenerGenerate(ctx context.Context, request *protos.GenerateRequest) (reply *protos.GenerateReply, err error) {
	reply = &protos.GenerateReply{}
	if request.Pcurl == "" && request.MobileUrl == "" && request.MobileClientUrl == "" {
		return reply, errors.New("At least a normal url")
	}

	reply, err = ShortGenerate(request.Pcurl, request.MobileUrl, request.MobileClientUrl)
	if err != nil {
		glog.V(0).Infof("ShortenerQuery %v", err)
		return reply, err
	}
	return reply, err
}

func (*ShortenerServer) ShortenerQuery(ctx context.Context, request *protos.QueryRequest) (reply *protos.QueryReply, err error) {
	reply = &protos.QueryReply{}

	reply, err = EnerQuery(request.Url, request.Shortener, request.Page, request.PageSize)
	if err != nil {
		glog.V(0).Infof("ShortenerQuery %v", err)
		return reply, err
	}
	return reply, err
}

func (*ShortenerServer) ShortenerDelete(ctx context.Context, request *protos.DelEnerRequest) (reply *protos.DelEnerReply, err error) {

	reply = &protos.DelEnerReply{}
	reply, err = EnerDelete(request.Url, request.Shortener)

	if err != nil {
		glog.V(0).Infof("ShortenerDelete %v", err)
		return reply, err
	}
	return reply, err
}
