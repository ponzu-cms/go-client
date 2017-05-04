# Ponzu HTTP Client - Go

## Work in Progress
- Currently no way to add headers or anything to the client request
- Needs a way to add Files to the multipartForm func
- API in flux, alpha software

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


    //------------------------------------------------------------------
    // GET content (single item)
    //------------------------------------------------------------------

    // fetch single Content item of type Blog by ID 1
    resp, err := cms.Content("Blog", 1)
    if err != nil {
        fmt.Println("Blog:1 error:", err)
        return
    }

    fmt.Println(resp.Data[0]["title"])
	fmt.Println(fieldName)



    //------------------------------------------------------------------
    // GET contents (multiple items)
    //------------------------------------------------------------------

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



    //------------------------------------------------------------------
    // SEARCH content
    //------------------------------------------------------------------

    // fetch the search results for a query "Steve" from Content items of type Blog
    resp, err = cms.Search("Blog", "Steve")
    if err != nil {
        fmt.Println("Blog:search error:", err)
        return
    }

    fmt.Println(resp.Data[0]["title"])
    fmt.Println(resp.Data[1]["title"])



    //------------------------------------------------------------------
    // GET File metadata (single item)
    //------------------------------------------------------------------

    // fetch file metadata for uploaded file with slug "mudcracks-mars.jpg" (slug is normalized filename)
    resp, err = cms.FileBySlug("mudcracks-mars.jpg")
    if err != nil {
        fmt.Println("File:slug error:", err)
        return
    }

    fmt.Println(resp.Data[0]["name"])
    fmt.Println(resp.Data[0]["content_type"])
    fmt.Println(resp.Data[0]["content_length"])



    //------------------------------------------------------------------
    // CREATE content (single item)
    //------------------------------------------------------------------

    // create Content item of type Blog with data
    data := make(url.Values)
    data.Set("title", "Added via API")
    data.Set("body", "<p>i'm not sure about this.</p>")
    data.Set("author", "Steve")

    // nil indicates no data params are filepaths, 
    // otherwise would be a []string of key names that are filepaths (docs coming)
    resp, err = cms.Create("Blog", data, nil)
    if err != nil {
        fmt.Println("Create:Blog error:", err)
        return
    }

    fmt.Println(resp.Data[0]["status"], resp.Data[0]["id"])
    id := int(resp.Data[0]["id"].(float64))



    //------------------------------------------------------------------
    // UPDATE content (single item)
    //------------------------------------------------------------------

    // update Content item of type Blog and ID {id} with data
    data = make(url.Values)
    data.Set("title", "Added then updated via API")
    data.Set("author", "API Steve")

    resp, err = cms.Update("Blog", id, data, nil)
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



    //------------------------------------------------------------------
    // DELETE content (single item)
    //------------------------------------------------------------------

    // delete Content item of type Blog with ID {id}
    resp, err = cms.Delete("Blog", id)
    if err != nil {
        fmt.Println("Delete:Blog:#id error:", err, id)
        return
    }

    fmt.Println(resp.Data[0]["status"], resp.Data[0]["id"])
}

```

Alternatively, the `resp` return value (which is an `*APIResponse` type) contains the response body's original
`[]byte` as `resp.JSON`, which you can use to unmarshal to content structs from
your Ponzu's `content` package.