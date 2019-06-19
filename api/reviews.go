package api


type ReviewCollection struct {
	Url 					   string				`json:"url"`
	TotalCount				   string				`json:"total_count"`
	Reviews 				   []ReviewData			`json:"data"`
	Pages struct {
		NextURL 			   string				`json:"next_url"`
	} `json:"pages"`

}

type ReviewData struct {
	Id    					   int 					`json:"id"`
	Object 					   string 				`json:"object"`
	Data struct {
		SubjectId			   int					`json:"subject_id"`
		SubjectType			   string				`json:"subject_type"`
		MeaningCorrect		   int					`json:"meaning_correct"`
		MeaningIncorrect	   int					`json:"meaning_incorrect"`
		MeaningMaxStreak	   int					`json:"meaning_max_streak"`
		MeaningCurrentStreak   int					`json:"meaning_current_streak"`
		ReadingCorrect		   int					`json:"reading_correct"`
		ReadingIncorrect	   int					`json:"reading_incorrect"`
		ReadingMaxStreak	   int					`json:"reading_max_streak"`
		ReadingCurrentStreak   int					`json:"reading_current_streak"`
		PercentageCorrect	   int					`json:"percentage_correct"`
		Hidden				   bool					`json:"hidden"`
	} `json:"data"`
}
