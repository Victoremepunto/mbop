create table if not exists public.allowlist(
    ip_block varchar not null,
    org_id varchar not null,
    created_at timestamp default now() not null
);

alter table if exists public.allowlist
    add constraint allowlist_unique_cidr
        primary key (ip_block);
