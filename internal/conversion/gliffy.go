package conversion

import (
	internal "diagram-converter/internal"
	datastr "diagram-converter/internal/datastructures"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

var xOffset float64
var yOffset float64

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

	xOffset, yOffset = GetXYOffset(input)

	var output datastr.GliffyScene
	var objects []datastr.GliffyObject
	objectIDs := map[string]int{}

	objects, objectIDs, err = AddElements(false, input, objects, objectIDs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to add element(s): %s\n", err)
		os.Exit(1)
	}

	objects, _, err = AddElements(true, input, objects, objectIDs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to add element(s) with parent(s): %s\n", err)
		os.Exit(1)
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

func AddElements(addChildren bool, input datastr.ExcalidrawScene, objects []datastr.GliffyObject, objectIDs map[string]int) ([]datastr.GliffyObject, map[string]int, error) {
	graphics := internal.MapGraphics()

	for i, element := range input.Elements {
		if len(element.ContainerId) > 0 && !addChildren {
			continue
		}
		if len(element.ContainerId) == 0 && addChildren {
			continue
		}

		var object datastr.GliffyObject
		var shape datastr.GliffyShape
		var text datastr.GliffyText
		var line datastr.GliffyLine
		var image datastr.GliffyImage

		object.X = element.X - xOffset
		object.Y = element.Y - yOffset
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
			shape.StrokeWidth = int64(math.Round(element.StrokeWidth))
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

				element.Text = strings.ReplaceAll(element.Text, "\n", "<br>")

				text.HTML = "<p style=\"text-align: " + element.TextAlign + ";\"><span style=\"font-family: " + fontFamily + "; font-size: " + fontSize + "px;\"><span style=\"\"><span style=\"color: " + fontColor + "; font-size: " + fontSize + "px; line-height: 16.5px;\">" + element.Text + "</span><br></span></span></p>"
				text.Valign = element.VerticalAlign
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
				line.StrokeWidth = int64(math.Round(element.StrokeWidth))
				line.FillColor = "none"
				line.StartArrowRotation = "auto"
				line.EndArrowRotation = "auto"
				line.InterpolationType = "linear"
				line.CornerRadius = 10
				line.Ortho = true
				line.ControlPath = element.Points
				line.StartArrow = ArrowheadConvExcGliffy(element.StartArrowhead)
				line.EndArrow = ArrowheadConvExcGliffy(element.EndArrowhead)

				object.Graphic.Line = &line
			}
		}

		for _, id := range graphics.Image.Excalidraw {
			if element.Type == id {
				object.UID = graphics.Image.Gliffy[0]
				object.Graphic.Type = "Image"

				dataUrl, err := EmbeddedImgConvExcGliffy(input, element.FileId)
				if err != nil {
					return nil, nil, err
				}

				image.Url = dataUrl
				image.StrokeColor = FillColorConvExcGliffy(element.StrokeColor)
				image.StrokeWidth = int64(math.Round(element.StrokeWidth))

				object.Graphic.Image = &image
			}
		}

		if object.Graphic.Type == "" {
			continue
		}

		object.ID = i
		objectIDs[element.ID] = object.ID

		fmt.Printf("  Adding object: %s (%s,%d,%d)\n", object.UID, element.ID, object.ID, object.Order)

		if len(element.ContainerId) > 0 {
			var parent int = 999999
			for obj_k, obj := range objects {
				if obj.ID == objectIDs[element.ContainerId] {
					parent = obj_k
				}
			}

			if parent == 999999 {
				return nil, nil, errors.New("unable to find object parent")
			}

			object.X = 2
			object.Y = 0
			object.Rotation = 0
			object.UID = ""
			object.Width = objects[parent].Width - 4
			object.Height = objects[parent].Height - 4

			fmt.Printf("  - Adding as child of %d\n", parent)

			children := append(objects[parent].Children, object)
			objects[parent].Children = children

			continue
		}

		objects = append(objects, object)
	}

	return objects, objectIDs, nil
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

func ArrowheadConvExcGliffy(head string) int {
	arrowHead := 0

	switch head {
	case "arrow":
		arrowHead = 1
	case "dot":
		arrowHead = 2
	}

	return arrowHead
}

func EmbeddedImgConvExcGliffy(input datastr.ExcalidrawScene, fileId string) (string, error) {
	file, ok := input.Files[fileId]
	if !ok {
		return "", fmt.Errorf("unable to find embedded file with id %s", fileId)
	}

	return file.DataUrl, nil
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

func GetXYOffset(input datastr.ExcalidrawScene) (float64, float64) {
	var xMin float64 = 0
	var yMin float64 = 0

	for _, element := range input.Elements {
		if element.X > xMin {
			xMin = element.X
		}

		if element.Y > yMin {
			yMin = element.Y
		}
	}

	for _, element := range input.Elements {
		if element.X < xMin {
			xMin = element.X
		}

		if element.Y < yMin {
			yMin = element.Y
		}
	}

	fmt.Printf("  Offset X: %f, Offset Y: %f\n", xMin, yMin)

	return xMin, yMin
}
