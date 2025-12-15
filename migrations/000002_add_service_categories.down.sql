ALTER TABLE services
DROP FOREIGN KEY fk_services_category,
DROP COLUMN category_id;

DROP TABLE service_categories;