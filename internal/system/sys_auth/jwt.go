package sys_auth

import (
	"github.com/Gary-Yez/go-admin/internal/state"
	"github.com/Gary-Yez/go-admin/internal/utils"
)

func generateLoginToken(userID, roleID uint, loginVersion uint64) (string, error) {
	lifetime, err := utils.ReadJWTLifetime()
	if err != nil {
		return "", err
	}
	return utils.NewJwt(state.Config().JWT.Secret).Generate(userID, roleID, loginVersion, lifetime)
}
