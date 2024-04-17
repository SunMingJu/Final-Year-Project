package common

//PageInfo 分页
type PageInfo struct {
	Page    int    `json:"page" form:"page"`                 // pagination
	Size    int    `json:"size" form:"size"`                 // page size
	Keyword string `json:"keyword,omitempty" form:"keyword"` //keywords
}

func (p *PageInfo) Init() {
	if p.Size == 0 {
		p.Size = 20
	}
	if p.Page == 0 {
		p.Page = 1
	}
}
