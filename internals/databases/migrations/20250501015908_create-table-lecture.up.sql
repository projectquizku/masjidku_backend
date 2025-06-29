CREATE TABLE IF NOT EXISTS lectures (
    lecture_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lecture_title VARCHAR(255) NOT NULL,
    lecture_description TEXT,
    total_lecture_sessions INTEGER,
    lecture_status BOOLEAN DEFAULT FALSE, -- FALSE = ongoing, TRUE = finished
    lecture_certificate_id UUID REFERENCES certificates(certificate_id) ON DELETE SET NULL,
    lecture_image_url TEXT,
    lecture_teachers JSONB,
    lecture_masjid_id UUID REFERENCES masjids(masjid_id) ON DELETE CASCADE,

    -- Penyesuaian baru
    lecture_is_registration_required BOOLEAN DEFAULT FALSE,
    lecture_is_paid BOOLEAN DEFAULT FALSE,
    lecture_price INT,
    lecture_payment_deadline TIMESTAMP,
    lecture_payment_scope VARCHAR(10) DEFAULT 'lecture', -- 'lecture' or 'session'

     -- Tambahan dari lecture_sessions
    lecture_capacity INT,
    lecture_is_public BOOLEAN DEFAULT TRUE,

    lecture_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP -- ⬅️ Soft delete support

);

-- Indexing
CREATE INDEX IF NOT EXISTS idx_lecture_masjid_id 
  ON lectures(lecture_masjid_id);

CREATE INDEX IF NOT EXISTS idx_lecture_created_at 
  ON lectures(lecture_created_at DESC);

CREATE INDEX IF NOT EXISTS idx_lecture_masjid_created_at 
  ON lectures(lecture_masjid_id, lecture_created_at DESC);

CREATE INDEX IF NOT EXISTS idx_lecture_status 
  ON lectures(lecture_status);

CREATE INDEX IF NOT EXISTS idx_lecture_certificate_id 
  ON lectures(lecture_certificate_id);

CREATE INDEX IF NOT EXISTS idx_lecture_payment_scope 
  ON lectures(lecture_payment_scope);

-- Index deleted_at untuk optimasi soft delete
CREATE INDEX IF NOT EXISTS idx_lectures_deleted_at 
  ON lectures(deleted_at);



CREATE TABLE IF NOT EXISTS user_lectures (
    user_lecture_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Relasi
    user_lecture_lecture_id UUID NOT NULL REFERENCES lectures(lecture_id) ON DELETE CASCADE,
    user_lecture_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Progres dan nilai
    user_lecture_grade_result INT,
    user_lecture_total_completed_sessions INT DEFAULT 0,

    -- Pendaftaran dan pembayaran (jika level lecture)
    user_lecture_is_registered BOOLEAN DEFAULT FALSE,
    user_lecture_has_paid BOOLEAN DEFAULT FALSE,
    user_lecture_paid_amount INT,
    user_lecture_payment_time TIMESTAMP,

    -- Metadata
    user_lecture_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Unik: satu user hanya boleh daftar satu kali ke satu lecture
    UNIQUE(user_lecture_lecture_id, user_lecture_user_id)
);

CREATE INDEX IF NOT EXISTS idx_user_lecture_lecture_id 
  ON user_lectures(user_lecture_lecture_id);

CREATE INDEX IF NOT EXISTS idx_user_lecture_user_id 
  ON user_lectures(user_lecture_user_id);
