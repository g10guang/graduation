# -*- coding: utf-8 -*-
# author: Xiguang Liu<g10guang@foxmail.com>
# 2019-05-09 10:52
# test read_api /get interface
import time
import requests
import asyncio
import random

consume = []
fail = 0
qps = dict()
fids = []


async def get(lock):
    global consume, fail
    start = time.time()
    url = "http://127.0.0.1:10002/get"
    fid = random_fid()
    rsp = requests.post(url, data={"fid": fid})
    # print('text: {}'.format(rsp.text))
    # print('header: {}'.format(rsp.headers))
    await lock.acquire()
    try:
        if rsp.status_code != 200:
            fail += 1
        consume.append(time.time() - start)
        if int(start) not in qps:
            qps[int(start)] = 1
        else:
            qps[int(start)] += 1
    finally:
        lock.release()


def main():
    print("test read_api /get start")
    loop = asyncio.get_event_loop()
    lock = asyncio.Lock()
    tasks = []
    for i in range(10000):
        tasks.append(get(lock))
    start_ts = time.time()
    loop.run_until_complete(asyncio.wait(tasks))
    end_ts = time.time()
    print("test read_api /get end")
    print("consume: {}s".format(end_ts - start_ts))


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


def load_fids():
    with open("./fid.txt") as f:
        for line in f.readlines():
            fids.append(int(line))


def random_fid():
    global fids
    index = random.randint(0, len(fids) - 1)
    return fids[index]


if __name__ == '__main__':
    load_fids()
    main()
    statistic()
