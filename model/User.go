package model

type User struct {
	Base
	Name     string
	Email    string
	Password uint
}