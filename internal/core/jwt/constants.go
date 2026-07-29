package jwt

import (
	"time"
)

const COOKIE_KEY_REFRESH_TOKEN string = "refresh_token"
const HEADER_KEY_ACCESS_TOKEN string = "X-Access-Token"
const CONTEXT_KEY_PAYLOAD string = "payload"
const TOKEN_TYPE_ACCESS string = "access"
const TOKEN_TYPE_REFRESH string = "refresh"
const ACCESS_TOKEN_EXPIRES time.Duration = 15 * 60
const REFRESH_TOKEN_EXPIRES time.Duration = 14 * 24 * 60 * 60
