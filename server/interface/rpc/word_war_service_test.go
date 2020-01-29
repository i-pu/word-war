package rpc

import (
	"context"
	"os"
	"testing"

	pb "github.com/i-pu/word-war/server/interface/rpc/pb"
	"github.com/i-pu/word-war/server/repository/memory"
	"github.com/i-pu/word-war/server/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GRPCTestHealthCheckSuite struct {
	suite.Suite
	grpcServer *wordWarService
}

const (
	serverVersion = "1.0.10"
)

func (suite *GRPCTestHealthCheckSuite) SetupTest() {
	if err := os.Setenv("SERVER_VERSION", serverVersion); err != nil {
		suite.Error(err)
	}
	messageRepo := memory.NewMessageRepository()

	counterRepo := memory.NewCounterRepository()
	counterUsecase := usecase.NewCounterUsecase(counterRepo)

	gameRepo := memory.NewGameStateRepository()
	gameUsecase := usecase.NewGameUsecase(gameRepo, messageRepo, counterRepo)

	resultRepo := memory.NewResultRepository()
	resultUsecase := usecase.NewResultUsecase(resultRepo, gameRepo)

	matchingUsecase := usecase.NewMatchingUsecase(gameRepo)
	suite.grpcServer = newWordWarService(gameUsecase, counterUsecase, resultUsecase, matchingUsecase)
}

func (suite *GRPCTestHealthCheckSuite) TearDownTest() {
	// defer goleak.VerifyNone(suite.T())
	if err := os.Unsetenv("SERVER_VERSION"); err != nil {
		suite.Error(err)
	}
}

func (suite *GRPCTestHealthCheckSuite) TestWordWarService_HealthCheck() {
	// FIREBASE_CREDENTIALS: /go/src/app/serviceAccount.json
	// SERVER_VERSION: 1.2.1

	in := pb.HealthCheckRequest{}
	res, err := suite.grpcServer.HealthCheck(context.Background(), &in)
	if err != nil {
		suite.Error(err)
	}
	assert.Equalf(suite.T(), serverVersion, res.ServerVersion, "should be equal.")
	assert.True(suite.T(), res.Active, "should be true.")
}

func TestRunGRPCTestHealthCheckSuite(t *testing.T) {
	suite.Run(t, new(GRPCTestHealthCheckSuite))
}
