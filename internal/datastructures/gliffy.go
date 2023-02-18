package internal

type GliffyScene struct {
	ContentType       string `json:"contentType"`
	EmbeddedResources struct {
		Index     int64    `json:"index"`
		Resources []string `json:"resources"`
	} `json:"embeddedResources"`
	Metadata struct {
		AutosaveDisabled bool        `json:"autosaveDisabled"`
		EditorVersion    interface{} `json:"editorVersion"`
		ExportBorder     bool        `json:"exportBorder"`
		LastSerialized   int64       `json:"lastSerialized"`
		Libraries        []string    `json:"libraries"`
		LoadPosition     string      `json:"loadPosition"`
		Revision         int64       `json:"revision"`
		Title            string      `json:"title"`
	} `json:"metadata"`
	Stage struct {
		AutoFit         bool           `json:"autoFit"`
		Background      string         `json:"background"`
		DrawingGuidesOn bool           `json:"drawingGuidesOn"`
		ExportBorder    bool           `json:"exportBorder"`
		GridOn          bool           `json:"gridOn"`
		Height          int64          `json:"height"`
		ImageCache      struct{}       `json:"imageCache"`
		Layers          []GliffyLayer  `json:"layers"`
		MaxHeight       int64          `json:"maxHeight"`
		MaxWidth        int64          `json:"maxWidth"`
		NodeIndex       int64          `json:"nodeIndex"`
		Objects         []GliffyObject `json:"objects"`
		PrintModel      struct {
			DisplayPageBreaks bool   `json:"displayPageBreaks"`
			FitToOnePage      bool   `json:"fitToOnePage"`
			PageSize          string `json:"pageSize"`
			Portrait          bool   `json:"portrait"`
		} `json:"printModel"`
		ShapeStyles  struct{}    `json:"shapeStyles"`
		SnapToGrid   bool        `json:"snapToGrid"`
		ThemeData    interface{} `json:"themeData"`
		ViewportType string      `json:"viewportType"`
		Width        float64     `json:"width"`
	} `json:"stage"`
	Version string `json:"version"`
}

type GliffyObject struct {
	Children interface{} `json:"children"`
	Graphic  struct {
		Shape *GliffyShape `json:",omitempty"`
		Text  *GliffyText  `json:",omitempty"`
		Line  *GliffyLine  `json:",omitempty"`
		Type  string       `json:"type"`
	} `json:"graphic"`
	Height          float64  `json:"height"`
	Hidden          bool     `json:"hidden"`
	ID              int      `json:"id"`
	LayerID         string   `json:"layerId"`
	LinkMap         []string `json:"linkMap"`
	LockAspectRatio bool     `json:"lockAspectRatio"`
	LockShape       bool     `json:"lockShape"`
	Order           int      `json:"order"`
	Rotation        float64  `json:"rotation"`
	UID             string   `json:"uid"`
	Width           float64  `json:"width"`
	X               float64  `json:"x"`
	Y               float64  `json:"y"`
}

type GliffyShape struct {
	DashStyle   string  `json:"dashStyle"`
	DropShadow  bool    `json:"dropShadow"`
	FillColor   string  `json:"fillColor"`
	Gradient    bool    `json:"gradient"`
	Opacity     float64 `json:"opacity"`
	ShadowX     int64   `json:"shadowX"`
	ShadowY     int64   `json:"shadowY"`
	State       int64   `json:"state"`
	StrokeColor string  `json:"strokeColor"`
	StrokeWidth int64   `json:"strokeWidth"`
	Tid         string  `json:"tid"`
}

type GliffyText struct {
	CalculatedHeight   float64     `json:"calculatedHeight"`
	CalculatedWidth    int64       `json:"calculatedWidth"`
	Hposition          string      `json:"hposition"`
	HTML               string      `json:"html"`
	OuterPaddingBottom int64       `json:"outerPaddingBottom"`
	OuterPaddingLeft   int64       `json:"outerPaddingLeft"`
	OuterPaddingRight  int64       `json:"outerPaddingRight"`
	OuterPaddingTop    int64       `json:"outerPaddingTop"`
	Overflow           string      `json:"overflow"`
	PaddingBottom      int64       `json:"paddingBottom"`
	PaddingLeft        int64       `json:"paddingLeft"`
	PaddingRight       int64       `json:"paddingRight"`
	PaddingTop         int64       `json:"paddingTop"`
	Tid                interface{} `json:"tid"`
	Valign             string      `json:"valign"`
	Vposition          string      `json:"vposition"`
}

type GliffyLine struct {
	ControlPath        [][]float64 `json:"controlPath"`
	CornerRadius       int64       `json:"cornerRadius"`
	DashStyle          interface{} `json:"dashStyle"`
	EndArrow           int         `json:"endArrow"`
	EndArrowRotation   string      `json:"endArrowRotation"`
	FillColor          string      `json:"fillColor"`
	HopType            interface{} `json:"hopType"`
	InterpolationType  string      `json:"interpolationType"`
	LockSegments       struct{}    `json:"lockSegments"`
	Ortho              bool        `json:"ortho"`
	StartArrow         int         `json:"startArrow"`
	StartArrowRotation string      `json:"startArrowRotation"`
	StrokeColor        string      `json:"strokeColor"`
	StrokeWidth        int64       `json:"strokeWidth"`
}

type GliffyLayer struct {
	Active    bool   `json:"active"`
	GUID      string `json:"guid"`
	Locked    bool   `json:"locked"`
	Name      string `json:"name"`
	NodeIndex int64  `json:"nodeIndex"`
	Order     int64  `json:"order"`
	Visible   bool   `json:"visible"`
}
