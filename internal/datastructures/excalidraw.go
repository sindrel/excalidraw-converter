package internal

type ExcalidrawElementBinding struct {
	ElementID string  `json:"elementId"`
	Focus     float64 `json:"focus"`
	Gap       float64 `json:"gap"`
}

type ExcalidrawSceneElement struct {
	Angle           float64 `json:"angle"`
	BackgroundColor string  `json:"backgroundColor"`
	Baseline        float64 `json:"baseline"`
	BoundElements   []struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"boundElements"`
	ContainerId  string                   `json:"containerId"`
	EndArrowhead string                   `json:"endArrowhead"`
	EndBinding   ExcalidrawElementBinding `json:"endBinding,omitempty"`
	FillStyle    string                   `json:"fillStyle"`
	FontFamily   int64                    `json:"fontFamily"`
	FontSize     float64                  `json:"fontSize"`
	GroupIds     []interface{}            `json:"groupIds"`
	Height       float64                  `json:"height"`
	ID           string                   `json:"id"`
	FileId       string                   `json:"fileId"`
	IsDeleted    bool                     `json:"isDeleted"`
	Link         string                   `json:"link"`
	Opacity      float64                  `json:"opacity"`
	Points       [][]float64              `json:"points"`
	Pressures    []float64                `json:"pressures"`
	Roughness    int64                    `json:"roughness"`
	Roundness    struct {
		Type int64 `json:"type"`
	} `json:"roundness"`
	Seed             int64                    `json:"seed"`
	SimulatePressure bool                     `json:"simulatePressure"`
	StartArrowhead   string                   `json:"startArrowhead"`
	StartBinding     ExcalidrawElementBinding `json:"startBinding,omitempty"`
	StrokeColor      string                   `json:"strokeColor"`
	StrokeSharpness  string                   `json:"strokeSharpness"`
	StrokeStyle      string                   `json:"strokeStyle"`
	StrokeWidth      float64                  `json:"strokeWidth"`
	Text             string                   `json:"text"`
	TextAlign        string                   `json:"textAlign"`
	Type             string                   `json:"type"`
	Version          int64                    `json:"version"`
	VersionNonce     int64                    `json:"versionNonce"`
	VerticalAlign    string                   `json:"verticalAlign"`
	Width            float64                  `json:"width"`
	X                float64                  `json:"x"`
	Y                float64                  `json:"y"`
}

type ExcalidrawScene struct {
	AppState struct {
		GridSize            int64  `json:"gridSize"`
		ViewBackgroundColor string `json:"viewBackgroundColor"`
	} `json:"appState"`
	Elements []ExcalidrawSceneElement `json:"elements"`
	Files    map[string]struct {
		DataURL  string `json:"dataURL"`
		ID       string `json:"id"`
		MimeType string `json:"mimeType"`
	} `json:"files"`
	Source  string `json:"source"`
	Type    string `json:"type"`
	Version int64  `json:"version"`
}
