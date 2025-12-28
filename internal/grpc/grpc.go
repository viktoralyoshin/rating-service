package grpc

import (
	"rating-service/internal/config"

	"github.com/rs/zerolog/log"
	gamepb "github.com/viktoralyoshin/playhub-proto/gen/go/games"
	socialpb "github.com/viktoralyoshin/playhub-proto/gen/go/social"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var GamesClient gamepb.GameServiceClient
var SocialClient socialpb.SocialServiceClient

func Connect(cfg *config.Config) {
	gamesServiceConn := connect(cfg.GameServiceAddr)
	socialServiceConn := connect(cfg.SocialServiceAddr)

	GamesClient = gamepb.NewGameServiceClient(gamesServiceConn)
	SocialClient = socialpb.NewSocialServiceClient(socialServiceConn)
}

func connect(addr string) *grpc.ClientConn {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().
			Err(err).
			Str("target_addr", addr).
			Msg("failed to initialize grpc connection")
	}

	return conn
}
