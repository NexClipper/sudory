USE sudory;

CREATE TABLE `service_command_v1` (
	`id` INT(11) NOT NULL AUTO_INCREMENT COMMENT '아이디',
	`uuid` CHAR(32) NOT NULL COMMENT 'uuid' COLLATE 'utf8mb4_general_ci',
	`created_by` VARCHAR(255) NULL DEFAULT NULL COMMENT '생성자' COLLATE 'utf8mb4_general_ci',
	`created_at` DATETIME NULL DEFAULT NULL COMMENT '생성시간',
	`updated_by` VARCHAR(255) NULL DEFAULT NULL COMMENT '수정자' COLLATE 'utf8mb4_general_ci',
	`updated_at` DATETIME NULL DEFAULT NULL COMMENT '수정시간',
	`deleted_by` VARCHAR(255) NULL DEFAULT NULL COLLATE 'utf8mb4_general_ci',
	`deleted_at` DATETIME NULL DEFAULT NULL COMMENT '삭제시간',
	`name` VARCHAR(255) NOT NULL COLLATE 'utf8mb4_general_ci',
	`summary` VARCHAR(255) NULL DEFAULT NULL COLLATE 'utf8mb4_general_ci',
	`api_version` VARCHAR(255) NOT NULL COLLATE 'utf8mb4_general_ci',
	`template_uuid` CHAR(32) NOT NULL DEFAULT '' COMMENT 'template 테이블 uuid' COLLATE 'utf8mb4_general_ci',
	`method` VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'method ex) kube.node.get, kube.pods.list' COLLATE 'utf8mb4_general_ci',
	`args` VARCHAR(255) NULL DEFAULT '' COMMENT 'kubernetes method arguments' COLLATE 'utf8mb4_general_ci',
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE INDEX `uuid` (`uuid`) USING BTREE
)
COMMENT='commandtable'
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
