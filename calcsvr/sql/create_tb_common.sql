CREATE DATABASE IF NOT EXISTS db_zbjs_common DEFAULT CHARSET=utf8 collate utf8_unicode_ci ;
USE db_zbjs_common;

DROP TABLE IF EXISTS tb_lastcalcinfo;
CREATE TABLE `tb_lastcalcinfo` (
	`zb_type` VARCHAR(50) NOT NULL DEFAULT '' COLLATE 'utf8_unicode_ci' COMMENT '指标类型',
	`calc_time` INT(11) NOT NULL DEFAULT '0' COMMENT '指标最后计算的时间',
	`zb_value` TEXT NOT NULL COLLATE 'utf8_unicode_ci',
	`update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (`zb_type`)
)ENGINE=INNODB DEFAULT CHARSET=utf8 COLLATE utf8_unicode_ci COMMENT='指标最后计算时间表';