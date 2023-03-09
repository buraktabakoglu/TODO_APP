package models

import "time"

type Activation_links struct {
    ID        int    `gorm:"primary_key;auto_increment" json:"id"`
    Token     string `gorm:"not null" json:"token"`
    Is_used   bool   `gorm:"default:false" json:"is_used"`
    UserID    int    `gorm:"not null" json:"user_id"`
    CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}
