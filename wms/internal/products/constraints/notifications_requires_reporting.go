//go:build notifications && !reporting

package constraints

var _ = doesNotCompileNotificationsRequiresReporting
