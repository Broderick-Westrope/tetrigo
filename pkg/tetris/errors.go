package tetris

import "errors"

var (
	ErrBufferZoneTooSmall = errors.New("matrix height must be greater than 20 to allow for a buffer zone of 20 lines")
)
