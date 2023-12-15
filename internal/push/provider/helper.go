package provider

import "allaboutapps.dev/aw/go-starter/internal/push"

// DEPRECATED: sendMulticastWithProvider
// Allows to send same notification to multiple receivers.
//
// This helper function is deprecated and might be removed with future releases.
// Please use sendMulticastWithProvider instead defined in push package.
func sendMulticastWithProvider(p push.Provider, tokens []string, title string, message string, data map[string]string, silent bool, collapseKey ...string) []push.ProviderSendResponse {
	responseSlice := make([]push.ProviderSendResponse, 0)

	for _, token := range tokens {
		responseSlice = append(responseSlice, p.Send(token, title, message, data, silent, collapseKey...))
	}

	return responseSlice
}
