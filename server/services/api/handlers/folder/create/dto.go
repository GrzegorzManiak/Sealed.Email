package create

import "github.com/GrzegorzManiak/NoiseBackend/services/api/session"

type Input struct {
	DomainID            string `json:"domainID" validate:"required,PublicID"`
	EncryptedFolderName string `json:"encryptedFolderName"`
}

type Output struct {
	Folder
}

type Folder struct {
	FolderID            string `json:"folderID"`
	EncryptedFolderName string `json:"encryptedFolderName"`
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: true,
}
