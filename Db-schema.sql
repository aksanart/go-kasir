CREATE TABLE `category`  (
  `category_id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`category_id`)
);