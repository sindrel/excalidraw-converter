package conversion

import (
	internal "diagram-converter/internal"
	datastr "diagram-converter/internal/datastructures"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

func ConvertExcalidrawDiagramToMermaidAndSaveToFile(importPath string, exportPath string) error {
	fmt.Printf("Parsing input file: %s\n", importPath)

	data, err := os.ReadFile(importPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "File reading failed. %s\n", err)
		os.Exit(1)
	}

	output, err := ConvertExcalidrawDiagramToMermaid(string(data))
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

// ConvertExcalidrawDiagramToMermaid converts an Excalidraw diagram to a Mermaid flowchart string.
func ConvertExcalidrawDiagramToMermaid(data string) (string, error) {
	var input datastr.ExcalidrawScene
	err := json.Unmarshal([]byte(data), &input)
	if err != nil {
		return "", errors.New("Unable to parse input: " + err.Error())
	}
	return BuildMermaidFromScene(input)
}

// Helper to format a node definition for Mermaid
func constructMermaidNodeDef(name, label, shape string) string {
	// Mermaid node shapes: [rectangle], ((circle)), (round), {diamond}
	switch shape {
	case "{":
		return fmt.Sprintf("%s{\"%s\"}", name, label)
	case "((":
		return fmt.Sprintf("%s((\"%s\"))", name, label)
	case "(":
		return fmt.Sprintf("%s(\"%s\")", name, label)
	default:
		return fmt.Sprintf("%s[\"%s\"]", name, label)
	}
}

// Helper to map Excalidraw edge type and stroke style to Mermaid arrow
func constructMermaidEdgeArrow(elType, endArrowhead, strokeStyle string) string {
	arrow := "--"
	if elType == "arrow" || endArrowhead == "arrow" {
		arrow = "-->"
	}
	if strokeStyle == "dashed" {
		if arrow == "-->" {
			arrow = "-.->"
		} else {
			arrow = "-.-"
		}
	} else if strokeStyle == "dotted" {
		if arrow == "-->" {
			arrow = "==>"
		} else {
			arrow = "==="
		}
	}
	return arrow
}

// Helper to extract edge label
func constructMermaidEdgeLabel(linkText, elText string) string {
	label := linkText
	if label == "" && elText != "" {
		label = strings.ReplaceAll(elText, "\n", " ")
	}
	if label != "" {
		return fmt.Sprintf("|%s|", label)
	}
	return ""
}

// Helper to format style string for Mermaid
func constructMermaidStyleString(style string) string {
	style = strings.ReplaceAll(style, ";", ",")
	style = strings.TrimSuffix(style, ",")
	style = strings.TrimSpace(style)
	if style != "" && !strings.HasSuffix(style, ";") {
		style += ";"
	}
	return style
}

// BuildMermaidFromScene converts an ExcalidrawScene struct to a Mermaid flowchart string.
func BuildMermaidFromScene(input datastr.ExcalidrawScene) (string, error) {
	nodeMap := make(map[string]string) // Excalidraw ID -> Mermaid node name
	nodeLabels := make(map[string]string)
	nodeShapes := make(map[string]string)
	nodeStyles := make(map[string]string)
	nodeCount := 0

	// First, collect text for each containerId (for nodes) and for links (for edges)
	containerText := make(map[string]string)
	linkText := make(map[string]string)
	for _, el := range input.Elements {
		if el.IsDeleted {
			continue
		}
		if el.Type == "text" && el.ContainerId != "" {
			for _, parent := range input.Elements {
				if parent.ID == el.ContainerId && (parent.Type == "rectangle" || parent.Type == "diamond" || parent.Type == "ellipse" || parent.Type == "roundRectangle") {
					if containerText[el.ContainerId] != "" {
						containerText[el.ContainerId] += " "
					}
					containerText[el.ContainerId] += strings.ReplaceAll(el.Text, "\n", " ")
				}
				if parent.ID == el.ContainerId && (parent.Type == "arrow" || parent.Type == "line") {
					if linkText[el.ContainerId] != "" {
						linkText[el.ContainerId] += " "
					}
					linkText[el.ContainerId] += strings.ReplaceAll(el.Text, "\n", " ")
				}
			}
		}
	}

	// Assign node names and gather node info
	for _, el := range input.Elements {
		if el.IsDeleted {
			continue
		}
		if el.Type == "rectangle" || el.Type == "diamond" || el.Type == "ellipse" || el.Type == "roundRectangle" {
			name := fmt.Sprintf("N%d", nodeCount)
			nodeMap[el.ID] = name
			label := containerText[el.ID]
			if label == "" {
				label = name
			}
			nodeLabels[el.ID] = label
			// Shape mapping
			shape := "[" // default rectangle
			switch el.Type {
			case "rectangle":
				if el.Roundness.Type > 0 {
					shape = "("
				} else {
					shape = "["
				}
			case "roundRectangle":
				shape = "("
			case "ellipse":
				shape = "(("
			case "diamond":
				shape = "{" // Mermaid diamond
			}
			nodeShapes[el.ID] = shape
			// Style mapping (stroke, fill, etc.)
			style := ""
			if el.StrokeStyle == "dashed" {
				style += "stroke-dasharray: 5 5;"
			} else if el.StrokeStyle == "dotted" {
				style += "stroke-dasharray: 2 2;"
			}
			if el.StrokeColor != "" {
				style += fmt.Sprintf("stroke:%s;", el.StrokeColor)
			}
			if el.BackgroundColor != "transparent" && el.BackgroundColor != "" {
				style += fmt.Sprintf("fill:%s;", el.BackgroundColor)
			}
			if el.Opacity < 100 {
				style += fmt.Sprintf("opacity:%.2f;", el.Opacity/100.0)
			}
			nodeStyles[el.ID] = style
			nodeCount++
		}
	}

	orientation := getFlowchartOrientation(input)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("flowchart %s\n", orientation))

	// Output all nodes equally
	nodeIDs := make([]string, 0, len(nodeMap))
	for id := range nodeMap {
		nodeIDs = append(nodeIDs, id)
	}
	sort.Strings(nodeIDs)
	for _, id := range nodeIDs {
		name := nodeMap[id]
		label := nodeLabels[id]
		shape := nodeShapes[id]
		nodeDef := constructMermaidNodeDef(name, label, shape)
		sb.WriteString(nodeDef + "\n")
	}

	// Output edges (arrows/lines)
	for _, el := range input.Elements {
		if el.IsDeleted {
			continue
		}
		if el.Type == "arrow" || el.Type == "line" {
			startID := el.StartBinding.ElementID
			endID := el.EndBinding.ElementID
			if startID == "" || endID == "" {
				continue
			}
			startNode, ok1 := nodeMap[startID]
			endNode, ok2 := nodeMap[endID]
			if !ok1 || !ok2 {
				continue
			}

			arrow := constructMermaidEdgeArrow(el.Type, el.EndArrowhead, el.StrokeStyle)
			labelStr := constructMermaidEdgeLabel(linkText[el.ID], el.Text)
			edgeDef := fmt.Sprintf("%s %s%s %s\n", startNode, arrow, labelStr, endNode)
			sb.WriteString(edgeDef)
		}
	}

	// Output style blocks for nodes
	sortedStyleIDs := make([]string, 0, len(nodeMap))
	for id := range nodeMap {
		sortedStyleIDs = append(sortedStyleIDs, id)
	}
	sort.Strings(sortedStyleIDs)
	for _, id := range sortedStyleIDs {
		name := nodeMap[id]
		style := constructMermaidStyleString(nodeStyles[id])
		if style != "" {
			sb.WriteString(fmt.Sprintf("style %s %s\n", name, style))
		}
	}

	return sb.String(), nil
}

func getFlowchartOrientation(input datastr.ExcalidrawScene) string {
	var minX, minY, maxX, maxY float64
	first := true
	for _, el := range input.Elements {
		if el.IsDeleted {
			continue
		}
		if el.Type == "rectangle" || el.Type == "diamond" || el.Type == "ellipse" || el.Type == "roundRectangle" {
			x1 := el.X
			y1 := el.Y
			x2 := el.X + el.Width
			y2 := el.Y + el.Height
			if first {
				minX, maxX = x1, x2
				minY, maxY = y1, y2
				first = false
			} else {
				if x1 < minX {
					minX = x1
				}
				if y1 < minY {
					minY = y1
				}
				if x2 > maxX {
					maxX = x2
				}
				if y2 > maxY {
					maxY = y2
				}
			}
		}
	}
	width := maxX - minX
	height := maxY - minY
	orientation := "TD"
	if width > height {
		orientation = "LR"
	}

	return orientation
}
