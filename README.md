# Ponzu HTTP Client - Go

### Usage
```go
package main

import (
	"fmt"
	"net/url"

	"github.com/ponzu-cms/go-client"
)

func main() {
    // configure the http client
    cfg := client.Config{
        Host:         "http://localhost:8080",
        DisableCache: false, // defaults to false, here for documentation
    }
    cms := client.New(cfg)

    // fetch single Content item of type Blog by ID 1
    resp, err := cms.Content("Blog", 1)
    if err != nil {
        fmt.Println("Blog:1 error:", err)
        return
    }

    fmt.Println(resp.Data[0]["title"])

    // fetch multiple Content items of type Blog, using the default QueryOptions
    // Count: 10
    // Offset: 0
    // Order: DESC
    resp, err = cms.Contents("Blog", client.QueryOptions{})
    if err != nil {
        fmt.Println("Blog:multi error:", err)
        return
    }

    for i := 0; i < len(resp.Data); i++ {
        fmt.Println(resp.Data[i]["title"])
    }

    // fetch the search results for a query "Steve" from Content items of type Blog
    resp, err = cms.Search("Blog", "Steve")
    if err != nil {
        fmt.Println("Blog:search error:", err)
        return
    }

    fmt.Println(resp.Data[0]["title"])
    fmt.Println(resp.Data[1]["title"])

    // fetch file metadata for uploaded file with slug "mudcracks-mars.jpg" (slug is normalized filename)
    resp, err = cms.UploadBySlug("mudcracks-mars.jpg")
    if err != nil {
        fmt.Println("Uploads:slug error:", err)
        return
    }

    fmt.Println(resp.Data[0]["name"])
    fmt.Println(resp.Data[0]["content_type"])
    fmt.Println(resp.Data[0]["content_length"])

    // create Content item of type Blog with data
    data := make(url.Values)
    data.Set("title", "Added via API")
    data.Set("body", "<p>i'm not sure about this.</p>")
    data.Set("author", "Steve")

    resp, err = cms.Create("Blog", data)
    if err != nil {
        fmt.Println("Create:Blog error:", err)
        return
    }

    fmt.Println(resp.Data[0]["status"], resp.Data[0]["id"])
    id := int(resp.Data[0]["id"].(float64))

    // update Content item of type Blog and ID 9 with data
    data = make(url.Values)
    data.Set("title", "Added then updated via API")
    data.Set("body", "<p>i'm not sure about this.</p>")
    data.Set("author", "API Steve")

    resp, err = cms.Update("Blog", 9, data)
    if err != nil {
        fmt.Println("Create:Blog error:", err)
        return
    }

    resp, err = cms.Search("Blog", `"API Steve"`)
    if err != nil {
        fmt.Println("Blog:search error:", err)
        return
    }

    fmt.Println(resp.Data[0]["title"])

    // delete Content item of type Blog with ID {id}
    resp, err = cms.Delete("Blog", id)
    if err != nil {
        fmt.Println("Delete:Blog:#id error:", err, id)
        return
    }

    fmt.Println(resp.Data[0]["status"], resp.Data[0]["id"])
}

```
