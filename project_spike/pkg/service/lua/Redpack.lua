local redpackId = ARGV[1]
local userId = ARGV[2]
local redpackKey = KEYS[1] .. redpackId
local redpackList = KEYS[2] .. redpackId
local redpackSet = KEYS[3] .. redpackId
local res = redis.call('llen', redpackKey)
if (tonumber(res) <= 0) then
    -- 超卖了，没有库存了
    return 1
end
if redis.call("sismember", redpackSet, userId) == 1 then
    -- 用户已经买过了
    return 2
end
-- 记录用户已经买过的信息
redis.call('sadd', redpackSet, userId)
-- 给用户发一个红包并且记录下来
local money = redis.call('lpop', redpackKey)
redis.call('lpush', redpackList, tostring(money))
return 0