package db

var cOrderRecordTable = `
	CREATE TABLE if not exists %s (
		order_id           VARCHAR(128)  NOT NULL UNIQUE,
		user_id            VARCHAR(128)  NOT NULL,
		value              VARCHAR(32)   DEFAULT 0,
		created_time       DATETIME      DEFAULT CURRENT_TIMESTAMP,
		cycle_time         VARCHAR(64)   DEFAULT "",
		state              INT           DEFAULT 0,
		done_state         INT           DEFAULT 0,
		done_time          DATETIME      DEFAULT CURRENT_TIMESTAMP,
		vps_id             BIGINT(20)    NOT NULL,
		msg                VARCHAR(2048) DEFAULT "",
		order_type         INT           DEFAULT 0,
		expiration         DATETIME      DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (order_id),
		KEY idx_user (user_id)
	) ENGINE=InnoDB COMMENT='order record';`

var cInstanceDetailsTable = `
	CREATE TABLE if not exists %s (
		id          		 BIGINT(20) NOT NULL AUTO_INCREMENT,
		region_id          	 VARCHAR(128)  NOT NULL,
		instance_id          VARCHAR(128)  DEFAULT '',
		instance_name        VARCHAR(128)  DEFAULT '',
		user_id          	 VARCHAR(128)  NOT NULL,
		order_id          	 VARCHAR(128)  NOT NULL,
		instance_type        VARCHAR(128)  DEFAULT '',
		image_id      		 VARCHAR(128)  DEFAULT '',
		security_group_id    VARCHAR(128)  DEFAULT '',
		instance_charge_type VARCHAR(128)  DEFAULT '',
		internet_charge_type VARCHAR(128)  DEFAULT '',
		period_unit          VARCHAR(128)  DEFAULT '',
		period          	 INT           DEFAULT 0,
		bandwidth_out        INT           DEFAULT 0,
		bandwidth_in         INT           DEFAULT 0,
		system_disk_size     INT           DEFAULT 0,
	    ip_address 	         VARCHAR(16)   DEFAULT '',
	    value       		 VARCHAR(32)   DEFAULT 0,
	    memory  		     FLOAT         DEFAULT 0,
	    memory_used  		 FLOAT         DEFAULT 0,  
	    cores                INT           DEFAULT 0,
	    auto_renew           INT           DEFAULT 0,
	    cores_used  		 FLOAT         DEFAULT 0,
	    system_disk_category VARCHAR(128)  DEFAULT '',
	    os_type 		     VARCHAR(128)  DEFAULT '',
	    expired_time 		 VARCHAR(128)  DEFAULT '',
	    data_disk 		     VARCHAR(1028) DEFAULT '',
		created_time         DATETIME      DEFAULT CURRENT_TIMESTAMP,
		access_key           VARCHAR(32)   DEFAULT '',
		renew                VARCHAR(16)   DEFAULT '',
		state                VARCHAR(16)   DEFAULT '',
		update_time          DATETIME      DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id),
		KEY idx_user (user_id),
		KEY idx_instance (instance_id)
	) ENGINE=InnoDB COMMENT='vps instance';`

var cRechargeTable = `
	CREATE TABLE if not exists %s (
		order_id           VARCHAR(128) NOT NULL UNIQUE,
		from_addr          VARCHAR(128) DEFAULT "",
		to_addr            VARCHAR(128) NOT NULL,
		user_id            VARCHAR(128) DEFAULT "",
		value              VARCHAR(32)  DEFAULT 0,
		created_time       DATETIME     DEFAULT CURRENT_TIMESTAMP,
		state              INT          DEFAULT 0,
		done_time          DATETIME     DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (order_id),
		KEY idx_user (user_id),
		KEY idx_to (to_addr)
	) ENGINE=InnoDB COMMENT='recharge info';`

