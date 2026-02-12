#!/usr/bin/php
<?php

require __DIR__ . '/vendor/autoload.php';

// load .env
$dotenv = Dotenv\Dotenv::createImmutable(__DIR__);
$dotenv->load();

// set from env variable
$host = $_ENV['MYSQL_HOST'];
$port = $_ENV['MYSQL_PORT'] ?? 3306;
$db = $_ENV['MYSQL_NAME'];
$user = $_ENV['MYSQL_USER'];
$pass = $_ENV['MYSQL_PASS'];

$folder = __DIR__ . "/sql/mysql";

try {
    $pdo = new PDO(
        "mysql:host=$host;port=$port;dbname=$db;charset=utf8mb4",
        $user,
        $pass,
        [
            PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION,
            PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC,
        ]
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
        | Retrieve MySQL column metadata
        |--------------------------------------------------------------------------
        */
        $metaStmt = $pdo->prepare("
            SELECT column_name, data_type, column_key
            FROM information_schema.columns
            WHERE table_schema = :db
              AND table_name   = :table
        ");
        $metaStmt->execute([
            "db" => $db,
            "table" => $table
        ]);

        $columnTypes = [];
        $primaryKey = 'id';
        foreach ($metaStmt->fetchAll() as $col) {
            $columnTypes[$col['COLUMN_NAME']] = $col['DATA_TYPE'];
            if ($col['COLUMN_KEY'] == 'PRI') {
                $primaryKey = $col['COLUMN_NAME'];
            }
        }

        if (!$columnTypes) {
            throw new Exception("Table '$table' not found in database.");
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
            'INSERT INTO `%s` (`%s`) VALUES (%s) ON DUPLICATE KEY UPDATE `' . $primaryKey . '` = `' . $primaryKey . '`',
            $table,
            implode('`,`', $columns),
            implode(',', $placeholders)
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
                elseif (in_array($type, ['int', 'bigint', 'smallint', 'mediumint', 'tinyint'])) {
                    $value = (int) $value;
                }
                // decimal / float
                elseif (in_array($type, ['decimal', 'numeric', 'float', 'double'])) {
                    $value = (float) $value;
                }
                // boolean (MySQL usually tinyint(1))
                elseif (in_array($type, ['boolean', 'tinyint'])) {
                    $value = in_array(strtolower((string) $value), ['1', 'true', 't', 'yes'], true) ? 1 : 0;
                }
                // timestamp / datetime
                elseif (in_array($type, ['timestamp', 'datetime'])) {
                    $date = new DateTime($value);
                    $value = $date->format('Y-m-d H:i:s');
                }

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
