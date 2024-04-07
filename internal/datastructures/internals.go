package internal

type GraphicTypes struct {
	Rectangle struct {
		Excalidraw []string `json:"excalidraw"`
		Gliffy     []string `json:"gliffy"`
	} `json:"rectangle"`
	Text struct {
		Excalidraw []string `json:"excalidraw"`
		Gliffy     []string `json:"gliffy"`
	} `json:"text"`
	Line struct {
		Excalidraw []string `json:"excalidraw"`
		Gliffy     []string `json:"gliffy"`
	} `json:"line"`
	Ellipse struct {
		Excalidraw []string `json:"excalidraw"`
		Gliffy     []string `json:"gliffy"`
	} `json:"ellipse"`
	Diamond struct {
		Excalidraw []string `json:"excalidraw"`
		Gliffy     []string `json:"gliffy"`
	} `json:"diamond"`
	Image struct {
		Excalidraw []string `json:"excalidraw"`
		Gliffy     []string `json:"gliffy"`
	} `json:"image"`
	Freedraw struct {
		Excalidraw []string `json:"excalidraw"`
		Gliffy     []string `json:"gliffy"`
	} `json:"freedraw"`
}

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

type ElementSizeOffset struct {
	Width  float64
	Height float64
}

type ElementPositionOffset struct {
	X float64
	Y float64
}
