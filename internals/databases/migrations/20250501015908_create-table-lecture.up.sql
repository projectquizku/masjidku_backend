-- Tabel lectures: berisi informasi kajian/ceramah
CREATE TABLE IF NOT EXISTS lectures (
    lecture_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lecture_title VARCHAR(255) NOT NULL,
    lecture_description TEXT,
    total_lecture_sessions INTEGER, -- ❗️Tambahan: untuk skenario yang memiliki batas sesi
    lecture_is_recurring BOOLEAN DEFAULT FALSE, -- ✅ Apakah ini kajian berulang?
    lecture_recurrence_interval INTEGER, -- ✅ Jumlah hari antara kajian berulang
    lecture_masjid_id UUID REFERENCES masjids(masjid_id) ON DELETE CASCADE,
    lecture_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk pencarian cepat berdasarkan masjid
CREATE INDEX idx_lecture_masjid_id ON lectures(lecture_masjid_id);


-- Tabel user_lectures: relasi user mengikuti kajian
CREATE TABLE IF NOT EXISTS user_lectures (
    user_lecture_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_lecture_grade_result INT, -- nilai hasil jika ada evaluasi
    user_lecture_lecture_id UUID NOT NULL REFERENCES lectures(lecture_id) ON DELETE CASCADE,
    user_lecture_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_lecture_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_lecture_lecture_id, user_lecture_user_id) -- satu user tidak bisa dua kali ikut satu kajian
);

-- Index untuk pencarian cepat
CREATE INDEX idx_user_lecture_lecture_id ON user_lectures(user_lecture_lecture_id);
CREATE INDEX idx_user_lecture_user_id ON user_lectures(user_lecture_user_id);
