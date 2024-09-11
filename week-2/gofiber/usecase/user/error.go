package user

type UserNotFoundError struct{}

func (e *UserNotFoundError) Error() string {
	return "user not found"
}
