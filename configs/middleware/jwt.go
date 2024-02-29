package middleware

import (
	"encoding/json"
	"fmt"
	config "payment-service/configs"
	userDto "payment-service/internal/modules/user/models/dto"
	userEntity "payment-service/internal/modules/user/models/entity"
	userRepoQueries "payment-service/internal/modules/user/repositories/queries"
	"payment-service/internal/pkg/constants"
	"payment-service/internal/pkg/databases/mongodb"
	"payment-service/internal/pkg/errors"
	helpers "payment-service/internal/pkg/helpers"
	"payment-service/internal/pkg/log"
	"payment-service/internal/pkg/redis"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

type Middlewares struct {
	redisClient redis.Collections
}

func NewMiddlewares(redis redis.Collections) Middlewares {
	return Middlewares{
		redisClient: redis,
	}
}

func (m Middlewares) VerifyBasicAuth() fiber.Handler {
	return basicauth.New(basicauth.Config{
		Users: map[string]string{
			config.GetConfig().UsernameBasicAuth: config.GetConfig().PasswordBasicAuth,
		},
	})
}

func (m Middlewares) VerifyBearer() fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger := log.GetLogger()
		redisClient := m.redisClient
		helperImpl := &helpers.JwtImpl{}
		parseToken, err := helperImpl.JWTAuthorization(c.Request())
		if err != nil {
			return helpers.RespError(c, logger, err)
		}
		token := strings.Split(string(c.Request().Header.Peek("Authorization")), " ")
		if len(token) != 2 || (token[0] != "Bearer" && token[0] != "bearer") {
			logger.Error(c.Context(), "Invalid token format", token[1])
			return helpers.RespError(c, logger, errors.ForbiddenError("Invalid token format"))
		}

		blocklist, _ := redisClient.Get(c.Context(), fmt.Sprintf("%s:%s", constants.RedisKeyBlockListJwt, token[1])).Result()
		if blocklist != "" {
			logger.Error(c.Context(), "Access token expired!", "Token blocklist")
			return helpers.RespError(c, logger, errors.UnauthorizedError("Access token expired!"))
		}
		result, _ := redisClient.Get(c.Context(), fmt.Sprintf("%s:%s", constants.RedisKeyGetProfileUser, parseToken.UserId)).Result()
		if result == "" {
			userQueryMongodbRepo := userRepoQueries.NewQueryMongodbRepository(mongodb.NewMongoDBLogger(mongodb.GetSlaveConn(), mongodb.GetSlaveDBName(), logger), logger)
			resp := <-userQueryMongodbRepo.FindOneUserId(c.Context(), parseToken.UserId)
			if resp.Error != nil {
				return helpers.RespError(c, logger, resp.Error)
			}
			if resp.Data == nil {
				return helpers.RespError(c, logger, errors.ForbiddenError("Invalid token!"))
			}
			convert, ok := resp.Data.(*userEntity.User)
			if !ok {
				return helpers.RespError(c, logger, errors.UnauthorizedError("Access token expired!"))
			}
			dataUser, _ := json.Marshal(userDto.UserData{
				Data: userDto.UserResp{
					FullName:     convert.FullName,
					Email:        convert.Email,
					MobileNumber: convert.MobileNumber,
					Role:         convert.Role,
					UserId:       convert.UserId,
					CreatedAt:    convert.CreatedAt,
					UpdatedAt:    convert.UpdatedAt,
				},
			})
			redisClient.Set(c.Context(), fmt.Sprintf("%s:%s", constants.RedisKeyGetProfileUser, parseToken.UserId), dataUser, 20*time.Minute)
		}
		c.Locals("userId", parseToken.UserId)
		c.Locals("userRole", parseToken.Role)
		return c.Next()
	}

}
