SET FOREIGN_KEY_CHECKS = 0;

-- USERS
CREATE TABLE users (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(150) NOT NULL,
  username VARCHAR(100) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  role ENUM('owner','kasir','staff','courier') NOT NULL,
  phone VARCHAR(30),
  is_active TINYINT(1) DEFAULT 1,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- SERVICES
CREATE TABLE services (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  code VARCHAR(50) NOT NULL UNIQUE,
  name VARCHAR(150) NOT NULL,
  unit VARCHAR(20) NOT NULL,
  unit_price BIGINT NOT NULL,
  is_active TINYINT(1) DEFAULT 1,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- ORDERS
CREATE TABLE orders (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  kode_nota VARCHAR(50) NOT NULL UNIQUE,
  customer_name VARCHAR(150),
  customer_phone VARCHAR(30),
  customer_address TEXT,
  is_delivery TINYINT(1) DEFAULT 0,
  total_price BIGINT DEFAULT 0,
  payment_method ENUM('paid','unpaid','cod') DEFAULT 'unpaid',
  payment_status ENUM('unpaid','paid','cod_pending') DEFAULT 'unpaid',
  status_internal ENUM(
    'pending',
    'in-progress',
    'ready-pickup',
    'ready-delivery',
    'being-delivered',
    'finished-delivery',
    'picked-up',
    'cancelled'
  ) DEFAULT 'pending',
  estimated_ready_at DATETIME DEFAULT NULL,
  notes TEXT,
  created_by BIGINT DEFAULT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_orders_status (status_internal),
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB;

-- ORDER ITEMS
CREATE TABLE order_items (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  order_id BIGINT NOT NULL,
  service_id BIGINT DEFAULT NULL,
  description VARCHAR(255),
  quantity INT DEFAULT 1,
  weight_kg DECIMAL(8,2) DEFAULT 0,
  unit_price BIGINT NOT NULL,
  subtotal BIGINT NOT NULL,
  FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
  FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE SET NULL
) ENGINE=InnoDB;

-- PAYMENTS
CREATE TABLE payments (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  order_id BIGINT NOT NULL,
  method ENUM('cash','cod') NOT NULL,
  amount BIGINT NOT NULL,
  status ENUM('pending','confirmed') DEFAULT 'pending',
  collected_by BIGINT DEFAULT NULL,
  collected_at DATETIME DEFAULT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
  FOREIGN KEY (collected_by) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB;

-- STATUS HISTORY
CREATE TABLE status_history (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  order_id BIGINT NOT NULL,
  from_status VARCHAR(50),
  to_status VARCHAR(50) NOT NULL,
  actor_id BIGINT DEFAULT NULL,
  actor_role VARCHAR(50),
  note VARCHAR(255),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
  FOREIGN KEY (actor_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB;

-- COURIER ASSIGNMENTS
CREATE TABLE courier_assignments (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  order_id BIGINT NOT NULL UNIQUE,
  courier_id BIGINT DEFAULT NULL,
  picked_at DATETIME DEFAULT NULL,
  delivered_at DATETIME DEFAULT NULL,
  cod_collected_amount BIGINT DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
  FOREIGN KEY (courier_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB;

-- RECEIPTS
CREATE TABLE receipts (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  order_id BIGINT NOT NULL,
  file_path VARCHAR(255),
  generated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
) ENGINE=InnoDB;

SET FOREIGN_KEY_CHECKS = 1;


CREATE TABLE `users` (

	`id` BIGINT(19) NOT NULL AUTO_INCREMENT,

	`full_name` VARCHAR(150) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',

	`username` VARCHAR(100) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',

	`email` VARCHAR(150) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',

	`password_hash` VARCHAR(255) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',

	`role` ENUM('owner','cashier','staff','courier') NOT NULL COLLATE 'utf8mb4_0900_ai_ci',

	`phone_number` VARCHAR(30) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',

	`is_active` TINYINT(1) NULL DEFAULT '1',

	`last_login_at` DATETIME NULL DEFAULT NULL,

	`created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,

	`updated_at` TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,

	PRIMARY KEY (`id`) USING BTREE,

	UNIQUE INDEX `username` (`username`) USING BTREE,

	UNIQUE INDEX `email` (`email`) USING BTREE,

	UNIQUE INDEX `phone_number` (`phone_number`) USING BTREE,

	INDEX `idx_users_full_name` (`full_name`) USING BTREE

)

COLLATE='utf8mb4_0900_ai_ci'

ENGINE=InnoDB

AUTO_INCREMENT=5

;



CREATE TABLE `service_categories` (

	`id` BIGINT(19) NOT NULL AUTO_INCREMENT,

	`category_name` VARCHAR(150) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',

	`description` VARCHAR(255) NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',

	`is_active` TINYINT(1) NULL DEFAULT '1',

	`created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,

	`updated_at` TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,

	PRIMARY KEY (`id`) USING BTREE,

	UNIQUE INDEX `idx_category_name` (`category_name`) USING BTREE

)

COLLATE='utf8mb4_0900_ai_ci'

ENGINE=InnoDB

AUTO_INCREMENT=3

;



CREATE TABLE `services` (

	`id` BIGINT(19) NOT NULL AUTO_INCREMENT,

	`code` VARCHAR(50) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',

	`service_name` VARCHAR(150) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',

	`unit` ENUM('kg','pcs') NOT NULL COLLATE 'utf8mb4_0900_ai_ci',

	`price` DECIMAL(15,2) NOT NULL DEFAULT '0.00',

	`is_active` TINYINT(1) NULL DEFAULT '1',

	`created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,

	`updated_at` TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,

	`category_id` BIGINT(19) NOT NULL,

	`duration_hours` INT(10) NOT NULL DEFAULT '72',

	PRIMARY KEY (`id`) USING BTREE,

	UNIQUE INDEX `code` (`code`) USING BTREE,

	INDEX `fk_services_category` (`category_id`) USING BTREE,

	INDEX `idx_services_name` (`service_name`) USING BTREE,

	CONSTRAINT `fk_services_category` FOREIGN KEY (`category_id`) REFERENCES `service_categories` (`id`) ON UPDATE NO ACTION ON DELETE RESTRICT

)

COLLATE='utf8mb4_0900_ai_ci'

ENGINE=InnoDB

AUTO_INCREMENT=3

;



CREATE TABLE `orders` (

	`id` BIGINT(19) NOT NULL AUTO_INCREMENT,

	`invoice_number` VARCHAR(50) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',

	`customer_id` BIGINT(19) NULL DEFAULT NULL,

	`customer_name` VARCHAR(150) NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',

	`customer_phone` VARCHAR(30) NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',

	`customer_address` TEXT NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',

	`is_delivery` TINYINT(1) NULL DEFAULT '0',

	`total_price` DECIMAL(15,2) NULL DEFAULT '0.00',

	`payment_status` ENUM('unpaid','paid','cod_pending') NULL DEFAULT 'unpaid' COLLATE 'utf8mb4_0900_ai_ci',

	`status_internal` ENUM('pending','in-progress','ready-pickup','ready-delivery','being-delivered','finished-delivery','picked-up','cancelled') NULL DEFAULT 'pending' COLLATE 'utf8mb4_0900_ai_ci',

	`estimated_ready_at` TIMESTAMP NULL DEFAULT NULL,

	`notes` TEXT NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',

	`created_by` BIGINT(19) NULL DEFAULT NULL,

	`created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,

	`updated_at` TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,

	PRIMARY KEY (`id`) USING BTREE,

	UNIQUE INDEX `unique_invoice_number` (`invoice_number`) USING BTREE,

	INDEX `idx_orders_status` (`status_internal`) USING BTREE,

	INDEX `created_by` (`created_by`) USING BTREE,

	INDEX `fk_orders_customer` (`customer_id`) USING BTREE,

	INDEX `idx_orders_created_at` (`created_at`) USING BTREE,

	CONSTRAINT `fk_orders_customer` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL,

	CONSTRAINT `orders_ibfk_1` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL

)

COLLATE='utf8mb4_0900_ai_ci'

ENGINE=InnoDB

AUTO_INCREMENT=6

;



CREATE TABLE `order_items` (

	`id` BIGINT(19) NOT NULL AUTO_INCREMENT,

	`order_id` BIGINT(19) NOT NULL,

	`service_id` BIGINT(19) NULL DEFAULT NULL,

	`item_notes` VARCHAR(255) NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',

	`quantity` INT(10) NULL DEFAULT NULL,

	`qty_pieces` INT(10) NULL DEFAULT NULL,

	`weight_kg` DECIMAL(8,2) NULL DEFAULT NULL,

	`unit_price` DECIMAL(15,2) NOT NULL,

	`subtotal` DECIMAL(15,2) NOT NULL,

	PRIMARY KEY (`id`) USING BTREE,

	INDEX `order_id` (`order_id`) USING BTREE,

	INDEX `order_items_ibfk_2` (`service_id`) USING BTREE,

	CONSTRAINT `order_items_ibfk_1` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,

	CONSTRAINT `order_items_ibfk_2` FOREIGN KEY (`service_id`) REFERENCES `services` (`id`) ON UPDATE NO ACTION ON DELETE RESTRICT

)

COLLATE='utf8mb4_0900_ai_ci'

ENGINE=InnoDB

AUTO_INCREMENT=5

;



CREATE TABLE `customers` (

	`id` BIGINT(19) NOT NULL AUTO_INCREMENT,

	`full_name` VARCHAR(150) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',

	`phone_number` VARCHAR(30) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',

	`address` TEXT NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',

	`is_active` TINYINT(1) NULL DEFAULT '1',

	`created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,

	`updated_at` TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,

	PRIMARY KEY (`id`) USING BTREE,

	UNIQUE INDEX `phone` (`phone_number`) USING BTREE,

	INDEX `idx_customers_name` (`full_name`) USING BTREE

)

COLLATE='utf8mb4_0900_ai_ci'

ENGINE=InnoDB

AUTO_INCREMENT=2

;



CREATE TABLE `payments` (

	`id` BIGINT(19) NOT NULL AUTO_INCREMENT,

	`order_id` BIGINT(19) NOT NULL,

	`method` ENUM('cash','transfer','qris','ewallet') NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',

	`amount` DECIMAL(15,2) NOT NULL,

	`amount_received` DECIMAL(15,2) NOT NULL DEFAULT '0.00',

	`amount_change` DECIMAL(15,2) NOT NULL DEFAULT '0.00',

	`reference_no` VARCHAR(100) NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',

	`status` ENUM('pending','confirmed','void') NOT NULL DEFAULT 'pending' COLLATE 'utf8mb4_0900_ai_ci',

	`created_by` BIGINT(19) NOT NULL,

	`collected_by` BIGINT(19) NULL DEFAULT NULL,

	`collected_at` TIMESTAMP NULL DEFAULT NULL,

	`created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,

	`updated_at` TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,

	PRIMARY KEY (`id`) USING BTREE,

	INDEX `collected_by` (`collected_by`) USING BTREE,

	INDEX `payments_ibfk_1` (`order_id`) USING BTREE,

	INDEX `fk_payments_creator` (`created_by`) USING BTREE,

	INDEX `idx_payments_created_at` (`created_at`) USING BTREE,

	CONSTRAINT `fk_payments_creator` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE RESTRICT,

	CONSTRAINT `payments_ibfk_1` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON UPDATE NO ACTION ON DELETE RESTRICT,

	CONSTRAINT `payments_ibfk_2` FOREIGN KEY (`collected_by`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL

)

COLLATE='utf8mb4_0900_ai_ci'

ENGINE=InnoDB

AUTO_INCREMENT=5

;



CREATE TABLE `deliveries` (

	`id` BIGINT(19) NOT NULL AUTO_INCREMENT,

	`order_id` BIGINT(19) NOT NULL,

	`delivery_status` ENUM('ready-delivery','being-delivered','finished-delivery') NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',

	`shipping_cost` DECIMAL(15,2) NOT NULL DEFAULT '0.00',

	`courier_id` BIGINT(19) NULL DEFAULT NULL,

	`courier_departed_at` TIMESTAMP NULL DEFAULT NULL,

	`courier_arrived_at` TIMESTAMP NULL DEFAULT NULL,

	`receiver_name` VARCHAR(100) NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',

	`cod_collected_amount` DECIMAL(15,2) NULL DEFAULT '0.00',

	`created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,

	`updated_at` TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,

	PRIMARY KEY (`id`) USING BTREE,

	UNIQUE INDEX `order_id` (`order_id`) USING BTREE,

	INDEX `courier_id` (`courier_id`) USING BTREE,

	INDEX `delivery_status_idx` (`delivery_status`) USING BTREE,

	CONSTRAINT `deliveries_ibfk_1` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,

	CONSTRAINT `deliveries_ibfk_2` FOREIGN KEY (`courier_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL

)

COLLATE='utf8mb4_0900_ai_ci'

ENGINE=InnoDB

AUTO_INCREMENT=4

;



CREATE TABLE `status_history` (

	`id` BIGINT(19) NOT NULL AUTO_INCREMENT,

	`order_id` BIGINT(19) NOT NULL,

	`previous_status` VARCHAR(50) NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',

	`new_status` VARCHAR(50) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',

	`actor_id` BIGINT(19) NULL DEFAULT NULL,

	`actor_role` VARCHAR(50) NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',

	`notes` TEXT NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',

	`created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,

	PRIMARY KEY (`id`) USING BTREE,

	INDEX `order_id` (`order_id`) USING BTREE,

	INDEX `actor_id` (`actor_id`) USING BTREE,

	INDEX `idx_status_history_created_at` (`created_at`) USING BTREE,

	CONSTRAINT `status_history_ibfk_1` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,

	CONSTRAINT `status_history_ibfk_2` FOREIGN KEY (`actor_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL

)

COLLATE='utf8mb4_0900_ai_ci'

ENGINE=InnoDB

AUTO_INCREMENT=15

;

