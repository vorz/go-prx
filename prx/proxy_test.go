package prx

import "testing"

func ProxyTest(t *testing.T) {
	t.Error("YEAH, YOU FAIL")
}

func init() {
	//	http.Handle("/bobo", ConstantHanlder("bobo"))
	//	http.Handle("/query", QueryHandler{})
}
