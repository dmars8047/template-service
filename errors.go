package main

import "errors"

var ErrNotFound = errors.New("record not found")

var ErrNameConflict = errors.New("name conflict found")
