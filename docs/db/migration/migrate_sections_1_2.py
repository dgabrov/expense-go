"""
Throwaway migration: olddata -> expense (spec_data.md, sections 1 and 2).

  Section 1: accounts               (distinct transaction.account_id -> expense.account)
  Section 2: groups and categories  (olddata.category.category_cd -> expense.expense_group / expense_category)

Section 2 rules:
  "group1:category1" -> expense_group    (group_id='group1', group_name='Group1')        once per group
                        expense_category (category_id='group1:category1', group_id='group1',
                                          category_name='Category1')
  "donations"        -> expense_group    (group_id='none', group_name='None')            once
                        expense_category (category_id='donations', group_id='none',
                                          category_name='Donations')
  NULL / blank category_cd is skipped.

Requires: pip install pymysql
Run this BEFORE migrate_transactions.py (transaction rows reference account and expense_category).

Everything runs in ONE transaction on the new database: any error rolls back all inserts.
Re-running on a database that already has these rows fails on duplicate keys by design.
"""

import pymysql

HOST = "localhost"
PORT = 3306
USER = "expense"
PASSWORD = "expense"
OLD_DB = "olddata"
NEW_DB = "expense"

NO_GROUP_ID = "none"
NO_GROUP_NAME = "None"


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


def capitalize_first(text):
    # 'group1' -> 'Group1', 'personalcare' -> 'Personalcare'. Rest of the text is kept as-is.
    return text[:1].upper() + text[1:]


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


def build_groups_categories(codes):
    # Pure transformation, no database access.
    # Returns:
    #   groups     : dict group_id -> group_name (each group inserted once)
    #   categories : list of (category_id, group_id, category_name)
    groups = {}
    categories = []
    for code in codes:
        if ":" in code:
            group_id, sub = code.split(":", 1)
            groups.setdefault(group_id, capitalize_first(group_id))
            categories.append((code, group_id, capitalize_first(sub)))
        else:
            groups.setdefault(NO_GROUP_ID, NO_GROUP_NAME)
            categories.append((code, NO_GROUP_ID, capitalize_first(code)))
    return groups, categories


def migrate_groups_categories(old, new):
    # Section 2: groups and categories, driven by olddata.category.category_cd.
    # NULL or blank rows are skipped.
    with old.cursor() as cur:
        cur.execute("SELECT category_cd FROM category")
        codes = [
            row[0].strip()
            for row in cur.fetchall()
            if row[0] is not None and row[0].strip()
        ]

    groups, categories = build_groups_categories(codes)

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
