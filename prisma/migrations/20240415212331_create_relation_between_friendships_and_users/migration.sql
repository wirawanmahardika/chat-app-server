-- AddForeignKey
ALTER TABLE `friendships` ADD CONSTRAINT `friendships_id_user_1_fkey` FOREIGN KEY (`id_user_1`) REFERENCES `users`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `friendships` ADD CONSTRAINT `friendships_id_user_2_fkey` FOREIGN KEY (`id_user_2`) REFERENCES `users`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;
