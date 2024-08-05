package repository

import "e-commerce/internal/models"

// Save Blacklist token in the blacklistToken collection
func (p *Postgres) BlacklistToken(token *models.BlacklistTokens) error {
	if err := p.DB.Create(token).Error; err != nil {
		return err
	}
	return nil
}

// TokenInBlacklist checks if token is already in the blacklist collection
func (p *Postgres) TokenInBlacklist(token *string) bool {
	blacklistToken := &models.BlacklistTokens{}
	if err := p.DB.Where("token = ?", token).First(&blacklistToken).Error; err != nil {
		return false
	}
	return true
}
