CREATE TABLE user_module_attempts (
  user_module_attempt_id SERIAL PRIMARY KEY,
  user_module_attempt_user_id UUID NOT NULL,
  user_module_attempt_target_type INTEGER NOT NULL,        -- 1 = reading, 2 = quiz, 3 = evaluation, 4 = exam
  user_module_attempt_target_id INTEGER NOT NULL,
  user_module_attempt_percentage_grade INTEGER,            -- nilai akhir (0-100)
  user_module_attempt_time_duration INTEGER,               -- durasi dalam detik
  user_module_attempt_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  user_module_attempt_batch_id VARCHAR(100),              -- ✅ Baru
  user_module_attempt_submitted_at TIMESTAMP    
);


CREATE INDEX idx_user_module_lookup
  ON user_module_attempts (
    user_module_attempt_user_id,
    user_module_attempt_target_type,
    user_module_attempt_target_id
  );

CREATE INDEX idx_user_module_batch ON user_module_attempts(user_module_attempt_batch_id);


CREATE TABLE user_answer_attempts (
  user_answer_attempt_id SERIAL PRIMARY KEY,
  user_answer_attempt_user_id UUID NOT NULL,
  user_answer_attempt_target_type INTEGER NOT NULL,         -- 1=reading, 2=quiz, 3=evaluation, 4=exam
  user_answer_attempt_target_id INTEGER NOT NULL,
  user_answer_attempt_question_id INTEGER NOT NULL,
  user_answer_attempt_answer VARCHAR(1) NOT NULL,
  user_answer_attempt_is_correct BOOLEAN DEFAULT FALSE,
  user_answer_attempt_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  user_answer_attempt_batch_id VARCHAR(100),                -- ✅ Tambahkan ini
  user_answer_attempt_submitted_at TIMESTAMP                -- ✅ Tambahkan ini
);


CREATE INDEX idx_user_answer_lookup
  ON user_answer_attempts (
    user_answer_attempt_user_id,
    user_answer_attempt_target_type,
    user_answer_attempt_target_id,
    user_answer_attempt_question_id
  );


CREATE INDEX idx_user_answer_batch ON user_answer_attempts(user_answer_attempt_batch_id);