-- KEYS[1] = sorted set active key
-- KEYS[2] = hold key

-- ARGV[1] = current timestamp
-- ARGV[2] = expire timestamp
-- ARGV[3] = limit
-- ARGV[4] = session id
-- ARGV[5] = ttl seconds

local activeKey = KEYS[1]
local holdKey = KEYS[2]

local now = tonumber(ARGV[1])
local expireAt = tonumber(ARGV[2])
local limit = tonumber(ARGV[3])
local sessionId = ARGV[4]
local ttl = tonumber(ARGV[5])

-- cleanup expired user
redis.call('ZREMRANGEBYSCORE', activeKey, '-inf', now)

-- hitung active user
local current = redis.call('ZCARD', activeKey)

if current >= limit then
    return {
        0,
        current
    }
end

-- cek apakah user sudah punya hold
local exists = redis.call('EXISTS', holdKey)

if exists == 1 then
    return {
        2,
        current
    }
end

-- simpan hold
redis.call('SET', holdKey, 1, 'EX', ttl)

-- masukkan ke sorted set
redis.call('ZADD', activeKey, expireAt, sessionId)

return {
    1,
    current + 1
}