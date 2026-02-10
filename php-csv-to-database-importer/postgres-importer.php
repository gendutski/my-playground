#!/usr/bin/php
<?php

require __DIR__ . '/vendor/autoload.php';

// load .env
$dotenv = Dotenv\Dotenv::createImmutable(__DIR__);
$dotenv->load();

// set from env variable
$host = $_ENV['POSTGRESQL_HOST'];
$port = $_ENV['POSTGRESQL_PORT'];
$db = $_ENV['POSTGRESQL_NAME'];
$user = $_ENV['POSTGRESQL_USER'];
$pass = $_ENV['POSTGRESQL_PASS'];

$folder = __DIR__ . "/sql/postgresql";

try {
    $pdo = new PDO(
        "pgsql:host=$host;port=$port;dbname=$db",
        $user,
        $pass,
        [PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION]
    );

    $files = glob($folder . "/*.csv");

    if (!$files) {
        throw new Exception("There is no CSV file in the {$folder} folder");
    }

    foreach ($files as $csvFile) {

        $table = strtolower(pathinfo($csvFile, PATHINFO_FILENAME));
        $table = preg_replace('/^\d+-|\.csv$/', '', $table);

        echo "Importing $table...\n";

        /*
        |--------------------------------------------------------------------------
        | Retrieve PostgreSQL column metadata
        |--------------------------------------------------------------------------
        */
        $metaStmt = $pdo->prepare("
            SELECT column_name, data_type
            FROM information_schema.columns
            WHERE table_name = :table
        ");
        $metaStmt->execute(["table" => $table]);

        $columnTypes = [];
        foreach ($metaStmt->fetchAll(PDO::FETCH_ASSOC) as $col) {
            $columnTypes[$col["column_name"]] = $col["data_type"];
        }

        if (!$columnTypes) {
            throw new Exception("Tabel '$table' not found in database.");
        }

        /*
        |--------------------------------------------------------------------------
        | Open CSV & read header
        |--------------------------------------------------------------------------
        */
        $handle = fopen($csvFile, "r");
        $columns = fgetcsv($handle);

        $placeholders = array_map(fn($c) => ":" . $c, $columns);

        $sql = sprintf(
            'INSERT INTO "%s" ("%s") VALUES (%s) ON CONFLICT DO NOTHING',
            $table,
            implode('","', $columns),
            implode(",", $placeholders)
        );

        $stmt = $pdo->prepare($sql);
        $pdo->beginTransaction();

        $rowCount = 0;

        while (($row = fgetcsv($handle)) !== false) {

            // skip blank lines
            if (count(array_filter($row)) === 0) {
                continue;
            }

            $data = [];

            foreach ($columns as $i => $col) {

                $value = $row[$i] ?? null;
                $type = $columnTypes[$col] ?? null;

                if (is_string($value)) {
                    $value = trim($value);
                }

                // empty → NULL
                if ($value === "" || $value === null) {
                    $value = null;
                }
                // integer
                elseif (in_array($type, ["integer", "bigint", "smallint"])) {
                    $value = (int) $value;
                }
                // numeric / decimal
                elseif (in_array($type, ["numeric", "double precision", "real"])) {
                    $value = (float) $value;
                }
                // boolean
                elseif ($type === "boolean") {
                    $value = in_array(strtolower($value), ["true", "1", "t", "yes"]) ? 'true' : 'false';
                }
                // timestamp/date → leave as a string (PostgreSQL will parse)
                // json/jsonb → leave as a JSON string

                $data[":" . $col] = $value;
            }

            $stmt->execute($data);

            $rowCount++;
            if ($rowCount % 1000 === 0) {
                echo "  $rowCount rows imported...\n";
            }
        }

        $pdo->commit();
        fclose($handle);

        echo "✔ Done $table ($rowCount rows)\n\n";
    }

    echo "=== ALL IMPORTS COMPLETED ===\n";

} catch (Exception $e) {

    if (isset($pdo) && $pdo->inTransaction()) {
        $pdo->rollBack();
    }

    echo "ERROR: " . $e->getMessage() . "\n";
}
