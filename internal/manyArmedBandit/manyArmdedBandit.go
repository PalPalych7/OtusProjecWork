package manyArmedBandit

type BanditConfig struct {
	FullLearnigCount     int // количество запросов в режиме "полного обучения"
	PartialLearningCount int // количество запросов в режиме "чатичного обучения"
	FinalRandomPecent    int // вероятность случайного выбора после обучения (в процентах)
}

type BannerStruct struct {
	Id         int
	ShowCount  int
	ClickCount int
}

type MyBandit struct {
	BanditConfig BanditConfig
}

func (b MyBandit) GetBannerId(arrStruct []BannerStruct) int {
	return 1
}

func New(BC BanditConfig) MyBandit {
	return MyBandit{BC}
}
