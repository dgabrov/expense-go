CREATE TABLE session
(
    session_id VARCHAR(64)  NOT NULL PRIMARY KEY,
    token      VARCHAR(255) NOT NULL,
    user_id    VARCHAR(64)  NOT NULL,
    login      VARCHAR(64)  NOT NULL,
    name       VARCHAR(64)  NOT NULL DEFAULT '',
    expiry_dt  DATETIME     NOT NULL,
    expired_ind varchar(1) not null default 'N' -- *ind fields are either Y or N, and transpose to bool true if Y (case-insensitive) otherwise false
) CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci;

CREATE TABLE account
(
    account_id varchar(64) not null primary key,
    account_name varchar(64) not null,
    cash varchar(1) not null default 'N'
) CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci;

CREATE TABLE expense_group
(
    group_id varchar(64) not null primary key,
    group_name varchar(64) not null
) CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci;

create table expense_category (
    category_id varchar(64) not null primary key,
    group_id varchar(64) not null,
    category_name varchar(64) not null
) CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci;

alter table expense_category add constraint fk_category_group foreign key (group_id) references expense_group(group_id);

create table transaction (
    transaction_id varchar(64) not null primary key,
    name varchar(255) not null,
    transaction_date date not null,
    amount number(14,2) not null,
    account_id varchar(64) not null,
    category_id varchar(64) null
) CHARACTER SET utf8mb4 COLLATE utf8mb4_uca1400_ai_ci;

alter table transaction add constraint fk_transaction_account foreign key (account_id) references account(account_id);
alter table transaction add constraint fk_transaction_category foreign key (category_id) references expense_category(category_id);
