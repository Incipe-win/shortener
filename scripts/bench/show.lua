-- wrk Lua script for /:short_url redirect endpoint
-- Requires SHORT_URLS environment variable or hardcoded URLs
-- Usage: SHORT_URLS="abc,def,xyz" wrk -t4 -c100 -d30s -s scripts/bench/show.lua http://localhost:8888

-- 预填充的短 URL（需确保数据库中存在）
local short_urls = {}

init = function(args)
    -- 从环境变量或默认值读取短 URL 列表
    local env_urls = os.getenv("SHORT_URLS")
    if env_urls and env_urls ~= "" then
        for s in string.gmatch(env_urls, "([^,]+)") do
            table.insert(short_urls, s)
        end
    else
        -- 如果为空，wrk 会在 done 中报告 404
        table.insert(short_urls, "test1")
    end
end

request = function()
    local surl = short_urls[math.random(#short_urls)]
    return wrk.format("GET", "/" .. surl)
end
