/* 
    角色基础信息表
*/
CREATE TABLE `role_base` (
  `role_id` bigint(19) NOT NULL COMMENT '角色id',
  `user_id` bigint(19) NOT NULL COMMENT '账号id',
  `role_name` varchar(255) NOT NULL DEFAULT '' COMMENT '角色名字',
  `create_time` bigint(19) NOT NULL DEFAULT '0' COMMENT '创建时间',
  `level` int(11) NOT NULL DEFAULT '0' COMMENT '等级',
  `login_time` bigint(19) NOT NULL COMMENT '最近登录时间',
  `offline_time` bigint(19) NOT NULL DEFAULT '0' COMMENT '最近离线时间',
  `money` bigint(19) NOT NULL DEFAULT '0' COMMENT '资金',
  PRIMARY KEY (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;