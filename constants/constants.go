package constants

const AuthSignInPath = "/auth/sign-in"
const AuthSignUpPath = "/auth/sign-up"
const AuthSignOutPath = "/auth/sign-out"
const AuthSubmitSignInPath = "/auth/submit-sign-in"
const AuthSubmitSignUpPath = "/auth/submit-sign-up"
const AuthEnterCodePath = "/auth/enter-code"
const AuthSubmitCodePath = "/auth/submit-code"
const AuthResubmitCodePath = "/auth/resubmit-code"
const AuthCancelPath = "/auth/cancel-sign-in"

const IndexPath = "/"
const LayoutsMainPath = "/layouts/main"
const ProtectedPath = "/protected"

const CameFromKey = "came-from-key"
const EmailKey = "email-key"
const ErrorKey = "error-key"
const ExpectedCodeKey = "expected-code-key"
const IsSignedInKey = "is-signed-in-key"
const MessageKey = "message-key"
const SubmitTimeKey = "submit-time-key"
const UrlToReturnToKey = "url-to-return-to-key"
const WrongCodeEnteredCount = "wrong-code-entered-count"

const IsSignedInValue = "TrUe"

const CodeExpireTimeInSeconds = 20 * 60
const WrongCodeFailureCount = 3
