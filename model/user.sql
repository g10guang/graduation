CREATE TABLE `user` (
  `uid` bigint(20) NOT NULL,
  `status` tinyint(4) NOT NULL DEFAULT '0',
  `extra` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci