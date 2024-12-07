package response

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/lascape/gopkg/response/ecode"
)

type InterceptorFunc func(ctx *gin.Context, p *Response)

var (
	respPool = &sync.Pool{
		New: func() interface{} {
			return &Response{}
		},
	}
)

func getResponseInstance(ctx *gin.Context, interceptors ...InterceptorFunc) {
	p := respPool.Get().(*Response)
	for _, interceptor := range interceptors {
		interceptor(ctx, p)
	}
	returnValue(ctx, p)
	p.init()
	respPool.Put(p)
}

// 实现默认的返回值拦截器
var returnValue = InterceptorFunc(func(ctx *gin.Context, p *Response) {
	ctx.JSON(200, p)
})

func traceInterceptor() InterceptorFunc {
	return func(ctx *gin.Context, p *Response) {
		traceId, ok := ctx.Get("trace_id")
		if ok {
			p.TraceId = traceId.(string)
		}
	}
}

// 实现接口返回success拦截器
func successInterceptor(data interface{}) InterceptorFunc {
	return InterceptorFunc(func(ctx *gin.Context, p *Response) {
		p = p.withData(data).withResult(true).withError(ecode.Wrap(ecode.Success, nil))
	})
}

// 实现接口返回error拦截器，注入错误值
func errorInterceptor(err error, data ...interface{}) InterceptorFunc {
	return InterceptorFunc(func(ctx *gin.Context, p *Response) {
		p = p.withResult(false)
		if len(data) == 1 {
			p = p.withData(data[0])
		} else {
			p = p.withData(data)
		}

		if e1, ok := err.(*ecode.ErrorX); ok {
			p = p.withError(e1)
		} else if e2, ok := err.(*ecode.Errno); ok {
			p = p.withError(ecode.Wrap(e2, nil))
		} else {
			p = p.withError(ecode.Wrap(ecode.Error, err))
		}
	})
}
