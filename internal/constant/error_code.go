package constant

const (
	CodeSuccess = 0

	CodeInvalidQueryParams = 40000
	CodeInvalidUUID        = 40002
	CodeVisitorHashMissing = 40003

	CodePhotoNotFound = 40401

	CodeTooManyBehaviorRequests = 42901
	CodeSuspiciousBehavior      = 42902

	CodeInternalServer = 50000
	CodePhotosList     = 50001
	CodePhotoDetail    = 50002
	CodePhotoLike      = 50003
	CodePhotoDownload  = 50004
	CodeTagList        = 50005
	CodeFilterList     = 50006
)

const (
	MsgSuccess                 = "success"
	MsgInvalidQueryParams      = "invalid query params"
	MsgInvalidUUID             = "invalid uuid"
	MsgVisitorHashMissing      = "visitor hash missing"
	MsgPhotoNotFound           = "photo not found"
	MsgTooManyBehaviorRequests = "too many behavior requests"
	MsgSuspiciousBehavior      = "suspicious behavior blocked"
	MsgInternalServerError     = "internal server error"
)
