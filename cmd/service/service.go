package service

type Service struct {
	Id         string
	Name       string
	Depends_on []string
}

func (s *Service) IsDependent(id string) bool {
	for i := 0; i < len(s.Depends_on); i++ {
		if s.Depends_on[i].Id == id {
			return true
		}
	}

	return false
}
