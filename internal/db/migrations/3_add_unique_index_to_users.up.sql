ALTER TABLE `users` ADD UNIQUE INDEX `idx_users_username` (`username`);
ALTER TABLE `users` ADD UNIQUE INDEX `idx_users_email` (`email`);