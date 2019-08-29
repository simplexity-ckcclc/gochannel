DROP TABLE IF EXISTS `click_info`;
CREATE TABLE `click_info` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `app_key` varchar(45) NOT NULL DEFAULT '' COMMENT 'App应用名称',
  `channel_id` varchar(55) NOT NULL COMMENT '渠道标识',
  `device_id` varchar(255) NOT NULL DEFAULT '' COMMENT '设备唯一标识，imei或idfa',
  `os_type` varchar(255) NOT NULL DEFAULT '' COMMENT '系统，ios/android',
  `click_time` bigint(20) NOT NULL COMMENT '点击时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


DROP TABLE IF EXISTS `app_channel`;
CREATE TABLE `app_channel` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `app_key` varchar(45) NOT NULL DEFAULT '' COMMENT 'App应用标识',
  `channel_id` varchar(45) NOT NULL DEFAULT '' COMMENT '渠道标识',
  `public_key` varchar(1024) NOT NULL DEFAULT '' COMMENT '公钥',
  `private_key` varchar(1024) NOT NULL DEFAULT '' COMMENT '私钥',
  UNIQUE KEY `uniq_app_key_&_channel_id` (`app_key`, `channel_id`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


DROP TABLE IF EXISTS `sdk_device_report`;
CREATE TABLE `sdk_device_report` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `imei` varchar(45) DEFAULT '' COMMENT 'Android IMEI',
  `idfa` varchar(45) DEFAULT '' COMMENT 'IOS IDFA',
  `app_key` varchar(45) NOT NULL DEFAULT '' COMMENT 'App应用标识',
  `channel_id` varchar(45) DEFAULT '' COMMENT '渠道',
  `resolution` varchar(45) DEFAULT '' COMMENT '设备分辨率',
  `language` varchar(45) DEFAULT '' COMMENT '设备语言',
  `os_type` varchar(45) NOT NULL DEFAULT '' COMMENT '操作系统',
  `os_version` varchar(45) DEFAULT '' COMMENT '操作系统版本',
  `activate_time` varchar(45) NOT NULL DEFAULT '' COMMENT '接收激活时间',
  `source_ip` varchar(45) DEFAULT '' COMMENT '源IP',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `app_process_info`;
CREATE TABLE `app_process_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `app_key` varchar(45) NOT NULL DEFAULT '' COMMENT 'App应用标识',
  `process_time` bigint(20) NOT NULL DEFAULT 0 COMMENT '已完成处理的最新时间戳',
  UNIQUE KEY `uniq_app_key` (`app_key`),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `callback_info`;
CREATE TABLE `callback_info` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `app_key` varchar(45) NOT NULL DEFAULT '' COMMENT 'App应用名称',
  `channel_id` varchar(55) NOT NULL COMMENT '渠道标识',
  `device_id` varchar(255) NOT NULL DEFAULT '' COMMENT '设备唯一标识，imei或idfa',
  `os_type` varchar(255) NOT NULL DEFAULT '' COMMENT '系统，ios/android',
  `click_time` bigint(20) NOT NULL COMMENT '点击时间',
  `activate_time` bigint(20) NOT NULL COMMENT '激活时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
