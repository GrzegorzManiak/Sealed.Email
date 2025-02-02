package email

type EncryptedInbox struct {
	DisplayName       string `json:"displayName" validate:"lte=200"`
	EmailHash         string `json:"emailHash" validate:"required,email"`
	PublicKey         string `json:"publicKey" validate:"Encrypted-B64-Key"`
	EncryptedEmailKey string `json:"encryptedEmailKey" validate:"Encrypted-B64-Key"`
	Nonce             string `json:"nonce" validate:"required,base64,gte=40,lte=200"`
}
