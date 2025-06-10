CREATE TABLE lecture_sessions (
  lecture_session_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  lecture_session_title VARCHAR(255) NOT NULL,
  lecture_session_description TEXT,
  lecture_session_teacher_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  lecture_session_scheduled_time TIMESTAMP NOT NULL,
  lecture_session_place TEXT,
  lecture_session_lecture_id UUID REFERENCES lectures(lecture_id) ON DELETE CASCADE,
  lecture_session_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexing
CREATE INDEX idx_lecture_sessions_teacher ON lecture_sessions(lecture_session_teacher_id);
CREATE INDEX idx_lecture_sessions_lecture ON lecture_sessions(lecture_session_lecture_id);
CREATE INDEX idx_lecture_sessions_schedule ON lecture_sessions(lecture_session_scheduled_time);


CREATE TABLE user_lecture_sessions (
  user_lecture_session_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_lecture_session_status_attendance VARCHAR(50), -- e.g., 'hadir', 'tidak_hadir', 'izin'
  user_lecture_session_grade_result FLOAT,
  user_lecture_session_lecture_session_id UUID NOT NULL REFERENCES lecture_sessions(lecture_session_id) ON DELETE CASCADE,
  user_lecture_session_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  user_lecture_session_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexing
CREATE INDEX idx_user_lecture_sessions_user ON user_lecture_sessions(user_lecture_session_user_id);
CREATE INDEX idx_user_lecture_sessions_lecture_session ON user_lecture_sessions(user_lecture_session_lecture_session_id);
CREATE INDEX idx_user_lecture_sessions_status ON user_lecture_sessions(user_lecture_session_status_attendance);
