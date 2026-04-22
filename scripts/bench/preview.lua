-- wrk Lua script for /api/preview/:short_url endpoint
-- Usage: SHORT_URLS="abc,def,xyz" wrk -t4 -c100 -d30s -s scripts/bench/preview.lua http://localhost:8888

local short_urls = {}

init = function(args)
    local env_urls = os.getenv("SHORT_URLS")
    if env_urls and env_urls ~= "" then
        for s in string.gmatch(env_urls, "([^,]+)") do
            table.insert(short_urls, s)
        end
    else
        table.insert(short_urls, "test1")
    end
end

request = function()
    local surl = short_urls[math.random(#short_urls)]
    return wrk.format("GET", "/api/preview/" .. surl)
end
