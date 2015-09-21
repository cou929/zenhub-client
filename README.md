# zenhub-client

[![wercker status](https://app.wercker.com/status/9e8d950f6fd2826cae1d779858a85b79/s/master "wercker status")](https://app.wercker.com/project/bykey/9e8d950f6fd2826cae1d779858a85b79)

Unofficial client library of [ZenHub.io](https://www.zenhub.io/) written in golang.

## Usage

```go
import zenhub "github.com/cou929/zenhub-client"

authToken := "your auth token"
organization := "your_org"
repository := "your_repo"
client := zenhub.NewClient(authToken, organization, repository)

issueNumber := 1234
pipelines, err := client.GetPipelines(issueNumber)

dstPipelineName := "Todo"
res, err := client.UpdatePipeline(issueNumber, dstPipelineName)
```

## CAUTION

This library is very alpha version because ZenHub.io has no public API. The interfaces would have breaking chages in the future.

## License

MIT
