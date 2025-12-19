package valueobjects

type PasswordHash struct {
	value string
}

func NewPasswordHash(value string) (pass PasswordHash, err error) {
	// Aquí podrías agregar lógica para validar o procesar el hash
	return PasswordHash{value: value}, nil
}

func (p PasswordHash) Value() string {
	return p.value
}
