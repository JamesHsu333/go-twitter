package utils

import (
	"context"

	"github.com/JamesHsu333/go-twitter/pkg/httpErrors"
	"github.com/JamesHsu333/go-twitter/pkg/logger"
)

// Validate is user from owner of content
func ValidateIsOwner(ctx context.Context, creatorID string, logger logger.Logger) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return err
	}

	if user.UserID.String() != creatorID {
		logger.Errorf(
			"ValidateIsOwner, userID: %v, creatorID: %v",
			user.UserID.String(),
			creatorID,
		)
		return httpErrors.Forbidden
	}

	return nil
}
