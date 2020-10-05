package gsm7bit

import "errors"

var (
	ErrInvalidCharacter = errors.New("gsm7bit: invalid gsm7 character")
	ErrInvalidByte      = errors.New("gsm7bit: invalid gsm7 byte")
)
