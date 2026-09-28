# GoFrame Template For SingleRepo

Quick Start: 
- https://goframe.org/quick

## 准备

- 下载镜像
docker pull alpine:3.20

- 构建与发送
  `gf build; scp v1.0.0\linux_amd64\main sin@192.168.3.22:~/dcodes/backend/`

- docker
  `docker build -t dcenter-backend .`

  `docker build -t dcenter-backend . && docker run -d -p 8000:8000 --restart unless-stopped --name dcenter-backend dcenter-backend`

## 数据表

### 用户表 `dc_user`

id
user_name
email
password
inviter_id
recently_login_at
created_at
deleted_at
```sql
CREATE TABLE `dc_user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `user_name` varchar(50) NOT NULL COMMENT '用户名',
  `email` varchar(128) NOT NULL COMMENT '邮箱',
  `phone` varchar(11) DEFAULT NULL COMMENT '手机号',
  `password` varchar(60) NOT NULL COMMENT '密码哈希(bcrypt)',
  `inviter_id` bigint unsigned DEFAULT NULL COMMENT '直接邀请人ID，关联dc_user.id，无邀请人为null',
  `user_rank` tinyint NOT NULL DEFAULT 1 COMMENT '用户权限等级'，
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '账号状态：1正常 0禁用 2锁定',
  `last_login_at` datetime(3) DEFAULT NULL COMMENT '最近登录时间',
  `last_login_ip` varchar(45) DEFAULT NULL COMMENT '最近登录IP',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '软删除时间，NULL=未删除',
  PRIMARY KEY (`id`),
  -- 软删除场景联合唯一：已删除账号的用户名/邮箱可以复用
  UNIQUE KEY `uk_username` (`user_name`,`deleted_at`),
  UNIQUE KEY `uk_email` (`email`,`deleted_at`),
  KEY `idx_inviter_id` (`inviter_id`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
```

### 菜单表 `dc_menu`
id
title
icon
sort
parent_id 0 表示为一级菜单
rank int 表示菜单等级
status 是否停用
created_at
delete_at

```sql
CREATE TABLE `dc_menu` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '菜单ID',
  `title` varchar(100) NOT NULL COMMENT '菜单显示名称',
  `icon` varchar(100) DEFAULT NULL COMMENT '菜单图标',
  `sort` int NOT NULL DEFAULT 0 COMMENT '菜单排序，数值越小越靠前',
  `parent_id` bigint unsigned DEFAULT 0 COMMENT '父菜单ID，0代表一级菜单',
  `visible_rank` tinyint NOT NULL DEFAULT 1 COMMENT '菜单可见等级，数值越高权限要求越高；用于判断普通用户是否可查看',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '菜单状态：1启用，0停用',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_visible_rank` (`visible_rank`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='菜单表';
```

### 用户权限表`dc_auth`
id
user_id
menu_id
created_at
deleted_at
creator_id
```sql
CREATE TABLE `dc_auth` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID，关联dc_user.id',
  `menu_id` bigint unsigned NOT NULL COMMENT '菜单ID，关联dc_menu.id',
  `creator_id` bigint unsigned NOT NULL COMMENT '创建人ID，操作授权的管理员',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  -- 同一个用户同一个菜单不能重复授权（软删除联合唯一）
  UNIQUE KEY `uk_user_menu` (`user_id`,`menu_id`,`deleted_at`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户单独授权菜单表';
```