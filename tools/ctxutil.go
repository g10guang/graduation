package tools

import (
	"context"
	"math/rand"
	"strconv"
	"time"
)

// log id in context
const CTX_LOG_ID = "log_id"

// create log id
func MakeLogID() string {
	return time.Now().Format("20060102150405") + strconv.Itoa(int(GetLocalIpInt())) + strconv.Itoa(int(rand.Int31()))
}

// create a context with log id
func NewCtxWithLogID() context.Context {
	ctx := context.Background()
	context.WithValue(ctx, CTX_LOG_ID, MakeLogID())
	return ctx
}
