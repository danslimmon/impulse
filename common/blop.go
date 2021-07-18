package common

type Blip struct {
	Text     string
	Children []*Blip
}

type Blop struct {
	Root *Blip
}

func (b *Blop) Name() string {
	return b.Root.Text
}
