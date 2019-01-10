DROP TABLE IF EXISTS `click_info`;
CREATE TABLE `click_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `app_key` varchar(45) DEFAULT NULL COMMENT '点击App Key',
  `device_id` varchar(255) NOT NULL COMMENT '设备唯一标识',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


DROP TABLE IF EXISTS `appkey_sig`;
CREATE TABLE `appkey_sig` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `app_key` varchar(45) DEFAULT NULL COMMENT '点击App Key',
  `public_key` varchar(1024) NOT NULL COMMENT '公钥',
  `private_key` varchar(1024) NOT NULL COMMENT '私钥',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
