package setup

import "mime"

func MimeTypes() {
	_ = mime.AddExtensionType(".css", "text/css")
	_ = mime.AddExtensionType(".ico", "image/x-icon")
	_ = mime.AddExtensionType(".gif", "image/gif")
	_ = mime.AddExtensionType(".jpeg", "image/jpeg")
	_ = mime.AddExtensionType(".jpg", "image/jpeg")
	_ = mime.AddExtensionType(".js", "application/javascript")
	_ = mime.AddExtensionType(".json", "application/json")
	_ = mime.AddExtensionType(".htm", "text/html")
	_ = mime.AddExtensionType(".html", "text/html")
	_ = mime.AddExtensionType(".png", "image/png")
}
