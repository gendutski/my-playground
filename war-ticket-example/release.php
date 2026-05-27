<?php

session_start();

$redis = require 'redis.php';

$sessionId = session_id();

$redis->del("flash:hold:$sessionId");

$redis->zRem('flash:active', $sessionId);

echo "Slot dilepas";