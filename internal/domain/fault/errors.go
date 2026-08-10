package fault

/********Repository Errors********/
func SaveAccountError(err error) *DomainError {
	return &DomainError{
		Code:            CodeDatabaseError,
		FriendlyMessage: "Unexpected server error; please try again later.",
		Cause:           err,
		Origin:          CallerName(2),
	}
}

func FindAccountError(err error) *DomainError {
	return &DomainError{
		Code:            CodeDatabaseError,
		FriendlyMessage: "Unexpected server error; please try again later.",
		Cause:           err,
		Origin:          CallerName(2),
	}
}

func FindAccountNotFoundError(err error) *DomainError {
	return &DomainError{
		Code:            CodeNotFound,
		FriendlyMessage: "Account not found.",
		Cause:           err,
		Origin:          CallerName(2),
	}
}

func SaveAccountAlreadyExistsError(err error) *DomainError {
	return &DomainError{
		Code:            CodeAlreadyExists,
		FriendlyMessage: "Account already exists.",
		Cause:           err,
		Origin:          CallerName(2),
	}
}

/***************Domain Errors****************/
func InvalidEntityError(err error, fields map[string]any) *DomainError {
	return &DomainError{
		Code:            CodeInvalidEntity,
		FriendlyMessage: "Some of the information entered is incorrect; please review it and try again.",
		Cause:           err,
		Origin:          CallerName(2),
		Fields:          fields,
	}
}

/***************UOW Errors****************/
func BeginTransactionError(err error) *DomainError {
	return &DomainError{
		Code:            CodeDatabaseError,
		FriendlyMessage: "Unexpected server error; please try again later.",
		Cause:           err,
	}
}

func CommitTransactionError(err error) *DomainError {
	return &DomainError{
		Code:            CodeDatabaseError,
		FriendlyMessage: "Unexpected server error; please try again later.",
		Cause:           err,
	}
}

func RollbackTransactionError(err error) *DomainError {
	return &DomainError{
		Code:            CodeDatabaseError,
		FriendlyMessage: "Unexpected server error; please try again later.",
		Cause:           err,
	}
}

/***********JSON Decoder Error *****************/
func ErrorJSONSyntaxError(err error) *DomainError {
	return &DomainError{
		Code:            CodeBadRequest,
		FriendlyMessage: "json unmarshal type error",
		Cause:           err,
	}
}

func ErrorJSONUnmarshalTypeError(err error) *DomainError {
	return &DomainError{
		Code:            CodeBadRequest,
		FriendlyMessage: "json unmarshal type error",
		Cause:           err,
	}
}

func ErrorJSON(err error) *DomainError {
	return &DomainError{
		Code:            CodeBadRequest,
		FriendlyMessage: "json unmarshal type error",
		Cause:           err,
	}
}
