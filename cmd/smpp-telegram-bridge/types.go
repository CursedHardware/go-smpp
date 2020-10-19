package main

type Configuration struct {
	Token          string  `json:"token"`
	AllowedUserIDs []int64 `json:"allowed_user_ids"`
}

func (c Configuration) isAllowedUserID(id int64) bool {
	for _, userId := range c.AllowedUserIDs {
		if userId == id {
			return true
		}
	}
	return false
}
