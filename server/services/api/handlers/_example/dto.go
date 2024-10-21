package register

type Input struct {
	// -- Required fields
	User string `json:"u" validate:"required"`
	PI   string `json:"π" validate:"required"`
	T    string `json:"t" validate:"required"`
	TOS  bool   `json:"tos" validate:"required"`

	// -- Optional fields
	RecoveryEmail string `json:"recoveryEmail" validate:"omitempty,email"`
}

type Output struct {
	X3  []byte `json:"x3"`
	PI3 string `json:"π3"`
}
