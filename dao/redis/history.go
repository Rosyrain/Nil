package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"nil/models"
	"strconv"
	"time"
)

var (
	history = models.History{Capacity: 3, Size: 0}
	head    = &models.HistoryNode{PostID: "", Time: time.Now()}
	tail    = &models.HistoryNode{PostID: "", Time: time.Now()}
)

func InitUserHistory(uid int64) error {
	// 初始化双向链表
	head.Pre = nil
	head.Next = tail
	tail.Pre = head
	tail.Next = nil
	history.Cache = make(map[string]*models.HistoryNode)

	ukey := GetRedisKey(KeyUserHistoryPrefix + strconv.Itoa((int(uid))))
	zsetCmd := client.ZRangeWithScores(ukey, 0, -1)
	if zsetCmd.Err() != nil {
		return zsetCmd.Err()
	}
	zsetValues := zsetCmd.Val()

	for _, zsetValue := range zsetValues {
		member := zsetValue.Member.(string)
		score := zsetValue.Score
		unixTimestamp := int64(score)
		// 根据 Unix 时间戳创建 time.Time 对象
		t := time.Unix(unixTimestamp, 0)
		historyNode := &models.HistoryNode{
			PostID: member,
			Time:   t,
			Pre:    nil,
			Next:   nil,
		}
		history.Cache[member] = historyNode
		addToHead(historyNode)
		history.Size += 1
	}
	return nil
}

func InsertHistory(pid string, uid int64) error {
	if err := InitUserHistory(uid); err != nil {
		return err
	}
	ukey := GetRedisKey(KeyUserHistoryPrefix + strconv.Itoa((int(uid))))
	_, ok := history.Cache[pid]
	fmt.Println(ok)
	if !ok {
		node := &models.HistoryNode{
			PostID: pid,
			Time:   time.Now(),
			Pre:    nil,
			Next:   nil,
		}
		addToHead(node)
		history.Cache[pid] = node
		history.Size += 1
		if history.Size > history.Capacity {
			removeNode := removeTail()
			delete(history.Cache, removeNode.PostID)
			client.ZRem(ukey, pid)
			history.Size -= 1
		}
	} else {
		node, _ := history.Cache[pid]
		node.Time = time.Now()
		moveToHead(node)
	}
	//无论是新加还是修改,都是ZAdd
	client.ZAdd(ukey, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: pid,
	})
	return nil
}

func addToHead(node *models.HistoryNode) {
	next := head.Next
	node.Pre = head
	node.Next = next
	head.Next = node
	next.Pre = node
}

func removeNode(node *models.HistoryNode) {
	node.Pre.Next = node.Next
	node.Next.Pre = node.Pre
}

func moveToHead(node *models.HistoryNode) {
	removeNode(node)
	addToHead(node)
}

func removeTail() (node *models.HistoryNode) {
	node = tail.Pre
	removeNode(node)
	return node
}

func GetUserHistoryList(p *models.ParamHistoryList) (data []string, err error) {
	ukey := GetRedisKey(KeyUserHistoryPrefix + strconv.Itoa((int(p.UserID))))
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	return client.ZRevRange(ukey, start, end).Result()
}
