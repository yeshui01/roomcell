package timeutil

import "time"

// 零点时间
func ZeroTime(t int64) int64 {
	tmObj := time.Unix(t, 0)
	return t - (int64(tmObj.Hour())*3600 + int64(tmObj.Minute())*60 + int64(tmObj.Second()))
}

// 当前时间:秒
func NowTime() int64 {
	return time.Now().Unix()
}

// 当前时间:毫秒
func NowTimeMs() int64 {
	return time.Now().UnixNano() / 1000000
}

// 判断是同一周
func IsSameWeek(start, end time.Time) bool {
	// >=7天，一定是两个不同的周
	if end.Sub(start).Seconds() >= (7 * 24 * time.Hour).Seconds() {
		return false
	}
	w1 := start.Weekday()
	w2 := end.Weekday()
	switch time.Sunday {
	case w1:
		w1 = 7
	case w2:
		w2 = 7
	}
	return w2 >= w1
}

// 时间字符串转换成时间戳("2019-09-16" 或者 "2019-09-16 00:00:00")
func TimeStrToTimestamp(timeStr string) int64 {
	var t int64
	loc, _ := time.LoadLocation("Local")
	if len(timeStr) == 19 {
		t1, _ := time.ParseInLocation("2006/01/02/15/04/05", timeStr, loc)
		t = t1.Unix()
	} else if len(timeStr) == 10 {
		t1, _ := time.ParseInLocation("2006-01-02", timeStr, loc)
		t = t1.Unix()
	} else if len(timeStr) == 8 {
		t1, _ := time.ParseInLocation("20060102", timeStr, loc)
		t = t1.Unix()
	}
	return t
}
