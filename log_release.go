//go:build !debug

package geom2d

// Debug is a no-op in production builds.
func logDebugf(
	_ string,
	_ ...interface{},
) {
}
