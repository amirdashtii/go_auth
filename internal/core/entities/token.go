package entities

// TokenPair ساختاری برای نگهداری توکن‌های دسترسی و رفرش
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
