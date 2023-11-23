package internal

import (
	datastr "diagram-converter/internal/datastructures"
)

func MapGraphics() datastr.GraphicTypes {
	var graphics datastr.GraphicTypes

	graphics.Rectangle.Gliffy = []string{
		"com.gliffy.shape.basic.basic_v1.default.rectangle",
		"com.gliffy.shape.basic.basic_v1.default.round_rectangle",
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

	graphics.Image.Gliffy = []string{
		"com.gliffy.shape.basic.basic_v1.default.image",
	}
	graphics.Image.Excalidraw = []string{
		"image",
	}

	return graphics
}
