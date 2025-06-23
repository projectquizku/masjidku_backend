-- =============================
-- TABLE: lectures
-- =============================
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
    lecture_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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

-- Index baru untuk pencarian by certificate
CREATE INDEX IF NOT EXISTS idx_lecture_certificate_id 
  ON lectures(lecture_certificate_id);


-- =============================
-- TABLE: user_lectures
-- =============================
CREATE TABLE IF NOT EXISTS user_lectures (
    user_lecture_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_lecture_grade_result INT,
    user_lecture_lecture_id UUID NOT NULL REFERENCES lectures(lecture_id) ON DELETE CASCADE,
    user_lecture_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_lecture_total_completed_sessions INT DEFAULT 0,
    user_lecture_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_lecture_lecture_id, user_lecture_user_id)
);

-- Indexing
CREATE INDEX IF NOT EXISTS idx_user_lecture_lecture_id ON user_lectures(user_lecture_lecture_id);
CREATE INDEX IF NOT EXISTS idx_user_lecture_user_id ON user_lectures(user_lecture_user_id);
