CREATE TABLE `users` (
	`ID` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
	`Email` varchar(128) NOT NULL COMMENT 'Email',
	`Password` varchar(128) NOT NULL COMMENT '密码',
	`PasswordSalt` varchar(128) DEFAULT NULL COMMENT '密码SALT',
	`Nickname` varchar(32) NOT NULL COMMENT '昵称',
	`Avatar` varchar (128) DEFAULT NULL COMMENT '头像URL',
	`IsApproved` tinyint(1) NOT NULL DEFAULT '1' COMMENT '审核状态',
	`IsLocked` tinyint(1) NOT NULL DEFAULT '0' COMMENT '锁定状态',
	`IsDeleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT '删除状态',
	`LastLoginDate` datetime DEFAULT NULL COMMENT '最后登录时间',
	`LastPasswordChangedDate` datetime DEFAULT NULL COMMENT '最后密码修改时间',
	`LastLockoutDate` datetime DEFAULT NULL COMMENT '最后锁定时间',
	`Remark` varchar(256) DEFAULT NULL COMMENT '备注',
	`CreatedOn` datetime NOT NULL COMMENT '创建时间',
	`UpdatedOn` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (`ID`),
	UNIQUE KEY `idx_users_email` (`Email`)
) ENGINE=InnoDB AUTO_INCREMENT=989187 DEFAULT CHARSET=utf8 COMMENT='用户表';