-- +migrate Up
CREATE TABLE IF NOT EXISTS donations (
    donation_id SERIAL PRIMARY KEY,
    donation_user_id UUID REFERENCES users(id) ON DELETE SET NULL, 
    donation_amount INTEGER NOT NULL CHECK (donation_amount > 0), 
    donation_message TEXT, 
    donation_status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (
        donation_status IN ('pending', 'paid', 'expired', 'canceled')
    ), 
    donation_order_id VARCHAR(100) NOT NULL UNIQUE CHECK (
        char_length(donation_order_id) <= 100
    ), 
    donation_payment_token TEXT, 
    donation_payment_gateway VARCHAR(50) DEFAULT 'midtrans',
    donation_payment_method VARCHAR,
    donation_paid_at TIMESTAMP, 
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,                       -- Timestamp pembuatan data
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,                       -- Timestamp terakhir kali data diperbarui
    deleted_at TIMESTAMP                                                           -- Soft delete (jika ada penghapusan)
);

-- 🔍 Index untuk pencarian cepat order_id (case-insensitive)
CREATE INDEX IF NOT EXISTS idx_donations_order_id_lower 
    ON donations (LOWER(donation_order_id));

-- 🔍 Index umum untuk pencarian berdasarkan user
CREATE INDEX IF NOT EXISTS idx_donations_user_id 
    ON donations (donation_user_id);



-- +migrate Up
CREATE TABLE IF NOT EXISTS donation_questions (
    donation_question_id SERIAL PRIMARY KEY,  -- ID unik untuk setiap entri pertanyaan dari donasi
    donation_question_donation_id BIGINT NOT NULL REFERENCES donations(donation_id) ON DELETE CASCADE,
    -- Relasi ke tabel 'donations'. Jika donasi dihapus, semua entri pertanyaannya ikut terhapus
    donation_question_question_id BIGINT NOT NULL REFERENCES questions(question_id) ON DELETE CASCADE,
    -- Relasi ke tabel 'questions'. Menunjukkan pertanyaan mana yang dihubungkan ke donasi ini
    donation_question_user_progress_id BIGINT REFERENCES user_progress(user_progress_id) ON DELETE SET NULL,
    -- Opsional: digunakan untuk melacak progress user pada pertanyaan ini, bisa NULL jika belum ada
    donation_question_user_message TEXT,
    -- Opsional: pesan personal dari user/donatur yang ingin dikaitkan ke pertanyaan
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Waktu saat entri ini dibuat
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    -- Waktu terakhir entri ini diperbarui
);

-- Index tambahan untuk pencarian cepat berdasarkan donation_id
CREATE INDEX IF NOT EXISTS idx_donation_question_donation_id ON donation_questions(donation_question_donation_id);

-- Index tambahan untuk pencarian cepat berdasarkan question_id
CREATE INDEX IF NOT EXISTS idx_donation_question_question_id ON donation_questions(donation_question_question_id);
