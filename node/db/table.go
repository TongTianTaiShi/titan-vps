package db

var cOrderRecordTable = `
	CREATE TABLE if not exists %s (
		order_id           VARCHAR(128) NOT NULL UNIQUE,
		from_addr          VARCHAR(128) DEFAULT "",
		to_addr            VARCHAR(128) DEFAULT "",
		user_id            VARCHAR(128) NOT NULL,
		tx_hash            VARCHAR(128) DEFAULT "",
		value              VARCHAR(32)  DEFAULT 0,
		created_height     INT          DEFAULT 0,
		created_time       DATETIME     DEFAULT CURRENT_TIMESTAMP,
		state              INT          DEFAULT 0,
		done_state         INT          DEFAULT 0,
		done_time          DATETIME     DEFAULT CURRENT_TIMESTAMP,
		done_height        INT          DEFAULT 0,
		vps_id             BIGINT(20)   NOT NULL,
		msg                VARCHAR(2048) DEFAULT "",
	    trade_price  	   VARCHAR(128)  DEFAULT "",
		expiration         DATETIME     DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (order_id),
		KEY idx_user (user_id)
	) ENGINE=InnoDB COMMENT='order record';`

var cInstanceDetailsTable = `
	CREATE TABLE if not exists %s (
		id          		BIGINT(20) NOT NULL AUTO_INCREMENT,
		region_id          	VARCHAR(128) ,
		instance_id          	VARCHAR(128) ,
		user_id          	VARCHAR(128) ,
		order_id          	VARCHAR(128) ,
		instance_type       VARCHAR(128) NOT NULL,
		dry_run             TINYINT(1) NOT NULL DEFAULT 0,
		image_id      		VARCHAR(128) NOT NULL,
		security_group_id   VARCHAR(128) NOT NULL,
		instance_charge_type VARCHAR(128) NOT NULL,
		internet_charge_type VARCHAR(128) NOT NULL,
		period_unit         VARCHAR(128) NOT NULL,
		period          	INT          DEFAULT 0,
		bandwidth_out       INT          DEFAULT 0,
		bandwidth_in        INT          DEFAULT 0,
		system_disk_size        Float  NOT NULL DEFAULT 0,
	    ip_address 	VARCHAR(128) NOT NULL,
	    trade_price  		Float   NOT NULL DEFAULT 0,
	    memory  		Float   NOT NULL DEFAULT 0,
	    memory_used  		Float   NOT NULL DEFAULT 0,
	    cores        INT          DEFAULT 0,
	    cores_used  		Float   NOT NULL DEFAULT 0,
	    system_disk_category  		VARCHAR(128) NOT NULL,
	    os_type 		VARCHAR(128) NOT NULL,
		created_time       DATETIME     DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
	) ENGINE=InnoDB COMMENT='vps instance';`

// var cMyServersTable = `
//	CREATE TABLE if not exists %s (
//		id          		BIGINT(20) NOT NULL AUTO_INCREMENT,
//		server_name         VARCHAR(128) NOT NULL,
//		system       		VARCHAR(128) NOT NULL,
//		location             VARCHAR(128) NOT NULL,
//		price      		VARCHAR(128) NOT NULL,
//		status  		TINYINT(1) NOT NULL DEFAULT 0,
//		internet_charge_type 	TINYINT(1) NOT NULL DEFAULT 0,
//		created_time       DATETIME     DEFAULT CURRENT_TIMESTAMP,
//		PRIMARY KEY (id)
//	) ENGINE=InnoDB COMMENT='instance info';`

var cMyInstancesTable = `
	CREATE TABLE if not exists %s (
		id          		BIGINT(20) NOT NULL AUTO_INCREMENT,
		instance_id         VARCHAR(128) NOT NULL,
		order_id       		VARCHAR(128) NOT NULL,
		user_id       		VARCHAR(128) NOT NULL,
		private_key_status    TINYINT(1) NOT NULL DEFAULT 0,
	    instance_name         VARCHAR(128) NOT NULL,
		instance_system       		VARCHAR(128) NOT NULL,
		location             VARCHAR(128) NOT NULL,
		price      		float NOT NULL DEFAULT 0,
		state  		 VARCHAR(28) NOT NULL,
		internet_charge_type 	VARCHAR(128) NOT NULL,
		created_time       DATETIME     DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
	) ENGINE=InnoDB COMMENT='my instance';`

