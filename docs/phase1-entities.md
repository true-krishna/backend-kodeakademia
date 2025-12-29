# Phase 1 MVP Entities

This document summarizes all core entities for the KodeAkademia MVP, based on the detailed requirements in `docs/phase1.txt`.

---

## 1. User
- **id**: UUID (PK)
- **name**: string
- **email**: string (unique)
- **avatar_url**: string (nullable)
- **created_at**: timestamp
- **updated_at**: timestamp
- **google_id**: string (nullable, for Google OAuth)

---

## 2. Course
- **id**: UUID (PK)
- **title**: string
- **slug**: string (unique)
- **description**: text
- **cover_image_url**: string (nullable)
- **price**: integer (in cents)
- **is_published**: boolean
- **created_at**: timestamp
- **updated_at**: timestamp

---

## 3. Section
- **id**: UUID (PK)
- **course_id**: UUID (FK → Course)
- **title**: string
- **order**: integer (for ordering sections)

---

## 4. Lesson
- **id**: UUID (PK)
- **section_id**: UUID (FK → Section)
- **title**: string
- **video_url**: string
- **order**: integer (for ordering lessons)
- **duration_seconds**: integer

---

## 5. Enrollment
- **id**: UUID (PK)
- **user_id**: UUID (FK → User)
- **course_id**: UUID (FK → Course)
- **enrolled_at**: timestamp
- **is_active**: boolean

---

## 6. LessonProgress
- **id**: UUID (PK)
- **user_id**: UUID (FK → User)
- **lesson_id**: UUID (FK → Lesson)
- **completed_at**: timestamp (nullable)
- **last_watched_seconds**: integer (nullable)

---

## 7. Order
- **id**: UUID (PK)
- **user_id**: UUID (FK → User)
- **course_id**: UUID (FK → Course)
- **amount_paid**: integer (in cents)
- **status**: enum (pending, paid, failed, refunded)
- **created_at**: timestamp
- **payment_reference**: string (nullable)

---

## 8. Coupon
- **id**: UUID (PK)
- **code**: string (unique)
- **discount_percent**: integer
- **max_uses**: integer
- **expires_at**: timestamp (nullable)
- **created_at**: timestamp

---

## 9. CouponUsage
- **id**: UUID (PK)
- **coupon_id**: UUID (FK → Coupon)
- **user_id**: UUID (FK → User)
- **used_at**: timestamp

---

## 10. Certificate
- **id**: UUID (PK)
- **user_id**: UUID (FK → User)
- **course_id**: UUID (FK → Course)
- **issued_at**: timestamp
- **certificate_url**: string

---

## 11. UserActivity
- **id**: UUID (PK)
- **user_id**: UUID (FK → User)
- **activity_type**: string (e.g., login, enroll, complete_lesson)
- **activity_data**: jsonb (nullable)
- **created_at**: timestamp

---

## Relationships Overview
- **User** can enroll in many **Courses** (via **Enrollment**)
- **Course** has many **Sections**; **Section** has many **Lessons**
- **User** can track progress in **Lessons** (via **LessonProgress**)
- **User** can place **Orders** for **Courses**
- **Coupon** can be used by many **Users** (via **CouponUsage**)
- **User** can earn **Certificates** for **Courses**
- **UserActivity** logs all key user actions

---

## Notes
- All timestamps are in UTC
- All UUIDs are v4
- All monetary values are in cents (integer)
- All foreign keys are indexed
- All entities have `created_at` and `updated_at` where relevant

---

*Generated from `docs/phase1.txt` on 2025-12-29.*
