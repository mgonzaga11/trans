#!/bin/sh

sqlite3 db1 <<EOF
CREATE TABLE IF NOT EXISTS matches (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  player1 TEXT NOT NULL,
  player2 TEXT NOT NULL,
  winner TEXT,
  duration INTEGER NOT NULL,
  created_at TEXT NOT NULL
);
EOF