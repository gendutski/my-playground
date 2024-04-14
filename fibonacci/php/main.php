<?php

// https://www.youtube.com/watch?v=Zo1wu6tO0_g
function fibonacciFormula($n)
{
    $satu = 1;
    $akarLima = sqrt(5);
    $dua = 2;

    $a = ($satu + $akarLima) / $dua;
    $b = ($satu - $akarLima) / $dua;
    $result = (pow($a, $n) - pow($b, $n)) / $akarLima;
    return round($result);
}

function fibonacciLoop($n)
{
    $f = array(0, 1);
    for ($i = 2; $i <= $n; $i++) {
        $f[$i] = $f[$i - 1] + $f[$i - 2];
    }
    return $f[$n];
}

// run
for ($n = 1; $n <= 100; $n++) {
    $formula = fibonacciFormula($n);
    $loop = fibonacciLoop($n);
    $isEqual = $formula == $loop;
    printf("Nth = %d, Formula = %d, Loop = %d, Is equal = %s\n", $n, $formula, $loop, $isEqual ? 'true' : 'false');
    if (!$isEqual) {
        break;
    }
}
