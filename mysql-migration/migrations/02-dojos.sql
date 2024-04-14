CREATE TABLE IF NOT EXISTS `dojos` (
  `id` bigint UNSIGNED NOT NULL,
  `serial` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `region_id` int NOT NULL,
  `sport_id` int NOT NULL COMMENT 'ref to params, type: sport-type',
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text COLLATE utf8mb4_unicode_ci,
  `contact_person` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contact_no` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `address` text COLLATE utf8mb4_unicode_ci,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,

  PRIMARY KEY (`id`),
  UNIQUE KEY `dojos_serial_unique` (`serial`),
  INDEX `dojos_region_id_sport_id_name_index` (`region_id`,`sport_id`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;