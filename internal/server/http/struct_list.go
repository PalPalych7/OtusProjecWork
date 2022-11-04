package internalhttp

type SlotBanner struct {
	SlotId   int
	BannerId int
}

type ForBannerClick struct {
	SlotId     int
	BannerId   int
	SocGroupId int
}

type ForGetBanner struct {
	SlotId     int
	SocGroupId int
}
