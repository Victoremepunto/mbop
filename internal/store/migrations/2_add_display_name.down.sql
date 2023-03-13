alter table public.registrations
    drop column if exists display_name;

alter table public.registrations
    drop constraint if exists display_name_unique;
