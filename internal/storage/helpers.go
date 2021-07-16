package storage

import "regexp"

var Rgx = regexp.MustCompile(`(C|H|T|DP|DV|DO|DS|DU|DF|DA|DB)M(\d{6})T(\d{2})_([1-9]\d{3})(\d{1,10})\.(ZIP|zip)$`)
