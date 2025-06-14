CREATE TABLE lecture_sessions_questions (
  lecture_sessions_question_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  lecture_sessions_question TEXT NOT NULL,
  lecture_sessions_question_answer TEXT NOT NULL,
  lecture_sessions_question_correct CHAR(1) NOT NULL CHECK (lecture_sessions_question_correct IN ('A', 'B', 'C', 'D')),
  lecture_sessions_question_explanation TEXT,
  lecture_sessions_question_lecture_session_id UUID REFERENCES lecture_sessions(lecture_session_id) ON DELETE SET NULL,
  lecture_sessions_question_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexing
CREATE INDEX idx_lecture_sessions_questions_lecture_session_id ON lecture_sessions_questions(lecture_sessions_question_lecture_session_id);
CREATE INDEX idx_lecture_sessions_questions_created_at ON lecture_sessions_questions(lecture_sessions_question_created_at);


CREATE TABLE lecture_sessions_question_links (
  lecture_sessions_question_link_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  question_id UUID NOT NULL REFERENCES lecture_sessions_questions(lecture_sessions_question_id) ON DELETE CASCADE,
  exam_id UUID REFERENCES lecture_sessions_exams(lecture_sessions_exam_id) ON DELETE CASCADE,
  quiz_id UUID REFERENCES lecture_sessions_quiz(lecture_sessions_quiz_id) ON DELETE CASCADE,
  question_order INTEGER,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  CHECK (exam_id IS NOT NULL OR quiz_id IS NOT NULL)
);

-- Indexing
CREATE INDEX idx_question_links_question_id ON lecture_sessions_question_links(question_id);
CREATE INDEX idx_question_links_exam_id ON lecture_sessions_question_links(exam_id);
CREATE INDEX idx_question_links_quiz_id ON lecture_sessions_question_links(quiz_id);



CREATE TABLE lecture_sessions_user_questions (
  lecture_sessions_user_question_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  lecture_sessions_user_question_answer CHAR(1) NOT NULL CHECK (lecture_sessions_user_question_answer IN ('A', 'B', 'C', 'D')),
  lecture_sessions_user_question_is_correct BOOLEAN NOT NULL,
  lecture_sessions_user_question_question_id UUID NOT NULL REFERENCES lecture_sessions_questions(lecture_sessions_question_id) ON DELETE CASCADE,
  lecture_sessions_user_question_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexing (untuk efisiensi kueri per user/question)
CREATE INDEX idx_lecture_sessions_user_questions_question_id ON lecture_sessions_user_questions(lecture_sessions_user_question_question_id);
CREATE INDEX idx_lecture_sessions_user_questions_created_at ON lecture_sessions_user_questions(lecture_sessions_user_question_created_at);