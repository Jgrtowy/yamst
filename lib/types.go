package lib

type ReleaseInfo struct {
	Latest struct {
		Release string `json:"release"`
	} `json:"latest"`
}
