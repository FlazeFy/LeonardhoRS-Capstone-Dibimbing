package config

type Config struct {
	MaxSizeFile     int64
	AllowedFileType []string
}

var ResponseMessages = map[string]string{
	"post":        "created",
	"put":         "updated",
	"hard delete": "permanentally deleted",
	"soft delete": "deleted",
	"recover":     "recovered",
	"get":         "fetched",
	"login":       "login",
	"sign out":    "signed out",
}
var Departments = []string{"IT", "Human Resource", "Finance & Risk Management", "Marketing", "Sales", "Planning & Transformation", "Network"}
var Floors = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13"}
var AssetStatus = []string{"available", "in-use", "maintenance"}
var FindingCategories = []string{"broken", "missing", "upgrade", "feedback"}
var Days = []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
var ConfigFile = Config{
	MaxSizeFile:     10000000, // 10 MB
	AllowedFileType: []string{"jpg", "jpeg"},
}
