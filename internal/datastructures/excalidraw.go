package internal

type ExcalidrawScene struct {
	AppState struct {
		GridSize            interface{} `json:"gridSize"`
		ViewBackgroundColor string      `json:"viewBackgroundColor"`
	} `json:"appState"`
	Elements []struct {
		Angle           float64 `json:"angle"`
		BackgroundColor string  `json:"backgroundColor"`
		Baseline        float64 `json:"baseline"`
		BoundElements   []struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		} `json:"boundElements"`
		ContainerId  string        `json:"containerId"`
		EndArrowhead *string       `json:"endArrowhead"`
		FillStyle    string        `json:"fillStyle"`
		FontFamily   int64         `json:"fontFamily"`
		FontSize     float64       `json:"fontSize"`
		GroupIds     []interface{} `json:"groupIds"`
		Height       float64       `json:"height"`
		ID           string        `json:"id"`
		FileId       string        `json:"fileId"`
		IsDeleted    bool          `json:"isDeleted"`
		Opacity      float64       `json:"opacity"`
		Points       [][]float64   `json:"points"`
		Roughness    int64         `json:"roughness"`
		Roundness    struct {
			Type int64 `json:"type"`
		} `json:"roundness"`
		Seed            int64   `json:"seed"`
		StartArrowhead  *string `json:"startArrowhead"`
		StrokeColor     string  `json:"strokeColor"`
		StrokeSharpness string  `json:"strokeSharpness"`
		StrokeStyle     string  `json:"strokeStyle"`
		StrokeWidth     float64 `json:"strokeWidth"`
		Text            string  `json:"text"`
		TextAlign       string  `json:"textAlign"`
		Type            string  `json:"type"`
		Version         int64   `json:"version"`
		VersionNonce    int64   `json:"versionNonce"`
		VerticalAlign   string  `json:"verticalAlign"`
		Width           float64 `json:"width"`
		X               float64 `json:"x"`
		Y               float64 `json:"y"`
	} `json:"elements"`
	Files map[string]struct {
		DataUrl string `json:"dataUrl"`
	} `json:"files"`
	Source  string `json:"source"`
	Type    string `json:"type"`
	Version int64  `json:"version"`
}
