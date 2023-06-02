
## zero-example-base

-------

go-zero示例微服务的基础库

有自定义的表单验证方法，拦截器，以及异常处理

## 错误处理

在业务编码中，需要抛出err时，必须要抛出自定义错误`BusinessError`,此结构包含`code`,`msg`,`code`定义了
几个通用的编码，业务错误码必须大于200000，大于200000的公用编码有：

const ErrorBusiness uint32 = 200000     用于普通业务错误，http状态码为400

const ErrorNotFound = 200001            用于表示未找到对象等，http状态码为404

const ErrorRpcOther uint32 = 200002     用于表示RPC服务的系统错误（非业务错误），http状态码为400


如果有必要，可以自行表达错误码，无需写入基础库，但必须按要求大于200000。

logic异常时返回error：
```go

//返回不带堆栈跟踪的错误
return nil, xerr.NewBusinessError(xerr.SetCode(xerr.ErrorNotFound), xerr.SetMsg("订单不存在"))

//返回带堆栈跟踪的错误
return nil, errors.Wrapf(xerr.NewBusinessError(xerr.SetCode(xerr.ErrorNotFound), xerr.SetMsg("订单不存在")), "在查询订单数据库时错误 %+v", err)

```

业务代码返回的err，作为`response.Response(r, w, resp, err)`的入参，`response.Response`会自行把结果转换成`application/json`编码的http响应。

##处理http应答的方法
1. response.ParseParamErrResponse 参数解析错误时使用

```go
var req types.LoginReq
if err := httpx.Parse(r, &req); err != nil {
    response.ParseParamErrResponse(r, w, err)
    return
}
```
2. response.ValidateErrResponse 验证错误时使用
```go
if err := svcCtx.Validator.Validate.StructCtx(r.Context(), req); err != nil {
    response.ValidateErrResponse(r, w, err, svcCtx.Validator.Trans)
    return
}
```
3. 处理logic返回值时使用
```go
resp, err := l.Login(&req)
response.Response(r, w, resp, err)
```

##模板 

api handler.tpl
示例项目使用了第三方验证库，并封装成方法`custom_validate.InitValidator()`，在NewServiceContext时实例化，如果不喜欢这种方式，可以在需要的Handler中实例化。
```
package {{.PkgName}}

import (
	"net/http"
    "go-zero-base/utils/response"
	"github.com/zeromicro/go-zero/rest/httpx"
	{{.ImportPackages}}
)

func {{.HandlerName}}(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		{{if .HasRequest}}var req types.{{.RequestType}}
		if err := httpx.Parse(r, &req); err != nil {
		    response.ParseParamErrResponse(r, w, err)
			return
		}

		if err := svcCtx.Validator.Validate.StructCtx(r.Context(), req); err != nil {
            response.ValidateErrResponse(r, w, err, svcCtx.Validator.Trans)
            return
        }

		{{end}}l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx)
		{{if .HasResp}}resp, {{end}}err := l.{{.Call}}({{if .HasRequest}}&req{{end}})
		{{if .HasResp}}response.Response(r, w, resp, err){{else}}response.Response(r, w, nil, err){{end}}
	}
}
```

api context.tpl
示例项目使用的`gorm`作为数据库连接及交互框架，所以模板如下：
```
package svc

import (
    "go-zero-base/custom_validate"
	{{.configImport}}
)

type ServiceContext struct {
	Config {{.config}}
	DbEngine *gorm.DB
	Validator custom_validate.Validator //验证器
	{{.middleware}}
}

func NewServiceContext(c {{.config}}) *ServiceContext {
    db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{
        NamingStrategy: schema.NamingStrategy{
            TablePrefix:   "greet_", // 表名前缀，`User` 的表名应该是 `t_users`
            SingularTable: true,    // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
        },
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        panic(err)
    }


	return &ServiceContext{
		Config: c,
		DbEngine: db,
		Validator: custom_validate.InitValidator(),
		{{.middlewareAssignment}}
	}
}

```
