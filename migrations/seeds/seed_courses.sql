-- seeds/seed_courses.sql
-- Insert seed roles, users (owner), courses, purchases and progress for development/testing

INSERT INTO roles (name) VALUES ('owner') ON CONFLICT (name) DO NOTHING;
INSERT INTO roles (name) VALUES ('learner') ON CONFLICT (name) DO NOTHING;

-- Create an owner user
WITH owner AS (
  INSERT INTO users (email, name, avatar_url, google_sub, role_id)
  VALUES ('owner@example.com', 'Owner User', NULL, 'owner-sub-1', (SELECT id FROM roles WHERE name='owner'))
  ON CONFLICT (email) DO UPDATE SET name = EXCLUDED.name
  RETURNING id
)
INSERT INTO courses (owner_id, title, description, category, price)
SELECT owner.id, 'Intro to Go', 'Learn Go from scratch', 'programming', 29.99 FROM owner;

INSERT INTO courses (owner_id, title, description, category, price)
SELECT owner.id, 'Advanced Go Concurrency', 'Concurrency patterns', 'programming', 49.99 FROM owner;

-- Create a sample learner and purchases
WITH learner AS (
  INSERT INTO users (email, name, avatar_url, google_sub, role_id)
  VALUES ('student1@example.com', 'Student One', NULL, 'student-sub-1', (SELECT id FROM roles WHERE name='learner'))
  ON CONFLICT (email) DO UPDATE SET name = EXCLUDED.name
  RETURNING id
), course1 AS (
  SELECT id FROM courses WHERE title = 'Intro to Go' LIMIT 1
), course2 AS (
  SELECT id FROM courses WHERE title = 'Advanced Go Concurrency' LIMIT 1
)
INSERT INTO purchases (user_id, course_id, amount)
SELECT learner.id, course1.id, 29.99 FROM learner, course1
ON CONFLICT (user_id, course_id) DO NOTHING;

INSERT INTO course_progress (user_id, course_id, completed)
SELECT learner.id, course1.id, true FROM learner, course1
ON CONFLICT (user_id, course_id) DO UPDATE SET completed = EXCLUDED.completed;
