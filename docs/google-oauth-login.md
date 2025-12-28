# Google OAuth Login Flow in Kodeakademia Backend

## Overview
This document explains how the Google OAuth login feature works in the Kodeakademia backend. It covers the authentication flow, user persistence, JWT issuance, and protected endpoint access.

---

## 1. OAuth Login Flow

1. **User Initiates Login**
   - The user visits `/login` or triggers the login flow from the frontend.
   - The backend redirects the user to Google's OAuth 2.0 authorization endpoint.

2. **Google Authentication**
   - The user authenticates with Google and grants requested permissions.
   - Google redirects the user back to the backend's `/auth/google/callback` endpoint with an authorization code.

3. **Token Exchange**
   - The backend exchanges the authorization code for access and ID tokens using Google's token endpoint.

4. **User Info Fetch**
   - The backend fetches the user's profile from Google's userinfo endpoint using the access token.

---

## 2. User Persistence Logic

- The backend checks the database for a user with the Google `sub` (unique subject ID).
- **If the user does not exist:**
  - A new user record is created in the database with information from Google (email, name, avatar, etc.).
- **If the user exists:**
  - The existing user record is used.
- This ensures each Google account maps to a single user in your system.

---

## 3. JWT Issuance

- After user lookup/creation, the backend generates a signed JWT (JSON Web Token) containing:
  - `uid`: The user's database ID
  - `email`, `name`, `picture`: User info
  - Standard claims: `iat` (issued at), `exp` (expiration), etc.
- The JWT is signed using the `JWT_SECRET` from environment variables.
- The JWT is returned to the client for use in authenticated requests.

---

## 4. Protected Endpoints

- Endpoints like `/me` are protected by JWT middleware.
- The client must send the JWT in the `Authorization: Bearer <token>` header.
- The middleware verifies the token and attaches user claims to the request context.
- Handlers can access user info from the context.

---

## 5. Security Notes

- Always use HTTPS in production.
- Store `JWT_SECRET` securely and never commit it to version control.
- Validate all tokens and handle errors gracefully.
- Consider adding CSRF protection (OAuth `state` parameter) and refresh token support for production.

---

## 6. Example .env Configuration

```
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
OAUTH_REDIRECT_URI=http://localhost:8080/auth/google/callback
OAUTH_SCOPES=openid email profile
JWT_SECRET=your-very-strong-secret
```

---

## 7. Example API Usage

- **Login:**
  - Open `/login` in a browser to start the flow.
- **Get JWT:**
  - After Google login, the backend returns `{ "token": "<jwt>", "user": { ... } }`.
- **Access Protected Endpoint:**
  - Send `Authorization: Bearer <jwt>` header to `/me` to get user info.

---

## 8. References
- [Google OAuth 2.0 Documentation](https://developers.google.com/identity/protocols/oauth2)
- [JWT Introduction](https://jwt.io/introduction)

---

For further details, see the code in `cmd/server/main.go` and the user usecase/repository logic.
