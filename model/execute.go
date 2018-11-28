package model

type Execute interface {
	Execute() ([]byte, Error)
}