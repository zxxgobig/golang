CREATE DATABASE IF NOT EXISTS cpgj DEFAULT CHARSET=utf8 collate utf8_unicode_ci;
use cpgj;


# 狙击镜 买卖点
DROP TABLE IF EXISTS tb_sighting_v2;
CREATE TABLE tb_sighting_v2(
    stockcode   VARCHAR(50)   NOT NULL COMMENT '股票代码',
    tradday     INT(11)       NOT NULL COMMENT '日期：eg:20190524',
    stype       INT(11)       NOT NULL COMMENT '建仓1，加仓2，建仓3，清仓4',
    open        INT(11)       NOT NULL COMMENT '开盘价N(3)',
    close       INT(11)       NOT NULL COMMENT '收盘价N(3)',
    reason      VARCHAR(100)  NOT NULL COMMENT '触发买卖点的原因',
    instime     TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '插入时间戳',
    PRIMARY KEY(stockcode,tradday)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 collate utf8_unicode_ci COMMENT='股票均价统计表';


# 公共表
CREATE TABLE tb_common_v2(
    stype           VARCHAR(50) NOT NULL COMMENT '指明压力支撑点、狙击镜买卖点 还是其它数据',
    start_date      INT(11)     NOT NULL COMMENT '开始日期, 如狙击镜买卖点开始计算的日期',
    last_calc_date  INT(11)     NOT NULL DEFAULT 0 COMMENT '最后计算日期，供服务启动时读取',
    PRIMARY KEY(stype)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 collate utf8_unicode_ci COMMENT='记录各种统计数据的起始计算时间、最后计算时间';

# 狙击镜 买卖点回测结果表
CREATE TABLE sighting_back_test_v2(
    stockcode   VARCHAR(50) NOT NULL COMMENT '股票代码',
    buyday      INT(11)     NOT NULL COMMENT '建仓日期：eg:20190524',
    clearday    INT(11)     NOT NULL COMMENT '清仓日期：eg:20190524',
    holddays    INT(11)     NOT NULL COMMENT '持有天数：[建仓--清仓]之间的交易日个数',
    bclose      FLOAT(4)    NOT NULL COMMENT '建仓日收盘价',
    cclose      FLOAT(4)    NOT NULL COMMENT '清仓日收盘价',
    srate       FLOAT(4)    NOT NULL COMMENT '建仓日 涨跌幅',
    brate       FLOAT(4)    NOT NULL COMMENT '加仓日 涨跌幅',
    instime     TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '插入时间戳',
    PRIMARY KEY(stockcode,buyday)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 collate utf8_unicode_ci COMMENT='狙击镜 买卖点回测结果表';


# 狙击战法 最大收益率信息表
CREATE TABLE `sighting_profit_max` (
  `stockcode` varchar(50) COLLATE utf8_unicode_ci NOT NULL COMMENT '股票代码',
  `rate` double(10,4) NOT NULL COMMENT '收益率',
  `initdate` int(11) NOT NULL COMMENT '建仓时间',
  `cleardate` int(11) NOT NULL COMMENT '清仓时间',
  `upserttime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`stockcode`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='记录狙击战法计算的各个股票的 最大收益周期信息';


# 狙击战法 所有收益率信息表
CREATE TABLE `sighting_profit` (
  `stockcode` varchar(50) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL COMMENT '股票代码',
  `rate` double(10,4) NOT NULL COMMENT '收益率',
  `initdate` date NOT NULL COMMENT '建仓时间',
  `cleardate` date NOT NULL COMMENT '清仓时间',
  `upserttime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`stockcode`,`initdate`,`cleardate`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='记录狙击战法计算的各个股票的 所有收益周期信息';
