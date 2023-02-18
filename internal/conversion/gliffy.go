package conversion

import (
	internal "diagram-converter/internal"
	datastr "diagram-converter/internal/datastructures"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

func ConvertExcalidrawToGliffy(importPath string, exportPath string) error {
	fmt.Printf("Parsing input file: %s\n", importPath)

	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	data, err := os.ReadFile(importPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read file: %s\n", err)
		os.Exit(1)
	}

	var input datastr.ExcalidrawScene
	err = json.Unmarshal(data, &input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse input: %s\n", err)
		os.Exit(1)
	}

	var output datastr.GliffyScene
	var objects []datastr.GliffyObject

	graphics := internal.MapGraphics()

	for i, element := range input.Elements {
		var object datastr.GliffyObject
		var shape datastr.GliffyShape
		var text datastr.GliffyText
		var line datastr.GliffyLine

		object.X = element.X
		object.Y = element.Y
		object.Width = element.Width
		object.Height = element.Height
		object.Rotation = internal.NormalizeRotation(element.Angle)
		object.LayerID = "dR5PnMr9lIuu"
		object.Order = i

		for _, id := range graphics.Rectangle.Excalidraw {
			if element.Type == id && element.Roundness.Type == 0 {
				object.UID = graphics.Rectangle.Gliffy[0]
				object.Graphic.Type = "Shape"
				shape.Tid = "com.gliffy.stencil.rectangle.basic_v1"
			}

			if element.Type == id && element.Roundness.Type > 0 {
				object.UID = graphics.Rectangle.Gliffy[1]
				object.Graphic.Type = "Shape"
				shape.Tid = "com.gliffy.stencil.round_rectangle.basic_v1"
			}
		}

		for _, id := range graphics.Ellipse.Excalidraw {
			if element.Type == id {
				object.UID = graphics.Ellipse.Gliffy[0]
				object.Graphic.Type = "Shape"
				shape.Tid = "com.gliffy.stencil.ellipse.basic_v1"
			}
		}

		for _, id := range graphics.Diamond.Excalidraw {
			if element.Type == id {
				object.UID = graphics.Diamond.Gliffy[0]
				object.Graphic.Type = "Shape"
				shape.Tid = "com.gliffy.stencil.diamond.basic_v1"
			}
		}

		if object.Graphic.Type == "Shape" {
			shape.DashStyle = StrokeStyleConvExcGliffy(element.StrokeStyle)
			shape.FillColor = FillColorConvExcGliffy(element.BackgroundColor)
			shape.StrokeColor = element.StrokeColor
			shape.StrokeWidth = element.StrokeWidth
			shape.Opacity = element.Opacity * 0.01

			if element.FillStyle != "solid" {
				shape.Gradient = true
			}

			object.Graphic.Shape = &shape
		}

		for _, id := range graphics.Text.Excalidraw {
			if element.Type == id {
				object.UID = graphics.Text.Gliffy[0]
				object.Graphic.Type = "Text"

				object.Width = element.Width * 1.2

				fontSize := strconv.FormatFloat(element.FontSize, 'f', 0, 64)
				fontColor := element.StrokeColor
				fontFamily := "Arial"
				if element.FontFamily == 3 {
					fontFamily = "Courier"
				}

				text.HTML = "<p style=\"text-align: " + element.TextAlign + ";\"><span style=\"font-family: " + fontFamily + "; font-size: " + fontSize + "px;\"><span style=\"\"><span style=\"color: " + fontColor + "; font-size: " + fontSize + "px; line-height: 16.5px;\">" + element.Text + "</span><br></span></span></p>"
				text.Valign = "middle"
				text.Overflow = "none"
				text.Vposition = "none"
				text.Hposition = "none"

				object.Graphic.Text = &text
			}
		}

		for _, id := range graphics.Line.Excalidraw {
			if element.Type == id {
				object.UID = graphics.Line.Gliffy[0]
				object.Graphic.Type = "Line"

				line.DashStyle = StrokeStyleConvExcGliffy(element.StrokeStyle)
				line.StrokeColor = element.StrokeColor
				line.StrokeWidth = element.StrokeWidth
				line.FillColor = "none"
				line.StartArrowRotation = "auto"
				line.EndArrowRotation = "auto"
				line.InterpolationType = "linear"
				line.CornerRadius = 10
				line.Ortho = true
				line.ControlPath = element.Points
				line.StartArrow = ArrowheadConvExGliffy(element.StartArrowhead)
				line.EndArrow = ArrowheadConvExGliffy(element.EndArrowhead)

				object.Graphic.Line = &line
			}
		}

		if object.Graphic.Type == "" {
			continue
		}

		fmt.Printf("  Adding object: %s\n", object.UID)

		object.ID = i
		objects = append(objects, object)
	}

	priorityGraphics := []string{
		//"Line",
		//"Text",
	}

	objects = OrderGliffyObjectsByPriority(objects, priorityGraphics)

	var layer datastr.GliffyLayer
	layer.Active = true
	layer.GUID = "dR5PnMr9lIuu"
	layer.Name = "Layer 0"
	layer.NodeIndex = 11
	layer.Visible = true

	output.ContentType = "application/gliffy+json"
	output.EmbeddedResources.Resources = []string{}
	output.Version = "1.3"
	output.Metadata.LastSerialized = timestamp
	output.Metadata.Libraries = []string{
		"com.gliffy.libraries.basic.basic_v1.default",
		"com.gliffy.libraries.flowchart.flowchart_v1.default",
	}
	output.Metadata.LoadPosition = "default"
	output.Metadata.Revision = 0
	output.Metadata.Title = "Import"
	output.Stage.Background = input.AppState.ViewBackgroundColor
	output.Stage.DrawingGuidesOn = true
	output.Stage.GridOn = true
	output.Stage.Height = 1024
	output.Stage.Layers = append(output.Stage.Layers, layer)
	output.Stage.MaxWidth = 5000
	output.Stage.MaxHeight = 5000
	output.Stage.Objects = objects
	output.Stage.PrintModel.PageSize = "Letter"
	output.Stage.PrintModel.Portrait = true
	output.Stage.SnapToGrid = true
	output.Stage.ViewportType = "default"
	output.Stage.Width = 1024

	outputJson, err := json.Marshal(output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occured during JSON marshaling: %s", err)
		os.Exit(1)
	}

	err = internal.WriteToFile(exportPath, string(outputJson))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to write diagram to file: %s", err)
		os.Exit(1)
	}

	fmt.Printf("Converted diagram saved to file: %s\n", exportPath)

	return nil
}

func StrokeStyleConvExcGliffy(style string) string {
	switch style {
	case "dashed":
		style = "8,8"
	case "dotted":
		style = "2,2"
	default:
		return ""
	}

	return style
}

func FillColorConvExcGliffy(color string) string {
	switch color {
	case "transparent":
		color = "none"
	}

	return color
}

func ArrowheadConvExGliffy(head string) int {
	arrowHead := 0

	switch head {
	case "arrow":
		arrowHead = 1
	case "dot":
		arrowHead = 2
	}

	return arrowHead
}

func OrderGliffyObjectsByPriority(objects []datastr.GliffyObject, prioritized []string) []datastr.GliffyObject {
	prioritySlot := len(objects)

	for i, object := range objects {
		prioritySlot++

		for _, t := range prioritized {
			if object.Graphic.Type == t {
				objects[i].Order = prioritySlot
			}
		}
	}

	return objects
}
