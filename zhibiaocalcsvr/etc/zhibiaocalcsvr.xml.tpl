<?xml version="1.0" encoding="UTF-8"?>
<zhibiaocalcsvr>
    <timeout>20</timeout>
    <alarm_url>{{ALARM_URL}}</alarm_url>
    <id>{{SYSTEM_ID}}</id>
    <!--日志相关配置-->
    <log>
        <file>../log/zhibiaocalcsvr.log</file>
        <level>DEBUG</level>
        <size>2048000000</size>
    </log>
    <!--与计算业务相关的配置-->
    <calc>
        <common_start_date>20180101</common_start_date>
        <sighting_start_date>20160101</sighting_start_date>
        <!-- 每隔xx秒打印一次服务状态 -->
        <print_status>30</print_status>
        <!-- 均价，单位是交易日-->
        <avgs>10,20,30,60,120,250</avgs>
        <!-- 最高最低价，单位是交易日 -->
        <lowest_highest>40,60</lowest_highest>
        <black_list>../etc/blacklist.txt</black_list>
    </calc>
    <!--数据库相关-->
    <mysql>
        <master>
            <ip>{{MASTER_MYSQL_IP}}</ip>
            <port>{{MASTER_MYSQL_PORT}}</port>
            <user>{{MASTER_MYSQL_USER}}</user>
            <passwd>{{MASTER_MYSQL_PWD}}</passwd>
            <db>cpgj</db>
            <max_idle>30</max_idle>
            <max_open>100</max_open>
        </master>
    </mysql>
    <redis>
        <hangqing>
            <ip>{{HQ_REDIS_IP}}</ip>
            <port>{{HQ_REDIS_PORT}}</port>
            <passwd>{{HQ_REDIS_PWD}}</passwd>
        </hangqing>
        <zixun>
            <ip>{{ZX_REDIS_IP}}</ip>
            <port>{{ZX_REDIS_PORT}}</port>
            <passwd>{{ZX_REDIS_PWD}}</passwd>
        </zixun>
    </redis>
    <mongo>
        <ip>{{MONGO_IP}}</ip>
        <user>{{MONGO_USER}}</user>
        <passwd>{{MONGO_PWD}}</passwd>
    </mongo>
</zhibiaocalcsvr>