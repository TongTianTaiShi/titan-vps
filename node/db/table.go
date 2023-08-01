package db

var cAssetRecordTable = `
	CREATE TABLE if not exists %s (
		hash               VARCHAR(128) NOT NULL UNIQUE,
		scheduler_sid      VARCHAR(128) NOT NULL,    
		cid                VARCHAR(128) NOT NULL,
		total_size         BIGINT       DEFAULT 0,
		total_blocks       INT          DEFAULT 0,
		edge_replicas      INT          DEFAULT 0,
		candidate_replicas INT          DEFAULT 0,
		expiration         DATETIME     NOT NULL,
		created_time       DATETIME     DEFAULT CURRENT_TIMESTAMP,
		end_time           DATETIME     DEFAULT CURRENT_TIMESTAMP,
		bandwidth          INT          DEFAULT 0,
		PRIMARY KEY (hash)
	) ENGINE=InnoDB COMMENT='asset record';`
