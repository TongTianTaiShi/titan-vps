package db

var cOrderRecordTable = `
	CREATE TABLE if not exists %s (
		hash               VARCHAR(128) NOT NULL UNIQUE,
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
		PRIMARY KEY (hash)
	) ENGINE=InnoDB COMMENT='asset record';`
