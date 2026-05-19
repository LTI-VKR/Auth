package application

type OAuthProvider string

const (
	GoogleProvider OAuthProvider = "google"
	VkProvider     OAuthProvider = "vk"
)

// IsValid проверяет, валиден ли провайдер
func (p OAuthProvider) IsValid() bool {
	switch p {
	case GoogleProvider, VkProvider:
		return true
	default:
		return false
	}
}

func (p OAuthProvider) String() string {
	return string(p)
}

func AllProviders() []OAuthProvider {
	return []OAuthProvider{
		GoogleProvider, VkProvider,
	}
}
