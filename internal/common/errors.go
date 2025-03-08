package common

import "errors"

var ErrIngestNoSuchSection = errors.New("unable to find section in trace file")
var ErrDataWithNoSectionInTraceFile = errors.New("data found outside of sections in trace file")
