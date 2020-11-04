#!/bin/bash

colorStr()
{
    if [[ "$1" == 'red' ]];then
        echo -e '\033[1;31m'$2'\033[0m'
    elif [[ "$1" == 'green' ]];then
        echo -e '\033[1;32m'$2'\033[0m'
    elif [[ "$1" == 'yellow' ]];then
        echo -e '\033[33m'$2'\033[0m'
    elif [[ "$1" == 'blue' ]];then
        echo -e '\033[34m'$2'\033[0m'
    elif [[ "$1" == 'bold' ]];then
        echo -e '\033[1m'$2'\033[0m'
    fi
}

errMsg()
{
    colorStr red "$1"
    msgDingtalk red "${GITLAB_USER_NAME}" "${CI_PROJECT_PATH}" "${CI_PIPELINE_URL}" "$2"
}

rmDockerImageAndContainer()
{
    readonly image_name=$1

    if [[ $# == 2 ]];then
        host="-H $2"
    else
        host=''
    fi

    container_ids=$(docker $host ps -a | grep "${image_name}:" | awk {'print $1'})
    if [[ -n "$container_ids" ]];then
        docker $host stop $container_ids
        docker $host rm -f $container_ids
    fi

    image_ids=$(docker $host image ls | grep "${image_name} " | awk {'print $3'})
    if [[ -n "$image_ids" ]];then
        docker $host rmi ${image_ids[@]}
    fi
}

msgDingtalk()
{
    if [[ "$1" == 'red' ]];then
        color="#FF3030"
    elif [[ "$1" == 'green' ]];then
        color="#43CD80"
    elif [[ "$1" == 'yellow' ]];then
        color="#EEEE00"
    elif [[ "$1" == 'blue' ]];then
        color="#87CEFF"
    elif [[ "$1" == 'bold' ]];then
        color="#8B7765"
    fi

    readonly who=$2
    readonly title=$3
    readonly title_link=$4
    readonly text=$(echo "$5" | sed 's/\\/\\\\/g' | sed 's/\"/\\"/g' | sed 's/&/ and /g' )

    readonly url='https://oapi.dingtalk.com/robot/send?access_token=3d4b816664f06d810adc201b0d6004a117b4b766c75ca1b280622bc8d12c69dc'

    data='{
        "msgtype": "markdown",
        "markdown": {
             "title":"'$title'",
             "text": "### ['$title']('$title_link')\n'$text'",
        },
        "at": {
            "atMobiles": [
              "'$who'"
            ],
            "isAtAll": false
        }
    }'

    curl -X POST -H 'Content-Type: application/json' -d "$data" "$url"
}

msgMattermost()
{
    if [[ "$1" == 'red' ]];then
        color="#FF3030"
    elif [[ "$1" == 'green' ]];then
        color="#43CD80"
    elif [[ "$1" == 'yellow' ]];then
        color="#EEEE00"
    elif [[ "$1" == 'blue' ]];then
        color="#87CEFF"
    elif [[ "$1" == 'bold' ]];then
        color="#8B7765"
    fi

    readonly who=$2
    readonly title=$3
    readonly title_link=$4
    readonly text=$(echo "$5" | sed 's/\\/\\\\/g' | sed 's/\"/\\"/g' | sed 's/&/ and /g' )

    readonly url='http://192.168.1.116:8065/hooks/7cthjo41ttgizysaeu93m5ypjh'

    payload='payload={
        "username":"Docker小女仆",
        "text":"@'$who'",
        "attachments":[
            {
                "color":"'$color'",
                "title":"'$title'",
                "title_link":"'$title_link'/branches/",
                "text":"'"$text"'"
            }
        ]
    }'

    curl -X POST -d "$payload" "$url"
}