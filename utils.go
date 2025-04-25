package main

import "time"

// 返回几天前的 UTC 零点时间戳
func GetTimestampDaysAgo(days int) int64 {
	now := time.Now().UTC()
	daysAgo := now.AddDate(0, 0, -days)
	zero := time.Date(daysAgo.Year(), daysAgo.Month(), daysAgo.Day(), 0, 0, 0, 0, time.UTC)
	return zero.Unix()
}
