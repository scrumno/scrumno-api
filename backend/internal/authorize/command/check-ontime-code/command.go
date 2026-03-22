package check_ontime_code

import entity "github.com/scrumno/scrumno-api/internal/authorize/entity/codes"

type Command struct {
	Phone    string
	Code     string
	CodeType entity.CodesType
}
