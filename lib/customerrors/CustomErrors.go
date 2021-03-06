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

// ErrInvalidIPAddress complains about an invalid ip address
var ErrInvalidIPAddress = errors.New("IP Address passed is invalid")

// ErrGeneratePassword complains about not being able to generate a key
var ErrGeneratePassword = errors.New("Couldn't generate api key. Handle access appropriately")
