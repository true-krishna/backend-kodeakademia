package domain

type User struct {
    ID        int64  `json:"id"`
    Email     string `json:"email"`
    Name      string `json:"name"`
    AvatarURL string `json:"avatar_url"`
    GoogleSub string `json:"google_sub"`
    Role      string `json:"role"`
}
