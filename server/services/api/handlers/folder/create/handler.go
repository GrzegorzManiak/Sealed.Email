package create

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
)

func Handler(input *Input, data *services.Handler) (*Output, errors.AppError) {
	ownsDomain, domainId := CheckDomainOwnership(data.DatabaseConnection, input.DomainID, data.User.ID)
	if !ownsDomain {
		return nil, errors.User(
			"Could not create the folder.",
			"Folder creation failed!",
		)
	}

	folder, err := createFolder(input, data, domainId)
	if err != nil {
		return nil, err
	}

	return &Output{
		Folder: *folder,
	}, nil
}
