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
