CREATE DATABASE IF NOT EXISTS db_zbjs_dcxg DEFAULT CHARSET=utf8 collate utf8_unicode_ci ;
USE db_zbjs_dcxg;


#黄蓝区间--------------------------------------
DROP TABLE IF EXISTS tb_hlqj_dot;
CREATE TABLE tb_hlqj_dot(
    stockcode VARCHAR(50) NOT NULL COMMENT '股票代码',
    sdata INT(11) NOT NULL COMMENT '序列化日期 eg:20180917',
    upserttime	INT(11) NOT NULL COMMENT '客户端发送的检索主键同mongo',
    rsttype INT(11) NOT NULL COMMENT '买卖标志 0买, 1卖',
    instime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '插入时间戳',
    uptime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间戳',
    PRIMARY KEY(stockcode, sdata,upserttime,rsttype)
)ENGINE=innodb DEFAULT CHARSET=utf8 collate utf8_unicode_ci COMMENT='黄蓝区间结果统计表';
#DESC tb_hlqj_dot;

DROP TABLE IF EXISTS tb_hlqj_column;
CREATE TABLE tb_hlqj_column(
    stockcode VARCHAR(50) NOT NULL COMMENT '股票代码',
    sdata INT(11) NOT NULL COMMENT '序列化日期 eg:20180917',
    upserttime	INT(11) NOT NULL COMMENT '客户端发送的检索主键同mongo',
    rsttype INT(11) NOT NULL COMMENT '黄蓝标志 0黄, 1蓝',
    high INT(11) NOT NULL COMMENT '柱形图上坐标',
    low INT(11) NOT NULL COMMENT '柱形图下坐标',
    instime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '插入时间戳',
    uptime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间戳',
    PRIMARY KEY(stockcode, sdata,upserttime,rsttype)
)ENGINE=innodb DEFAULT CHARSET=utf8 collate utf8_unicode_ci COMMENT='黄蓝区间结果统计表';
#DESC tb_hlqj_column;


#操盘提醒-------------------------------------
DROP TABLE IF EXISTS tb_cptx_dot;
CREATE TABLE tb_cptx_dot(
    stockcode VARCHAR(50) NOT NULL COMMENT '股票代码',
    sdata INT(11) NOT NULL COMMENT '序列化日期 eg:20180917',
    upserttime	INT(11) NOT NULL COMMENT '客户端发送的检索主键同mongo',
    rsttype INT(11) NOT NULL COMMENT '进场标志 0波段进场, 1反弹进场，2超跌进场',
    instime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '插入时间戳',
    uptime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间戳',
    PRIMARY KEY(stockcode, sdata,upserttime,rsttype)
)ENGINE=innodb DEFAULT CHARSET=utf8 collate utf8_unicode_ci COMMENT='操盘提醒结果统计表';
#DESC tb_cptx_dot;

DROP TABLE IF EXISTS tb_cptx_column;
CREATE TABLE tb_cptx_column(
    stockcode VARCHAR(50) NOT NULL COMMENT '股票代码',
    sdata INT(11) NOT NULL COMMENT '序列化日期 eg:20180917',
    upserttime	INT(11) NOT NULL COMMENT '客户端发送的检索主键同mongo',
    rsttype INT(11) NOT NULL COMMENT '柱形颜色 0红, 1绿',
    high INT(11) NOT NULL COMMENT '柱形图上坐标',
    low INT(11) NOT NULL COMMENT '柱形图下坐标',
    instime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '插入时间戳',
    uptime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间戳',
    PRIMARY KEY(stockcode, sdata,upserttime,rsttype)
)ENGINE=innodb DEFAULT CHARSET=utf8 collate utf8_unicode_ci COMMENT='操盘提醒结果统计表';
#DESC tb_cptx_column;


