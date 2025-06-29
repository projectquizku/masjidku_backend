-- =============================
-- TABLE: lecture_sessions_questions
-- =============================
CREATE TABLE IF NOT EXISTS lecture_sessions_questions (
  lecture_sessions_question_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  lecture_sessions_question TEXT NOT NULL,
  lecture_sessions_question_answer TEXT NOT NULL,
  lecture_sessions_question_correct CHAR(1) NOT NULL CHECK (lecture_sessions_question_correct IN ('A', 'B', 'C', 'D')),
  lecture_sessions_question_explanation TEXT,
  lecture_sessions_question_quiz_id UUID,
  lecture_sessions_question_exam_id UUID,
  lecture_sessions_question_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT fk_question_quiz FOREIGN KEY (lecture_sessions_question_quiz_id)
    REFERENCES lecture_sessions_quiz (lecture_sessions_quiz_id) ON DELETE SET NULL,

  CONSTRAINT fk_question_exam FOREIGN KEY (lecture_sessions_question_exam_id)
    REFERENCES lecture_sessions_exams (lecture_sessions_exam_id) ON DELETE SET NULL
);

-- Index untuk pencarian berdasarkan quiz atau exam
CREATE INDEX IF NOT EXISTS idx_lecture_sessions_questions_quiz_id
  ON lecture_sessions_questions (lecture_sessions_question_quiz_id);

CREATE INDEX IF NOT EXISTS idx_lecture_sessions_questions_exam_id
  ON lecture_sessions_questions (lecture_sessions_question_exam_id);

CREATE INDEX IF NOT EXISTS idx_lecture_sessions_questions_created_at
  ON lecture_sessions_questions (lecture_sessions_question_created_at);


-- =============================
-- TABLE: user_lecture_sessions_questions
-- =============================
CREATE TABLE IF NOT EXISTS user_lecture_sessions_questions (
  lecture_sessions_user_question_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  lecture_sessions_user_question_answer CHAR(1) NOT NULL CHECK (lecture_sessions_user_question_answer IN ('A', 'B', 'C', 'D')),
  lecture_sessions_user_question_is_correct BOOLEAN NOT NULL,
  lecture_sessions_user_question_question_id UUID NOT NULL,
  lecture_sessions_user_question_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT fk_user_question_to_question FOREIGN KEY (lecture_sessions_user_question_question_id)
    REFERENCES lecture_sessions_questions (lecture_sessions_question_id) ON DELETE CASCADE
);

-- Indexing untuk efisiensi kueri
CREATE INDEX IF NOT EXISTS idx_user_lecture_sessions_questions_question_id
  ON user_lecture_sessions_questions (lecture_sessions_user_question_question_id);

CREATE INDEX IF NOT EXISTS idx_user_lecture_sessions_questions_created_at
  ON user_lecture_sessions_questions (lecture_sessions_user_question_created_at);
