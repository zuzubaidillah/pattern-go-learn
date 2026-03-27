-- 000001_create_users_table.down.sql
-- Digunakan untuk membongkar dan menghapus bersih seluruh isi tabel saat proses rollback/down (-1 step)

DROP TABLE IF EXISTS `users`;
