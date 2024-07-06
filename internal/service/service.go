//go:generate protoc --proto_path=../../api --go_out=api --go_opt=paths=source_relative --go-grpc_out=api --go-grpc_opt=paths=source_relative ../../api/api.proto
//go:generate mockery --name App
package service

import (
	"context"
	"net"
	"strconv"

	"github.com/trb1maker/gobannercircle/internal/service/api"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type App interface {
	BannerOn(ctx context.Context, slotID, bannerID int) error
	BannerOff(ctx context.Context, slotID, bannerID int) error
	Banner(ctx context.Context, slotID, groupID int) (bannerID int, err error)
	Click(ctx context.Context, slotID, bannerID, groupID int) error
}

func NewService(app App, host string, port uint16) *Service {
	return &Service{
		app:  app,
		addr: net.JoinHostPort(host, strconv.Itoa(int(port))),
	}
}

type Service struct {
	api.UnimplementedBannerRotationServer

	app  App
	srv  *grpc.Server
	addr string
}

func (s *Service) Start() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer lis.Close()

	s.srv = grpc.NewServer(grpc.UnaryInterceptor(loggingInterceptor))
	api.RegisterBannerRotationServer(s.srv, s)

	if err := s.srv.Serve(lis); err != nil {
		return err
	}
	return nil
}

func (s *Service) AddBanner(ctx context.Context, r *api.SlotConfigRequest) (*emptypb.Empty, error) {
	if err := s.app.BannerOn(ctx, int(r.GetSlotId()), int(r.GetBannerId())); err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *Service) RemoveBanner(ctx context.Context, r *api.SlotConfigRequest) (*emptypb.Empty, error) {
	if err := s.app.BannerOff(ctx, int(r.GetSlotId()), int(r.GetBannerId())); err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *Service) Find(ctx context.Context, r *api.BannerRequest) (*api.BannerResponse, error) {
	bannerID, err := s.app.Banner(ctx, int(r.GetSlotId()), int(r.GetGroupId()))
	if err != nil {
		return nil, err
	}
	return &api.BannerResponse{
		BannerId: int64(bannerID),
	}, nil
}

func (s *Service) Success(ctx context.Context, r *api.SuccessRequest) (*emptypb.Empty, error) {
	if err := s.app.Click(ctx, int(r.GetSlotId()), int(r.GetGroupId()), int(r.GetGroupId())); err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *Service) Stop(ctx context.Context) {
	<-ctx.Done()
	s.srv.Stop()
}
