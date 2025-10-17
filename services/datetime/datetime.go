package datetime

import "time"

const MinTimeString = "1970-01-01 00:00:00"

const MaxTimeString = "9999-12-31 23:59:59"

var UtcMinTime, _ = time.ParseInLocation("2006-01-02 15:04:05", MinTimeString, time.UTC)
var UtcMaxTime, _ = time.ParseInLocation("2006-01-02 15:04:05", MaxTimeString, time.UTC)
var NullTime = time.Time{}
