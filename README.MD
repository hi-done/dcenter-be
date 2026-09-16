# GoFrame Template For SingleRepo

Quick Start: 
- https://goframe.org/quick

## 准备

- 下载镜像
docker pull alpine:3.20

- 构建与发送
  `scp v1.0.0\linux_amd64\main sin@192.168.3.22:~/dcodes/backend/`

- docker
  `docker build -t dcenter-backend .`

  `docker build -t dcenter-backend . && docker run -d -p 8000:8000 --restart unless-stopped --name dcenter-backend dcenter-backend`