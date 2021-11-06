package users

import (
	"fmt"
	"github.com/aveloper/blog/internal/http/response"
)

type UserNotFound struct {
	ID 		int32 `json:"id"`
}

func (d *UserNotFound) Message() string {
	return fmt.Sprintf("User not Found for the id %d", d.ID)
}

func (d *UserNotFound) Code() response.ErrorCode {
	return response.DefaultErrorCode
}

func (d *UserNotFound) Data() interface{} {
	return nil
}