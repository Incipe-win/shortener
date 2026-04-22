-- wrk Lua script for /api/convert endpoint
-- Usage: wrk -t4 -c100 -d30s -s scripts/bench/convert.lua http://localhost:8888

request = function()
    local urls = {
        "https://github.com/zeromicro/go-zero",
        "https://golang.org/doc/install",
        "https://react.dev/learn",
        "https://docs.docker.com/compose/",
        "https://redis.io/docs/latest/",
        "https://kafka.apache.org/documentation/",
        "https://prometheus.io/docs/introduction/overview/",
        "https://grafana.com/docs/grafana/latest/",
    }
    local url = urls[math.random(#urls)]
    local body = '{"long_url":"' .. url .. '"}'

    return wrk.format("POST", "/api/convert", {
        ["Content-Type"] = "application/json",
    }, body)
end
