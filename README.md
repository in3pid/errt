Embeds stack traces to errors and simplifies debugging code. Eg.:
```
func (s *Service) Method(arg int) (foo *Foo, err error) {
	defer func() {
		if err != nil {
			glog.Error(errt.TraceDeferred(err)) // line 19
		}
	}()
	if err = a(); err != nil {
		return nil, err
	}
	if err = b(); err != nil {
		return nil, err // line 26
	}
	// more errors to check...
}
```
Error output:
```
E0827 19:20:12.147815    3733 foo.go:19] some error
foo.go:26 foo.(*Service).Method
foo:67 foo.(*Service).StateLoop
asm_amd64.s:2233 runtime.goexit
```
