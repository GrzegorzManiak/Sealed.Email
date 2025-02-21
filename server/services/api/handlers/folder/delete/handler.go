package delete

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/folder/create"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
)

func Handler(input *Input, data *services.Handler) (*Output, errors.AppError) {
	ownsDomain, domainId := create.CheckDomainOwnership(data.DatabaseConnection, input.DomainID, data.User.ID)
	if !ownsDomain {
		return nil, errors.User(
			"Could not create the folder.",
			"Folder creation failed!",
		)
	}

	if err := deleteFolder(data.DatabaseConnection, domainId, data.User.ID, input.FolderID); err != nil {
		return nil, errors.User("Sorry, we couldn't delete the folder.", "Folder deletion failed!")
	}

	return &Output{}, nil
}
