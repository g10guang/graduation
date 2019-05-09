# -*- coding: utf-8 -*-
# author: Xiguang Liu<g10guang@foxmail.com>
# 2019-05-09 15:14
# test write_api /delete
import requests
import asyncio
import time

consume = []
fail = 0
qps = dict()
fids = []


async def delete(lock, fid):
    global consume, fail, qps
    start = time.time()
    url = "http://127.0.0.1:10003/delete"
    rsp = requests.post(url, data={"uid": 1, "fid": fid})
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
    print("test read_api /get start")
    loop = asyncio.get_event_loop()
    lock = asyncio.Lock()
    tasks = []
    for fid in fids:
        tasks.append(delete(lock, fid))
    start_ts = time.time()
    loop.run_until_complete(asyncio.wait(tasks))
    end_ts = time.time()
    print("test read_api /get end")
    print("consume: {}s".format(end_ts - start_ts))


def load_fids():
    with open("./fid.txt") as f:
        for line in f.readlines():
            fids.append(int(line))


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
    load_fids()
    main()
    statistic()
