package attributes

type ErrorListDataContainer struct {
   Errors []ErrorData `json:"errors"`
}

type ErrorData struct {
   Status int               `json:"status"`
   Code   string            `json:"code"`
   Title  string            `json:"title"`
   Detail string            `json:"detail"`
   Meta   map[string]string `json:"meta"`
}
