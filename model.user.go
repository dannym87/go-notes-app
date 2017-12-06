package main

type User struct {
	BaseModel
	Email     string `json:"email"`
	Password  string `json:"-"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Scope     string `json:"scope"`
}
