package lib

import (
	"time"
)

const DateTimeFormat = "2006-01-02 15:04:05"
const DateTimeFormatNoSeparator = "20060102150405"

var NowFunc = time.Now
