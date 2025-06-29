CREATE TABLE IF NOT EXISTS lecture_sessions (
    lecture_session_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lecture_session_title VARCHAR(255) NOT NULL,
    lecture_session_description TEXT,
    lecture_session_teacher JSONB NOT NULL,
    lecture_session_start_time TIMESTAMP NOT NULL,
    lecture_session_end_time TIMESTAMP NOT NULL,
    lecture_session_place TEXT,
    lecture_session_image_url TEXT,

    lecture_session_lecture_id UUID REFERENCES lectures(lecture_id) ON DELETE CASCADE,
    lecture_session_certificate_id UUID REFERENCES certificates(certificate_id) ON DELETE SET NULL,

    -- Metadata
    lecture_session_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index utama untuk mencari sesi berdasarkan lecture
CREATE INDEX IF NOT EXISTS idx_lecture_sessions_lecture_id 
  ON lecture_sessions(lecture_session_lecture_id);

-- Index untuk pencarian berdasarkan waktu mulai (misal untuk jadwal mendatang)
CREATE INDEX IF NOT EXISTS idx_lecture_sessions_start_time 
  ON lecture_sessions(lecture_session_start_time);

-- Index untuk pencarian berdasarkan waktu selesai (jika butuh range waktu)
CREATE INDEX IF NOT EXISTS idx_lecture_sessions_end_time 
  ON lecture_sessions(lecture_session_end_time);

-- Index untuk pencarian berdasarkan ID teacher (dalam JSON)
CREATE INDEX IF NOT EXISTS idx_lecture_sessions_teacher_id 
  ON lecture_sessions ((lecture_session_teacher->>'id'));

-- Index untuk pencarian berdasarkan nama teacher (dalam JSON)
CREATE INDEX IF NOT EXISTS idx_lecture_sessions_teacher_name 
  ON lecture_sessions ((lecture_session_teacher->>'name'));

-- Index untuk pencarian berdasarkan sertifikat per sesi
CREATE INDEX IF NOT EXISTS idx_lecture_sessions_certificate_id 
  ON lecture_sessions (lecture_session_certificate_id);

-- Index untuk sorting sesi terbaru
CREATE INDEX IF NOT EXISTS idx_lecture_sessions_created_at 
  ON lecture_sessions (lecture_session_created_at DESC);



CREATE TABLE IF NOT EXISTS user_lecture_sessions (
  user_lecture_session_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  -- Kehadiran dan nilai
  user_lecture_session_attendance_status INT, -- 0 = tidak hadir, 1 = hadir, 2 = hadir online
  user_lecture_session_grade_result FLOAT,

  -- Relasi
  user_lecture_session_lecture_session_id UUID NOT NULL REFERENCES lecture_sessions(lecture_session_id) ON DELETE CASCADE,
  user_lecture_session_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

  -- Metadata
  user_lecture_session_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- Index by user untuk ambil daftar kehadiran user
CREATE INDEX IF NOT EXISTS idx_user_lecture_sessions_user 
  ON user_lecture_sessions(user_lecture_session_user_id);

-- Index by sesi untuk ambil siapa aja yang hadir
CREATE INDEX IF NOT EXISTS idx_user_lecture_sessions_lecture_session 
  ON user_lecture_sessions(user_lecture_session_lecture_session_id);

-- Index by status kehadiran (untuk statistik absensi)
CREATE INDEX IF NOT EXISTS idx_user_lecture_sessions_attendance_status 
  ON user_lecture_sessions(user_lecture_session_attendance_status);
