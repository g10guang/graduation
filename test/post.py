# -*- coding: utf-8 -*-
# author: Xiguang Liu<g10guang@foxmail.com>
# 2019-05-09 10:30
# test write_api /post interface
import time
import requests
import asyncio

consume = []
fail = 0
qps = {}


async def post(lock):
    global consume, fail
    start = time.time()
    url = "http://127.0.0.1:10003/post"
    with open("./test.png", mode="rb") as f:
        rsp = requests.post(url, data={"uid": 1}, files={"file": f})
    await lock.acquire()
    try:
        consume.append(time.time() - start)
        if rsp.status_code != 200:
            fail += 1
        if int(start) not in qps:
            qps[int(start)] = 1
        else:
            qps[int(start)] += 1
    finally:
        lock.release()


def main():
    start_ts = time.time()
    print('test write_api /post start')
    print('test write_api /post end')
    loop = asyncio.get_event_loop()
    lock = asyncio.Lock()
    task = []
    for i in range(1000):
        task.append(post(lock))
    loop.run_until_complete(asyncio.wait(task))
    end_ts = time.time()
    print('consume: {}s'.format(end_ts - start_ts))


def statistic():
    global consume, fail, qps
    print('statistic')
    total = sum(consume)
    average = sum(consume) / len(consume)
    print('total cost: {}s'.format(total))
    print('average cost: {}s'.format(average))
    print('fail: {}'.format(fail))
    total = 0
    for second, n in qps.items():
        total += n
    print("avg_qps: {}".format(total / len(qps)))


if __name__ == '__main__':
    main()
    statistic()
