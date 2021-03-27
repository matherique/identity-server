package service

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
