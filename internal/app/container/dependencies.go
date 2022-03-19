package container

import (
	"github.com/bhankey/pharmacy-automatization-user/internal/adapter/repository/emailrepo"
	"github.com/bhankey/pharmacy-automatization-user/internal/adapter/repository/onetimecodesrepo"
	"github.com/bhankey/pharmacy-automatization-user/internal/adapter/repository/userrepo"
	"github.com/bhankey/pharmacy-automatization-user/internal/delivery/grpc/user"
	"github.com/bhankey/pharmacy-automatization-user/internal/service/userservice"
	"time"
)

func (c *Container) GetUserGRPCHandler() *user.GRPCHandler {
	const key = "UserGRPCHandler"

	dependency, ok := c.dependencies[key]
	if ok {
		typedDependency, ok := dependency.(*user.GRPCHandler)
		if ok {
			return typedDependency
		}
	}

	typedDependency := user.NewUserGRPCHandler(c.getUserSrv())

	c.dependencies[key] = typedDependency

	return typedDependency
}

func (c *Container) getUserSrv() *userservice.UserService {
	const key = "UserSrv"

	dependency, ok := c.dependencies[key]
	if ok {
		typedDependency, ok := dependency.(*userservice.UserService)
		if ok {
			return typedDependency
		}
	}

	typedDependency := userservice.NewUserService(
		c.getUserStorage(),
		c.getEmailStorage(),
		c.getOneTimeCodesPasswordStorage(),
	)

	c.dependencies[key] = typedDependency

	return typedDependency
}

func (c *Container) getUserStorage() *userrepo.Repository {
	const key = "UserStorage"

	dependency, ok := c.dependencies[key]
	if ok {
		typedDependency, ok := dependency.(*userrepo.Repository)
		if ok {
			return typedDependency
		}
	}

	typedDependency := userrepo.NewUserRepo(c.masterPostgresDB, c.slavePostgresDB)

	c.dependencies[key] = typedDependency

	return typedDependency
}

func (c *Container) getEmailStorage() *emailrepo.EmailRepo {
	const key = "EmailStorage"

	dependency, ok := c.dependencies[key]
	if ok {
		typedDependency, ok := dependency.(*emailrepo.EmailRepo)
		if ok {
			return typedDependency
		}
	}

	typedDependency := emailrepo.NewEmailRepo(c.smtpClient, c.smtpMessageFrom)

	c.dependencies[key] = typedDependency

	return typedDependency
}

func (c *Container) getOneTimeCodesPasswordStorage() *onetimecodesrepo.ResetCodesRepo {
	const key = "OneTimeCodesPasswordStorage"

	dependency, ok := c.dependencies[key]
	if ok {
		typedDependency, ok := dependency.(*onetimecodesrepo.ResetCodesRepo)
		if ok {
			return typedDependency
		}
	}

	const timeOfLife = time.Second * 15 // TODO move to config or something else

	typedDependency := onetimecodesrepo.NewResetCodesRepo(c.redisConnection, timeOfLife)

	c.dependencies[key] = typedDependency

	return typedDependency
}
