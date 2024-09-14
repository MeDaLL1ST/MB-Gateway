package pkg

import "mbgateway/config"

func validateToken(token string) bool {
	if token == config.Cfg.APIKey {
		return true
	} else {
		return false
	}
}
