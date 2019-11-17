package zaplog

type OpError struct {
	Op  string
	Err error
}
