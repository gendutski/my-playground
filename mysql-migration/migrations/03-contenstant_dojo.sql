CREATE TABLE IF NOT EXISTS `contestant_dojo` (
  `contestant_id` bigint NOT NULL,
  `dojo_id` bigint NOT NULL,
  `is_active` tinyint(1) NOT NULL,
  `updated_at` timestamp NOT NULL,

  PRIMARY KEY (`contestant_id`, `dojo_id`),
  INDEX `contestant_dojo_search_idx1` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;