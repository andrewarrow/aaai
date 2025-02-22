package diff

type Change struct {
	File   string  `json:"file"`
	Ranges []Range `json:"ranges"`
}

type Range struct {
	Start  int      `json:"s"`
	End    int      `json:"e"`
	Before []string `json:"b"`
	After  []string `json:"a"`
}
