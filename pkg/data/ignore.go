package data

/*
Why this interface?
Later it would be helpful to ignore on the basis of prefix/suffix etc
by having different implementations
Or maybe we would like to ignore based on REGEX
*/
type ScanIgnorable interface {
	ShouldIgnore(string) (bool, error)
}
