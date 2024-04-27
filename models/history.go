package models

import "time"

type HistoryNode struct {
	PostID string
	Time   time.Time
	Pre    *HistoryNode
	Next   *HistoryNode
}

type History struct {
	Cache    map[string]*HistoryNode
	Capacity int64
	Size     int64
}
