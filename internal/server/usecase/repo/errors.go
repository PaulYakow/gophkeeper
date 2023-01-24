package repo

const (
	ErrUserNotExist = RepoErr("user already exists")
)

type RepoErr string

func (e RepoErr) Error() string {
	return string(e)
}
