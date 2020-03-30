package command

// Login command accepts firebase authentication payload
type Login struct {
	User struct {
		UID           string      `json:"uid"`
		DisplayName   string      `json:"displayName"`
		PhotoURL      string      `json:"photoURL"`
		Email         string      `json:"email"`
		EmailVerified bool        `json:"emailVerified"`
		PhoneNumber   interface{} `json:"phoneNumber"`
		IsAnonymous   bool        `json:"isAnonymous"`
		TenantID      interface{} `json:"tenantId"`
		ProviderData  []struct {
			UID         string      `json:"uid"`
			DisplayName string      `json:"displayName"`
			PhotoURL    string      `json:"photoURL"`
			Email       string      `json:"email"`
			PhoneNumber interface{} `json:"phoneNumber"`
			ProviderID  string      `json:"providerId"`
		} `json:"providerData"`
		APIKey          string `json:"apiKey"`
		AppName         string `json:"appName"`
		AuthDomain      string `json:"authDomain"`
		StsTokenManager struct {
			APIKey         string `json:"apiKey"`
			RefreshToken   string `json:"refreshToken"`
			AccessToken    string `json:"accessToken"`
			ExpirationTime int64  `json:"expirationTime"`
		} `json:"stsTokenManager"`
		RedirectEventID interface{} `json:"redirectEventId"`
		LastLoginAt     string      `json:"lastLoginAt"`
		CreatedAt       string      `json:"createdAt"`
		MultiFactor     struct {
			EnrolledFactors []interface{} `json:"enrolledFactors"`
		} `json:"multiFactor"`
	} `json:"user"`
	Credential struct {
		ProviderID       string `json:"providerId"`
		SignInMethod     string `json:"signInMethod"`
		OauthIDToken     string `json:"oauthIdToken"`
		OauthAccessToken string `json:"oauthAccessToken"`
	} `json:"credential"`
	OperationType      string `json:"operationType"`
	AdditionalUserInfo struct {
		ProviderID string `json:"providerId"`
		IsNewUser  bool   `json:"isNewUser"`
		Profile    struct {
			Name          string `json:"name"`
			GrantedScopes string `json:"granted_scopes"`
			ID            string `json:"id"`
			VerifiedEmail bool   `json:"verified_email"`
			GivenName     string `json:"given_name"`
			Locale        string `json:"locale"`
			FamilyName    string `json:"family_name"`
			Email         string `json:"email"`
			Picture       string `json:"picture"`
		} `json:"profile"`
	} `json:"additionalUserInfo"`
}
