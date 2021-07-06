package main

type Configuration struct {
	Token          string  `json:"token"`
	AllowedUserIDs []int `json:"allowed_user_ids"`
}

func (c Configuration) isAllowedUserID(id int) bool {
	for _, userId := range c.AllowedUserIDs {
		if userId == id {
			return true
		}
	}
	return false
}
