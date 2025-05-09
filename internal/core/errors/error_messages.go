package errors

var (
	// User related errors
	ErrUserNotFound         = New(NotFoundError, "User not found", "کاربر یافت نشد", nil)
	ErrDuplicatePhoneNumber = New(ValidationError, "Phone number already exists", "این شماره تلفن قبلاً ثبت شده است", nil)
	ErrDuplicateEmail       = New(ValidationError, "Email already exists", "این ایمیل قبلاً ثبت شده است", nil)
	ErrUpdateUser           = New(InternalError, "Failed to update user", "خطا در به\u200cروزرسانی کاربر", nil)
	ErrCreateUser           = New(InternalError, "Failed to create user", "خطا در ایجاد کاربر", nil)

	// Database related errors
	ErrDatabaseInit = New(InternalError, "Failed to initialize database", "خطا در راه\u200cاندازی پایگاه داده", nil)
	ErrRedisInit    = New(InternalError, "Failed to initialize redis", "خطا در راه\u200cاندازی redis", nil)
	ErrGetUsers     = New(InternalError, "Failed to get users", "خطا در دریافت لیست کاربران", nil)
	ErrGetUser      = New(InternalError, "Failed to get user", "خطا در دریافت اطلاعات کاربر", nil)

	// Admin related errors
	ErrChangeRole   = New(InternalError, "Failed to change user role", "خطا در تغییر نقش کاربر", nil)
	ErrChangeStatus = New(InternalError, "Failed to change user status", "خطا در تغییر وضعیت کاربر", nil)
	ErrDeleteUser   = New(InternalError, "Failed to delete user", "خطا در حذف کاربر", nil)
	ErrForbidden    = New(AuthorizationError, "Access denied", "شما دسترسی لازم برای انجام این عملیات را ندارید", nil)

	// General errors
	ErrInternalServer    = New(InternalError, "Internal server error", "خطای داخلی سرور", nil)
	ErrInvalidRequest    = New(ValidationError, "Invalid request", "درخواست نامعتبر است", nil)
	ErrInvalidUserID     = New(ValidationError, "Invalid user ID", "شناسه کاربر نامعتبر است", nil)
	ErrInvalidUserIDType = New(ValidationError, "Invalid user ID format", "فرمت شناسه کاربر نامعتبر است", nil)

	// Authentication related errors
	ErrInvalidCredentials   = New(AuthenticationError, "Invalid credentials", "نام کاربری یا رمز عبور اشتباه است", nil)
	ErrAccountDeactivated   = New(AuthenticationError, "Account is deactivated", "حساب کاربری غیرفعال است", nil)
	ErrUserNotAuthenticated = New(AuthenticationError, "Authentication required", "لطفاً ابتدا وارد حساب کاربری خود شوید", nil)

	// Token related errors
	ErrInvalidToken       = New(AuthenticationError, "Invalid token", "توکن نامعتبر است", nil)
	ErrTokenCreation      = New(InternalError, "Failed to create token", "خطا در ایجاد توکن", nil)
	ErrRemoveToken        = New(InternalError, "Failed to remove token", "خطا در حذف توکن", nil)
	ErrGetToken           = New(InternalError, "Failed to get token", "خطا در دریافت توکن", nil)
	ErrTokenNotFound      = New(NotFoundError, "Token not found", "توکن یافت نشد", nil)
	ErrAddToken           = New(InternalError, "Failed to add token", "خطا در اضافه کردن توکن", nil)
	ErrRefreshToken       = New(InternalError, "Failed to refresh token", "خطا در تجدید توکن", nil)
	ErrMissingAuthHeader  = New(AuthenticationError, "Authorization header is required", "هدر احراز هویت الزامی است", nil)
	ErrParseToken         = New(AuthenticationError, "Failed to parse token", "خطا در تجزیه توکن", nil)
	ErrInvalidTokenClaims = New(AuthenticationError, "Invalid token claims", "اطلاعات توکن نامعتبر است", nil)
	ErrInvalidTokenType   = New(AuthenticationError, "Invalid token type", "نوع توکن نامعتبر است", nil)

	// User operation errors
	ErrLogin          = New(AuthenticationError, "Failed to login", "خطا در ورود", nil)
	ErrLogout         = New(InternalError, "Failed to logout", "خطا در خروج", nil)
	ErrChangePassword = New(InternalError, "Failed to change password", "خطا در تغییر رمز عبور", nil)

	// Configuration related errors
	ErrLoadConfig = New(InternalError, "Failed to load configuration", "خطا در بارگذاری تنظیمات", nil)

	// Validation errors
	ErrInvalidSortField    = New(ValidationError, "Sort field is invalid", "فیلد مرتب\u200cسازی نامعتبر است", nil)
	ErrInvalidRoleField    = New(ValidationError, "Role field is invalid", "فیلد نقش نامعتبر است", nil)
	ErrInvalidStatusField  = New(ValidationError, "Status field is invalid", "فیلد وضعیت نامعتبر است", nil)
	ErrInvalidOrderField   = New(ValidationError, "Order field is invalid", "فیلد ترتیب نامعتبر است", nil)
	ErrInvalidPhoneNumber  = New(ValidationError, "Phone number must start with 09 and be 11 digits", "شماره موبایل باید با 09 شروع شده و 11 رقم باشد", nil)
	ErrInvalidFirstName    = New(ValidationError, "First name field is invalid", "فیلد نام کوچک نامعتبر است", nil)
	ErrInvalidLastName     = New(ValidationError, "Last name field is invalid", "فیلد نام خانوادگی نامعتبر است", nil)
	ErrInvalidEmail        = New(ValidationError, "Email field is invalid", "فیلد ایمیل نامعتبر است", nil)
	ErrInvalidPassword     = New(ValidationError, "Password must be at least 8 characters and include uppercase, lowercase, and a number", "رمز عبور باید حداقل ۸ کاراکتر و شامل حروف بزرگ، کوچک و عدد باشد", nil)
	ErrInvalidOldPassword  = New(ValidationError, "Old password is invalid", "رمز عبور قدیمی نامعتبر است", nil)
	ErrInvalidNewPassword  = New(ValidationError, "New password must be at least 8 characters and include uppercase, lowercase, and a number", "رمز عبور جدید باید حداقل ۸ کاراکتر و شامل حروف بزرگ، کوچک و عدد باشد", nil)
	ErrInvalidRefreshToken = New(ValidationError, "Refresh token is invalid", "توکن بروزرسانی نامعتبر است", nil)
)
