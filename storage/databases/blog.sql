CREATE TABLE `articles`
(
    `id`         int(11) unsigned    NOT NULL AUTO_INCREMENT,
    `title`      varchar(100)        NOT NULL COMMENT '标题',
    `content`    varchar(2000)       NOT NULL COMMENT '内容',
    `user_id`    int(11) unsigned    NOT NULL COMMENT '用户id',
    `view_num`   int(11)             NOT NULL COMMENT '浏览次数',
    `created_at` int(11) unsigned    NOT NULL COMMENT '创建时间',
    `updated_at` int(11) unsigned    NOT NULL COMMENT '更新时间',
    `status`     tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '状态，0正常，1删除',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE `comments`
(
    `id`         int(10) unsigned    NOT NULL AUTO_INCREMENT,
    `user_id`    int(11)             NOT NULL COMMENT '用户id',
    `article_id` int(10) unsigned    NOT NULL COMMENT '文字id',
    `content`    varchar(500)        NOT NULL COMMENT '评论内容',
    `created_at` int(11) unsigned    NOT NULL COMMENT '创建时间',
    `updated_at` int(11) unsigned    NOT NULL COMMENT '更新时间',
    `status`     tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '状态，0正常，1删除',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE `users`
(
    `id`         int(10) unsigned NOT NULL AUTO_INCREMENT,
    `name`       varchar(50)      NOT NULL COMMENT '用户名',
    `email`      varchar(50)      NOT NULL COMMENT '邮箱',
    `password`   varchar(20)      NOT NULL COMMENT '密码',
    `created_at` int(11) unsigned NOT NULL COMMENT '创建时间',
    `updated_at` int(11) unsigned NOT NULL COMMENT '更新时间',
    `status`     int(11) unsigned NOT NULL COMMENT '状态，0正常，1删除',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

