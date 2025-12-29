# Domain Entities in Kodeakademia Backend

This document describes the main domain entities used in the Kodeakademia backend, their purpose, and key fields. These entities form the core of the business logic and database schema.

---

## 1. User

Represents a learner or platform owner.

**Fields:**
- `id`: Unique user ID (primary key)
- `email`: User's email address (unique)
- `google_sub`: Google OAuth subject (unique, nullable)
- `name`: Full name
- `avatar_url`: Profile picture URL
- `role_id`: Foreign key to roles table
- `created_at`: Timestamp

**Purpose:**  
Tracks authentication, authorization, and user profile information.

---

## 2. Course

Represents a course created by the owner.

**Fields:**
- `id`: Unique course ID (primary key)
- `title`: Course title
- `category`: Course category (e.g., Programming)
- `owner_id`: Reference to the user who owns the course
- `price`: Course price (in smallest currency unit)
- `created_at`, `updated_at`: Timestamps

**Purpose:**  
Defines the learning content available for purchase and completion.

---

## 3. Purchase

Represents a learner's purchase of a course.

**Fields:**
- `id`: Unique purchase ID (primary key)
- `user_id`: Reference to the purchasing user
- `course_id`: Reference to the purchased course
- `amount`: Purchase amount
- `purchased_at`: Timestamp

**Purpose:**  
Tracks which users have access to which courses.

---

## 4. CourseProgress

Represents a learner's progress in a course.

**Fields:**
- `id`: Unique progress ID (primary key)
- `user_id`: Reference to the learner
- `course_id`: Reference to the course
- `completed`: Boolean flag for completion
- `updated_at`: Timestamp

**Purpose:**  
Tracks course completion status for each learner.

---

## 5. Role

Represents user roles for authorization.

**Fields:**
- `id`: Unique role ID (primary key)
- `name`: Role name (e.g., `owner`, `learner`)

**Purpose:**  
Defines access levels and permissions in the system.

---

## Notes

- All entities use `created_at` and/or `updated_at` for auditability.
- Foreign keys (e.g., `owner_id`, `user_id`, `course_id`, `role_id`) enforce relationships.
- See the `migrations/` directory for exact SQL schema definitions.

---

For more details, refer to the Go structs in `internal/domain/entity/` and the repository interfaces in `internal/domain/repository/`.
