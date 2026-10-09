package setup

// Version is the backend's release version. Release builds set it with
//
//	-ldflags "-X bitbucket.org/daforester/go-audience/apps/fantasy/setup.Version=0.0.1"
//
// from the fantasy-backend/vX.Y.Z git tag; anything else reports "dev".
var Version = "dev"
