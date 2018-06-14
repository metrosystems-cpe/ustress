#!/usr/bin/env bash
# Simple script to build golang microservice binary using metro cloud pipeline (2tier)
# Author: ionut.ilie@metrosystems.net


# get https://github.com/Masterminds/glide
curl https://glide.sh/get | sh
if [[ $? -ne 0 ]] ;then 
    printf "Failed to install glide"
    exit 1;
fi

# static variable
GIT_DIR=${GOPATH}/src/git.metrosystems.net/reliability-engineering
APP_DIR=rest-monkey
WRK_DIR=$(pwd)/
printf "GIT DIR: %s\n" ${GIT_DIR} 
printf "APP DIR: %s\n" ${APP_DIR}
printf "WRK DIR: %s\n" ${WRK_DIR}

# prepare git folder structure in GOPATH
mkdir -p ${GIT_DIR}/
if [[ $? -ne 0 ]] ;then 
    printf "Failed to create DIR: %s\n" ${GIT_DIR} 
    exit 1;
  else 
    printf "Succsesfully created DIR: %s \n" ${GIT_DIR}
fi

# symlink mnt/workspace content to git folder structure in GOPATH
ln -s ${WRK_DIR} ${GIT_DIR}/${APP_DIR}
if [[ $? -ne 0 ]] ;then 
    printf "Failed to symlink: %s/ %s \n" ${WRK_DIR} ${GIT_DIR}/${APP_DIR} 
    exit 1;
  else 
    printf "Succsesfully created symlink: %s/ %s \n" ${WRK_DIR} ${GIT_DIR}/${APP_DIR}
fi

# Change Directory
cd ${GIT_DIR}/${APP_DIR}
if [[ $? -ne 0 ]] ;then 
    printf "Failed to cd: %s \n" ${GIT_DIR}/${APP_DIR}
    exit 1;
fi
# remove web-ui folder as it is not necessary for the build 
rm -rf ./web-ui
if [[ $? -ne 0 ]] ;then 
    printf "Failed to rm -rf ./web-ui \n"
fi
# Use glide to rezolve dependencies
glide up
# make linux"
make linux




