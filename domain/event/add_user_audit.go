package event

import (
	"github.com/reecerussell/distro-blog/domain/dto"
	"time"
)

type AddUserAudit struct {
	Message string
	Date time.Time
	UserID string
	PerformingUserID string
	Before *dto.User
	After *dto.User
}