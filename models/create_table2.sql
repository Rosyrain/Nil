DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
                      `id` int NOT NULL AUTO_INCREMENT,
                      `user_id` bigint(20) NOT NULL,
                      `username` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
                      `password` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
                      `email` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
                      `gender` tinyint(4) NOT NULL DEFAULT '0',
                      `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                      `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                      `status` tinyint(4) NOT NULL DEFAULT '0',
                      PRIMARY KEY (`id`),
                      UNIQUE KEY `idx_username` (`username`) USING BTREE,
                      UNIQUE KEY `idx_user_id` (`user_id`) USING BTREE
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

DROP TABLE IF EXISTS `chunk`;
CREATE TABLE `chunk` (
                         `id` int NOT NULL AUTO_INCREMENT,
                         `chunk_id` bigint(20) unsigned NOT NULL,
                         `chunk_name` varchar(128) COLLATE utf8mb4_general_ci NOT NULL,

                         `introduction` varchar(256) COLLATE utf8mb4_general_ci NOT NULL,
                         `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                         PRIMARY KEY (`id`),
                         UNIQUE KEY `idx_chunk_id` (`chunk_id`),
                         UNIQUE KEY `idx_chunk_name` (`chunk_name`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

DROP TABLE IF EXISTS `post`;
CREATE TABLE `post` (
                        `id` bigint(20) NOT NULL AUTO_INCREMENT,
                        `post_id` bigint(20) NOT NULL COMMENT '帖子id',
                        `title` varchar(128) COLLATE utf8mb4_general_ci NOT NULL COMMENT '标题',
                        `content` varchar(8192) COLLATE utf8mb4_general_ci NOT NULL COMMENT '内容',
                        `author_id` bigint(20) NOT NULL COMMENT '作者的用户id',
                        `chunk_id` bigint(20) NOT NULL COMMENT '所属板块',
                        `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '帖子状态',
                        `vote_number` bigint(20) NOT NULL DEFAULT '0' COMMENT '投票情况',
                        `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                        `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                        PRIMARY KEY (`id`),
                        UNIQUE KEY `idx_post_id` (`post_id`),
                        KEY `idx_author_id` (`author_id`),
                        KEY `idx_chunk_id` (`chunk_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

DROP TABLE IF EXISTS `comment`;

CREATE TABLE `comment` (
                         `id` bigint(20) NOT NULL AUTO_INCREMENT,
                         `comment_id` bigint(20) not null,
                         `author_id` bigint(20) NOT NULL COMMENT '作者的用户id',
                         `post_id` bigint(20) NOT NULL COMMENT '帖子id',
                         `content` VARCHAR(255) COLLATE utf8mb4_general_ci NOT NULL ,
                         `vote_num` bigint(20) default 0,
                         `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                       PRIMARY KEY (`id`),
                       UNIQUE KEY `idx_comment_id`(`comment_id`),
                        KEY `idx_author_id`(`author_id`),
                        KEY `idx_post_id`(`post_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


DROP TABLE IF EXISTS `superuser`;
CREATE TABLE `superuser` (
                        `id` int NOT NULL AUTO_INCREMENT,
                        `user_id` bigint(20) NOT NULL,
                        `chunk_id` bigint(20)  NOT NULL,
                        `username` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
                        `password` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
                        `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                        `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                        PRIMARY KEY (`id`),
                        UNIQUE KEY `idx_username` (`username`) USING BTREE,
                        UNIQUE KEY `idx_user_id` (`user_id`) USING BTREE
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

insert into `superuser`(username,password,chunk_id) values ("rosyrain","jx20031002",0);

show tables;
