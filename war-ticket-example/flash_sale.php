<?php

session_start();

$redis = require 'redis.php';

$sessionId = session_id();

$limit = 1000;
$ttl = 900;

$now = time();
$expireAt = $now + $ttl;

$lua = file_get_contents('flash_sale.lua');

$result = $redis->eval(
    $lua,
    [
        'flash:active',
        "flash:hold:$sessionId",

        $now,
        $expireAt,
        $limit,
        $sessionId,
        $ttl
    ],
    2
);

$status = $result[0];
$count = $result[1];

switch ($status) {

    case 0:
        echo "Flash sale penuh. Active: $count";
        break;

    case 1:
        echo "Berhasil masuk flash sale. Active: $count";
        break;

    case 2:
        echo "Kamu sudah punya slot.";
        break;
}