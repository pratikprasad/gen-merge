# gen-merge

## Overview 
A library to generate merging functions for structs in golang. WIP.

For creating a test suite of variations on a base configuration, it's convenient to use the prototype pattern.

This library generates a Merge and MergeOverride function for each struct. 

For each unspecified field on source, copy the value from target
  
```
  func (source <T>) Merge(target <T>) <T>
```
For each specified field on target, copy the value to source
```
  func (source <T>) MergeOverride(target <T>) <T>
```

Right now, a field is considered unspecified if the value of the field corresponds to the zero value of that field's type.

## Final State

Given a project, for every package, create a `genmerge.go` file containing merge functions for every exported type.

If the fields for any particular struct implement the merge functions, then the merging should rely not on the zero value, but call the corresponding merge function directly.

## Current State

Given a package, generate merge functions for all (including unexported) structs that rely on zero-value as opposed to "recursive" merging. 
