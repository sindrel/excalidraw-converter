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
}

func MapGraphics() GraphicTypes {
	var graphics GraphicTypes

	graphics.Rectangle.Gliffy = []string{
		"com.gliffy.shape.basic.basic_v1.default.rectangle",
	}
	graphics.Rectangle.Excalidraw = []string{
		"rectangle",
	}

	graphics.Ellipse.Gliffy = []string{
		"com.gliffy.shape.basic.basic_v1.default.ellipse",
	}
	graphics.Ellipse.Excalidraw = []string{
		"ellipse",
	}

	graphics.Diamond.Gliffy = []string{
		"com.gliffy.shape.flowchart.flowchart_v1.default.decision",
	}
	graphics.Diamond.Excalidraw = []string{
		"diamond",
	}

	graphics.Text.Gliffy = []string{
		"com.gliffy.shape.basic.basic_v1.default.text",
	}
	graphics.Text.Excalidraw = []string{
		"text",
	}

	graphics.Line.Gliffy = []string{
		"com.gliffy.shape.basic.basic_v1.default.line",
	}
	graphics.Line.Excalidraw = []string{
		"line",
		"arrow",
	}

	return graphics
}
