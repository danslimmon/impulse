package common

type Blip struct {
	Title string
}

func (blip *Blip) GetTitle() string {
	return blip.Title
}

func (blip *Blip) SetTitle(s string) {
	blip.Title = s
}

type Blop struct {
	Title string
}

func (blop *Blop) GetTitle() string {
	return blop.Title
}

func (blop *Blop) SetTitle(s string) {
	blop.Title = s
}
