package service

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func loggingInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	start := time.Now()
	address := "unknown"
	fullMethod := strings.Split(info.FullMethod, "/")

	p, ok := peer.FromContext(ctx)
	if ok {
		address = p.Addr.String()
	}

	res, err := handler(ctx, req)
	slog.LogAttrs(
		ctx,
		slog.LevelInfo,
		"bannerRotation",
		slog.String("from", address),
		slog.String("method", fullMethod[len(fullMethod)-1]),
		slog.Any("err", err),
		slog.String("dur", time.Since(start).String()),
	)

	return res, err
}
