I need to do the data migration from old database having a structure to new database
for this application for expense management. 

Parameters for old database: localhost:3306, mariadb, user expense, password expense, database olddata
ddl is in @docs/db/migration/old_ddl.sql

Parameters for new database: localhost:3306, mariadb, user expense, password expense, database expense
ddl is in @docs/db/db_ddl.sql

I need a python script that provides migration as follows:

**section 1**

1. Accounts
- there is no table accounts in olddata
- go in transactions and get distinct account_id
- each of these, insert in expense database account table same value in account_id and account_name
- if account_id or name is "cash" case insensitive, please put cash field to 'Y' otherwise do not populate it and will get value 'N' automatically

2. Groups and categories
- in olddata there is a table called category with only one field, category_cd
- there is one with empty value either null or empty string, that one will not be migrated
- there are two possibilities, the category_cd contains : (colon) or not

if contains colon
    example
    group1:category1 
    - one time only for that group1, you insert in expense_group entry with group_id group1 and group_name Group1
    in expense_category you insert as follows:
    group_id group1
    category_id group1:category1
    category_name Category1

    better I gave sample rather than wasting time with generic description
else 
    example: none
    enter in expense_group group_id none and group_name None
    enter in category categoryId none, group_id none and category_name None
end if

python script you create: same folder, and you hardcode everything as it will be a throwaway
