-- 0001_init.down.sql
-- Drops schema created in 0001_init.up.sql (reverse order)

DROP INDEX IF EXISTS idx_users_google_sub;
DROP INDEX IF EXISTS idx_courses_created_at;
DROP INDEX IF EXISTS idx_courses_category;

DROP TABLE IF EXISTS course_progress;
DROP TABLE IF EXISTS purchases;
DROP TABLE IF EXISTS courses;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;
