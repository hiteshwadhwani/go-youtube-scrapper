package youtubeservice

type Service interface {
	Get()
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{
		repository: repository,
	}
}

func (y *service) Get() {
	y.repository.Get()
}
