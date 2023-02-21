package lib

import "errors"

var (
	ClassNotFoundError   = errors.New("Class not found!")
	ObjectNotFoundError  = errors.New("Object not found!")
	UninitializedKBError = errors.New("Uninitialized KB!")
	InvalidClassError    = errors.New("Invalid class!")
	CompilerError        = errors.New("Compiler error in statement!")
	LiteralNotFoundError = errors.New("Literal not found!")
)
