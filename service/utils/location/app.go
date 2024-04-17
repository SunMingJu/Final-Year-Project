package location

import "path"

type AppConfigStruct struct {
	ImagePath struct {
		SystemHeadPortrait string //System Avatar Path
		//UserHeadPortrait   string //User avatar path
		//LiveCover          string //Live Cover
	}
}

var AppConfig *AppConfigStruct

func init() {
	AppConfig = &AppConfigStruct{
		ImagePath: struct {
			SystemHeadPortrait string
			//UserHeadPortrait   string
			//LiveCover          string
		}{
			SystemHeadPortrait: path.Clean("assets/static/img/users/headPortrait/system"),
		},
	}
}
