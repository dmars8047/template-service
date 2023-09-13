package api

import "errors"

// ErrNotFound is returned when a record is not found
var ErrNotFound = errors.New("record not found")

// ErrNameConflict is returned when a name conflict is found
var ErrNameConflict = errors.New("name conflict found")
