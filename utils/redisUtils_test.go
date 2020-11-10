package utils

import (
	"testing"
)

//连接redis
func TestRedisInit(t *testing.T) {
	RedisInit()
}

//测试RedisSelectDB
func TestRedisSelectDB(t *testing.T) {
	RedisSelectDB(RedisInit())
}

func TestRedisExample(t *testing.T) {
	RedisExample()
}

//RedisSet
func TestRedisSet(t *testing.T) {
	RedisSet(RedisInit(), "k1abc123", "1234")
}

//RedisGet
func TestRedisGet(t *testing.T) {
	RedisGet(RedisInit(), "k1abc123")
}

func TestRedisDelete(t *testing.T) {
	RedisDelete(RedisInit(), "kabc123")
}

//RedisExpireSet
func TestRedisExpireSet(t *testing.T) {
	RedisExpireSet(RedisInit(), "k1abc123", 10)
}

//RedisHSet  覆盖set
func TestRedisHSet(t *testing.T) {
	//RedisHSet(RedisInit(), "tingabc12", "k1", "v1")
	RedisHSet(RedisInit(), "tingabc12", "k1", "v2")
}

//RedisHGet
func TestRedisHGet(t *testing.T) {
	RedisHGet(RedisInit(), "kabc12", "k1")
}

//RedisHMSet
func TestRedisHMSet(t *testing.T) {
	clear := make(map[string]string, 0)
	clear["field"] = "abc"
	RedisHMSet(RedisInit(), "clearlingTJabc", clear)
}

//RedisHMGet
func TestRedisHMGet(t *testing.T) {
	RedisHMGet(RedisInit(), "clear", OldData(30))
}
