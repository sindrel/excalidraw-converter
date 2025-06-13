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

type edgeStyleInfo struct {
	index int
	style string
}
type edge struct {
	startNode string
	arrow     string
	labelStr  string
	endNode   string
}

// Calls for conversion to Mermaid flowchart and saves it to a file
func ConvertExcalidrawDiagramToMermaidAndSaveToFile(importPath string, exportPath string, flowDirection string) error {
	output, err := ConvertExcalidrawDiagramToMermaidAndOutputAsString(importPath, exportPath, flowDirection)
	if err != nil {
		return err
	}

	err = internal.WriteToFile(exportPath, output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Saving diagram failed. %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Converted diagram saved to file: %s\n", exportPath)

	return nil
}

// Converts an Excalidraw diagram to a Mermaid flowchart and returns it as a string
func ConvertExcalidrawDiagramToMermaidAndOutputAsString(importPath string, exportPath string, flowDirection string) (string, error) {
	fmt.Printf("Parsing input file: %s\n", importPath)

	data, err := os.ReadFile(importPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "File reading failed. %s\n", err)
		os.Exit(1)
	}

	output, err := ConvertExcalidrawDiagramToMermaid(string(data), flowDirection)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Diagram conversion failed. %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Diagram successfully converted to Mermaid format\n")

	return string(output), nil
}

// Reads and unmarshals an Excalidraw diagram and calls for conversion to a Mermaid flowchart string
func ConvertExcalidrawDiagramToMermaid(data string, flowDirection string) (string, error) {
	var input datastr.ExcalidrawScene
	err := json.Unmarshal([]byte(data), &input)
	if err != nil {
		return "", errors.New("Unable to parse input: " + err.Error())
	}
	return BuildMermaidFromScene(input, flowDirection)
}

