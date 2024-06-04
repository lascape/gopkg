## HTTPX

这是通用的请求库，沿用大部分之前版本的逻辑，移除了代码中强制依赖的组件。开发者如果需要提供公共的插件能力，应该将代码写到httpx/plugins中，通过
OnBeforeRequest、OnAfterResponse函数注入到主流程。

## Usage

### GET

```go
resp := New(server.URL + "/update").SetBodyJson(s).Post(context.Background())
```

### POST

```go
resp := New(server.URL + "/update").SetBodyJson(s).Post(context.Background())

```