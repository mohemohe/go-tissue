package go_tissue

import "time"

// NumericString は JSON の文字列・数値のどちらとしてエンコードされていても
// 受け入れ可能な数値表現。値は文字列として保持される。
type NumericString string

func (n *NumericString) UnmarshalJSON(data []byte) error {
	s := string(data)
	if s == "null" {
		*n = ""
		return nil
	}
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		*n = NumericString(s[1 : len(s)-1])
		return nil
	}
	*n = NumericString(s)
	return nil
}

func (n NumericString) MarshalJSON() ([]byte, error) {
	if n == "" {
		return []byte("null"), nil
	}
	return []byte(n), nil
}

type User struct {
	ID                  int64  `json:"id"`
	Name                string `json:"name"`
	DisplayName         string `json:"display_name"`
	IsProtected         bool   `json:"is_protected"`
	PrivateLikes        bool   `json:"private_likes"`
	ProfileImageURL     string `json:"profile_image_url"`
	ProfileMiniImageURL string `json:"profile_mini_image_url"`
	Bio                 string `json:"bio"`
	URL                 string `json:"url"`
}

type CheckinSummary struct {
	CurrentSessionElapsed int64   `json:"current_session_elapsed"`
	TotalCheckins         int64   `json:"total_checkins"`
	TotalTimes            int64   `json:"total_times"`
	AverageInterval       float64 `json:"average_interval"`
	MedianInterval        NumericString `json:"median_interval"`
	LongestInterval       int64   `json:"longest_interval"`
	ShortestInterval      int64   `json:"shortest_interval"`
}

type Me struct {
	User
	CheckinSummary CheckinSummary `json:"checkin_summary"`
}

type Information struct {
	ID        int64     `json:"id"`
	Category  string    `json:"category"`
	Pinned    bool      `json:"pinned"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type DailyCheckinCount struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type TagCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type Checkin struct {
	ID                 int64     `json:"id"`
	CheckedInAt        time.Time `json:"checked_in_at"`
	Note               string    `json:"note"`
	Link               string    `json:"link"`
	Tags               []string  `json:"tags"`
	Source             string    `json:"source"`
	IsPrivate          bool      `json:"is_private"`
	IsTooSensitive     bool      `json:"is_too_sensitive"`
	DiscardElapsedTime bool      `json:"discard_elapsed_time"`
	User               User      `json:"user"`
	IsLiked            int       `json:"is_liked,omitempty"`
	LikesCount         int       `json:"likes_count,omitempty"`
	IsMuted            int       `json:"is_muted,omitempty"`
}

type UserCheckin struct {
	Checkin
	CheckinInterval     int64  `json:"checkin_interval"`
	PreviousCheckedInAt string `json:"previous_checked_in_at"`
}

type Collection struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	UserName  string    `json:"user_name"`
	User      User      `json:"user"`
	Title     string    `json:"title"`
	IsPrivate bool      `json:"is_private"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CollectionItem struct {
	ID           int64      `json:"id"`
	CollectionID int64      `json:"collection_id"`
	Collection   Collection `json:"collection"`
	UserID       int64      `json:"user_id"`
	UserName     string     `json:"user_name"`
	Link         string     `json:"link"`
	Note         string     `json:"note"`
	Tags         []string   `json:"tags"`
}
