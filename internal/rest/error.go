package rest

import "github.com/oligzeev/host-manager/internal/domain"

type Error struct {
	Ops      []domain.ErrOp `json:"ops"`
	Messages []string       `json:"messages"`
}

func E(err error) *Error {
	return &Error{
		Ops:      domain.EOps(err),
		Messages: domain.EMsgs(err),
	}
}
