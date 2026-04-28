package models

import "errors"

var UnrecognisedAction error = errors.New("Action not recognised")
var TooFewArguments error = errors.New("Too few arguments")
var AddressFormattedWrong error = errors.New("Address improperly formatted")
var MXNotFound error = errors.New("MX not found")
