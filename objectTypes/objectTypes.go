package objectTypes

import "chatting-system-backend/model"

type LoginCredentials struct {
	*model.User
}
