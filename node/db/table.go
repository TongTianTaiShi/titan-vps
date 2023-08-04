package db

var cOrderRecordTable = `
	CREATE TABLE if not exists %s (
		order_id           VARCHAR(128) NOT NULL UNIQUE,
		from_addr          VARCHAR(128) ,
		to_addr            VARCHAR(128) NOT NULL,
		value              BIGINT       DEFAULT 0,
		created_height     INT          DEFAULT 0,
		created_time       DATETIME     DEFAULT CURRENT_TIMESTAMP,
		state              INT          DEFAULT 0,
		done_state         INT          DEFAULT 0,
		done_time          DATETIME     DEFAULT CURRENT_TIMESTAMP,
		done_height        INT          DEFAULT 0,
		vps_id             VARCHAR(128) NOT NULL,
		PRIMARY KEY (order_id)
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
