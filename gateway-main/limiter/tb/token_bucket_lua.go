package tb

import "github.com/redis/go-redis/v9"

var AllowN = redis.NewScript(AllowNScript)

var AllowNScript = `
local key = KEYS[1]
local rate=tonumber(ARGV[1])
local burst=tonumber(ARGV[2])
local count=tonumber(ARGV[3])
local expires=tonumber(ARGV[4])

redis.replicate_commands()
local t=redis.call('TIME')
local now = t[1] + t[2]/1000000
local exists=redis.call('EXISTS',key)
if exists == 0 then
    redis.call('hset',key,"rate",rate)
    redis.call('hset',key,"burst",burst)
    redis.call('hset',key,"lastBucket",0)
    -- 默认‘初次发放’是在10秒前
    redis.call('hset',key,"lastTime",now-10)
	redis.call('hset',key,'bucket',0)
end

redis.call('expire',key,expires)

local lastTime=redis.call('hget',key,"lastTime")
local lastBucket=redis.call('hget',key,"lastBucket")
local elasped=now - lastTime
local delta= elasped * rate 
local bucket= redis.call('hget',key,"bucket")

-- 桶里现在有几个ticket？
if ( delta + bucket  < burst) 
then
    bucket=delta+bucket
else
	bucket=burst
end

if ( count <= bucket )then
--	local msg="[yes]elasped:" .. elasped .. ",delta:" .. delta .. ",bucket:" .. bucket
--	redis.log(redis.LOG_WARNING, msg)
    lastBucket=lastBucket+count
    lastTime=now
    bucket=bucket-count
	redis.call('hset',key,"lastBucket",lastBucket)
	redis.call('hset',key,"lastTime",lastTime)
	redis.call('hset',key,"bucket",bucket) -- 剩余 token
    return {count, bucket}
else
    -- 不能发
--	local msg="[no]elasped:" .. elasped .. ",delta:" .. delta .. ",bucket:" .. bucket
--	redis.log(redis.LOG_WARNING, msg)
    return {-1, bucket} 
end
`
