
package curl

/*
#cgo linux pkg-config: libcurl
#include <curl/curl.h>
static CURLSHcode curl_share_setopt_int(CURLSH *handle, CURLSHoption option, int parameter) {
  return curl_share_setopt(handle, option, parameter);
}
static CURLSHcode curl_share_setopt_pointer(CURLSH *handle, CURLSHoption option, void *parameter) {
  return curl_share_setopt(handle, option, parameter);
}
*/
import "C"
import (
	"unsafe"
	"os"
)

// implement os.Error interface
type CurlShareError C.CURLMcode

func (e CurlShareError) String() string {
	// ret is const char*, no need to free
	ret := C.curl_share_strerror(C.CURLSHcode(e))
	return C.GoString(ret)
}


func newCurlShareError(errno C.CURLSHcode) os.Error {
	if errno == C.CURLSHE_OK {		// if nothing wrong
		return nil
	}
	return CurlShareError(errno)
}


type CURLSH struct {
	handle unsafe.Pointer
}

func ShareInit() *CURLSH {
	p := C.curl_share_init()
	return &CURLSH{p}
}

func (shcurl *CURLSH) Cleanup() os.Error {
	p := shcurl.handle
	return newCurlShareError(C.curl_share_cleanup(p))
}

func (shcurl *CURLSH) Setopt(opt int, param interface{}) os.Error {
	p := shcurl.handle
	if param == nil {
		return newCurlShareError(C.curl_share_setopt_pointer(p, C.CURLSHoption(opt), nil))
	}
	switch opt {
//	case SHOPT_LOCKFUNC, SHOPT_UNLOCKFUNC, SHOPT_USERDATA:
//		panic("not supported")
	case SHOPT_SHARE, SHOPT_UNSHARE:
		if val, ok := param.(int); ok {
			return newCurlShareError(C.curl_share_setopt_int(p, C.CURLSHoption(opt), C.int(val)))
		}
	}
	panic("not supported CURLSH.Setopt opt or param")
	return nil
}
