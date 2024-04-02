package define

const (
	ExistUserName = iota + 4001
	ShortPasswordLength
	PasswordInconsistency
	VerificationError
	EmptyMail
	CodeSendFailed
	EmptyAccountOrPassword
	AccountNotFind
	ErrorPassword
	SignatureTooLong
	NotMatchMail
	ImageFormatError
	CodeExpired
	SameNameFavorite
	NotFindFavorite
	ProhibitFavoritesNameEmpty
	ErrorVideoFormat
)
const (
	QueryUserError = iota + 5001
	PasswordEncryptionError
	ObtainUserInformationFailed
	CreateTokenError
	CreateUserFailed
	UploadUserAvatarFailed
	CreateWebSocketFailed
	LogoutUserFailed
	ModifySignatureFailed
	FollowUserFailed
	UnfollowUserFailed
	ModifyUserNameFailed
	ModifyPasswordFailed
	OpenFileFailed
	ReadFileFailed
	CreateVideoCoverFailed
	CreateFavoriteFailed
	DeleteFavoriteFailed
	ModifyFavoriteFailed
	UploadVideoFailed
	GetVideoInfoFailed
	SendVideoFailed
	CreateSearchRecordFailed
	CreateCommentToVideoFailed
	CreateCommentToUserFailed
	CreateMessageFailed
)
