package models

type Source struct {
	Id        int
	MediaName string
	Url       string
}

func (src Source) SourceIdNewsApi() int64 {
	return 1
}

func (src Source) SourceIdReddit() int64 {
	return 2
}
