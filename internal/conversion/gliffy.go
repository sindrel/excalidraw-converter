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

func ConvertExcalidrawToGliffyFile(importPath string, exportPath string) error {
	fmt.Printf("Parsing input file: %s\n", importPath)

	data, err := os.ReadFile(importPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "File reading failed. %s\n", err)
		os.Exit(1)
	}

	output, err := ConvertExcalidrawToGliffy(string(data))
	if err != nil {
		fmt.Fprintf(os.Stderr, "File parsing failed. %s\n", err)
		os.Exit(1)
	}

	err = internal.WriteToFile(exportPath, string(output))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Saving diagram failed. %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Converted diagram saved to file: %s\n", exportPath)

	return nil
}

func ConvertExcalidrawToGliffy(data string) (string, error) {
	fmt.Printf("Converting to Gliffy format...\n")

	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	var input datastr.ExcalidrawScene
	err := json.Unmarshal([]byte(data), &input)
	if err != nil {
		return "", errors.New("Unable to parse input: " + err.Error())
	}

	xOffset, yOffset = GetXYOffset(input)

	var output datastr.GliffyScene
	var objects []datastr.GliffyObject

	output.EmbeddedResources.Resources = []datastr.GliffyEmbeddedResource{}

	objectIDs := map[string]int{}

	objects, output, objectIDs, err = AddElements(false, input, output, objects, objectIDs)
	if err != nil {
		return "", errors.New("Unable to add element(s): " + err.Error())
	}

	objects, output, _, err = AddElements(true, input, output, objects, objectIDs)
	if err != nil {
		return "", errors.New("Unable to add element(s) with parent(s): " + err.Error())
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
		return "", errors.New("Error occurred during JSON marshaling + " + err.Error())
	}

	return string(outputJson), nil
}

func AddElements(addChildren bool, input datastr.ExcalidrawScene, scene datastr.GliffyScene, objects []datastr.GliffyObject, objectIDs map[string]int) ([]datastr.GliffyObject, datastr.GliffyScene, map[string]int, error) {
	graphics := internal.MapGraphics()

	for i, element := range input.Elements {
		if len(element.ContainerId) > 0 && !addChildren {
			continue
		}
		if len(element.ContainerId) == 0 && addChildren {
			continue
		}

		var hasParent bool = false
		var parent int = 999999
		if len(element.ContainerId) > 0 {
			for obj_k, obj := range objects {
				if obj.ID == objectIDs[element.ContainerId] {
					parent = obj_k
				}
			}

			if parent == 999999 {
				return nil, scene, nil, errors.New("unable to find object parent")
			}

			hasParent = true
		}

		var object datastr.GliffyObject
		var shape datastr.GliffyShape
		var text datastr.GliffyText
		var line datastr.GliffyLine
		var image datastr.GliffyImage
		var svg datastr.GliffySvg
		var embedded datastr.GliffyEmbeddedResource

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

		for _, id := range graphics.Freedraw.Excalidraw {
			if element.Type == id {
				object.UID = graphics.Freedraw.Gliffy[0]
				object.Graphic.Type = "Svg"

				var embeddedResourceId = len(scene.EmbeddedResources.Resources) + 1

				element.Width = element.Width + 8
				element.Height = element.Height + 8

				svg.EmbeddedResourceID = embeddedResourceId
				svg.StrokeColor = element.StrokeColor
				svg.StrokeWidth = int64(FreedrawStrokeWidthConvExcGliffy(element.StrokeWidth))
				svg.DropShadow = false
				svg.ShadowX = 0
				svg.ShadowY = 0

				var svgFill = "none"
				if element.BackgroundColor != "transparent" {
					svgFill = element.BackgroundColor
				}

				xMin, yMin, points := AddPointsOffset(element.Points)
				var svgPath = ConvertPointsToSvgPath(points, element.Width, element.Height, svg.StrokeColor, svgFill, svg.StrokeWidth)
				svg.Svg = svgPath

				object.X = object.X + xMin
				object.Y = object.Y + yMin

				embedded.ID = embeddedResourceId
				embedded.MimeType = "image/svg+xml"
				embedded.Data = svgPath
				embedded.X = 1
				embedded.Y = 1
				embedded.Width = element.Width
				embedded.Height = element.Height
				scene.EmbeddedResources.Resources = append(scene.EmbeddedResources.Resources, embedded)

				object.Graphic.Svg = &svg
			}
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

				if hasParent && objects[parent].Graphic.Line != nil {
					text.Overflow = "both"
				}

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
				line.FillColor = FillColorConvExcGliffy(element.BackgroundColor)
				line.StartArrowRotation = "auto"
				line.EndArrowRotation = "auto"
				line.InterpolationType = "linear"
				line.CornerRadius = 10
				line.Ortho = true
				line.ControlPath = element.Points

				if element.StartArrowhead == nil {
					line.StartArrow = 0
				} else {
					line.StartArrow = ArrowheadConvExcGliffy(*element.StartArrowhead)
				}

				if element.EndArrowhead == nil {
					line.EndArrow = 0
				} else {
					line.EndArrow = ArrowheadConvExcGliffy(*element.EndArrowhead)
				}

				object.Graphic.Line = &line
			}
		}

		for _, id := range graphics.Image.Excalidraw {
			if element.Type == id {
				object.UID = graphics.Image.Gliffy[0]
				object.Graphic.Type = "Image"

				DataURL, err := EmbeddedImgConvExcGliffy(input, element.FileId)
				if err != nil {
					return nil, scene, nil, err
				}

				image.Url = DataURL
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

		if hasParent {
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

	return objects, scene, objectIDs, nil
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

	return file.DataURL, nil
}

func FreedrawStrokeWidthConvExcGliffy(strokeWidth float64) float64 {
	switch strokeWidth {
	case 1:
		strokeWidth = 2
	case 2:
		strokeWidth = 2
	case 4:
		strokeWidth = 6
	}

	return strokeWidth
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

	// TODO: When snapping to grid, this should be the grid size
	// or: Default to 20 (and update tests)
	xMin -= 10
	yMin -= 10

	fmt.Printf("  Canvas Offset X: %f, Offset Y: %f\n", xMin, yMin)

	return xMin, yMin
}

func AddPointsOffset(points [][]float64) (float64, float64, [][]float64) {
	var xMin float64 = 0
	var yMin float64 = 0
	var output [][]float64

	for _, point := range points {
		if point[0] < xMin {
			xMin = point[0]
		}

		if point[1] < yMin {
			yMin = point[1]
		}
	}

	for _, point := range points {
		x := point[0] + math.Abs(xMin)
		y := point[1] + math.Abs(yMin)

		output = append(output, []float64{x, y})
	}

	return xMin, yMin, output
}

func ConvertPointsToSvgPath(points [][]float64, width float64, height float64, stroke string, fill string, strokeWidth int64) string {
	var path string

	for i, point := range points {
		var pointX = point[0] + 3
		var pointY = point[1] + 3

		if i == 0 {
			path = fmt.Sprintf("M%.1f %.1f", pointX, pointY)
		} else {
			path = fmt.Sprintf("%sL%.1f %.1f", path, pointX, pointY)
		}
	}

	svg := fmt.Sprintf("<svg width=\"%.0f\" height=\"%.0f\" viewBox=\"0 0 %.0f %.0f\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\">\n<path fill=\"%s\" stroke=\"%s\" stroke-width=\"%d\"  stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"%s\"></path>\n</svg>", width, height, width, height, fill, stroke, strokeWidth, path)

	return svg
}
