#!/bin/bash
source /opt/docker-config/export_env_var
process=$(ps aux | grep "supervisord" | grep -v "grep")
if [[ -z "$process" ]];then
    /usr/bin/supervisord -c /opt/wwwroot/jsb-message-center-service/supervisor.conf
fi
source /opt/docker-config/export_env_var
process=$(ps aux | grep "supervisord" | grep -v "grep")
if [[ -z "$process" ]];then
    /usr/bin/supervisord -c /opt/wwwroot/jsb-message-center-service/supervisor.conf
fi