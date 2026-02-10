# PHP CSV to Database Importer

A simple PHP CLI tool to import CSV files into a database using **PDO**.
This repository supports **MySQL** and **PostgreSQL** with separate importer scripts.

## Features

- Import CSV files into database tables
- Supports:
  - ✅ MySQL
  - ✅ PostgreSQL
- Automatic column mapping based on CSV headers
- Uses database metadata to cast data types correctly
- Skips duplicate records safely
- Transaction-based import for better performance
- CLI-friendly (no web server required)

## Importer Scripts

| Script                  | Database        |
| ----------------------- | --------------- |
| `mysql-importer.php`    | MySQL / MariaDB |
| `postgres-importer.php` | PostgreSQL      |

## CSV file name format:

- Prefix number is optional and only used for ordering
- Table name is derived from the filename
  Example:

```
001-companies.csv -> companies table
```

## CSV Structure

- First row **must be the column headers**
- Column names **must match database column names**

## Example CSV

[sample.csv](sample.csv)

## CSV Folder Structure

Place CSV files in the correct folder depending on the database type:

```
project-root/
├── mysql-importer.php
├── postgres-importer.php
├── vendor/
├── .env
└── sql/
    ├── mysql/
    │   ├── 001-companies.csv
    │   └── 002-users.csv
    └── postgres/
        ├── 001-companies.csv
        └── 002-users.csv
```

- MySQL importer reads from: `sql/mysql/`
- PostgreSQL importer reads from: `sql/postgres/`

## Installation

### 1. Clone the repository

```sh
git clone https://github.com/gendutski/my-playground.git
cd php-csv-to-database-importer
```

### 2. Install dependencies

```sh
composer install
```

### 3. Create and configure .env

#### MySQL `.env` Example

```.env
MYSQL_HOST=127.0.0.1
MYSQL_PORT=3306
MYSQL_NAME=your_database_name
MYSQL_USER=your_username
MYSQL_PASS=your_password
```

#### PostgreSQL `.env` Example

```.env
POSTGRESQL_HOST=127.0.0.1
POSTGRESQL_PORT=5432
POSTGRESQL_NAME=your_database_name
POSTGRESQL_USER=your_username
POSTGRESQL_PASS=your_password
```

## Usage

### Import to MySQL

```sh
php mysql-importer.php
```

### Import to PostgreSQL

```sh
php postgres-importer.php
```

## Notes

- Tables must already exist in the database
- Duplicate records are ignored:
  - MySQL: INSERT IGNORE
  - PostgreSQL: ON CONFLICT DO NOTHING
- Empty CSV values are converted to NULL
- Import runs inside a transaction for each file
