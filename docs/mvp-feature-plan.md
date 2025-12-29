# MVP Feature Plan (Phase 1)

This plan details all features to be developed for the MVP, based on the finalized entity model. Each CRUD operation is a separate feature. Admin CRUD and soft-delete for Course and Lesson are mandatory and prioritized.

---

## 1. Admin Features (Priority)

### Course Management
- Create Course
- Read/List Courses
- Update Course
- Delete (Soft-delete) Course

### Section Management
- Create Section
- Read/List Sections
- Update Section
- Delete Section

### Lesson Management
- Create Lesson
- Read/List Lessons
- Update Lesson
- Delete (Soft-delete) Lesson

---

## 2. User Features

### Registration & Authentication
- Google OAuth Login
- User Profile (Read/Update)

### Course Browsing & Enrollment
- List Published Courses
- Enroll in Course (Create Enrollment)
- View Enrollments

### Lesson Progress
- Mark Lesson as Complete (Create/Update LessonProgress)
- Track Last Watched Position

### Orders & Payments
- Create Order (Purchase Course)
- View Orders

### Coupons
- Redeem Coupon (Create CouponUsage)
- View Available/Used Coupons

### Certificates
- Issue Certificate (on course completion)
- View Certificates

---

## 3. Cross-Entity Features
- Enrollment: User enrolls in Course
- Coupon Redemption: User applies Coupon to Order
- Certificate Issuance: System issues Certificate on course completion

---

## 4. Backend Components per Feature
For each feature:
- Migration (if schema change needed)
- Repository (data access)
- Usecase (business logic)
- Handler (HTTP endpoint)
- DTO (request/response)

---

## 5. Soft-Delete Implementation
- Course: Add `deleted_at` timestamp, filter out soft-deleted in queries
- Lesson: Add `deleted_at` timestamp, filter out soft-deleted in queries

---

## 6. Tracking & Documentation
- All features and progress to be tracked in this plan file
- Update as features are completed or requirements change

---

*Last updated: 2025-12-29*
