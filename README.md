# gen-merge
A library to generate merging functions for structs in golang. WIP.

For creating a test suite of variations on a base configuration, it's convenient to use the prototype pattern.

This library generates a Merge function per struct type to apply the values of one struct onto another. Doing this
without reflection minimizes runtime errors.
