package dto

import "time"

type UserAudit struct {
	Message string `json:"message"`
	Date time.Time `json:"date"`
	UserID string `json:"userId"`
	UserFullname string `json:"userFullname"`
}
