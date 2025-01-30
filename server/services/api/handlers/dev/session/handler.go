package devSession

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
)

func Handler(input *Input, data *services.Handler) (*Output, helpers.AppError) {
	return &Output{
		ContentSessionID: data.Session.Content.SessionID,
		ContentUserID:    data.Session.Content.UserID,

		Token:     data.Session.Token,
		Signature: data.Session.Signature,

		HeaderVersion: data.Session.Header.Ver,
		HeaderGroup:   data.Session.Header.Group,
		HeaderEat:     data.Session.Header.ExpiresAt,
		HeaderCat:     data.Session.Header.CreatedAt,
		HeaderRat:     data.Session.Header.RefreshAt,
	}, nil
}
