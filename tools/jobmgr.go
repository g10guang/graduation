package tools

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"runtime"
	"sync"
	"time"
)

type Job interface {
	Run() (interface{}, error)
	GetName() string
}

type Result struct {
	Result interface{}
	Err    error
}

// 用于执行并发任务的 JobMgr
type JobMgr struct {
	jobs      []Job
	timeout   time.Duration
	lock      sync.Mutex
	doneCh    chan bool
	errMap    map[string]error
	resultMap map[string]*Result
}

var ErrJobMgrTimeout = errors.New("JobMgr job timeout")
var ErrJobFail = errors.New("Some Job Fail. Please Check")

func NewJobMgr(timeout time.Duration) *JobMgr {
	return &JobMgr{
		timeout:   timeout,
		jobs:      make([]Job, 0),
		resultMap: make(map[string]*Result),
		doneCh:    make(chan bool),
	}
}

func (mgr *JobMgr) AddJob(j Job) {
	mgr.jobs = append(mgr.jobs, j)
}

// 只要有一个 Job 返回 Error 则返回 error
func (mgr *JobMgr) Start(ctx context.Context) error {
	mgr.parallel(ctx)
	return mgr.join(ctx)
}

// 每一个 job 对应一个 goroutine 执行
func (mgr *JobMgr) parallel(ctx context.Context) {
	startTime := time.Now().UnixNano()
	for _, job := range mgr.jobs {
		go func(j Job) {
			var err error
			var result interface{}
			// 不能让当天 job 影响到其他 goroutine
			defer func() {
				if e := recover(); e != nil {
					const size = 64 << 10
					buf := make([]byte, size)
					buf = buf[:runtime.Stack(buf, false)]
					logrus.Panicf("job: %s panic", j.GetName())
					e = fmt.Errorf("job: %s panic", j.GetName())
					mgr.doneCh <- false
				} else {
					// 通知 join 当天 job 结束
					mgr.doneCh <- err == nil
				}
			}()
			defer func() {
				logrus.Debugf("job: %s consume: %dns", j.GetName(), time.Now().UnixNano()-startTime)
			}()
			result, err = j.Run()
			mgr.SetResult(j.GetName(), &Result{
				Result: result,
				Err:    err,
			})
		}(job)
	}
}

func (mgr *JobMgr) join(ctx context.Context) error {
	timeout := time.After(mgr.timeout)
	errNum := 0
	for i := 0; i < len(mgr.jobs); i++ {
		// 等待所有 job 退出或者超时
		select {
		case msg := <-mgr.doneCh:
			if msg == false {
				errNum++
			}
		case <-timeout:
			logrus.Errorf("jobmgr timeout: %v finish job num: %d", mgr.timeout, i)
			return ErrJobMgrTimeout
		}
	}

	if errNum > 0 {
		return ErrJobFail
	}
	return nil
}

func (mgr *JobMgr) SetResult(jobName string, result *Result) {
	// golang map 并非并发安全
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	mgr.resultMap[jobName] = result
}

func (mgr *JobMgr) GetResult(jobName string) *Result {
	return mgr.resultMap[jobName]
}
