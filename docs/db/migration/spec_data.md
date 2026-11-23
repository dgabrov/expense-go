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
    left of colon is the group, right of colon is category
    insert once into expense_group table,  group_id and group_name same value
    insert in expense_category group_id as presented above, category_id and category_name same value
    category_id and category_name take the FULL category_cd (e.g. "services:insurance"), NOT only the part right of the colon
    (reason: category_id is the primary key and the right-hand parts collide, e.g. services:insurance / shelter:insurance / transportation:insurance all end in "insurance")

    example: services:insurance -> expense_group (group_id='services', group_name='services'); expense_category (category_id='services:insurance', group_id='services', category_name='services:insurance')

    attention, group shall be inserted only once, of course
else 
    insert into expense_group the value in both group_id and group_name
    then insert same values for all three fields in one expense_category entry

    example: none is both group and category and means transaction is not considered for statistics; very well, you create entry in expense_group with both group_id and group_name with value 'none' and one entry in expense_category with all three fields with value 'none', one for group_id and then two values for category_id and category_name
end if

python script you create: same folder, and you hardcode everything as it will be a throwaway
