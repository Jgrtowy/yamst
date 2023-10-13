package lib

type ReleaseInfo struct {
	Latest struct {
		Release  string `json:"release"`
		Snapshot string `json:"snapshot"`
	} `json:"latest"`
	Versions []VersionInfo `json:"versions"`
}

type VersionInfo struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	Url         string `json:"url"`
	Time        string `json:"time"`
	ReleaseTime string `json:"releaseTime"`
}

type PackageInfo struct {
	Downloads struct {
		Server struct {
			Url string `json:"url"`
		}
	}
}

type PaperBuilds struct {
	Builds []int32 `json:"builds"`
}

type PaperBuildInfo struct {
	Downloads struct {
		Application struct {
			Name string `json:"name"`
		} `json:"application"`
	} `json:"downloads"`
}

type ForgeBuilds struct {
	Builds []string `json:"promos"`
}
