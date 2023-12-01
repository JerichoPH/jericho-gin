package types

type TimeFormat string

const (
	TIME_FORMAT_DATETIME    TimeFormat = "2006-01-02 15:04:05"
	TIME_FORMAT_DATE        TimeFormat = "2006-01-02"
	TIME_FORMAT_TIME        TimeFormat = "15:04:05"
	TIME_FORMAT_YEAR        TimeFormat = "2006"
	TIME_FORMAT_MONTH       TimeFormat = "01"
	TIME_FORMAT_DAY         TimeFormat = "02"
	TIME_FORMAT_HOUR        TimeFormat = "15"
	TIME_FORMAT_MINUTE      TimeFormat = "04"
	TIME_FORMAT_SECOND      TimeFormat = "05"
	TIME_FORMAT_YEAR_MONTH  TimeFormat = "2006-01"
	TIME_FORMAT_HOUR_MINUTE TimeFormat = "15:04"
)

type AccountField string

const (
	ACCOUNT_ID                    AccountField = "***ACCOUNT.ID***"
	ACCOUNT_ACCOUNT               AccountField = "***ACCOUNT.USERNAME***"
	ACCOUNT_OPEN_ID               AccountField = "***ACCOUNT.UUID***"
	ACCOUNT_NICKNAME              AccountField = "***ACCOUNT.NICKNAME***"
	ACCOUNT_WORK_AREA_UNIQUE_CODE AccountField = "***ACCOUNT.WORK_AREA_UNIQUE_CODE***"
	ACCOUNT_AUTHORIZATION         AccountField = "***ACCOUNT.AUTHORIZATION***"
)

type Marker string

const (
	MARKER_ORIGINAL = "ORIGINAL"
	MARKER_FINISHED = "FINISHED"
)

type ResponseMethod string

const (
	RESPONSE_METHOD_OK                ResponseMethod = "OK"
	RESPONSE_METHOD_DATUM             ResponseMethod = "DATUM"
	RESPONSE_METHOD_PAGER             ResponseMethod = "PAGER"
	RESPONSE_METHOD_JQUERY_DATA_PAGER ResponseMethod = "JQUERY-DATA-PAGER"
	RESPONSE_METHOD_CREATED           ResponseMethod = "CREATED"
	RESPONSE_METHOD_UPDATED           ResponseMethod = "UPDATED"
	RESPONSE_METHOD_DELETED           ResponseMethod = "DELETED"
)
