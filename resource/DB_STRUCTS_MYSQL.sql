DROP TABLE IF EXISTS `click_info`;
CREATE TABLE `click_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `app_key` varchar(45) NOT NULL DEFAULT '' COMMENT 'App应用名称',
  `channel_id` varchar(55) DEFAULT NULL COMMENT '渠道标识',
  `device_id` varchar(255) NOT NULL DEFAULT '' COMMENT '设备唯一标识，imei或idfa',
  `click_time` bigint(20) DEFAULT NULL COMMENT '点击时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


DROP TABLE IF EXISTS `channel_sig`;
CREATE TABLE `channel_sig` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `app_key` varchar(45) NOT NULL DEFAULT '' COMMENT 'App应用标识',
  `channel_id` varchar(45) NOT NULL DEFAULT '' COMMENT '渠道标识',
  `public_key` varchar(1024) NOT NULL DEFAULT '' COMMENT '公钥',
  `private_key` varchar(1024) NOT NULL DEFAULT '' COMMENT '私钥',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


DROP TABLE IF EXISTS `sdk_device_report`;
CREATE TABLE `sdk_device_report` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `imei` varchar(45) DEFAULT '' COMMENT 'Android IMEI',
  `idfa` varchar(45) DEFAULT '' COMMENT 'IOS IDFA',
  `app_key` varchar(45) NOT NULL DEFAULT '' COMMENT 'App应用标识',
  `channel_id` varchar(45) DEFAULT '' COMMENT '渠道',
  `resolution` varchar(45) DEFAULT '' COMMENT '设备分辨率',
  `language` varchar(45) DEFAULT '' COMMENT '设备语言',
  `os_type` varchar(45) NOT NULL DEFAULT '' COMMENT '操作系统',
  `os_version` varchar(45) DEFAULT '' COMMENT '操作系统版本',
  `receive_time` varchar(45) NOT NULL DEFAULT '' COMMENT '接收上报时间',
  `source_ip` varchar(45) DEFAULT '' COMMENT '源IP',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
