CREATE DATABASE IF NOT EXISTS db_zbjs_hjjj DEFAULT CHARSET=utf8 collate utf8_unicode_ci ;
USE db_zbjs_hjjj;


#筹码分布
DROP TABLE IF EXISTS tb_cmfb;
CREATE TABLE tb_cmfb(
    stockcode VARCHAR(50) NOT NULL COMMENT '股票代码',
    sdata INT(11) NOT NULL COMMENT '序列化日期 eg:20180917',
    upserttime	INT(11) NOT NULL COMMENT '客户端发送的检索主键同mongo',
    rsttype INT(11) NOT NULL COMMENT '分集标志 0集, 1分',
    instime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '插入时间戳',
    uptime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间戳',
    PRIMARY KEY(stockcode, sdata,upserttime,rsttype)
)ENGINE=innodb DEFAULT CHARSET=utf8 collate utf8_unicode_ci COMMENT='筹码分布结果统计表';
#DESC tb_cmfb;

#见底出击
DROP TABLE IF EXISTS tb_jdcj;
CREATE TABLE tb_jdcj(
    stockcode VARCHAR(50) NOT NULL COMMENT '股票代码',
    sdata INT(11) NOT NULL COMMENT '序列化日期 eg:20180917',
    upserttime	INT(11) NOT NULL COMMENT '客户端发送的检索主键同mongo',
    rsttype INT(11) NOT NULL COMMENT '趋势类别1：短期,2中期,3长期',
    instime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '插入时间戳',
    uptime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间戳',
    PRIMARY KEY(stockcode, sdata,upserttime,rsttype)
)ENGINE=innodb DEFAULT CHARSET=utf8 collate utf8_unicode_ci COMMENT='见底出击结果统计表';
#DESC tb_jdcj;
