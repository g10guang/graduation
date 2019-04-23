存储层的抽象，为了后续从本地文件系统存储切换到 HDFS 存储

# hdfs 笔记

```bash
docker run -it -p 2122:2122 -p 8030:8030 -p 8031:8031 -p 8032:8032 -p 8083:8083 -p 8040:8040 -p 8088:8088 -p 19888:19888 -p 49707:49707 -p 50010:50010 -p 50020:50020 -p 50070:50070 -p 50075:50075 -p 50090:50090  sequenceiq/hadoop-docker:2.7.1 /etc/bootstrap.sh -bash
```

节点：

- namenode：9000
- datanode：50075
- webhdfs: 50070

web 界面：

```
http://10.8.118.15:50070/explorer.html#/user

http://10.8.118.15:8088/cluster/nodes
```

[webhdfs 文档](http://hadoop.apache.org/docs/stable/hadoop-project-dist/hadoop-hdfs/WebHDFS.html)

webhdfs

## mkdir

```bash
➜  tmp git:(b2) ✗ curl -i -X PUT "http://10.8.118.15:50070/webhdfs/v1/oss/image?op=MKDIRS&user.name=root"
HTTP/1.1 200 OK
Cache-Control: no-cache
Expires: Sun, 31 Mar 2019 14:05:50 GMT
Date: Sun, 31 Mar 2019 14:05:50 GMT
Pragma: no-cache
Expires: Sun, 31 Mar 2019 14:05:50 GMT
Date: Sun, 31 Mar 2019 14:05:50 GMT
Pragma: no-cache
Content-Type: application/json
Set-Cookie: hadoop.auth="u=root&p=root&t=simple&e=1554077150977&s=H714pJSxPaAlKec7tcgj7intmv8="; Path=/; Expires=Mon, 01-Apr-2019 00:05:50 GMT; HttpOnly
Transfer-Encoding: chunked
Server: Jetty(6.1.26)

{"boolean":true}%
```

## stat

```bash
➜  tmp git:(b2) ✗ curl -i  "http://10.8.118.15:50070/webhdfs/v1/user?op=GETFILESTATUS"
HTTP/1.1 200 OK
Cache-Control: no-cache
Expires: Sun, 31 Mar 2019 10:49:00 GMT
Date: Sun, 31 Mar 2019 10:49:00 GMT
Pragma: no-cache
Expires: Sun, 31 Mar 2019 10:49:00 GMT
Date: Sun, 31 Mar 2019 10:49:00 GMT
Pragma: no-cache
Content-Type: application/json
Transfer-Encoding: chunked
Server: Jetty(6.1.26)

{"FileStatus":{"accessTime":0,"blockSize":0,"childrenNum":1,"fileId":16386,"group":"supergroup","length":0,"modificationTime":1450036470505,"owner":"root","pathSuffix":"","permission":"755","replication":0,"storagePolicy":0,"type":"DIRECTORY"}}%
```

## write

首先访问 webhdfs，然后被 redirect 到对应的 data node

```bash
// 初始请求
curl -i -X PUT -T a "http://10.8.118.15:50070/webhdfs/v1/tmp/b?op=CREATE&user.name=root&overwrite=true&noredirect=true"

HTTP/1.1 307 TEMPORARY_REDIRECT
Cache-Control: no-cache
Expires: Sun, 31 Mar 2019 13:41:24 GMT
Date: Sun, 31 Mar 2019 13:41:24 GMT
Pragma: no-cache
Expires: Sun, 31 Mar 2019 13:41:24 GMT
Date: Sun, 31 Mar 2019 13:41:24 GMT
Pragma: no-cache
Content-Type: application/octet-stream
Set-Cookie: hadoop.auth="u=root&p=root&t=simple&e=1554075684267&s=Ra+RmicRTuVM1cyB7/9hPpQgVWc="; Path=/; Expires=Sun, 31-Mar-2019 23:41:24 GMT; HttpOnly
Location: http://7afa289e0e6e:50075/webhdfs/v1/tmp/b?op=CREATE&user.name=root&namenoderpcaddress=7afa289e0e6e:9000&overwrite=false
Content-Length: 0
Server: Jetty(6.1.26)

// redirect 后的请求
curl -i -X PUT -T a "http://10.8.118.15:50075/webhdfs/v1/tmp/a?op=CREATE&user.name=root&namenoderpcaddress=7afa289e0e6e:9000&blocksize=1048576&buffersize=64&overwrite=true&permission=755&replication=1"
```

## read

按照 http 响应 Location 来重定向到 DataNode 的 url

```bash
curl -i -L "http://10.8.118.15:50070/webhdfs/v1/tmp/a?op=OPEN"
```

## delete

```bash
curl -i -X DELETE "http://10.8.118.15:50070/webhdfs/v1/tmp/a?op=DELETE&user.name=root"
```

## FAQ

对于 webhdfs Redirect Location 中使用 Hostname 而不是 ip:port 的情况，应该添加 `/etc/hosts`