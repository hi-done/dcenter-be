// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"time"
)

type (
	IGold interface {
		GetGoldPrice() (float64, string, time.Time)
	}
)

var (
	localGold IGold
)

func Gold() IGold {
	if localGold == nil {
		panic("implement not found for interface IGold, forgot register?")
	}
	return localGold
}

func RegisterGold(i IGold) {
	localGold = i
}
