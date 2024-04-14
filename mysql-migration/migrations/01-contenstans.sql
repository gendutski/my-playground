CREATE TABLE IF NOT EXISTS `contestants` (
  `id` bigint UNSIGNED NOT NULL,
  `nik` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `birth_dt` date NOT NULL,
  `body_weight` double(8,2) NOT NULL,
  `gender` enum('male','female') COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `telp_no` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,

  PRIMARY KEY (`id`),
  UNIQUE KEY `contestants_nik_unique` (`nik`),
  INDEX `contestants_search_idx1` (`birth_dt`,`body_weight`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;