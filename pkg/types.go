package pkg

import "time"

type Tail struct {
	Pk                  string      `json:"PK"`
	Sk                  string      `json:"SK"`
	ApplicationIndexKey string      `json:"ApplicationIndexKey"`
	ApplicationLoaderID string      `json:"ApplicationLoaderId"`
	BuildLink           interface{} `json:"BuildLink"`
	Config              Config      `json:"Config"`
	ContentLength       int64       `json:"ContentLength"`
	Digest              string      `json:"Digest"`
	GitHash             string      `json:"GitHash"`
	GitPackage          string      `json:"GitPackage"`
	GitTags             string      `json:"GitTags"`
	Name                string      `json:"Name"`
	Signature           string      `json:"Signature"`
	TimestampCreated    time.Time   `json:"TimestampCreated"`
}

type Config struct {
	Command   []string `json:"command"`
	Env       []string `json:"env"`
	Mounts    []string `json:"mounts"`
	Signature string   `json:"signature"`
	Workdir   string   `json:"workdir"`
}
