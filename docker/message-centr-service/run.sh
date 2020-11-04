#!/bin/bash

# 配置host
cat /opt/docker-config/hosts/$ENV >> /etc/hosts

# 配置php
\cp -rf /opt/docker-config/php/php.ini $PHP_INI_DIR/conf.d/
\rm -rf $PHP_INI_DIR/../php-fpm.d
\cp -rf /opt/docker-config/php/php-fpm.d $PHP_INI_DIR/../

# 配置.env
if [[ ! -f "/opt/wwwroot/conf/.env" ]];then
    exit 1
fi
\cp -f /opt/wwwroot/conf/.env $APP_PATH
\cp -f /opt/wwwroot/conf/.env $APP_PATH/.env.$ENV

# 配置时区
\cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
echo 'Asia/Shanghai' >/etc/timezone

# 配置进程监控环境变量
echo 'source /opt/docker-config/export_env_var' >> /opt/docker-config/check_process.sh

# 生成环境变量声明文件
cat > /opt/docker-config/export_env_var << EOF
export ENV=$ENV
export APP_NAME=$APP_NAME
export APP_PATH=$APP_PATH
export FRAMEWORK=$FRAMEWORK
EOF

# 配置日志文件夹
if [[ ! -d "/opt/log/$APP_NAME" ]];then
    mkdir -p /opt/log/$APP_NAME
    chown -R nobody:nobody /opt/log/$APP_NAME
fi
chmod -R 777 /opt/log/$APP_NAME

if [[ -f "${APP_PATH}supervisor.conf" && ! -d "/opt/log/$APP_NAME/supervisor/run" ]];then
    mkdir -p /opt/log/$APP_NAME/supervisor/run
    chown -R nobody:nobody /opt/log/$APP_NAME/supervisor
    chmod -R 777 /opt/log/$APP_NAME/supervisor
fi

# 启动 laravel 缓存
if [[ "$FRAMEWORK" == "laravel" ]];then
    chown -R nobody:nobody ${APP_PATH}storage
    chmod -R 777 ${APP_PATH}storage
    $PHP_CLI ${APP_PATH}artisan config:cache
    $PHP_CLI ${APP_PATH}artisan route:cache
fi

check_process=0

# 启动php-fpm
if [[ "$PROCESS_FPM_WEB" == "1" || "$PROCESS_FPM_RPC" == "1" ]];then
    if [[ "$PROCESS_FPM_WEB" == "1" ]];then
        mv $PHP_INI_DIR/../php-fpm.d/web.conf.bak $PHP_INI_DIR/../php-fpm.d/web.conf
    fi

    if [[ "$PROCESS_FPM_RPC" == "1" ]];then
        mv $PHP_INI_DIR/../php-fpm.d/rpc.conf.bak $PHP_INI_DIR/../php-fpm.d/rpc.conf
    fi

    $PHP_FPM -D

    # 写入进程监控脚本
    cat >> /opt/docker-config/check_process.sh << EOF
process=\$(ps aux | grep "php-fpm: master" | grep -v "grep")
if [[ -z "\$process" ]];then
    $PHP_FPM -D
fi
EOF

    check_process=1
fi

# 启动队列
if [[ "$PROCESS_QUEUE" == "1" ]];then
    #启动supervisor
    if [[ -f "${APP_PATH}supervisor.conf" ]];then
        /usr/bin/supervisord -c ${APP_PATH}supervisor.conf
        # 写入进程监控脚本
        cat >> /opt/docker-config/check_process.sh << EOF
process=\$(ps aux | grep "supervisord" | grep -v "grep")
if [[ -z "\$process" ]];then
    /usr/bin/supervisord -c ${APP_PATH}supervisor.conf
fi
EOF

        check_process=1
    fi

    if [[ "$FRAMEWORK" == "ci" && -f "${APP_PATH}swoole.php" ]];then
        $PHP_CLI ${APP_PATH}swoole.php -s queue -d
        # 写入进程监控脚本
        cat >> /opt/docker-config/check_process.sh << EOF
process=\$(ps aux | grep "swoole_queue:master" | grep -v "grep")
if [[ -z "\$process" ]];then
    $PHP_CLI ${APP_PATH}swoole.php -s queue -d
fi
EOF

        check_process=1
    fi
fi

# 加入项目计划任务
if [[ "$PROCESS_CRONTAB" == "1" && -f "${APP_PATH}.crontab" ]];then
    crontab -u nobody ${APP_PATH}.crontab
fi

# 加入进程监控计划任务
if [[ "$check_process" == "1" ]];then
    crontab /opt/docker-config/.crontab
fi

# 启动计划任务
if [[ "$PROCESS_CRONTAB" == "1" || "$check_process" == "1" ]];then
    $CRONTAB_START
fi

# 启动sidecar
if [[ "$PROCESS_SIDECAR" == "1" ]];then
    /opt/docker-config/sidecar.sh
fi

# 再次修改日志目录权限，防止刚才的命令产生出权限不对的日志
chmod -R 777 /opt/log/$APP_NAME

tail -f /etc/hosts