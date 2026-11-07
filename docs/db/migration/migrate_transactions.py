"""
Throwaway migration: olddata.transaction -> expense.transaction.

Mapping (old -> new):
  id               -> transaction_id
  name             -> name
  date_posted      -> transaction_date   (date part only; new column is DATE)
  amount           -> amount
  account_id       -> account_id         (same value; account rows come from migrate_sections_1_2.py)
  category_cd      -> category_id        (same value; expense_category.category_id is the full category_cd)

Old columns with no new equivalent (user_id, bank_id, account_type, transaction_type,
bank_transaction_id, comments, hash) are not migrated.

Requires: pip install pymysql
Run migrate_sections_1_2.py first (accounts and categories must exist: foreign keys).

Everything runs in ONE transaction on the new database: any error rolls back all inserts.
Re-running without clearing the expense.transaction table fails on duplicate keys by design.
"""

import pymysql

HOST = "localhost"
PORT = 3306
USER = "expense"
PASSWORD = "expense"
OLD_DB = "olddata"
NEW_DB = "expense"
BATCH_SIZE = 1000


def connect(database):
    return pymysql.connect(
        host=HOST,
        port=PORT,
        user=USER,
        password=PASSWORD,
        database=database,
        charset="utf8mb4",
        autocommit=False,
    )


def main():
    old = connect(OLD_DB)
    new = connect(NEW_DB)
    try:
        with old.cursor() as cur:
            cur.execute(
                "SELECT id, name, DATE(date_posted), amount, account_id, category_cd "
                "FROM transaction ORDER BY id"
            )
            rows = cur.fetchall()

        total = 0
        with new.cursor() as cur:
            for start in range(0, len(rows), BATCH_SIZE):
                batch = rows[start:start + BATCH_SIZE]
                cur.executemany(
                    "INSERT INTO transaction "
                    "(transaction_id, name, transaction_date, amount, account_id, category_id) "
                    "VALUES (%s, %s, %s, %s, %s, %s)",
                    batch,
                )
                total += len(batch)

        new.commit()
        print(f"Committed: transactions={total}")
    except Exception:
        new.rollback()
        print("Rolled back: no changes were written to the expense database.")
        raise
    finally:
        old.close()
        new.close()


if __name__ == "__main__":
    main()
