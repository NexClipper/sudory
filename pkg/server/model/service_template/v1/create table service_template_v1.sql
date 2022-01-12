USE sudory;

CREATE TABLE `service_template_v1` (
	`id` INT(11) NOT NULL AUTO_INCREMENT COMMENT '아이디',
	`uuid` CHAR(32) NOT NULL COMMENT 'uuid' COLLATE 'utf8mb4_general_ci',
	`created_by` VARCHAR(255) NULL DEFAULT NULL COMMENT '생성자' COLLATE 'utf8mb4_general_ci',
	`created_at` DATETIME NULL DEFAULT NULL COMMENT '생성시간',
	`updated_by` VARCHAR(255) NULL DEFAULT NULL COMMENT '수정자' COLLATE 'utf8mb4_general_ci',
	`updated_at` DATETIME NULL DEFAULT NULL COMMENT '수정시간',
	`deleted_by` VARCHAR(255) NULL DEFAULT NULL COMMENT '삭제자' COLLATE 'utf8mb4_general_ci',
	`deleted_at` DATETIME NULL DEFAULT NULL COMMENT '삭제시간',
	`name` VARCHAR(255) NOT NULL COMMENT '템플릿 이름' COLLATE 'utf8mb4_general_ci',
	`summary` VARCHAR(255) NULL DEFAULT '' COMMENT '템플릿 설명' COLLATE 'utf8mb4_general_ci',
	`api_version` VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'api version' COLLATE 'utf8mb4_general_ci',
	`origin` VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'origin name  ex) predefined, userdefined' COLLATE 'utf8mb4_general_ci',
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE INDEX `name` (`name`) USING BTREE,
	UNIQUE INDEX `uuid` (`uuid`) USING BTREE
)
COMMENT='template table'
COLLATE='utf8mb4_general_ci'
ENGINE=InnoDB
;
