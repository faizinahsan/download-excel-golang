-- Optimasi Database - Tambahkan index yang diperlukan

-- Index untuk kolom created (jika sering digunakan untuk sorting berdasarkan tanggal)
CREATE INDEX IF NOT EXISTS idx_customer_update_profile_created ON customer_update_profile(created DESC);

-- Index untuk kombinasi id dan created untuk query yang lebih cepat
CREATE INDEX IF NOT EXISTS idx_customer_update_profile_id_created ON customer_update_profile(id DESC, created DESC);

-- Jika sering menggunakan filter berdasarkan customer_name atau email
CREATE INDEX IF NOT EXISTS idx_customer_update_profile_customer_name ON customer_update_profile(customer_name);
CREATE INDEX IF NOT EXISTS idx_customer_update_profile_customer_email ON customer_update_profile(customer_email);

-- Index untuk reference_number jika sering digunakan untuk pencarian
CREATE INDEX IF NOT EXISTS idx_customer_update_profile_reference_number ON customer_update_profile(reference_number);

-- Untuk query yang menggunakan range berdasarkan ID, buat partial index
-- Index ini hanya untuk data dengan ID yang besar (data terbaru)
CREATE INDEX IF NOT EXISTS idx_customer_update_profile_recent_id ON customer_update_profile(id DESC)
WHERE id > (SELECT MAX(id) - 2000000 FROM customer_update_profile);

-- Update statistics untuk query planner yang lebih baik
ANALYZE customer_update_profile;

-- Atur work_mem untuk session yang menjalankan query besar (optional)
-- SET work_mem = '256MB'; -- Sesuaikan dengan available memory
