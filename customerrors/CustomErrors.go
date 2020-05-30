package customerrors

import "errors"

// ErrDownloadFile complains about the inability to download a file
var ErrDownloadFile = errors.New("Couldn't download file")

// ErrCreateFile complains about the inability to create or save a file
var ErrCreateFile = errors.New("Couldn't create or save file")

// ErrHashChecksum complains that a hash check fails
var ErrHashChecksum = errors.New("Hash checksum doesn't match")

// ErrMMDBNotFound complains that the db is not found or unreadable
var ErrMMDBNotFound = errors.New("Unable to read Maxmind DB")
