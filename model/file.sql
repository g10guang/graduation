CREATE TABLE `file` (
  `fid` bigint(20) NOT NULL,
  `uid` bigint(20) NOT NULL,
  `name` varchar(255) NOT NULL DEFAULT '',
  `size` bigint(20) NOT NULL,
  `md5` char(32) NOT NULL,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `extra` varchar(255) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`fid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci