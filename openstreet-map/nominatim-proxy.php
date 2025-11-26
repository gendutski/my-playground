<?php

header('Content-Type: application/json');
header("Access-Control-Allow-Origin: *");

if (!isset($_GET['q'])) {
    echo json_encode(['error' => 'query param q is required']);
    exit;
}

$q = urlencode($_GET['q']);

$url = "https://nominatim.openstreetmap.org/search?format=jsonv2&q={$q}";
sleep(1);

// Gunakan curl
$ch = curl_init($url);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
curl_setopt($ch, CURLOPT_USERAGENT, 'MyApp/1.0'); // WAJIB
$response = curl_exec($ch);

if ($response === false) {
    http_response_code(500);
    echo json_encode(['error' => curl_error($ch)]);
    curl_close($ch);
    exit;
}

curl_close($ch);

// langsung echo hasil dari nominatim
echo $response;
