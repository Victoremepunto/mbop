alter table registrations
    drop constraint display_name_unique;

-- changing it to be uniqe per org_id instead of globally
alter table registrations
    add constraint display_name_unique
        unique (display_name, org_id);
