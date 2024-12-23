// user.go
package repo

// User represents the user table schema
type User struct {
	TlpUsername string `gorm:"column:tlp_username"`
	DiscordID   string `gorm:"column:discord_id;primaryKey"`
	DiscordName string `gorm:"column:discord_name"`
	IsApprover  bool   `gorm:"column:is_approver"`
}

// TableName overrides the table name
func (User) TableName() string {
	return "user"
}

// GetApprovers returns all users with is_approver = true
func GetApprovers() ([]User, error) {
	var approvers []User
	result := DB.Where("is_approver = ?", true).Find(&approvers)
	if result.Error != nil {
		return nil, result.Error
	}
	return approvers, nil
}

// GetUserByTlpUsername returns a user by their TLP username
func GetUserByTlpUsername(tlpUsername string) (*User, error) {
	var user User
	result := DB.Where("tlp_username = ?", tlpUsername).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
