DROP TABLE IF EXISTS `click_info`;
CREATE TABLE `click_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `app_key` varchar(45) NOT NULL DEFAULT '' COMMENT 'App应用名称',
  `channel_id` varchar(55) DEFAULT NULL COMMENT '渠道标识',
  `device_id` varchar(255) NOT NULL DEFAULT '' COMMENT '设备唯一标识，imei或idfa',
  `click_time` bigint(20) DEFAULT NULL COMMENT '点击时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=41 DEFAULT CHARSET=utf8;


DROP TABLE IF EXISTS `channel_sig`;
CREATE TABLE `channel_sig` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `app_key` varchar(45) NOT NULL DEFAULT '' COMMENT 'App应用标识',
  `channel_id` varchar(45) NOT NULL DEFAULT '' COMMENT '渠道标识',
  `public_key` varchar(1024) NOT NULL DEFAULT '' COMMENT '公钥',
  `private_key` varchar(1024) NOT NULL DEFAULT '' COMMENT '私钥',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8;
