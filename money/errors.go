package money

// defining error
type Error string

// Error implementing error interface
func (e Error) Error() string{
	return string(e)
}