// Converts an ExcalidrawScene struct to a Mermaid flowchart string
func BuildMermaidFromScene(input datastr.ExcalidrawScene, flowDirection string) (string, error) {
	nodeMap := make(map[string]string)
	nodeLabels := make(map[string]string)
	nodeShapes := make(map[string]string)
	nodeStyles := make(map[string]string)
	nodeCount := 0

	// First, collect text and fontSize for each containerId (for nodes) and for links (for edges)
	containerText := make(map[string]string)
	containerTextColor := make(map[string]string)
	containerFontSize := make(map[string]float64)
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
					containerText[el.ContainerId] += strings.ReplaceAll(el.Text, "\n", "<br>")
					// If text color is not default, record it
					if el.StrokeColor != "" && el.StrokeColor != "#1e1e1e" && el.StrokeColor != "black" {
						containerTextColor[el.ContainerId] = el.StrokeColor
					}
					// Record fontSize (use the first found for this container)
					if _, ok := containerFontSize[el.ContainerId]; !ok {
						containerFontSize[el.ContainerId] = el.FontSize
					}
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
	linkedNodes := make(map[string]string) // nodeName -> url
	for _, el := range input.Elements {
		if el.IsDeleted {
			continue
		}
		if el.Type == "rectangle" || el.Type == "diamond" || el.Type == "ellipse" || el.Type == "roundRectangle" {
			name := fmt.Sprintf("N%d", nodeCount)
			nodeMap[el.ID] = name
			label := containerText[el.ID]
			if label == "" {
				label = " "
			}
			nodeLabels[el.ID] = label
			// Shape mapping
			shape := "[" // Default rectangle
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
				shape = "{"
			}
			nodeShapes[el.ID] = shape
			nodeStyles[el.ID] = getMermaidNodeStyle(el, containerFontSize, containerTextColor)
			if el.Link != "" && el.Link != "null" {
				linkedNodes[name] = el.Link
			}
			nodeCount++
		}
	}

	direction := getMermaidFlowchartDirection(flowDirection, input)
	// curve := "basis"  // Default curve
	// renderer := "elk" // Default renderer

	var sb strings.Builder
	// sb.WriteString(fmt.Sprintf("%%%%{ init: { 'flowchart': { 'curve': '%s', 'defaultRenderer': '%s' } } }%%%%\n", curve, renderer))
	sb.WriteString(fmt.Sprintf("flowchart %s\n", direction))

	containedBy, children := detectMermaidSpatialContainment(input.Elements)

	// Recursively output nodes/subgraphs
	var writeNodeOrSubgraph func(id string, depth int)
	writeNodeOrSubgraph = func(id string, depth int) {
		name := nodeMap[id]
		label := nodeLabels[id]
		shape := nodeShapes[id]
		if len(children[id]) > 0 {
			sortedChildren := make([]string, len(children[id]))
			copy(sortedChildren, children[id])
			sort.Strings(sortedChildren)
			sb.WriteString(strings.Repeat("  ", depth) + "subgraph " + name + " [\"" + label + "\"]\n")
			for _, cid := range sortedChildren {
				writeNodeOrSubgraph(cid, depth+1)
			}
			sb.WriteString(strings.Repeat("  ", depth) + "end\n")
		} else {
			nodeDef := constructMermaidNodeDef(name, label, shape)
			sb.WriteString(strings.Repeat("  ", depth) + nodeDef + "\n")
		}
	}

	// Output all top-level nodes/subgraphs
	nodeIDs := make([]string, 0, len(nodeMap))
	for id := range nodeMap {
		nodeIDs = append(nodeIDs, id)
	}
	sort.Strings(nodeIDs)
	for _, id := range nodeIDs {
		if containedBy[id] == "" { // Only top-level
			writeNodeOrSubgraph(id, 0)
		}
	}

	// Build map of edges
	edgeStyleList := []edgeStyleInfo{}
	edges := []edge{}
	edgeCount := 0
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
			edges = append(edges, edge{
				startNode: startNode,
				arrow:     arrow,
				labelStr:  labelStr,
				endNode:   endNode,
			})

			// Build edge style string if needed
			style := ""
			if el.StrokeStyle == "dashed" {
				style += "stroke-dasharray: 5 5;"
			} else if el.StrokeStyle == "dotted" {
				style += "stroke-dasharray: 2 2;"
			}
			if el.StrokeColor != "" && el.StrokeColor != "#1e1e1e" {
				// Ensure the color starts with '#'
				color := el.StrokeColor
				if !strings.HasPrefix(color, "#") {
					color = "#" + color
				}
				style += fmt.Sprintf("stroke:%s,color:black;", color)
			}
			// In edge style mapping (replace the current stroke-width logic)
			if el.StrokeWidth == 4 {
				style += "stroke-width:2;"
			} else if el.StrokeWidth == 1 {
				style += "stroke-width:0.5;"
			}
			if el.Opacity < 100 {
				style += fmt.Sprintf("opacity:%.2f;", el.Opacity/100.0)
			}
			if style != "" {
				edgeStyleList = append(edgeStyleList, edgeStyleInfo{index: edgeCount, style: style})
			}
			edgeCount++
		}
	}

	// Identify pairs of edges that can be grouped and chained
	edgesByStart := make(map[string][]edge)
	for _, e := range edges {
		edgesByStart[e.startNode] = append(edgesByStart[e.startNode], e)
	}
	groupedIndices := make(map[int]bool)

	// Sort the startNode keys for consistent output
	startNodes := make([]string, 0, len(edgesByStart))
	for k := range edgesByStart {
		startNodes = append(startNodes, k)
	}
	sort.Strings(startNodes)

	for _, startNode := range startNodes {
		group := edgesByStart[startNode]
		if len(group) == 2 {
			// Find indices in the original edges slice
			var idx1, idx2 int = -1, -1
			for i := range edges {
				if edges[i] == group[0] && idx1 == -1 {
					idx1 = i
				} else if edges[i] == group[1] && idx2 == -1 {
					idx2 = i
				}
			}
			if idx1 != -1 {
				groupedIndices[idx1] = true
			}
			if idx2 != -1 {
				groupedIndices[idx2] = true
			}

			e1 := group[0]
			e2 := group[1]

			combined := fmt.Sprintf("%s %s%s %s %s%s %s\n", e1.endNode, e1.arrow, e1.labelStr, e1.startNode, e2.arrow, e2.labelStr, e2.endNode)
			sb.WriteString(combined)
		}
	}

	// Output all other edges not in a pair
	for i := range edges {
		if groupedIndices[i] {
			continue
		}
		e := &edges[i]
		edgeDef := fmt.Sprintf("%s %s%s %s\n", e.startNode, e.arrow, e.labelStr, e.endNode)
		sb.WriteString(edgeDef)
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

	// Output linkStyle blocks for edges (correct Mermaid syntax)
	for _, info := range edgeStyleList {
		styleStr := constructMermaidStyleString(info.style)
		if styleStr != "" {
			sb.WriteString(fmt.Sprintf("linkStyle %d %s\n", info.index, styleStr))
		}
	}

	// Output click links for nodes with links
	for node, url := range linkedNodes {
		sb.WriteString(fmt.Sprintf("click %s \"%s\" _blank\n", node, url))
	}

	return sb.String(), nil
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

	// Map endArrowhead to Mermaid edge types
	switch endArrowhead {
	case "arrow":
		arrow = "-->"
	case "circle_outline", "circle":
		arrow = "--o"
	case "arrow_bidirectional":
		arrow = "<-->"
	case "circle_outline_bidirectional":
		arrow = "o--o"
	// Add more mappings here as needed
	default:
		if elType == "arrow" {
			arrow = "-->"
		}
	}

	if strokeStyle == "dashed" || strokeStyle == "dotted" {
		if arrow == "-->" {
			arrow = "-.->"
		} else if arrow == "--o" {
			arrow = "-.o"
		} else if arrow == "<-->" {
			arrow = "<-.-.->"
		} else if arrow == "o--o" {
			arrow = "o-.-o"
		} else {
			arrow = "-.-"
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
		return fmt.Sprintf("|\"%s\"|", label)
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

// Helper function to find the index of an element by ID (for containedBy logic)
func getExcalidrawElementIndexByID(elements []datastr.ExcalidrawSceneElement, id string) int {
	for i, el := range elements {
		if el.ID == id {
			return i
		}
	}
	return -1
}

// Helper function to detect spatial containment for subgraphs
func detectMermaidSpatialContainment(elements []datastr.ExcalidrawSceneElement) (map[string]string, map[string][]string) {
	containedBy := make(map[string]string) // nodeID -> parentID
	children := make(map[string][]string)  // parentID -> []childID
	for idA, elA := range elements {
		if elA.IsDeleted || !(elA.Type == "rectangle" || elA.Type == "diamond" || elA.Type == "ellipse" || elA.Type == "roundRectangle") {
			continue
		}
		boxA := struct{ x1, y1, x2, y2 float64 }{elA.X, elA.Y, elA.X + elA.Width, elA.Y + elA.Height}
		for idB, elB := range elements {
			if idA == idB || elB.IsDeleted || !(elB.Type == "rectangle" || elB.Type == "diamond" || elB.Type == "ellipse" || elB.Type == "roundRectangle") {
				continue
			}
			boxB := struct{ x1, y1, x2, y2 float64 }{elB.X, elB.Y, elB.X + elB.Width, elB.Y + elB.Height}
			// B is inside A?
			if boxB.x1 >= boxA.x1 && boxB.y1 >= boxA.y1 && boxB.x2 <= boxA.x2 && boxB.y2 <= boxA.y2 {
				// Only set if not already contained by a smaller parent
				if parent, ok := containedBy[elB.ID]; !ok || (parent != "" && (elA.Width*elA.Height < elements[getExcalidrawElementIndexByID(elements, parent)].Width*elements[getExcalidrawElementIndexByID(elements, parent)].Height)) {
					containedBy[elB.ID] = elA.ID
				}
			}
		}
	}
	for child, parent := range containedBy {
		children[parent] = append(children[parent], child)
	}
	return containedBy, children
}

// Helper to map node style from Excalidraw element
func getMermaidNodeStyle(el datastr.ExcalidrawSceneElement, containerFontSize map[string]float64, containerTextColor map[string]string) string {
	style := ""
	if el.StrokeStyle == "dashed" {
		style += "stroke-dasharray: 5 5;"
	} else if el.StrokeStyle == "dotted" {
		style += "stroke-dasharray: 2 2;"
	}
	if el.StrokeColor != "" {
		style += fmt.Sprintf("stroke:%s;", el.StrokeColor)
	}
	// Map strokeWidth: 4 -> 2, 1 -> 0.5, otherwise omit
	if el.StrokeWidth == 4 {
		style += "stroke-width:2;"
	} else if el.StrokeWidth == 1 {
		style += "stroke-width:0.5;"
	}
	if el.BackgroundColor != "" {
		style += fmt.Sprintf("fill:%s;", el.BackgroundColor)
	}
	if el.Opacity < 100 {
		style += fmt.Sprintf("opacity:%.2f;", el.Opacity/100.0)
	}
	// Font size mapping: 16 -> 90%, 28 -> 110%, else omit (from associated text element)
	if fontSize, ok := containerFontSize[el.ID]; ok {
		if fontSize == 16 {
			style += "font-size:90%;"
		} else if fontSize == 28 {
			style += "font-size:110%;"
		}
	}
	// Add text color if found for this node
	if color, ok := containerTextColor[el.ID]; ok {
		if !strings.HasPrefix(color, "#") && color != "black" && color != "white" {
			color = "#" + color
		}
		style += fmt.Sprintf("color:%s;", color)
	}
	return style
}

// Helper to determine flowchart direction string
func getMermaidFlowchartDirection(flowDirection string, input datastr.ExcalidrawScene) string {
	switch strings.ToLower(flowDirection) {
	case "left-right", "lr":
		return "LR"
	case "right-left", "rl":
		return "RL"
	case "bottom-top", "bt":
		return "BT"
	case "top-down", "td":
		return "TD"
	default:
		return ""
	}
}
