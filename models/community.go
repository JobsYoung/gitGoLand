package models

import "time"

type Community struct {
	Id   int64  `json:"id,string" db:"community_id"`
	Name string `json:"community_name" db:"community_name"`
}

type CommunityDetail struct {
	Id           int64     `json:"id,string" db:"community_id"`
	Name         string    `json:"community_name" db:"community_name"`
	Introduction string    `json:"introduction" db:"introduction"`
	CreateTime   time.Time `json:"createTime" db:"create_time"`
}
