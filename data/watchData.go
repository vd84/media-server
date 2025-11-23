package data

type WatchData struct {
	UserID    uint `gorm:"primaryKey"`
	EpisodeID uint `gorm:"primaryKey"`

	User    User    `gorm:"constraint:OnDelete:CASCADE;"`
	Episode Episode `gorm:"constraint:OnDelete:CASCADE;"`

	Progress int64
}
