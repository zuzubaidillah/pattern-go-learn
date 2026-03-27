-- 000001_create_users_table.up.sql
-- Digunakan untuk membangun sturktur utama tabel users beserta kolom pelengkap autentikasi

CREATE TABLE IF NOT EXISTS `users` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT 'Primary key yang terhindar dari limitasi angka 32-bit',
    `name` VARCHAR(100) NOT NULL COMMENT 'Nama lengkap milik si pendaftar',
    `email` VARCHAR(100) NOT NULL UNIQUE COMMENT 'Validasi unik tingkat database agar tidak ada 2 akun dengan email sama',
    `password` VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'Menyimpan teks sandi yang sudah diasinkan dan dienkripsi kuat (Bcrypt)',
    `role` VARCHAR(50) NOT NULL DEFAULT 'user' COMMENT 'Peran kendali akses RBAC, cth: user, admin, super_admin',
    `status` VARCHAR(20) NOT NULL COMMENT 'Tanda apakah pengguna aktif, tersetuju, atau di-banned (active/inactive)',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'Kapan dia pertama kali mendaftar',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Kapan dia terakhir mengubah profil / ganti sandi'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
