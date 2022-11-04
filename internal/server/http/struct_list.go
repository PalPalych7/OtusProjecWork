package internalhttp

type SlotBanner struct {
	SlotID   int
	BannerID int
}

type ForBannerClick struct {
	SlotID     int
	BannerID   int
	SocGroupID int
}

type ForGetBanner struct {
	SlotID     int
	SocGroupID int
}
