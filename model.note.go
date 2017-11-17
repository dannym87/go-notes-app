package main

import "time"

type Note struct {
	Id      uint32    `jsonapi:"primary,notes"`
	Title   string    `jsonapi:"attr,title"`
	Created time.Time `jsonapi:"attr,created_at,iso8601"`
}
