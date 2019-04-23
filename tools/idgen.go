package tools

import (
	"github.com/bwmarrin/snowflake"
	"math/rand"
	"time"
)

// use snowflake lib: github.com/bwmarrin/snowflake

var idNode *snowflake.Node

func init() {
	var err error
	idNode, err = snowflake.NewNode(getNodeId())
	if err != nil {
		panic(err)
	}
}

func getNodeId() int64 {
	s := rand.NewSource(time.Now().UnixNano())
	return s.Int63() % 1024
}

func GenID() snowflake.ID {
	return idNode.Generate()
}
