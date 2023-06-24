package db

func GetQuestionMarks(num int) string {
	questionMarks := ""
	for i := 0; i < num; i++ {
		questionMarks += "?,"
	}
	return questionMarks[:len(questionMarks)-1]
}
