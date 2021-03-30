package service

import (
	"github.com/matherique/identity-service/lib/utils"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Service struct {
	Id         string
	Name       string
	Depends_on []string
	Secret     []byte
}

func (s *Service) IsDependent(id string) bool {
	for i := 0; i < len(s.Depends_on); i++ {
		if s.Depends_on[i] == id {
			return true
		}
	}

	return false
}

func (s *Service) GenerateToken(target *Service) (interface{}, error) {
	return utils.GenerateTokens(target.Id, target.Secret)
}
