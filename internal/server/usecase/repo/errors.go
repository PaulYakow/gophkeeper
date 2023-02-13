package repo

const (
	ErrUserExist = Err("user already exists")
)

type Err string

func (e Err) Error() string {
	return string(e)
}
