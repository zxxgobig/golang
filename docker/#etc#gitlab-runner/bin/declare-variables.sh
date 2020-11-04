if [[ "$CI_PROJECT_NAMESPACE" =~ ^bpc\/be ]];then
    readonly PROJECT_GROUP="bpc"
    readonly DOCKER_CONTAINER_NAME=${CI_PROJECT_PATH_SLUG#*-}
elif [[ "$CI_PROJECT_NAMESPACE" =~ ^tob-ai-platform ]]; then
    readonly PROJECT_GROUP="dsp"
    readonly DOCKER_CONTAINER_NAME=${CI_PROJECT_NAME}
    #readonly DOCKER_CONTAINER_NAME=${CI_PROJECT_PATH_SLUG:16}
else
    errMsg "${CI_PROJECT_NAMESPACE} 不是一个后端项目"
    exit 1
fi
readonly GIT_DOMAIN="https://oauth2:hqPTwde-n9YcdFsm2VVC@gitlab.ifchange.com"
readonly DOCKER_HUB=hub.ifchange.com/
#readonly PROJECT_NAME=${CI_PROJECT_PATH#*/}
readonly PROJECT_NAME=${CI_PROJECT_NAME}

readonly DOCKER_IMAGE_NAME="${DOCKER_HUB}${PROJECT_GROUP}/${DOCKER_CONTAINER_NAME}"
readonly DOCKER_IMAGE_NAME_AND_TAG="${DOCKER_IMAGE_NAME}:${CI_COMMIT_REF_SLUG}"
readonly PROJECT_REAL_PATH=/opt/wwwroot/${PROJECT_NAME}/