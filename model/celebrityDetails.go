package model

type Celebrity struct {
	CelebrityId   string
	CelebrityName string
	Image         string
	Gender        string
	Constellation string
	BirthDate     string
	Birthplace    string
	Jobs          string
	Nickname      string
	Family        string
	Introduction  string
	Photos        [5]string
	Rewords       [3]string
	Works         [5]RecentWorks
}

type RecentWorks struct {
	WorkId     string
	WorkImage  string
	WorkName   string
	WorkScores string
}
