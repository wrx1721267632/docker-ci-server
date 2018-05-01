DROP DATABASE IF EXISTS `deploy_web`;
CREATE DATABASE `deploy_web`;
USE `deploy_web`;


DROP TABLE IF EXISTS `account`;
CREATE TABLE `account` (
    `id` int(10) unsigned NOT NULL auto_increment COMMENT '用户ID',
    `name` varchar(30) NOT NULL COMMENT '用户名',
    `password` varchar(30) NOT NULL COMMENT '登录密码',
    PRIMARY KEY(`id`),
    KEY `name` (`name`)
)ENGINE=InnoDB default CHARSET=utf8 COMMENT '用户表';


DROP TABLE IF EXISTS `project`;
CREATE TABLE `project` (
	`id` int(10) unsigned NOT NULL auto_increment COMMENT '工程ID',
	`account_id` int(10) NOT NULL COMMENT '工程创建用户ID',
	`project_name` varchar(50) NOT NULL COMMENT '工程名称',
	`project_describe` varchar(500) default NULL COMMENT '工程描述',
	`git_docker_path` text default NULL COMMENT 'dockerfile文件git路径',
	`create_date` bigint(20) NOT NULL COMMENT '工程创建时间',
	`update_date` bigint(20) NOT NULL COMMENT '工程最后更新时间',
	`is_del` int(5) default '0' COMMENT '工程是否被删除',
	`project_member` varchar(500) default NULL COMMENT '工程成员列表',
	PRIMARY KEY(`id`),
	KEY `account_id` (`account_id`),
	KEY `project_name` (`project_name`),
	KEY `create_date` (`create_date`),
	KEY `update_date` (`update_date`)
)ENGINE=InnoDB default CHARSET=utf8 COMMENT '工程表';

-- DROP TABLE IF EXISTS `project_member`;
-- CREATE TABLE `project_member` (
--   `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '工程成员ID',
--   `project_id` int(10) NOT NULL COMMENT '工程ID',
--   `account_id` int(10) NOT NULL COMMENT '成员ID',
--   PRIMARY KEY (`id`),
--   KEY `project_id` (`project_id`),
--   KEY `account_id` (`account_id`)
-- ) ENGINE=InnoDB default CHARSET=utf8 COMMENT '工程成员表';


DROP TABLE IF EXISTS `construct_record`;
CREATE TABLE `construct_record` (
	`id` int(10) unsigned NOT NULL auto_increment COMMENT '工程构建记录ID',
	`account_id` int(10) NOT NULL COMMENT '用户成员ID',
	`project_id` int(10) NOT NULL COMMENT '工程ID',
	`mirror_id` int(10) default '0' COMMENT '构建产生的镜像ID',
	`construct_start`bigint(20) NOT NULL COMMENT '工程构建的开始时间点',
	`construct_end` bigint(20) NOT NULL COMMENT '工程构建的结束时间点',
	`construct_statu` int(5) default '0' COMMENT '工程构建状态',
	`construct_log` text default NULL COMMENT '工程构建日志',
	PRIMARY KEY(`id`),
	KEY `project_id` (`project_id`),
	KEY `account_id` (`account_id`),
	KEY `mirror_id` (`mirror_id`),
	KEY `construct_start` (`construct_start`)
)ENGINE=InnoDB default CHARSET=utf8 COMMENT '工程构建记录表';


DROP TABLE IF EXISTS `mirror`;
CREATE TABLE `mirror` (
	`id` int(10) unsigned NOT NULL auto_increment COMMENT '镜像ID',
	`mirror_name` varchar(50) NOT NULL COMMENT '镜像名称',
	`mirror_version` varchar(500) default NULL COMMENT '镜像版本',
	`mirror_describe` varchar(500) default NULL COMMENT '镜像描述',
	PRIMARY KEY(`id`),
	KEY `mirror_name` (`mirror_name`),
	KEY `mirror_version` (`mirror_version`)
)ENGINE=InnoDB default CHARSET=utf8 COMMENT '镜像表';


DROP TABLE IF EXISTS `host`;
CREATE TABLE `host` (
	`id` int(10) unsigned NOT NULL auto_increment COMMENT '主机ID',
	`host_name` varchar(50) NOT NULL COMMENT '主机名称',
	`ip` varchar(50) NOT NULL COMMENT '主机IP',
	PRIMARY KEY(`id`),
	KEY `host_name`(`host_name`),
	KEY `ip` (`ip`)
)ENGINE=InnoDB default CHARSET=utf8 COMMENT '主机表';


DROP TABLE IF EXISTS `service`;
CREATE TABLE `service` (
	`id` int(10) unsigned NOT NULL auto_increment COMMENT '服务ID',
	`account_id` int(10) NOT NULL COMMENT '创建用户ID',
	`service_name` varchar(50) NOT NULL COMMENT '服务名称',
	`service_describe` varchar(500) default NULL COMMENT '服务描述',
	`host_list` text default NULL COMMENT '机器列表',
	`mirror_list` varchar(500) default NULL COMMENT '镜像列表',
	`docker_config` text default NULL COMMENT 'docker config文件',
	`create_date` bigint(20) NOT NULL COMMENT '服务创建时间',
	`update_date` bigint(20) NOT NULL COMMENT '服务最后更新时间',
	`service_member` varchar(500) default NULL COMMENT '服务成员列表',
	`is_del` int(5) default '0' COMMENT '服务是否被删除',
	PRIMARY KEY(`id`),
	KEY `account_id` (`account_id`),
	KEY `service_name` (`service_name`),
	KEY `create_date` (`create_date`),
	KEY `update_date` (`update_date`)
)ENGINE=InnoDB default CHARSET=utf8 COMMENT '服务表';


-- DROP TABLE IF EXISTS `service_member`;
-- CREATE TABLE `service_member` (
--   `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '服务成员ID',
--   `service_id` int(10) NOT NULL COMMENT '服务ID',
--   `account_id` int(10) NOT NULL COMMENT '成员ID',
--   PRIMARY KEY (`id`),
--   KEY `service_id` (`service_id`),
--   KEY `account_id` (`account_id`)
-- ) ENGINE=InnoDB default CHARSET=utf8 COMMENT '服务成员表';


DROP TABLE IF EXISTS `deploy`;
CREATE TABLE `deploy` (
	`id` int(10) unsigned NOT NULL auto_increment COMMENT '部署ID',
	`service_id` int(10) NOT NULL COMMENT '所属服务ID',
	`account_id` int(10) NOT NULL COMMENT '部署用户ID',
	`deploy_start` bigint(20) NOT NULL COMMENT '部署的开始时间点',
	`deploy_end` bigint(20) NOT NULL COMMENT '部署的结束时间点',
	`host_list` text default NULL COMMENT '机器列表',
	`mirror_list` varchar(500) default NULL COMMENT '镜像列表',
	`docker_config` text default NULL COMMENT 'docker config文件',
	`deploy_statu` int(5) default '0' COMMENT '部署状态',
	`deploy_log` text default NULL COMMENT '部署日志',
	PRIMARY KEY(`id`),
	KEY `service_id` (`service_id`),
	KEY `account_id` (`account_id`),
	KEY	`deploy_start` (`deploy_start`)
)ENGINE=InnoDB default CHARSET=utf8 COMMENT '部署表';