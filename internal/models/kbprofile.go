package models

type KBProfile byte

const (
	Empty KBProfile = iota
	User
	Manager
	Admin
)
