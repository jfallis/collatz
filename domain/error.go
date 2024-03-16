package domain

type CollatzError string

func (ce CollatzError) Error() string {
	return string(ce)
}
