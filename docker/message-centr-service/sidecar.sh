#!/bin/bash

# 配置sidecar 日志目录
if [[ ! -d "/opt/log/$APP_NAME/sidecar" ]];then
    mkdir -p /opt/log/$APP_NAME/sidecar
    chown -R nobody:nobody /opt/log/$APP_NAME/sidecar
    chmod -R 777 /opt/log/$APP_NAME/sidecar
fi
if [[ ! -L "/opt/app/sidecar/logs" ]];then
    rm -rf /opt/app/sidecar/logs
    ln -s /opt/log/$APP_NAME/sidecar /opt/app/sidecar/logs
fi

if [[ -z "$EUREKA_URL" ]];then
    EUREKA_URL="http://eureka-center/eureka/"
fi

# 输出配置文件
cat > /opt/app/sidecar/discovery.conf << EOF
# 被代理的服务IP，与sidecar同host时可以用 127.0.0.1
proxy.host=127.0.0.1
# 被代理的服务端口
proxy.port=8080
# sidecar服务端口
server.port=8000
# 被代理的服务名
spring.application.name=${SPI_NAME}
# 注册中心服务分组
eureka.instance.app-group-name=PHP
# 注册中心地址
eureka.client.service-url.defaultZone=${EUREKA_URL}
#访问日志控制，true-全部输出、false-只输出接口名字不为空的日志
access-log-coercive=true
EOF

# 启动sidecar
${JAVA_HOME}/bin/java -server -Xms${JVM_XMS} -Xmx${JVM_XMX} -Xmn${JVM_XMN} -XX:+UseNUMA -XX:+UseParallelGC -jar /opt/app/sidecar/sidecar-serve.jar --boot.active-conf=file:/opt/app/sidecar/discovery.conf 2>&1 &
