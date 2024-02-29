package queries_test

import (
	"context"
	"payment-service/internal/modules/user"
	mongoRQ "payment-service/internal/modules/user/repositories/queries"
	mocks "payment-service/mocks/pkg/databases/mongodb"
	mocklog "payment-service/mocks/pkg/log"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CommandTestSuite struct {
	suite.Suite
	mockMongodb *mocks.Collections
	mockLogger  *mocklog.Logger
	repository  user.MongodbRepositoryQuery
	ctx         context.Context
}

func (suite *CommandTestSuite) SetupTest() {
	suite.mockMongodb = new(mocks.Collections)
	suite.mockLogger = &mocklog.Logger{}
	suite.repository = mongoRQ.NewQueryMongodbRepository(
		suite.mockMongodb,
		suite.mockLogger,
	)
	suite.ctx = context.WithValue(context.TODO(), "key", "value")
}

func TestCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CommandTestSuite))
}
