package conversion

import (
	datastr "diagram-converter/internal/datastructures"
	"encoding/json"
	"testing"
)

func TestConstructMermaidNodeDef(t *testing.T) {
	tests := []struct {
		id, text, shape, want string
	}{
		{"N1", "My Node", "[", `N1["My Node"]`},
		{"N2", "Decision", "{", `N2{"Decision"}`},
		{"N3", "Circle", "((", `N3(("Circle"))`},
		{"N4", "Round", "(", `N4("Round")`},
	}
	for _, tt := range tests {
		got := constructMermaidNodeDef(tt.id, tt.text, tt.shape)
		if got != tt.want {
			t.Errorf("id=%q, text=%q, shape=%q: got %q, want %q", tt.id, tt.text, tt.shape, got, tt.want)
		}
	}
}

func TestConstructMermaidEdgeArrow(t *testing.T) {
	tests := []struct {
		type_, endArrow, style, want string
	}{
		{"line", "", "", "--"},
		{"arrow", "", "", "-->"},
		{"line", "arrow", "", "-->"},
		{"arrow", "", "dashed", "-.->"},
		{"line", "", "dashed", "-.-"},
		{"arrow", "", "dotted", "==>"},
		{"line", "", "dotted", "==="},
	}
	for _, tt := range tests {
		got := constructMermaidEdgeArrow(tt.type_, tt.endArrow, tt.style)
		if got != tt.want {
			t.Errorf("type=%q, endArrow=%q, style=%q: got %q, want %q", tt.type_, tt.endArrow, tt.style, got, tt.want)
		}
	}
}

func TestConstructMermaidEdgeLabel(t *testing.T) {
	tests := []struct {
		label, elText, want string
	}{
		{"label", "", "|label|"},
		{"", "some\ntext", "|some text|"},
		{"", "", ""},
	}
	for _, tt := range tests {
		got := constructMermaidEdgeLabel(tt.label, tt.elText)
		if got != tt.want {
			t.Errorf("label=%q, elText=%q: got %q, want %q", tt.label, tt.elText, got, tt.want)
		}
	}
}

func TestConstructMermaidStyleString(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"stroke:red;fill:blue;", "stroke:red,fill:blue;"},
		{"stroke:red;", "stroke:red;"},
		{"", ""},
		{"stroke:red", "stroke:red;"},
	}
	for _, tt := range tests {
		got := constructMermaidStyleString(tt.in)
		if got != tt.want {
			t.Errorf("in=%q: got %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestBuildMermaidFromScene(t *testing.T) {
	// Two rectangles with text and an arrow between them
	scene := datastr.ExcalidrawScene{
		Elements: []datastr.ExcalidrawSceneElement{
			{ID: "1", Type: "rectangle", X: 0, Y: 0, Width: 100, Height: 50},
			{ID: "t1", Type: "text", ContainerId: "1", Text: "First"},
			{ID: "2", Type: "rectangle", X: 200, Y: 0, Width: 100, Height: 50},
			{ID: "t2", Type: "text", ContainerId: "2", Text: "Second"},
			{
				ID:   "3",
				Type: "arrow",
				StartBinding: datastr.ExcalidrawElementBinding{
					ElementID: "1",
					Focus:     0,
					Gap:       0,
				},
				EndBinding: datastr.ExcalidrawElementBinding{
					ElementID: "2",
					Focus:     0,
					Gap:       0,
				},
				EndArrowhead: "arrow",
			},
		},
	}
	b, _ := json.Marshal(scene)
	var realScene datastr.ExcalidrawScene
	_ = json.Unmarshal(b, &realScene)

	got, _ := BuildMermaidFromScene(realScene)
	want := `flowchart LR
N0["First"]
N1["Second"]
N0 --> N1
style N0 opacity:0.00;
style N1 opacity:0.00;
`
	if got != want {
		t.Errorf("BuildMermaidFromScene() = %q, want %q", got, want)
	}
}
