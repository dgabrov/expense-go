"""
Throwaway migration: olddata -> expense (spec_data.md, sections 1 and 2).

  Section 1: accounts       (distinct transaction.account_id -> expense.account)
  Section 2: groups/cats    (olddata.category -> expense.expense_group / expense_category)

Requires: pip install pymysql

Everything runs in ONE transaction on the new database: any error rolls back
all inserts. Re-running after a successful run fails on duplicate keys by design.
"""

import pymysql

HOST = "localhost"
PORT = 3306
USER = "expense"
PASSWORD = "expense"
OLD_DB = "olddata"
NEW_DB = "expense"


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


def migrate_accounts(old, new):
    # Section 1: accounts. The old DB has no account table, so take the distinct
    # account_id values used by transactions. account_name = account_id.
    # cash = 'Y' when the value is 'cash' (case-insensitive), otherwise 'N'.
    with old.cursor() as cur:
        cur.execute(
            "SELECT DISTINCT account_id FROM transaction WHERE account_id IS NOT NULL"
        )
        account_ids = [row[0] for row in cur.fetchall()]

    rows = []
    for account_id in account_ids:
        # account_name is the same value as account_id, so checking one checks both.
        is_cash = account_id.lower() == "cash"
        rows.append((account_id, account_id, "Y" if is_cash else "N"))

    with new.cursor() as cur:
        cur.executemany(
            "INSERT INTO account (account_id, account_name, cash) VALUES (%s, %s, %s)",
            rows,
        )
    return len(rows)


def migrate_groups_categories(old, new):
    # Section 2: groups and categories, driven by olddata.category.category_cd.
    # NULL or blank rows are skipped.
    #
    # With colon, e.g. "services:insurance":
    #   expense_group    (group_id='services', group_name='services')          - inserted once per group
    #   expense_category (category_id='services:insurance', group_id='services',
    #                     category_name='services:insurance')
    # category_id/category_name take the FULL category_cd, not the part after the colon,
    # because category_id is the primary key and the right-hand parts collide
    # (services:insurance / shelter:insurance / transportation:insurance).
    # Without colon, e.g. 'none': group and category are all 'none'.
    with old.cursor() as cur:
        cur.execute("SELECT category_cd FROM category")
        codes = [
            row[0].strip()
            for row in cur.fetchall()
            if row[0] is not None and row[0].strip()
        ]

    groups = {}          # group_id -> group_name (inserted once each)
    categories = []      # (category_id, group_id, category_name)
    for code in codes:
        if ":" in code:
            group = code.split(":", 1)[0]
        else:
            group = code  # no colon: group and category are both the whole value
        groups.setdefault(group, group)
        categories.append((code, group, code))

    with new.cursor() as cur:
        cur.executemany(
            "INSERT INTO expense_group (group_id, group_name) VALUES (%s, %s)",
            list(groups.items()),
        )
        cur.executemany(
            "INSERT INTO expense_category (category_id, group_id, category_name) "
            "VALUES (%s, %s, %s)",
            categories,
        )
    return len(groups), len(categories)


def main():
    old = connect(OLD_DB)
    new = connect(NEW_DB)
    try:
        n_accounts = migrate_accounts(old, new)
        n_groups, n_categories = migrate_groups_categories(old, new)
        new.commit()
        print(f"Committed: accounts={n_accounts}, groups={n_groups}, categories={n_categories}")
    except Exception:
        new.rollback()
        print("Rolled back: no changes were written to the expense database.")
        raise
    finally:
        old.close()
        new.close()


if __name__ == "__main__":
    main()
