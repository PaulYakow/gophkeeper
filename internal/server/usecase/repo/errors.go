package repo

const (
	ErrUserExist = RepoErr("user already exists")
)

type RepoErr string

func (e RepoErr) Error() string {
	return string(e)
}
