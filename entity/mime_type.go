package entity

func MimeType(fileExt string) string {
	switch fileExt {
	case "jpg", "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	default:
		return "application/octet-stream"
	}
}