var cRechargeTable = `
	CREATE TABLE if not exists %s (
		order_id           VARCHAR(128) NOT NULL UNIQUE,
		from_addr          VARCHAR(128) DEFAULT "",
		to_addr            VARCHAR(128) NOT NULL,
		user_id            VARCHAR(128) DEFAULT "",
		value              VARCHAR(32)  DEFAULT 0,
		created_height     INT          DEFAULT 0,
		created_time       DATETIME     DEFAULT CURRENT_TIMESTAMP,
		state              INT          DEFAULT 0,
		done_time          DATETIME     DEFAULT CURRENT_TIMESTAMP,
		done_height        INT          DEFAULT 0,
		PRIMARY KEY (order_id),
		KEY idx_user (user_id),
		KEY idx_to (to_addr)
	) ENGINE=InnoDB COMMENT='recharge info';`

var cWithdrawTable = `
	CREATE TABLE if not exists %s (
		order_id           VARCHAR(128) NOT NULL UNIQUE,
		from_addr          VARCHAR(128) DEFAULT "",
		to_addr            VARCHAR(128) NOT NULL,
		user_id            VARCHAR(128) DEFAULT "",
		withdraw_addr      VARCHAR(128) NOT NULL,
		withdraw_hash      VARCHAR(128) DEFAULT "",
		value              VARCHAR(32)  DEFAULT 0,
		created_height     INT          DEFAULT 0,
		created_time       DATETIME     DEFAULT CURRENT_TIMESTAMP,
		state              INT          DEFAULT 0,
		done_time          DATETIME     DEFAULT CURRENT_TIMESTAMP,
		done_height        INT          DEFAULT 0,
		executor           VARCHAR(128) DEFAULT "",
		PRIMARY KEY (order_id),
		KEY idx_user (user_id),
		KEY idx_to (to_addr)
	) ENGINE=InnoDB COMMENT='withdraw info';`

var cConfigTable = `
	CREATE TABLE if not exists %s (
		name       VARCHAR(16)  DEFAULT "",
		value      VARCHAR(32)  DEFAULT "",
		PRIMARY KEY (name)
	) ENGINE=InnoDB COMMENT='config info';`

var cUserTable = `
	CREATE TABLE if not exists %s (
		user_id      VARCHAR(128) NOT NULL UNIQUE,
		balance        VARCHAR(32)  DEFAULT 0,
		PRIMARY KEY (user_id)
	) ENGINE=InnoDB COMMENT='user info';`

var cRechargeAddressTable = `
	CREATE TABLE if not exists %s (
		addr      VARCHAR(128) NOT NULL UNIQUE,
		user_id   VARCHAR(128) DEFAULT "",
		PRIMARY KEY (addr)
	) ENGINE=InnoDB COMMENT='recharge address ';`

var cAdminTable = `
	CREATE TABLE if not exists %s (
		user_id      VARCHAR(128) NOT NULL UNIQUE,
		nick_name    VARCHAR(32)  DEFAULT 0,
		PRIMARY KEY (user_id)
	) ENGINE=InnoDB COMMENT='admin info';`

var cInstanceDefaultTable = `
	CREATE TABLE if not exists %s (
		region_id      VARCHAR(128) NOT NULL,
		instance_type_id    VARCHAR(128)  DEFAULT 0,
		memory_size    float  DEFAULT 0,
		cpu_architecture    VARCHAR(128)  DEFAULT 0,
		instance_category    VARCHAR(128)  DEFAULT 0,
		cpu_core_count    int  DEFAULT 0,
		available_zone    VARCHAR(128)  DEFAULT 0,
		instance_type_family    VARCHAR(128)  DEFAULT 0,
		physical_processor_model    VARCHAR(128)  DEFAULT 0,
		original_price    float  DEFAULT 0,
		price    float  DEFAULT 0,
	    created_time       DATETIME     DEFAULT CURRENT_TIMESTAMP,
	    updated_time       DATETIME     DEFAULT CURRENT_TIMESTAMP,
		UNIQUE KEY (region_id,instance_type_id)
	) ENGINE=InnoDB COMMENT='instance info';`
