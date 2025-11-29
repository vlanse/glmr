package editor

type Project struct {
	ID   int64
	Path string
}

type Settings struct {
	Cmd      string
	Projects []Project
}
