package port

import "fmt"

type DBFailureError struct { // UnWrappable
	Msg        string
	WrappedErr error
}

func (sc DBFailureError) Error() string {
	return fmt.Sprintf("%s. %s", sc.Msg, sc.WrappedErr.Error())
}

func (sc DBFailureError) Unwrap() error {
	return sc.WrappedErr
}

type ValidationError struct { // UnWrappable
	Msg        string
	WrappedErr error
}

func (sc ValidationError) Error() string {
	if sc.WrappedErr != nil {
		return fmt.Sprintf("%s, %s", sc.Msg, sc.WrappedErr.Error())
	}
	return sc.Msg
}

func (sc ValidationError) Unwrap() error {
	return sc.WrappedErr
}

type DependencyError struct { // UnWrappable
	Msg        string
	WrappedErr error
}

func (sc DependencyError) Error() string {
	if sc.WrappedErr != nil {
		return fmt.Sprintf("%s, %s", sc.Msg, sc.WrappedErr.Error())
	}
	return sc.Msg
}

func (sc DependencyError) Unwrap() error {
	return sc.WrappedErr
}