var cWithdrawTable = `
	CREATE TABLE if not exists %s (
		order_id           VARCHAR(128) NOT NULL UNIQUE,
		user_id            VARCHAR(128) DEFAULT "",
		withdraw_addr      VARCHAR(128) NOT NULL,
		withdraw_hash      VARCHAR(128) DEFAULT "",
		value              VARCHAR(32)  DEFAULT 0,
		created_time       DATETIME     DEFAULT CURRENT_TIMESTAMP,
		state              INT          DEFAULT 0,
		done_time          DATETIME     DEFAULT CURRENT_TIMESTAMP,
		executor           VARCHAR(128) DEFAULT "",
		PRIMARY KEY (order_id),
		KEY idx_user (user_id)
	) ENGINE=InnoDB COMMENT='withdraw info';`

var cConfigTable = `
	CREATE TABLE if not exists %s (
		name       VARCHAR(16)  DEFAULT "",
		value      VARCHAR(32)  DEFAULT "",
		PRIMARY KEY (name)
	) ENGINE=InnoDB COMMENT='config info';`

var cUserTable = `
	CREATE TABLE if not exists %s (
		user_id        VARCHAR(128) NOT NULL UNIQUE,
		balance        VARCHAR(32)  DEFAULT 0,
		created_time   DATETIME     DEFAULT CURRENT_TIMESTAMP,
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
		user_id       VARCHAR(128) NOT NULL UNIQUE,
		nick_name     VARCHAR(32)  DEFAULT 0,
		created_time  DATETIME     DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (user_id)
	) ENGINE=InnoDB COMMENT='admin info';`

var cInstanceDefaultTable = `
	CREATE TABLE if not exists %s (
		region_id                  VARCHAR(128)  NOT NULL,
		instance_type_id           VARCHAR(128)  DEFAULT 0,
		memory_size                FLOAT         DEFAULT 0,
		cpu_architecture           VARCHAR(128)  DEFAULT 0,
		instance_category          VARCHAR(128)  DEFAULT 0,
		cpu_core_count             INT           DEFAULT 0,
		available_zone             VARCHAR(128)  DEFAULT 0,
		instance_type_family       VARCHAR(128)  DEFAULT 0,
		physical_processor_model   VARCHAR(128)  DEFAULT 0,
		price                      FLOAT         DEFAULT 0,
		original_price             FLOAT         DEFAULT 0,
	    status                     VARCHAR(16)   DEFAULT 0,
	    created_time               DATETIME      DEFAULT CURRENT_TIMESTAMP,
	    updated_time               DATETIME      DEFAULT CURRENT_TIMESTAMP,
		UNIQUE KEY (region_id,instance_type_id)
	) ENGINE=InnoDB COMMENT='instance info';`

var cInstanceRefundTable = `
	CREATE TABLE if not exists %s (
		instance_id   VARCHAR(128) NOT NULL,
		executor      VARCHAR(128) DEFAULT "",
	    refund_time   VARCHAR(32)  DEFAULT "",
		PRIMARY KEY (instance_id)
	) ENGINE=InnoDB COMMENT='instance refund';`

var cInvitationTable = `
	CREATE TABLE if not exists %s (
		invitation_code   	VARCHAR(128) NOT NULL,
		id   	VARCHAR(128),
		PRIMARY KEY (invitation_code)
	) ENGINE=InnoDB COMMENT='invitation';`

var cProviderInfoTable = `
	CREATE TABLE if not exists %s (
		id VARCHAR(128) NOT NULL,
	    address VARCHAR(128),
	    email VARCHAR(128),
		passwd VARCHAR(128),
	    filecoin VARCHAR(128),
	    create_time BIGINT(20),
	    PRIMARY KEY (id)
	)ENGINE=InnoDB COMMENT='provider info';`

var cAccessKeyTable = `
    CREATE TABLE if not exists %s (
	    provider_id    VARCHAR(128) NOT NULL UNIQUE,
	    access_secret  VARCHAR(128) DEFAULT 0,
	    access_key     VARCHAR(128) DEFAULT 0,
		k_type         INT          DEFAULT 0,
		state          INT          DEFAULT 0,
		rebate         FLOAT        DEFAULT 0,
	    nick           VARCHAR(64)  DEFAULT 0,
	    created_time   DATETIME     DEFAULT CURRENT_TIMESTAMP,
	    PRIMARY KEY (access_secret),
		KEY idx_user (provider_id)
    ) ENGINE=InnoDB COMMENT='provider access key';`
