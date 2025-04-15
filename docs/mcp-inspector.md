# MCP inspector

ref: <https://modelcontextprotocol.io/docs/tools/inspector#type-script>

## inspect @modelcontextprotocol/server-postgres

```sh
npx -y @modelcontextprotocol/inspector npx @modelcontextprotocol/server-postgres postgresql://postgres:postgres@127.0.0.1:5432/mcpserver
```

## inspect @modelcontextprotocol/server-redis

Try fetch url: `https://jsonplaceholder.typicode.com/todos/1`

```sh
npx -y @modelcontextprotocol/inspector uvx mcp-server-fetch
```
