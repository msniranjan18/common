package postgres

import "strings"

// maskDatabaseURL masks the password in a PostgreSQL connection string for safe logging
func maskDatabaseURL(url string) string {
	// Standard format: postgres://username:password@localhost:5432/database
	if strings.Contains(url, "@") {
		parts := strings.Split(url, "@")
		if len(parts) == 2 {
			authPart := parts[0] // e.g., "postgres://username:password"
			hostPart := parts[1] // e.g., "localhost:5432/database"

			// Look for the password after the first colon in the auth part
			// Note: we start searching after "postgres://"
			prefix := ""
			if strings.Contains(authPart, "://") {
				schemeParts := strings.SplitN(authPart, "://", 2)
				prefix = schemeParts[0] + "://"
				authPart = schemeParts[1]
			}

			if strings.Contains(authPart, ":") {
				userPass := strings.SplitN(authPart, ":", 2)
				if len(userPass) == 2 {
					// userPass[0] is username, userPass[1] is password
					return prefix + userPass[0] + ":***@" + hostPart
				}
			}

			// If there's an @ but no :, just mask everything before @
			return prefix + "***@" + hostPart
		}
	}
	return url
}
