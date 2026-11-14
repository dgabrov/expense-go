create table account
(
    account_id   varchar(64)            not null
        primary key,
    account_name varchar(64)            not null,
    cash         varchar(1) default 'N' not null
)
    engine = InnoDB
    default charset = utf8mb4
    collate = utf8mb4_unicode_ci;

create table expense_group
(
    group_id   varchar(64) not null
        primary key,
    group_name varchar(64) not null
)
    engine = InnoDB
    default charset = utf8mb4
    collate = utf8mb4_unicode_ci;

create table expense_category
(
    category_id   varchar(64) not null
        primary key,
    group_id      varchar(64) not null,
    category_name varchar(64) not null,
    constraint fk_category_group
        foreign key (group_id) references expense_group (group_id)
)
    engine = InnoDB
    default charset = utf8mb4
    collate = utf8mb4_unicode_ci;

create table session
(
    session_id  varchar(64)             not null
        primary key,
    token       varchar(255)            not null,
    user_id     varchar(64)             not null,
    login       varchar(64)             not null,
    name        varchar(64) default ''  not null,
    expiry_dt   datetime                not null,
    expired_ind varchar(1)  default 'N' not null
)
    engine = InnoDB
    default charset = utf8mb4
    collate = utf8mb4_unicode_ci;

create table transaction
(
    transaction_id   varchar(64)    not null
        primary key,
    name             varchar(255)   not null,
    transaction_date date           not null,
    amount           decimal(14, 2) not null,
    account_id       varchar(64)    not null,
    category_id      varchar(64)    null,
    constraint fk_transaction_account
        foreign key (account_id) references account (account_id),
    constraint fk_transaction_category
        foreign key (category_id) references expense_category (category_id)
)
    engine = InnoDB
    default charset = utf8mb4
    collate = utf8mb4_unicode_ci;

