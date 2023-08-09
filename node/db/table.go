package db

var cOrderRecordTable = `
	CREATE TABLE if not exists %s (
		order_id           VARCHAR(128) NOT NULL UNIQUE,
		from_addr          VARCHAR(128) DEFAULT "",
		to_addr            VARCHAR(128) NOT NULL,
		user_addr          VARCHAR(128) DEFAULT "",
		tx_hash            VARCHAR(128) DEFAULT "",
		value              BIGINT       DEFAULT 0,
		created_height     INT          DEFAULT 0,
		created_time       DATETIME     DEFAULT CURRENT_TIMESTAMP,
		state              INT          DEFAULT 0,
		done_state         INT          DEFAULT 0,
		done_time          DATETIME     DEFAULT CURRENT_TIMESTAMP,
		done_height        INT          DEFAULT 0,
		vps_id             BIGINT(20)   NOT NULL,
		msg                VARCHAR(2048) DEFAULT "",
		PRIMARY KEY (order_id),
		KEY idx_user (user_addr),
		KEY idx_to (to_addr)
	) ENGINE=InnoDB COMMENT='order record';`

var cVpsInstanceTable = `
	CREATE TABLE if not exists %s (
		id          		BIGINT(20) NOT NULL AUTO_INCREMENT,
		region_id          	VARCHAR(128) ,
		instance_type       VARCHAR(128) NOT NULL,
		dry_run             TINYINT(1) NOT NULL DEFAULT 0,
		image_id      		VARCHAR(128) NOT NULL,
		security_group_id   VARCHAR(128) NOT NULL,
		instanceCharge_type VARCHAR(128) NOT NULL,
		period_unit         VARCHAR(128) NOT NULL,
		period          	INT          DEFAULT 0,
		bandwidth_out       INT          DEFAULT 0,
		bandwidth_in        INT          DEFAULT 0,
		PRIMARY KEY (id)
	) ENGINE=InnoDB COMMENT='vps instance';`

var cInstanceInfoTable = `
	CREATE TABLE if not exists %s (
		id          		BIGINT(20) NOT NULL AUTO_INCREMENT,
		instance_id         VARCHAR(128) NOT NULL,
		order_id       		VARCHAR(128) NOT NULL,
		dry_run             VARCHAR(128) NOT NULL,
		RequestId      		VARCHAR(128) NOT NULL,
		TradePrice  		VARCHAR(128) NOT NULL,
		PublicIpAddress 	VARCHAR(128) NOT NULL,
		PrivateKeyStatus    TINYINT(1) NOT NULL DEFAULT 0,
		PRIMARY KEY (id)
	) ENGINE=InnoDB COMMENT='instance info';`

var cRechargeTable = `
	CREATE TABLE if not exists %s (
		order_id           VARCHAR(128) NOT NULL UNIQUE,
		from_addr          VARCHAR(128) DEFAULT "",
		to_addr            VARCHAR(128) NOT NULL,
		user_addr          VARCHAR(128) DEFAULT "",
		tx_hash            VARCHAR(128) DEFAULT "",
		recharge_addr      VARCHAR(128) NOT NULL,
		recharge_hash      VARCHAR(128) DEFAULT "",
		value              VARCHAR(32)  DEFAULT 0,
		created_height     INT          DEFAULT 0,
		created_time       DATETIME     DEFAULT CURRENT_TIMESTAMP,
		state              INT          DEFAULT 0,
		done_time          DATETIME     DEFAULT CURRENT_TIMESTAMP,
		done_height        INT          DEFAULT 0,
		done_state         INT          DEFAULT 0,
		msg                VARCHAR(2048) DEFAULT "",
		PRIMARY KEY (order_id),
		KEY idx_user (user_addr),
		KEY idx_to (to_addr)
	) ENGINE=InnoDB COMMENT='recharge info';`

var cWithdrawTable = `
	CREATE TABLE if not exists %s (
		id                 VARCHAR(128) NOT NULL UNIQUE,
		from_addr          VARCHAR(128) DEFAULT "",
		to_addr            VARCHAR(128) NOT NULL,
		user_addr          VARCHAR(128) DEFAULT "",
		tx_hash            VARCHAR(128) DEFAULT "",
		recharge_addr      VARCHAR(128) NOT NULL,
		recharge_hash      VARCHAR(128) DEFAULT "",
		value              BIGINT       DEFAULT 0,
		created_height     INT          DEFAULT 0,
		created_time       DATETIME     DEFAULT CURRENT_TIMESTAMP,
		state              INT          DEFAULT 0,
		done_time          DATETIME     DEFAULT CURRENT_TIMESTAMP,
		done_height        INT          DEFAULT 0,
		msg                VARCHAR(2048) DEFAULT "",
		PRIMARY KEY (id),
		KEY idx_user (user_addr),
		KEY idx_to (to_addr)
	) ENGINE=InnoDB COMMENT='withdraw info';`
