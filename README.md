# Trace
Simple error tracing library for go

# Installation
```shell
go get -u github.com/gabehardgrave/trace
```

# About

An alternative to `fmt.Errorf` and `errors.Wrap`. 

```go
x, err := db.GetX()
if err != nil {
  return 0, trace.Wrap(err)
}
```

`trace.Wrap(err)` adds file name and line # information to the resulting `error`.

See `/tests` for more examples and use cases.
