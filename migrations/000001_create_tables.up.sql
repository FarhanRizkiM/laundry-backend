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
