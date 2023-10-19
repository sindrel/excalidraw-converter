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
	Freedraw struct {
		Excalidraw []string `json:"excalidraw"`
		Gliffy     []string `json:"gliffy"`
	} `json:"freedraw"`
}
