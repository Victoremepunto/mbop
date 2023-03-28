alter table registrations
    drop constraint display_name_unique;

alter table registrations
    add constraint display_name_unique
        unique (display_name);
