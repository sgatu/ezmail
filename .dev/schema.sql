CREATE TABLE `domain` (
  `id` varchar(64) COLLATE utf8mb4_bin NOT NULL,
  `domain_name` varchar(128) COLLATE utf8mb4_bin NOT NULL,
  `validated` tinyint NOT NULL DEFAULT '0',
  `records` varchar(1024) COLLATE utf8mb4_bin DEFAULT NULL,
  `created` datetime NOT NULL,
  `region` varchar(32) COLLATE utf8mb4_bin NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;



CREATE TABLE `email` (
  `id` bigint NOT NULL,
  `domain_id` bigint NOT NULL,
  `template_id` bigint DEFAULT NULL,
  `from` varchar(256) COLLATE utf8mb4_bin NOT NULL,
  `reply_to` varchar(256) COLLATE utf8mb4_bin DEFAULT NULL,
  `to` varchar(2048) COLLATE utf8mb4_bin NOT NULL,
  `bcc` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `created` datetime NOT NULL,
  `processed_date` datetime DEFAULT NULL,
  `context` varchar(2048) COLLATE utf8mb4_bin DEFAULT NULL,
  `expected_sent` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;



CREATE TABLE `template` (
  `id` bigint NOT NULL,
  `subject` varchar(256) COLLATE utf8mb4_bin NOT NULL,
  `html` text COLLATE utf8mb4_bin NOT NULL,
  `text` text COLLATE utf8mb4_bin NOT NULL,
  `created` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
