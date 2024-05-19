package check

import (
	"context"
	"testing"
)

// CheckerWhiteList 白名单
type CheckerWhiteList struct {
	preName string
	did     int64
}

func NewCheckerWhiteList(preName string, did int64) *CheckerWhiteList {
	return &CheckerWhiteList{
		preName: preName,
		did:     did,
	}
}

func (checker *CheckerWhiteList) DirectCheck() bool {
	// 这里是模拟根据 preName 从配置中心中读取白名单配置
	whiteList := []int64{123}
	for _, did := range whiteList {
		if checker.did == did {
			// log.info(ctx, "direct check success")
			return true
		}
	}
	return false
}

type CheckerReadTime struct {
	readTime int
}

func NewCheckerReadTime(readTime int) *CheckerReadTime {
	return &CheckerReadTime{
		readTime: readTime,
	}
}

func (c *CheckerReadTime) Check() bool {
	if c.readTime > 100 {
		return true
	}
	return false
}

// TestCheck 粗粒度白名单校验
func TestCheck(t *testing.T) {
	c1 := NewCheckerWhiteList("snack bar", 123)
	c2 := func() bool {
		// 做简单业务校验
		return false
	}
	c3 := NewCheckerReadTime(200)
	println(Check(context.Background(), c1, CheckerFun(c2), c3))
}

type DCheckerReadTime struct {
	CheckerWhiteList
	readTime int
}

func NewDCheckerReadTime(readTime int, preName string, did int64) *DCheckerReadTime {
	c := &DCheckerReadTime{
		readTime: readTime,
	}
	c.preName = preName
	c.did = did
	return c
}

func (c *DCheckerReadTime) Check() bool {
	if c.readTime > 100 {
		return true
	}
	return false
}

// TestCheck2 细粒度白名单校验
func TestCheck2(t *testing.T) {
	c1 := NewDCheckerReadTime(1, "test_check_2", 123) // 命中白名单，跳过检测
	c2 := NewDCheckerReadTime(1, "test_check_2", 111) // 没命中白名单，且业务校验不通过
	println(Check(context.Background(), c1, c2))
}
