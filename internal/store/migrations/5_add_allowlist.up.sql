create table if not exists public.allowlist(
    ip varchar not null,
    org_id varchar not null,
    created_at timestamp default now() not null
);
